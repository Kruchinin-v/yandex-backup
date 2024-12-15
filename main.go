package main

import (
	"fmt"
	"os"
	"yandex-backup/modules"
)

func main() {
	version := "1.0"
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-V" || arg == "--version" {
			fmt.Println("Version:", version)
			return
		}
	}
	var token string
	token = os.Getenv("YANDEX_TOKEN")
	dirFiles := os.Getenv("BACKUP_DIR")
	filePrefix := os.Getenv("FILE_PREFIX")
	backupDir := os.Getenv("YANDEX_DIR")
	notificationChatId := os.Getenv("NOTIFICATION_CHAT_ID")
	notificationBotToken := os.Getenv("NOTIFICATION_BOT_TOKEN")
	notificationSubjectLine := os.Getenv("NOTIFICATION_SUBJECT_LINE")
	if token == "" || dirFiles == "" {
		fmt.Println("Error: 48045")
		return
	}
	yandexDriveApiUrl := "https://cloud-api.yandex.net/v1/disk/resources"
	var yandexDriveClient = modules.YandexDrive{
		Token:                   token,
		DirFiles:                dirFiles,
		FilePrefix:              filePrefix,
		BackupDir:               backupDir,
		YandexDriveApiUrl:       yandexDriveApiUrl,
		NotificationChatId:      notificationChatId,
		NotificationBotToken:    notificationBotToken,
		NotificationSubjectLine: notificationSubjectLine,
	}
	yandexDriveClient.Backup()
}
