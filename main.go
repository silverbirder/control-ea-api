package main

import (
	"log"
	"net/http"

	"github.com/Silver-birder/control-ea-api/p"
)

func main() {
	http.HandleFunc("/GetVipData", p.GetVipData)
	http.HandleFunc("/SearchVipData", p.SearchVipData)
	http.HandleFunc("/UpdateVipData", p.UpdateVipData)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
