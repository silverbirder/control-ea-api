package p

import (
 "encoding/json"
 "net/http"
)

type Currency string

const (
 USD Currency = "USD"
 CHF Currency = "CHF"
 EUR Currency = "EUR"
 GBP Currency = "GBP"
 JPY Currency = "JPY"
)

type vip struct {
 StartDateTime string `json:"startDateTime"`
 EndDateTime string `json:"endDateTime"`
 Currency Currency `json:"currency"`
 RelatedCurrency []relatedCurrency `json:"relatedCurrency"`
 Title string `json:"title"`
 IsClose bool `json:"isClose"`
 IsDelete bool `json:"isDelete"`
}

type response struct {
 Status string `json:"status"`
 Vips []vip `json:"vips"`
}

type relatedCurrency struct {
 Currency Currency `json:"currency"`
}

func GetVipData(w http.ResponseWriter, r *http.Request) {
 if r.Method != "POST" {
  http.Error(w, "not post", http.StatusInternalServerError)
  return
 }
 if err := r.ParseForm(); err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 symbol := r.FormValue("sy")
 v1 := vip{
  StartDateTime: "2019.02.01 00:00",
  EndDateTime: "2019.02.01 02:00", // TODO: 時差を考慮. 日本時間 -7時間→世界標準時刻
  Currency: USD,
  RelatedCurrency: []relatedCurrency{{Currency:CHF}},
  Title: "big news, " + symbol,
  IsClose:true,
  IsDelete:true,
 }
 v2 := vip{
  StartDateTime: "2019.02.01 03:00",
  EndDateTime: "2019.02.01 05:00",
  Currency: USD,
  RelatedCurrency: []relatedCurrency{{Currency:USD}, {Currency:CHF}, {Currency:GBP}},
  Title: "big news2, " + symbol[3:],
  IsClose:true,
  IsDelete:false,
 }
 v3 := vip{
  StartDateTime: "2019.01.28 15:00",
  EndDateTime: "2019.01.28 20:00",
  Currency: USD,
  RelatedCurrency: []relatedCurrency{{Currency:EUR}, {Currency:JPY}},
  Title: "big news3, " + symbol[:3],
  IsClose:false,
  IsDelete:true,
 }
 vp := []vip{v1, v2, v3,}
 response := response{
  Status:"ok",
  Vips:vp,
 }
 res, err := json.Marshal(response)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 w.Header().Set("Content-Type", "application/json")
 w.Write(res)
}