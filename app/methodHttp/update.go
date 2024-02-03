package methodHttp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go_test/handlerDb"
	"go_test/share"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Update(output http.ResponseWriter, request *http.Request) {
	type responseJson struct {
		Result bool   `json:"result"`
		Info   string `json:"info"`
		Code   int    `json:"code"`
	}
	output.Header().Set("Content-Type", "application/json")
	updateFlagResult, _ := handlerDb.GetFlag(updateFlag)
	switch updateFlagResult {
	case updateEmpty, updateCompleted: // run new updating
		err = handlerDb.SetFlag(updateFlag, updateRunning)
		if err != nil {
			log.Println(err)
		}
		go updater()
		jsonText, err = json.Marshal(responseJson{true, "", 200})
	case 1: // update is going
		jsonText, err = json.Marshal(responseJson{true, "", 200})
	case 3: // last update has failed, this service is blocking for flagTTL seconds
		jsonText, err = json.Marshal(responseJson{false, "service unavailable", 503})
	default:
		jsonText, err = json.Marshal(responseJson{false, "unknown flag value", 500})
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

func updater() {
	defer func() {
		if issue := recover(); issue != nil {
			log.Println(issue)
			err := handlerDb.SetFlag(updateFlag, updateFailed)
			if err != nil {
				log.Println(err)
			}
		} else {
			err := handlerDb.SetFlag(updateFlag, updateCompleted)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	sourcePrepare()
	parser()
}

func sourcePrepare() {
	cacheStat, err := os.Stat(share.CacheFile)
	if err == nil {
		cacheLastModified := cacheStat.ModTime().UTC().Format(time.RFC1123)
		client := &http.Client{}
		request, err := http.NewRequest("HEAD", share.SourceUrl, nil)
		if err != nil {
			log.Println(err)
		} else {
			request.Header.Set("If-Modified-Since", cacheLastModified)
			response, err := client.Do(request)
			if err != nil || response.StatusCode != http.StatusNotModified {
				download(share.CacheFile, share.SourceUrl)
			}
		}
	} else {
		download(share.CacheFile, share.SourceUrl)
	}
}

func download(filePath string, url string) {
	fileHandler, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer fileHandler.Close()
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		panic(fmt.Errorf("HTTP error: %s", response.Status))
	}
	_, err = io.Copy(fileHandler, response.Body)
	if err != nil {
		panic(err)
	}
}

func parser() {
	fileHandler, err := os.Open(share.CacheFile)
	if err != nil {
		panic(err)
	}
	var entry share.XmlEntry
	defer fileHandler.Close()
	decoder := xml.NewDecoder(fileHandler)
	for token, _ := decoder.Token(); token != nil; token, _ = decoder.Token() {
		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "sdnEntry" {
				err = decoder.DecodeElement(&entry, &element)
				if err != nil {
					fmt.Println(err)
				} else {
					if entry.SdnType == "Individual" {
						err = handlerDb.InsertEntry(&entry)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}
}
