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
 Currency string `json: "currency"`
 Date time.Time `json: "date"`
 Title string `json: "title"`
 Volatility int `json: "volatility"`
 IsClose bool `json:"isClose"`
 IsDelete bool `json:"isDelete"`
 RelatedCurrency []relatedCurrency `json:"relatedCurrency"`
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
 gmt := now.Add(-9 * time.Hour)
 startGmt := time.Date(gmt.Year(), gmt.Month(), gmt.Day(), 0, 0, 0, 0, time.Local)
 endGmt := time.Date(gmt.Year(), gmt.Month(), gmt.Day(), 23, 59, 59, 0, time.Local)
 var f []vip
 q := datastore.
  NewQuery("VipData").
  Filter("currency =", symbol[3:]).
  Filter("volatility =", 3).
  Filter("date >", startGmt).
  Filter("date <", endGmt)
 _, err = dsClient.GetAll(ctx, q, &f)
 var b []vip
 q = datastore.
  NewQuery("VipData").
  Filter("currency =", symbol[:3]).
  Filter("volatility =", 3).
  Filter("date >", startGmt).
  Filter("date <", endGmt)
 _, err = dsClient.GetAll(ctx, q, &b)

 vp := append(f, b...)

 relatedCurrencyMap := map[string][]relatedCurrency{
  "EUR": {{"GBP"}, {"CHF"}},
  "GBP": {{"EUR"}, {"CHF"}},
  "CHF": {{"EUR"}, {"GBP"}},
 }
 for i:= 0; i < len(vp); i++ {
  vp[i].RelatedCurrency = relatedCurrencyMap[vp[i].Currency]
 }
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