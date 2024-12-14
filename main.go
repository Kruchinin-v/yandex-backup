package main

import (
	"os"
	"yandex-backup/modules"
)

func main() {
	var Token string
	Token = os.Getenv("YANDEX_TOKEN")
	var yandexDriveClient = modules.YandexDrive{Token: Token}
	yandexDriveClient.Backup()
}
