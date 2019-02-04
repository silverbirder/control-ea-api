package p

import (
 "context"
 "encoding/json"
 "fmt"
 "time"
 "net/http"

 "cloud.google.com/go/datastore"
)

type vip struct {
 StartDateTime string `json:"startDateTime"`
 EndDateTime string `json:"endDateTime"`
 Currency string `json:"currency"`
 RelatedCurrency []relatedCurrency `json:"relatedCurrency"`
 Title string `json:"title"`
 IsClose bool `json:"isClose"`
 IsDelete bool `json:"isDelete"`
}

type dataStore struct {
 Currency string `json: "currency"`
 Date time.Time `json: "date"`
 Title string `json: "title"`
 Volatility int `json: "volatility"`
}

type response struct {
 Status string `json:"status"`
 Vips []vip `json:"vips"`
}

type relatedCurrency struct {
 Currency string `json:"currency"`
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

 ctx := context.Background()
 dsClient, err := datastore.NewClient(ctx, "ma-web-tools")
 if err != nil {
  fmt.Println(err)
 }
 now := time.Now()
 var f []dataStore
 q := datastore.
  NewQuery("VipData").
  Filter("currency =", symbol[3:]).
  Filter("volatility =", 3).
  Filter("date >", now).
  Filter("date <", now.Add(24 * time.Hour))
 _, err = dsClient.GetAll(ctx, q, &f)
 if err != nil {
  fmt.Println(err)
 }
 var b []dataStore
 q = datastore.
  NewQuery("VipData").
  Filter("currency =", symbol[:3]).
  Filter("volatility =", 3).
  Filter("date >", now).
  Filter("date <", now.Add(24 * time.Hour))
 _, err = dsClient.GetAll(ctx, q, &b)
 if err != nil {
  fmt.Println(err)
 }
 res, err := json.Marshal(append(f, b...))

 w.Write(res)
 return

 v1 := vip{
  StartDateTime: "2019.02.01 00:00",
  EndDateTime: "2019.02.01 02:00", // TODO: 時差を考慮. 日本時間 -7時間→世界標準時刻
  Currency: "USD",
  RelatedCurrency: []relatedCurrency{{Currency:"CHF"}},
  Title: "big news, " + symbol,
  IsClose:true,
  IsDelete:true,
 }
 v2 := vip{
  StartDateTime: "2019.02.01 03:00",
  EndDateTime: "2019.02.01 05:00",
  Currency: "USD",
  RelatedCurrency: []relatedCurrency{{Currency:"USD"}, {Currency:"CHF"}, {Currency:"GBP"}},
  Title: "big news2, " + symbol[3:],
  IsClose:true,
  IsDelete:false,
 }
 v3 := vip{
  StartDateTime: "2019.01.28 15:00",
  EndDateTime: "2019.01.28 20:00",
  Currency: "USD",
  RelatedCurrency: []relatedCurrency{{Currency:"EUR"}, {Currency:"JPY"}},
  Title: "big news3, " + symbol[:3],
  IsClose:false,
  IsDelete:true,
 }
 vp := []vip{v1, v2, v3,}
 response := response{
  Status:"ok",
  Vips:vp,
 }
 res, err = json.Marshal(response)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 w.Header().Set("Content-Type", "application/json")
 w.Write(res)
}