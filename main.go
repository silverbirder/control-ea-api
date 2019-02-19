package main

import (
	"log"
	"net/http"

	"github.com/Silver-birder/control-ea-api/p"
)

func main() {
	http.HandleFunc("/", p.GetVipData)
	http.HandleFunc("/search", p.SearchVipData)
	http.HandleFunc("/update", p.UpdateVipData)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
