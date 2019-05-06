package Orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-booking/configure"
	"go-booking/model"
	"go-booking/model/Tokens"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var Token366 Tokens.TokenApteka366

type SixAndSixOrder struct {
	BasketSixAndSix BasketSixAndSix
	BaseStoreId     string `json:"baseStoreId"`
	Name            string
	Phone           string
}

type BasketSixAndSix struct {
	CodeGood string `json:"code"`
	Quantity []int  `json:"qty"`
}

type ResponseFromGetCart struct {
	Code string `json:"code"`
}

func (s *SixAndSixOrder) GetCart() (string, error) {
	url, _, _ := configure.GetURLData(configure.Server.Code_Farm.APTEKA366_CODE_FARM)
	url = url + "cart"
	response, err := s.DoReq(url, http.MethodGet, nil)
	if err != nil {
		return "", err
	}
	var cartNumber ResponseFromGetCart
	json.Unmarshal(response, &cartNumber)
	return cartNumber.Code, nil
}

func (s *SixAndSixOrder) AddToCart(cartId string) ([]byte, error) {
	url, _, _ := configure.GetURLData(configure.Server.Code_Farm.APTEKA366_CODE_FARM)
	path := url + "cart/" + cartId + "/entries"
	var qtyString string
	for i := 0; i < len(s.BasketSixAndSix.Quantity); i++ {
		qtyString = qtyString + strconv.Itoa(s.BasketSixAndSix.Quantity[i]) + ","
	}
	path = path + "?code=" + s.BasketSixAndSix.CodeGood + "&qty=" + qtyString
	response, err := s.DoReq(path, http.MethodPost, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *SixAndSixOrder) BookOrder() ([]byte, error) {
	urlConf, _, _ := configure.GetURLData(configure.Server.Code_Farm.APTEKA366_CODE_FARM)
	cartId, err := s.GetCart()
	if err != nil {
		return nil, err
	}
	_, err = s.AddToCart(cartId)
	if err != nil {
		return nil, err
	}

	//path := url + "orders/booking?cartId=" + cartId + "&pointOfServiceName=" + s.BaseStoreId + "&customerName=" + s.Name + "&customerPhone=" + s.Phone

	//path := "https://366.ru/rest/v2/reservation/77-res/orders/booking/?cartId="+cartId+"&pointOfServiceName="+s.BaseStoreId+"&customerName="+s.Name+"&customerPhone="+s.Phone

	var Url *url.URL
	Url, err = url.Parse(urlConf)
	if err != nil {
		fmt.Println()
	}

	Url.Path += "/orders/booking/"
	parameters := url.Values{}
	parameters.Add("cartId", cartId)
	parameters.Add("pointOfServiceName", s.BaseStoreId)
	parameters.Add("customerName", s.Name)
	parameters.Add("customerPhone", s.Phone)
	Url.RawQuery = parameters.Encode()

	req, err := http.NewRequest("POST", Url.String(), nil)
	if err != nil {
		return nil, err
	}

	if !Token366.CheckToken() {
		Token366.GetToken()
	}

	bearer := "bearer " + Token366.Access_token
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if strings.Contains(string(body), "error") {
		errorToLog := errors.New(string(body))
		return []byte("Error from api 366."), errorToLog
	}
	type answ struct {
		Code string `json:"code"`
	}

	var code answ
	json.Unmarshal(body, &code)

	done := fmt.Sprintf("\nSUCCESS: %s; Code: %s\n", string(body), code)
	return []byte(done), nil
}

func (s *SixAndSixOrder) MapData(order model.Order_Inter) {
	s.BaseStoreId = order.StoreId
	s.Name = order.Name
	s.Phone = order.Phone
	var strGood string
	for _, v := range order.Basket {
		s.BasketSixAndSix.Quantity = append(s.BasketSixAndSix.Quantity, v.Qty)
		strGood = strGood + v.Product_id + ","
	}
	s.BasketSixAndSix.CodeGood = strGood
}

func (s *SixAndSixOrder) DoReq(url string, method string, body io.Reader) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	if !Token366.CheckToken() {
		Token366.GetToken()
	}

	bearer := "bearer " + Token366.Access_token
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	response, err := ioutil.ReadAll(res.Body)

	err = res.Body.Close()
	return response, nil
}

func (s SixAndSixOrder) GetInformation() string {
	str := fmt.Sprintf("Name: %s;	Phone: %s;BaseStore: %s;	GoodCode: %s;	Qty: %v", s.Name, s.Phone, s.BaseStoreId, s.BasketSixAndSix.CodeGood, s.BasketSixAndSix.Quantity)
	return str
}
