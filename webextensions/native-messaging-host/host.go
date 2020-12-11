package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lhside/chrome-go"
	"log"
	"os"
)

const VERSION = "1.99";

type RequestParams struct {
	Id      string `json:id`
	Name    string `json:name`
	Version string `json:version`
}
type Request struct {
	Command string        `json:"command"`
	Params  RequestParams `json:"params"`
}

func main() {
	shouldReportVersion := flag.Bool("v", false, "v")
	flag.Parse()
	if *shouldReportVersion == true {
		fmt.Println(VERSION)
		return
	}

	log.SetOutput(os.Stderr)

	rawRequest, err := chrome.Receive(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	request := &Request{}
	if err := json.Unmarshal(rawRequest, request); err != nil {
		log.Fatal(err)
	}

	switch command := request.Command; command {
	case "register-addon":
		RegisterAndRespond(request.Params)
	case "unregister-addon":
		UnregisterAndRespond(request.Params)
	case "list-registered-addons":
		GetRegisteredIdsAndRespond()
	case "cleanup":
		CleanUpAndRespond()
	default: // just echo
		err = chrome.Post(rawRequest, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type BasicResponse struct {
	Error string `json:"error"`
}

func RegisterAndRespond(params RequestParams) {
	errorMessage := Register(params)
	response := &BasicResponse{errorMessage}
	body, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	err = chrome.Post(body, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func Register(params RequestParams) (errorMessage string) {
	return ""
}

func UnregisterAndRespond(params RequestParams) {
	errorMessage := Unregister(params)
	response := &BasicResponse{errorMessage}
	body, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	err = chrome.Post(body, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func Unregister(params RequestParams) (errorMessage string) {
	return ""
}

type ListRegisteredAddonsResponse struct {
	Ids   []string `json:"ids"`
	Error string   `json:"error"`
}

func GetRegisteredIdsAndRespond() {
	ids, errorMessage := GetRegisteredIds()
	response := &ListRegisteredAddonsResponse{ids, errorMessage}
	body, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	err = chrome.Post(body, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func GetRegisteredIds() (ids []string, errorMessage string) {
	return ids, ""
}

func CleanUpAndRespond() {
	errorMessage := CleanUp()
	response := &BasicResponse{errorMessage}
	body, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	err = chrome.Post(body, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func CleanUp() (errorMessage string) {
	return ""
}
