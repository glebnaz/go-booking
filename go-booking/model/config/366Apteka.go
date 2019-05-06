package config

import (
	"net/http"
	"time"
)

var Client = http.Client{
	Timeout: time.Second * 10, // Maximum of 10 secs
}

type Apteka366 struct {
	Client_Id     string
	Client_Secret string
	Grant_Type    string
	UserName      string
	Password      string
	UrlToken      string
	Url           string
	Url_Test      string
}

func (a Apteka366) GetUrlData(MockMode bool) string {
	if MockMode {
		return a.Url_Test
	} else {
		return a.Url_Test
	}
}
