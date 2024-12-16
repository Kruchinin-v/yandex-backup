package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type YandexDrive struct {
	Token                   string
	DirFiles                string
	FilePrefix              string
	BackupDir               string
	YandexDriveApiUrl       string
	Client                  *http.Client
	NotificationChatId      string
	NotificationBotToken    string
	NotificationSubjectLine string
	NotificationEnabled     string
	NotificationDebug       string
	numberOfElements        int
	listFiles               []string
	uploadUrl               string
}

type yandexResponse struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (yd *YandexDrive) findFiles() {
	files, err := os.ReadDir(yd.DirFiles)
	if err != nil {
		yd.sendMessageAdmin(fmt.Sprintf("Failed to browse the directory with backups:\n"+
			"%s", yd.DirFiles), "false")
		fmt.Println("Error: 80261")
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		suffix := ".zip"
		if strings.HasPrefix(file.Name(), yd.FilePrefix) && strings.HasSuffix(file.Name(), suffix) {
			yd.listFiles = append(yd.listFiles, file.Name())
			yd.numberOfElements += 1
		}
	}
	if len(yd.listFiles) == 0 {
		yd.sendMessageAdmin(fmt.Sprintf("No files found for backup in folder %s",
			yd.BackupDir), "false")
	}
	if yd.NotificationDebug == "True" {
		yd.sendMessageAdmin(fmt.Sprintf("<b>Debug</b>\nfiles:\n%s", yd.listFiles), "true")
	}
}

func (yd *YandexDrive) Backup() {
	yd.sendMessageAdmin(fmt.Sprintf("Debug is %s", yd.NotificationDebug), "true")

	currentDate := time.Now().Format("2006_01_02")
	yd.FilePrefix = fmt.Sprintf("%s_%s", yd.FilePrefix, currentDate)
	yd.findFiles()
	yd.createBackupDir()
	yd.runBackup()
}

func (yd *YandexDrive) runBackup() {
	for _, fileName := range yd.listFiles {
		uploadUrl := yd.getUploadUrl(fileName)
		if yd.NotificationDebug == "True" {
			yd.sendMessageAdmin(fmt.Sprintf("<b>Debug</b>\nuploadUrl:\n%s", uploadUrl), "true")
		}
		yd.uploadFile(uploadUrl, fileName)
	}
}

func (yd *YandexDrive) createBackupDir() {
	url := fmt.Sprintf("%s/?path=app:/%s", yd.YandexDriveApiUrl, yd.BackupDir)
	responseCode, _ := yd.makeRequest(url, "GET")
	if responseCode == 200 {
		return
	}
	responseCode, responseBody := yd.makeRequest(url, "PUT")
	if responseCode != 201 {
		yd.sendMessageAdmin(fmt.Sprintf("<b>Error creating directory on yandex disk</b>\n"+
			"<b>Code:</b> %d\n"+
			"<b>Response:</b>\n%s", responseCode, string(responseBody)), "false")
		fmt.Printf("Error: 4315%d\n", responseCode)
		os.Exit(1)
	}
}

func (yd *YandexDrive) uploadFile(yr yandexResponse, fileName string) {
	data, err := os.ReadFile(path.Join(yd.DirFiles, fileName))
	if err != nil {
		yd.sendMessageAdmin(fmt.Sprintf("Error reading backup file"), "false")
		fmt.Println("Error: 47046")
		os.Exit(1)
	}
	req, err := http.NewRequest(yr.Method, yr.Href, bytes.NewBuffer(data))
	if err != nil {
		yd.sendMessageAdmin(fmt.Sprintf("Failed to create a request\n<b>Error:</b>%s", err),
			"false")
		fmt.Println("Error: 11940")
		os.Exit(1)
	}
	req.Header.Set("Authorization", "OAuth "+yd.Token)

	client := &http.Client{}
	if yd.NotificationDebug == "True" {
		yd.sendMessageAdmin(fmt.Sprintf("<b>Debug</b>\nStart upload file"), "true")
	}
	resp, err := client.Do(req)
	if yd.NotificationDebug == "True" {
		yd.sendMessageAdmin(fmt.Sprintf("<b>Debug</b>\nEnd upload file"), "true")
	}
	if err != nil {
		yd.sendMessageAdmin(fmt.Sprintf("Failed to request\n<b>Error:</b>%s", err),
			"false")
		fmt.Println("Error: 24946")
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 201 {
		yd.sendMessageAdmin(fmt.Sprintf("Failed to load the file\n<b>Code: </b>%d\n <b>Response:</b>%s",
			resp.StatusCode, string(body)),
			"false")
		fmt.Println("Error: 53076")
		os.Exit(1)
	}
}

func (yd *YandexDrive) getUploadUrl(fileName string) yandexResponse {
	url := fmt.Sprintf("%s/upload/?path=app:/%s/%s&overwrite=false", yd.YandexDriveApiUrl, yd.BackupDir, fileName)
	respCode, respBody := yd.makeRequest(url, "GET")
	if respCode == 409 {
		yd.sendMessageAdmin(fmt.Sprintf("File %s already exist in disk", fileName), "true")
		os.Exit(1)
	}
	if respCode != 200 {
		yd.sendMessageAdmin(fmt.Sprintf("Error getting uploadUrl.\n"+
			"<b>Code</b>: %d, \n<b>Response</b>:\n%s", respCode,
			string(respBody)), "false")
		fmt.Printf("Error: 4316%d\n", respCode)
		os.Exit(1)
	}
	var configObject yandexResponse
	err := json.Unmarshal(respBody, &configObject)
	check(err)
	return configObject
}

func (yd *YandexDrive) makeRequest(url string, method string) (int, []byte) {
	request := yd.createRequest(url, method)
	resp, err := yd.Client.Do(request)
	check(err)
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	check(err)

	return resp.StatusCode, responseBody
}

func (yd *YandexDrive) createRequest(url string, method string) *http.Request {
	var err error
	Request, err := http.NewRequest(method, url, nil)
	check(err)
	Request.Header.Add("Authorization", yd.Token)
	yd.Client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	return Request
}

func (yd *YandexDrive) sendMessageAdmin(message string, disableNotification string) {
	if yd.NotificationEnabled != "True" {
		return
	}
	payload := map[string]interface{}{
		"chat_id":              yd.NotificationChatId,
		"text":                 fmt.Sprintf("<b>%s</b>\n%s", yd.NotificationSubjectLine, message),
		"disable_notification": disableNotification,
		"parse_mode":           "HTML",
	}
	// Преобразуем тело в JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error: 73121")
		os.Exit(1)
	}
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", yd.NotificationBotToken)
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: 73122")
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: 73123")
		os.Exit(1)
	}
	//body, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(body))
	resp.Body.Close()
}
