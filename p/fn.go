package p

import (
	"context"
	"encoding/json"
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
	ID int64 `json:"id"`
}

type response struct {
	Status string `json:"status"`
	Vips   []vip  `json:"vips"`
}

type relatedCurrency struct {
	Currency string `json:"currency"`
}

var relatedCurrencyMap = map[string][]relatedCurrency{
"EUR": {{"GBP"}, {"CHF"}},
"GBP": {{"EUR"}, {"CHF"}},
"CHF": {{"EUR"}, {"GBP"}},
}

var dangerZonePeriod = 1

func SearchVipData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "not get", http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, "ma-web-tools")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var vp []vip
	jst, _ := time.LoadLocation("Asia/Tokyo")
	jp := time.Now()
	date, ok := r.URL.Query()["date"]
	if ok {
		jp, _ = time.Parse("200601", date[0])
	}
	startJp := time.Date(jp.Year(), jp.Month(), jp.Day(), 0, 0, 0, 0, jst)
	adOneMonthTim := jp.AddDate(0, 1, 0)
	endJp := time.Date(adOneMonthTim.Year(), adOneMonthTim.Month(), adOneMonthTim.Day(), 0, 0, 0, 0, jst)
	q := datastore.
		NewQuery("VipData").
		Filter("date >=", startJp).
		Filter("date <", endJp)
	keys, err := dsClient.GetAll(ctx, q, &vp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(vp); i++ {
		vp[i].RelatedCurrency = relatedCurrencyMap[vp[i].Currency]
		if vp[i].RelatedCurrency == nil {
			vp[i].RelatedCurrency = []relatedCurrency{}
		}
		vp[i].StartDateTime = vp[i].Date.Add(time.Duration(-1*dangerZonePeriod) * time.Hour).Format("2006.01.02 15:04")
		vp[i].EndDateTime = vp[i].Date.Add(time.Duration(dangerZonePeriod) * time.Hour).Format("2006.01.02 15:04")
		vp[i].ID = keys[i].ID
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
func UpdateVipData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "not get", http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	_, err := datastore.NewClient(ctx, "ma-web-tools")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//if _, err = dsClient.Put(ctx, k, &e); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		Filter("date >=", startJp).
		Filter("date <=", endJp)
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

	for i := 0; i < len(vp); i++ {
		vp[i].RelatedCurrency = relatedCurrencyMap[vp[i].Currency]
		if vp[i].RelatedCurrency == nil {
			vp[i].RelatedCurrency = []relatedCurrency{}
		}
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
