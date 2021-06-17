package main

import (
	"encoding/json"
	"flag"
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/lhside/chrome-go"
	"github.com/mitchellh/go-ps"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

const VERSION = "2.0";

const BASE_PATH = `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`;
const APP_ID_FIREFOX = "{ec8030f7-c20a-464f-9b0e-13a3a9e97384}";
const APP_ID_THUNDERBIRD = "{3550f703-e582-4d05-9a08-453d09bdfdc6}";

type RequestParams struct {
	Id      string `json:id`
	Name    string `json:name`
	Version string `json:version`
	Creator string `json:creator`
}
type Request struct {
	Command string        `json:"command"`
	Params  RequestParams `json:"params"`
	Logging bool          `json:"logging"`
}

var appName = "";

func main() {
	shouldReportVersion := flag.Bool("v", false, "Output version")
	shouldListAddonIds := flag.Bool("l", false, "List registered addons")
	shouldCleanUp := flag.Bool("c", false, "Clean up registered addons")
	givenAppName := flag.String("a", "", "Name of the host application (Firefox or Thunderbird)")
	flag.Parse()
	if *givenAppName != "" {
		appName = *givenAppName
	} else {
		appName = GetAppName()
	}
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

	if request.Logging {
		logfileDir := os.ExpandEnv(`${temp}`)
		logRotationTime := time.Duration(24) * time.Hour
		maxAge := time.Duration(-1)
		rotateLog, err := rotatelogs.New(filepath.Join(logfileDir, "com.clear_code.list_addons_in_win_programs_we_host.log.%Y%m%d%H%M.txt"),
			rotatelogs.WithMaxAge(maxAge),
			rotatelogs.WithRotationTime(logRotationTime),
			rotatelogs.WithRotationCount(5),
		)
		if err != nil {
			log.Fatal(err)
		}
		defer rotateLog.Close()
		log.SetOutput(rotateLog)
		log.SetFlags(log.Ldate | log.Ltime)
	}
	log.Println("logging started");

	log.Println("command: " + request.Command);
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

func GetAppName() (appName string) {
	parentPID := os.Getppid()
	pidInfo, _ := ps.FindProcess(parentPID)
	if strings.Contains(pidInfo.Executable(), "firefox") {
		return "Firefox"
	}
	if strings.Contains(pidInfo.Executable(), "thunderbird") {
		return "Thunderbird"
	}
	return "Unknown"
}

func GetAppId() (appId string) {
	switch appName {
	case "Firefox":
		return APP_ID_FIREFOX
	case "Thunderbird":
		return APP_ID_THUNDERBIRD
	}
	return ""
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
	prefix := GetAppId() + "."
	parentKey, err := registry.OpenKey(registry.CURRENT_USER, BASE_PATH, registry.CREATE_SUB_KEY)
	if err != nil {
		return err.Error()
	}
	defer parentKey.Close()

	id := prefix + params.Id
	log.Println("register: id = " + id);
	addonKey, _, err := registry.CreateKey(parentKey, id, registry.SET_VALUE)
	if err != nil {
		return err.Error()
	}
	defer addonKey.Close()

	displayName := appName + ": " + params.Name
	log.Println("register : " + id + "/DisplayName = " + displayName);
	err = addonKey.SetStringValue("DisplayName", displayName)
	if err != nil {
		return err.Error()
	}
	log.Println("register : " + id + "/DisplayVersion = " + params.Version);
	err = addonKey.SetStringValue("DisplayVersion", params.Version)
	if err != nil {
		return err.Error()
	}
	appPath := GetAppPath()
	log.Println("register : " + id + "/UninstallString = " + appPath);
	err = addonKey.SetStringValue("UninstallString", appPath)
	if err != nil {
		return err.Error()
	}
	icon := appPath + ",0"
	log.Println("register : " + id + "/DisplayIcon = " + icon);
	err = addonKey.SetStringValue("DisplayIcon", icon)
	if err != nil {
		return err.Error()
	}
	log.Println("register : " + id + "/Publisher = " + params.Creator);
	err = addonKey.SetStringValue("Publisher", params.Creator)
	if err != nil {
		return err.Error()
	}
	log.Println("register: successfully registered " + id);

	return ""
}

func GetAppPath() (path string) {
	parentPID := os.Getppid()
	pidInfo, _ := ps.FindProcess(parentPID)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\` + pidInfo.Executable(),
		registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	path, _, err = key.GetStringValue("")
	if err != nil {
		log.Fatal(err)
	}
	return
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
	log.Println("unregister: id = " + prefix + id);
	err := registry.DeleteKey(registry.CURRENT_USER, BASE_PATH + `\` + prefix + id)
	if err != nil {
		return err.Error()
	}
	log.Println("successfully unregistered " + id);
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
	log.Println("get registered ids under " + BASE_PATH);
	parentKey, err := registry.OpenKey(registry.CURRENT_USER, BASE_PATH, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Fatal(err)
	}
	defer parentKey.Close()

	subkeyNames, err := parentKey.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatal(err)
	}

	prefix := GetAppId() + "."
	for _, name := range subkeyNames {
		log.Println("  name = " + name);
		if (strings.HasPrefix(name, prefix)) {
			id := strings.Replace(name, prefix, "", 1)
			log.Println("  => id = " + id);
			ids = append(ids, id)
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
