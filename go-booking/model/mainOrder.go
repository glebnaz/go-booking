package model

import (
	"io"
)

type Order interface {
	BookOrder() ([]byte, error)
	MapData(order Order_Inter)
	DoReq(url string, method string, body io.Reader) ([]byte, error)
	GetInformation() string
}
