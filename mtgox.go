/*
	Package mtgox provides a streaming implementation of Mt. Gox's bitcoin trading API.
*/
package mtgox

import (
	"code.google.com/p/go.net/websocket"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"fmt"
	"io/ioutil"
)

const (
	api_host      string = "wss://websocket.mtgox.com:443"
	api_path      string = "/mtgox"
	http_endpoint string = "http://mtgox.com/api/2"
	origin_url    string = "http://localhost"
)

type StreamingApi struct {
	ws     *websocket.Conn
	key    []byte
	secret []byte
    Errors chan error
    Ticker chan Ticker
    Info chan Info
    Depth chan Depth
    Trade chan Trade
    Orders chan []Order
}

type Config struct {
	Currencies []string
	Key        string
	Secret     string
}

func NewFromConfig(cfgfile string) (*StreamingApi, error) {
	file, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return nil, err
	}

	m := new(Config)
	err = json.Unmarshal(file, m)
	if err != nil {
		return nil, err
	}

	return New(m.Key, m.Secret, m.Currencies...)
}

func New(key, secret string, currencies ...string) (*StreamingApi, error) {
	url := fmt.Sprintf("%s%s?Currency=%s", api_host, api_path, strings.Join(currencies, ","))
	config, _ := websocket.NewConfig(url, origin_url)
	ws, err := websocket.DialConfig(config)

	if err != nil {
		return nil, err
	}

	api := &StreamingApi{
		ws:         ws,
        Ticker: make(chan Ticker),
        Info: make(chan Info),
        Depth: make(chan Depth),
        Trade: make(chan Trade),
        Orders: make(chan []Order),
	}

	api.key, err = hex.DecodeString(strings.Replace(key, "-", "", -1))
	if err != nil {
		return nil, err
	}

	api.secret, err = base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	return api, err
}

func (api *StreamingApi) Start(){
    var obj []byte
    go func(){
        for ;;{
            _ = websocket.Message.Receive(api.ws, &obj)
            api.handle(obj)
        }
    }()
}

func (api *StreamingApi) Close() error {
	return api.ws.Close()
}

func (api *StreamingApi) sign(body []byte) ([]byte, error) {
	mac := hmac.New(sha512.New, api.secret)
	_, err := mac.Write(body)
	if err != nil {
		return nil, err
	}

	return mac.Sum(nil), nil
}

func (api *StreamingApi) authenticatedSend(msg map[string]interface{}) error {
	if api.key == nil || api.secret == nil {
		return errors.New("API Key or secret is invalid or missing.")
	}

	req, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	signedReq, err := api.sign(req)
	if err != nil {
		return err
	}

	reqid := msg["id"]

	fullReq := append(append(api.key, signedReq...), req...)
	encodedReq := base64.StdEncoding.EncodeToString(fullReq)

	return websocket.JSON.Send(api.ws,
        map[string]interface{}{
            "op":      "call",
            "id":      reqid,
            "call":    encodedReq,
            "context": "mtgox.com",
        })
}

func (api *StreamingApi) handle(obj []byte){
    msg := make(map[string]interface{})
    err := json.Unmarshal(obj, &msg)
    if err != nil{
        fmt.Println(err)
    }
    var msgtype string
    var bytemsg []byte
    if channel_name, ok := msg["channel_name"]; ok{
        bits := strings.Split(channel_name.(string), ".")
        msgtype = bits[0]
        submsg := msg[msgtype].(map[string]interface{})
        bytemsg, err = json.Marshal(submsg)
    }else{
        submsg, ok := msg["result"].(map[string]interface{})
        if ok{
            bytemsg, err = json.Marshal(submsg)
            msgtype = "info"
        }else{
            submsg, ok := msg["result"].([]interface{})
            if ok{
                bytemsg, err = json.Marshal(submsg)
            }else{
                bytemsg, err = json.Marshal([]byte("[]"))
            }
            msgtype = "orders"
        }
    }
    if err != nil{
        select{
        case api.Errors<-err:
        default:
        }
    }

    switch(msgtype){

    case "ticker":
        var ticker Ticker
        err = json.Unmarshal(bytemsg, &ticker)
        if err != nil{
            select{
            case api.Errors<-err:
            default:
            }
        }else{
            select{
            case api.Ticker<-ticker:
            default:
            }
        }

    case "depth":
        var depth Depth
        err = json.Unmarshal(bytemsg, &depth)
        if err != nil{
            select{
            case api.Errors<-err:
            default:
            }
        }else{
            select{
            case api.Depth<-depth:
            default:
            }
        }

    case "trade":
        var trade Trade
        err = json.Unmarshal(bytemsg, &trade)
        if err != nil{
            select{
            case api.Errors<-err:
            default:
            }
        }else{
            select{
            case api.Trade<-trade:
            default:
            }
        }
    case "info":
        var info Info
        err = json.Unmarshal(bytemsg, &info)
        if err != nil{
            select{
            case api.Errors<-err:
            default:
            }
        }else{
            select{
            case api.Info<-info:
            default:
            }
        }
    case "orders":
        var orders []Order
        err = json.Unmarshal(bytemsg, &orders)
        if err != nil{
            select{
            case api.Errors<-err:
            default:
            }
        }else{
            select{
            case api.Orders<-orders:
            default:
            }
        }

    default:
        fmt.Println("handling ", msgtype)

    }
}
func (api *StreamingApi) RequestInfo() (c chan Info){
    api.call("private/info", nil)
    return api.Info
}

func (api *StreamingApi) RequestOrders() (c chan []Order){
    api.call("private/orders", nil)
    return api.Orders
}

func (api *StreamingApi) call(endpoint string, params map[string]interface{}) error{
    if params == nil {
        params = make(map[string]interface{})
    }
    msg := map[string]interface{}{
        "call":   endpoint,
        "item":   "BTC",
        "params": params,
        "id":     <-ids,
        "nonce":  <-nonces,
    }

    err := api.authenticatedSend(msg)
    return err
}

func (api *StreamingApi) HandleErrors(f func(error)){
    go func(){
        for err := range(api.Errors){
            f(err)
        }
    }()
}
func (api *StreamingApi) SubmitOrder(typ string, amount, price int64) chan []Order {
    api.call("order/add", map[string]interface{}{
        "type":       typ,
        "amount_int": amount,
        "interfaceprice_int":  price,
    })
    return api.Orders
}
