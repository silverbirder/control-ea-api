package p

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
)

type vip struct {
	Currency        string            `json: "currency"`
	Date            time.Time         `json: "date"`
	Title           string            `json: "title"`
	Volatility      int               `json: "volatility"`
	IsClose         bool              `json:"isClose"`
	IsDelete        bool              `json:"isDelete"`
	RelatedCurrency []relatedCurrency `json:"relatedCurrency"`
	StartDateTime string `json:"startDateTime"`
	EndDateTime string `json:"endDateTime"`
}

type response struct {
	Status string `json:"status"`
	Vips   []vip  `json:"vips"`
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

	jst, _ := time.LoadLocation("Asia/Tokyo")
	jp := time.Now()
	startJp := time.Date(jp.Year(), jp.Month(), jp.Day(), 0, 0, 0, 0, jst)
	endJp := time.Date(jp.Year(), jp.Month(), jp.Day(), 23, 59, 59, 0, jst)
	var f []vip
	q := datastore.
		NewQuery("VipData").
		Filter("currency =", symbol[3:]).
		Filter("volatility =", 3).
		Filter("date >", startJp).
		Filter("date <", endJp)
	_, err = dsClient.GetAll(ctx, q, &f)
	var b []vip
	q = datastore.
		NewQuery("VipData").
		Filter("currency =", symbol[:3]).
		Filter("volatility =", 3).
		Filter("date >", startJp).
		Filter("date <", endJp)
	_, err = dsClient.GetAll(ctx, q, &b)
	vp := append(f, b...)

	relatedCurrencyMap := map[string][]relatedCurrency{
		"EUR": {{"GBP"}, {"CHF"}},
		"GBP": {{"EUR"}, {"CHF"}},
		"CHF": {{"EUR"}, {"GBP"}},
	}
	dangerZonePeriod := 1
	for i := 0; i < len(vp); i++ {
		vp[i].RelatedCurrency = relatedCurrencyMap[vp[i].Currency]
		vp[i].StartDateTime = vp[i].Date.Add(time.Duration(-1*dangerZonePeriod) * time.Hour).Format("2006.01.02 15:04")
		vp[i].EndDateTime = vp[i].Date.Add(time.Duration(dangerZonePeriod) * time.Hour).Format("2006.01.02 15:04")
	}
	response := response{
		Status: "ok",
		Vips:   vp,
	}

	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
