package mtgox

import (
    "time"
    "strconv"
)

type Trade struct{
    Type string
    Trade_type string
    Properties string
    Now time.Time
    Amount float64
    Amount_int int64
    Primary string
    Price float64
    Price_int int64
    Item string
    Price_currency string
}

type Ticker struct{
    Vol Value // Volume
    Item string
    High Value //Highest value
    Low Value // Lowest Value
    Last Value // == Last_local
    Last_local Value // Last trade in auxilary currency
    Last_all Value // Last trade converted to auxilary currency
    Last_orig Value // Last trade in any currency
    Buy Value
    Sell Value
    VWap Value // Volume weighted average price
    Avg Value // Averaged price
    Now EpochTime `json:",string"`
}

type Value struct{
    Value float64 `json:",string"`
    Value_int int64 `json:",string"`
    Display string
    Display_short string
    Currency string
}

type Info struct{
    Created SimpleTime `json:",string"`
    Id string
    Index string
    Language string
    Last_login SimpleTime `json:",string"`
    Link string
    Login string
    Montly_Volume Value
    Trade_fee float64
    Rights []string
    Wallets map[string]Wallet
}

type Wallet struct{
    Balance Value
    Daily_Withdraw_Limit Value
    Max_Withdraw Value
    // Monthly_Withdraw_Limit nil
    Open_Orders Value
    Operations int64
}

type Rate struct{
    To string
    From string
    Rate float64
}

type Depth struct{
    Type int64
    Type_str string
    Volume float64 `json:",string"`
    Volume_int int64 `json:",string"`
    Now EpochTime `json:",string"`
    Price float64 `json:",string"`
    Price_int int64 `json:",string"`
    Item string
    Currency string
    Total_volume_int int64 `json:",string"`
}

type Order struct{
    Oid string
    Currency string
    Item string
    Type string
    Amount Value
    Effective_amount Value
    Price Value
    Status string
    Date EpochTime
    Priority EpochTime `json:",string"`
    Actions []string
}

type EpochTime struct{
    time.Time
}

type SimpleTime struct{
    time.Time
}

func (t *EpochTime) UnmarshalJSON(b []byte) error {
    result, err := strconv.ParseInt(string(b), 0, 64)
    if err != nil {
        return err
    }
    // convert the unix epoch to a Time object
    *t = EpochTime{time.Unix(0, result * 1000)}
    return nil
}
func (t *SimpleTime) UnmarshalJSON(b []byte) error {
    layout := "2006-01-02 15:04:05"
    time, err := time.Parse(layout, string(b))
    if err != nil {
        return err
    }
    *t = SimpleTime{time}
    return nil

}

