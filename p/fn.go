package p

import (
"encoding/json"
"net/http"
)

type vip struct {
 StartDateTime string `json:"startDateTime"`
 EndDateTime string `json:"endDateTime"`
 Currency string `json:"currency"`
 Title string `json:"title"`
}

func GetVipData(w http.ResponseWriter, r *http.Request) {
 v1 := vip{"2019/02/01 00:00:00", "2019/02/01 02:00:00", "USD", "なんか大変そうな発表1",}
 v2 := vip{"2019/02/01 04:00:00", "2019/02/01 08:00:00", "USD", "なんか大変そうな発表2",}
 v3 := vip{"2019/02/01 22:00:00", "2019/02/01 23:00:00", "USD","なんか大変そうな発表3",}
 vp := []vip{v1, v2, v3,}
 res, err := json.Marshal(vp)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 w.Header().Set("Content-Type", "application/json")
 w.Write(res)
}