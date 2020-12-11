package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lhside/chrome-go"
	"github.com/mitchellh/go-ps"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const VERSION = "1.99";

const BASE_PATH = "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall";
const APP_ID_FIREFOX = "{ec8030f7-c20a-464f-9b0e-13a3a9e97384}";
const APP_ID_THUNDERBIRD = "{3550f703-e582-4d05-9a08-453d09bdfdc6}";

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
	shouldListAddonIds := flag.Bool("l", false, "l")
	shouldCleanUp := flag.Bool("c", false, "c")
	flag.Parse()
	if *shouldReportVersion == true {
		fmt.Println(VERSION)
		return
	}
	if *shouldListAddonIds == true {
		ids, errorMessage := GetRegisteredIds()
		if errorMessage != "" {
			fmt.Println(errorMessage)
		}
		for _, id := range ids {
			fmt.Println(id)
		}
		return
	}
	if *shouldCleanUp == true {
		CleanUp()
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

func GetAppId() (appId string) {
	parentPID := os.Getppid()
	pidInfo, _ := ps.FindProcess(parentPID)
	if strings.Contains(pidInfo.Executable(), "firefox") {
		return APP_ID_FIREFOX
	}
	if strings.Contains(pidInfo.Executable(), "thunderbird") {
		return APP_ID_THUNDERBIRD
	}
	return APP_ID_FIREFOX
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
	errorMessage := Unregister(params.Id)
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

func Unregister(id string) (errorMessage string) {
	prefix := GetAppId() + "."
	err := registry.DeleteKey(registry.CURRENT_USER, BASE_PATH + "\\" + prefix + id)
	if err != nil {
		return err.Error()
	}
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
	key, err := registry.OpenKey(registry.CURRENT_USER, BASE_PATH, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	subkeyNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatal(err)
	}

	prefix := GetAppId() + "."
	for _, name := range subkeyNames {
		if (strings.HasPrefix(name, prefix)) {
			ids = append(ids, strings.Replace(name, prefix, "", 1))
		}
	}
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
	ids, errorMessage := GetRegisteredIds()
	for _, id := range ids {
		errorMessage = Unregister(id)
		if errorMessage != "" {
			return errorMessage
		}
	}
	return ""
}
