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
	Token             string
	DirFiles          string
	FilePrefix        string
	BackupDir         string
	YandexDriveApiUrl string
	Client            *http.Client
	listFiles         []string
	uploadUrl         string
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
		fmt.Printf("Ошибка при чтении директории: %v\n", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), yd.FilePrefix) {
			yd.listFiles = append(yd.listFiles, file.Name())
		}
	}
}

func (yd *YandexDrive) Backup() {
	currentDate := time.Now().Format("2006_01_02")
	yd.FilePrefix = fmt.Sprintf("%s_%s", yd.FilePrefix, currentDate)
	yd.findFiles()
	yd.createBackupDir()
	yd.runBackup()
}

func (yd *YandexDrive) runBackup() {
	for _, fileName := range yd.listFiles {
		uploadUrl := yd.getUploadUrl(fileName)
		yd.uploadFile(uploadUrl, fileName)
	}
}

func (yd *YandexDrive) createBackupDir() {
	url := fmt.Sprintf("%s/?path=app:/%s", yd.YandexDriveApiUrl, yd.BackupDir)
	responseCode, _ := yd.makeRequest(url, "GET")
	if responseCode == 200 {
		return
	}
	responseCode, _ = yd.makeRequest(url, "PUT")
	if responseCode != 201 {
		fmt.Printf("Error: 4315%d\n", responseCode)
		os.Exit(1)
	}
}

func (yd *YandexDrive) uploadFile(yr yandexResponse, fileName string) {
	data, err := os.ReadFile(path.Join(yd.DirFiles, fileName))
	if err != nil {
		fmt.Println("Error: 47046")
		os.Exit(1)
	}
	req, err := http.NewRequest(yr.Method, yr.Href, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error: 11940")
		os.Exit(1)
	}
	req.Header.Set("Authorization", "OAuth "+yd.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: 24946")
		os.Exit(1)
	}
	defer resp.Body.Close()
	//body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != 201 {
		fmt.Println("Error: 53076")
		os.Exit(1)
	}
}

func (yd *YandexDrive) getUploadUrl(fileName string) yandexResponse {
	url := fmt.Sprintf("%s/upload/?path=app:/%s/%s&overwrite=false", yd.YandexDriveApiUrl, yd.BackupDir, fileName)
	respCode, respBody := yd.makeRequest(url, "GET")
	if respCode != 200 {
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
