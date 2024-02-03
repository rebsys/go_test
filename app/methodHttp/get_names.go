package methodHttp

import (
	"encoding/json"
	"fmt"
	"go_test/handlerDb"
	"go_test/share"
	"log"
	"net/http"
	"strings"
)

const (
	queryType       = "type"
	queryTypeStrong = "strong"
	queryTypeWeak   = "weak"
	queryName       = "name"
)

func GetNames(output http.ResponseWriter, request *http.Request) {
	output.Header().Set("Content-Type", "application/json")
	defer func() {
		if issue := recover(); issue != nil {
			output.WriteHeader(400)
			_, err = fmt.Fprintf(output, "Invalid request: %s", issue)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	if !request.URL.Query().Has(queryName) {
		panic("Missing mandatory parameter '" + queryName + "'")
	}
	var searchType string = queryTypeWeak
	if strings.ToLower(request.URL.Query().Get(queryType)) == queryTypeStrong {
		searchType = queryTypeStrong
	}

	var responseArray []share.JsonEntry
	if searchType == queryTypeStrong {
		responseArray, err = handlerDb.SearchStrong(request.URL.Query().Get(queryName))
	} else {
		responseArray, err = handlerDb.SearchWeak(request.URL.Query().Get(queryName))
	}
	if err != nil {
		log.Println(err)
	}
	if len(responseArray) > 0 {
		jsonText, err = json.Marshal(responseArray)
		_, err = fmt.Fprint(output, string(jsonText))
	} else {
		_, err = fmt.Fprint(output, "[]")
	}
	if err != nil {
		log.Println(err)
	}
}
