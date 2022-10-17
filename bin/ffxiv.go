package bin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type FfCharacter struct {
	CharUrl string `json:"CharUrl"`
	ImgUrl string	`json:"ImgUrl"`
	Title string 	`json:"Title"`
	JobImg string `json:"JobImg"`
	Level string `json:"Level"`
	GrandCompany string `json:"GrandCompany"`
}

func GetFfCharacter(dataCenter string, name string)(FfCharacter,error) {
	var FfChar FfCharacter
	var err error
	url := fmt.Sprintf("http://valentintahon.com/apiFfxiv?name=%s&world=%s", name, dataCenter)
	urlSanitized := strings.Replace(url, " ", "%2B",-1)
	resp, err := http.Get(urlSanitized)
	bodyByte,err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	err = json.Unmarshal(bodyByte, &FfChar)
	return FfChar, err
}
