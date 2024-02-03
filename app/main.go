package main

import (
	_ "github.com/lib/pq"
	"go_test/handlerDb"
	"go_test/methodHttp"
	"log"
	"net/http"
)

const (
	tempFile  = "/tmp/temp.xml"
	sourceUrl = "https://www.treasury.gov/ofac/downloads/sdn.xml"
)

func main() {
	err := handlerDb.Init_db()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/update", methodHttp.Update)
	http.HandleFunc("/state", methodHttp.State)
	http.HandleFunc("/get_names", methodHttp.GetNames)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
