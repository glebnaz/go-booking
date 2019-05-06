package Orders

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-booking/configure"
	"go-booking/model"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type RostaOrder struct {
	Store_Id     string        `json:"storeId"`    //idСклада
	Ext_Id       string        `json:"extId"`      //КодЗаказаИсточника
	Ext_Date     string        `json:"extDate"`    //ДатаЗаказа
	Name_Client  string        `json:"clientName"` //имя клиента
	Comment      string        `json:"clientComment":`
	Phone_Client string        `json:"clientTel"`
	BasketArr    []BasketRosta `json:"basket"`
}

type BasketRosta struct {
	GoodsId  string `json:"goodsId"`
	Quantity string `json:"quantity"` //количество
}

type answer struct {
	ErrorsCode int    `json:"errorsCode"`
	Message    string `json:"message"`
}

var Client = http.Client{
	Timeout: time.Second * 10, // Maximum of 10 secs
}

func (r *RostaOrder) BookOrder() ([]byte, error) {
	response, err := r.DoReq("", http.MethodPost, nil)
	if err != nil {
		return nil, err
	}
	var answ answer

	err = json.Unmarshal(response, &answ)
	if answ.ErrorsCode != 0 {
		msg := fmt.Sprintf("Error from Api rigla. Message: %s", answ.Message)
		errorFromApi := errors.New(msg)
		return response, errorFromApi
	}

	return []byte("Success"), nil
}

func (r *RostaOrder) MapData(order model.Order_Inter){
	r.Store_Id = order.StoreId
	t1 := time.Now()
	strTìme, _ := t1.MarshalText()
	r.Ext_Date = string(strTìme)
	r.Ext_Id = "1"
	r.Name_Client = order.Name
	r.Phone_Client = order.Phone
	r.Comment = order.Comment
	for _, v := range order.Basket {
		var BasketNew BasketRosta
		BasketNew.GoodsId = v.Product_id
		BasketNew.Quantity = strconv.Itoa(v.Qty)
		r.BasketArr = append(r.BasketArr, BasketNew)
	}
}

func (r *RostaOrder) DoReq(url, method string, body io.Reader) ([]byte, error) {
	url, usr, pass := configure.GetURLData(configure.Server.Code_Farm.ROSTA_CODE_FARM)
	byteOrder, err := json.Marshal(r)
	body = bytes.NewReader(byteOrder)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.SetBasicAuth(usr, pass)

	res, err := Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	answer, err := ioutil.ReadAll(res.Body)

	return answer, nil
}

func (r RostaOrder) GetInformation() string {
	var str string
	str = fmt.Sprintf("Name: %s;	Phone %s;	Basket: %v", r.Name_Client, r.Phone_Client, r.BasketArr)
	return str
}
