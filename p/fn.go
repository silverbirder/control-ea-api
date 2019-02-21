package p

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
)

type saveVip struct {
	Currency   string    `json: "currency"`
	Date       time.Time `json: "date"`
	Title      string    `json: "title"`
	Volatility int       `json: "volatility"`
	IsClose    bool      `json:"isClose"`
	IsDelete   bool      `json:"isDelete"`
}

type vip struct {
	Currency        string            `json: "currency"`
	Date            time.Time         `json: "date"`
	Title           string            `json: "title"`
	Volatility      int               `json: "volatility"`
	IsClose         bool              `json:"isClose"`
	IsDelete        bool              `json:"isDelete"`
	RelatedCurrency []relatedCurrency `json:"relatedCurrency"`
	StartDateTime   string            `json:"startDateTime"`
	EndDateTime     string            `json:"endDateTime"`
	ID              int64             `json:"id"`
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
	filterDate := time.Now()
	date, ok := r.URL.Query()["date"]
	if ok {
		filterDate, _ = time.Parse("200601", date[0])
	}
	startJp := time.Date(filterDate.Year(), filterDate.Month(), filterDate.Day(), 0, 0, 0, 0, time.UTC)
	adOneMonthTim := filterDate.AddDate(0, 1, 0)
	endJp := time.Date(adOneMonthTim.Year(), adOneMonthTim.Month(), adOneMonthTim.Day(), 0, 0, 0, 0, time.UTC)
	q := datastore.
		NewQuery("FxVipData").
		Filter("Date >=", startJp).
		Filter("Date <", endJp).
		Order("Date")
	keys, err := dsClient.GetAll(ctx, q, &vp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if vp == nil {
		vp = []vip{}
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
	// for dev
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(res)
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
		NewQuery("FxVipData").
		Filter("Currency =", symbol[3:]).
		Filter("Volatility =", 3).
		Filter("Date >=", startJp).
		Filter("Date <=", endJp).
		Order("Date")
	_, err = dsClient.GetAll(ctx, q, &f)
	var b []vip
	q = datastore.
		NewQuery("FxVipData").
		Filter("Currency =", symbol[:3]).
		Filter("Volatility =", 3).
		Filter("Date >", startJp).
		Filter("Date <", endJp).
		Order("Date")
	_, err = dsClient.GetAll(ctx, q, &b)
	vp := append(f, b...)
	if vp == nil {
		vp = []vip{}
	}
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

func UpdateVipData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "not post", http.StatusInternalServerError)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	currency := r.Form.Get("currency")
	title := r.Form.Get("title")
	date, _ := time.Parse("2006-01-02 15:04:05", r.Form.Get("date"))
	volatility, _ := strconv.Atoi(r.Form.Get("volatility"))
	isClose, _ := strconv.ParseBool(r.Form.Get("isClose"))
	isDelete, _ := strconv.ParseBool(r.Form.Get("isDelete"))
	v := saveVip{
		Currency:   currency,
		Title:      title,
		Volatility: volatility,
		IsClose:    isClose,
		IsDelete:   isDelete,
		Date:       date,
	}
	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, "ma-web-tools")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	k := datastore.IDKey("FxVipData", id, nil)
	if _, err = dsClient.Put(ctx, k, &v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// for dev
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
