package methodHttp

import (
	"encoding/json"
	"fmt"
	"go_test/handlerDb"
	"log"
	"net/http"
)

func State(output http.ResponseWriter, request *http.Request) {
	type responseJson struct {
		Result bool   `json:"result"`
		Info   string `json:"info"`
	}
	output.Header().Set("Content-Type", "application/json")
	updateFlagResult, _ := handlerDb.GetFlag(updateFlag)
	switch updateFlagResult {
	case updateEmpty: // no update process
		jsonText, err = json.Marshal(responseJson{false, "empty"})
	case updateRunning: // update is going
		jsonText, err = json.Marshal(responseJson{false, "updating"})
	case updateCompleted: // update completed
		jsonText, err = json.Marshal(responseJson{true, "ok"})
	case updateFailed: // update failed
		jsonText, err = json.Marshal(responseJson{false, "error"})
	default:
		jsonText, err = json.Marshal(responseJson{false, "unknown flag value"})
		log.Printf("Unknown updateFlag: %d\n", updateFlagResult)
	}
	if err != nil {
		log.Println(err)
	}
	_, err = fmt.Fprint(output, string(jsonText))
	if err != nil {
		log.Println(err)
	}
}
