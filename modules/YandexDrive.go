package modules

import "fmt"

type YandexDrive struct {
	Token     string
	uploadUrl string
}

func (yd *YandexDrive) Backup() {
	fmt.Println(yd.Token)
}
