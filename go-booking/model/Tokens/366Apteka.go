package Tokens

import (
	"encoding/json"
	"fmt"
	"go-booking/configure"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var TimeTokenDie time.Time

var Client = http.Client{
	Timeout: time.Second * 10, // Maximum of 10 secs
}

type TokenApteka366 struct {
	Access_token  string `json:"access_token"`
	Token_Type    string `json:"token_type"`
	Refresh_token string `json:"refresh_token"`
	Expires_in    int    `json:"expires_in"`
	Scope         string `json:"Scope"`
}

func (t *TokenApteka366) GetToken() error {
	path := configure.Server.Apteka366.UrlToken + "?client_id=" + configure.Server.Apteka366.Client_Id +
		"&client_secret=" + configure.Server.Apteka366.Client_Secret + "&grant_type=" + configure.Server.Apteka366.Grant_Type + "&username=" + configure.Server.Apteka366.UserName + "&password=" + configure.Server.Apteka366.Password

	req, err := http.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	res, err := Client.Do(req)
	if err != nil {
		return err
	}

	timeGetToken := time.Now()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(response, t)
	if err != nil {
		return err
	}

	ex_string := "0h" + "0m" + strconv.Itoa(t.Expires_in) + "s"

	ex_duration, err := time.ParseDuration(ex_string)
	if err != nil {
		fmt.Println("Err when ParseDuration")
	}

	TimeTokenDie = timeGetToken.Add(ex_duration)
	LogToken(t)

	return nil
}

func (t *TokenApteka366) CheckToken() bool {
	return time.Now().Before(TimeTokenDie)
}

func LogToken(token *TokenApteka366) {
	timeToken := time.Now()
	str := fmt.Sprintf("Time: %v;	 Token: %s; Express_In: %v;\n", timeToken, token.Access_token, token.Expires_in)
	fmt.Println(str)
}
