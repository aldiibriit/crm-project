package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type CreateBRIVARequest struct {
	InstitutionCode string    `json:"institutionCode"`
	BrivaNo         int       `json:"brivaNo"`
	CustCode        string    `json:"custCode"`
	Nama            string    `json:"nama"`
	Amount          int       `json:"amount"`
	Keterangan      string    `json:"keterangan"`
	DateTime        time.Time `json:"expiredDate"`
}

func SignatureBRIVA() string {

	brivaRequest := CreateBRIVARequest{
		InstitutionCode: "777",
		BrivaNo:         9192,
		CustCode:        "1",
		Nama:            "Aldi Saputra",
		Amount:          20000,
		Keterangan:      "Testing create Briva",
		DateTime:        time.Now(),
	}

	jsonData, err := json.Marshal(brivaRequest)
	if err != nil {
		fmt.Println(err.Error())
	}

	// fmt.Println("JSON STRING", string(jsonData))

	currentTS := time.Now().Format(time.RFC3339)
	path := "path=" + "/v1/briva"
	verb := "&verb=" + "POST"
	token := "&token=" + "Bearer hi61CyLko5GFYVAFAXIhWn2DxqGS"
	timestamp := "&timestamp=" + currentTS
	body := "&body=" + string(jsonData)
	signingString := path + verb + token + timestamp + body

	digest := hmac.New(sha256.New, []byte("AtxaoZ0DnrbiTuPY"))

	digest.Write([]byte(signingString))

	signature := base64.StdEncoding.EncodeToString(digest.Sum(nil))

	fmt.Println(signature)

	return currentTS
}
