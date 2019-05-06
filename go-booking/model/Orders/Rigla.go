package Orders

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-booking/configure"
	"go-booking/controllers/log_controller"
	"go-booking/model"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var token string

type Order_Rigla struct {
	Name    string
	Phone   string
	Email   string
	Comment string
	StoreId int           `json:"store_id"`
	Basket  []model.Pharm `json:"basket"`
}

type body_Order struct {
	Order order `json:"order"`
}

type order struct {
	App_points string   `json:"app_point"`
	Customer   customer `json:"customer"`
	Comment    string   `json:"comment"`
	Shipment   shipment `json:"shipment"`
	Store_id   int      `json:"store_id"`
	Payment    payment  `json:"payment"`
}
type payment struct {
	PaymentMethod PaymentMethod `json:"paymentMethod"`
}

type PaymentMethod struct {
	Method          string  `json:"method"`
	Additional_data adddata `json:"additional_data"`
}

type adddata struct {
	Only_in_stock bool `json:"only_in_stock"`
}

type customer struct {
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email"`
	Telephone  string `json:"telephone"`
}

type shipment struct {
	Shipping_carrier_code string `json:"shipping_carrier_code"`
	Shipping_method_code  string `json:"shipping_method_code"`
}

type Bodyrigla struct {
	Cart CartItem `json:"cartItem"`
}

type CartItem struct {
	PharmId string `json:"sku"`
	Qty     int    `json:"qty"`
	CartId  string `json:"quote_id"`
}

func (r Order_Rigla) GetCart() (string, error) {
	url := configure.Server.Rigla.URL + "guest-carts"
	answer, err := r.DoReq(url, http.MethodPost, nil)
	if err != nil {
		return "", err
	}
	cartId := strings.Split(string(answer), `"`)

	return cartId[1], err
}

func (r Order_Rigla) AddToCart(cartId string) error {
	token = configure.Server.Rigla.Token
	url := configure.Server.Rigla.URL + "guest-carts/" + cartId + "/items"
	for _, u := range r.Basket {
		b := Bodyrigla{
			CartItem{
				PharmId: u.Product_id,
				Qty:     u.Qty,
				CartId:  cartId,
			},
		}
		byteBody, err := json.Marshal(b)

		if err != nil {
			fmt.Println(err)
			return err
		}
		body := bytes.NewReader(byteBody)
		_, err = r.DoReq(url, http.MethodPost, body)
		if err != nil {

			return err
		}
	}
	return nil
}

func (r *Order_Rigla) BookOrder() ([]byte, error) {

	token = configure.Server.Rigla.Token
	cartId, err := r.GetCart()
	if err != nil {
		return nil, err
	}
	err = r.AddToCart(string(cartId))
	if err != nil {
		return nil, err
	}
	url := configure.Server.Rigla.URL + "cart/" + cartId + "/place-order"
	store := fmt.Sprintf("%d", r.StoreId)
	b := body_Order{
		Order: order{
			App_points: "vseapteki.ru",
			Comment:    r.Comment,
			Customer: customer{
				r.Name,
				"-",
				"-",
				r.Email,
				r.Phone,
			},
			Store_id: 1,
			Shipment: shipment{
				"slpickup",
				store,
			},
			Payment: payment{
				PaymentMethod: PaymentMethod{
					Method: "checkmo",
					Additional_data: adddata{
						Only_in_stock: true,
					},
				},
			},
		},
	}

	byteBody, err := json.Marshal(b)

	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(byteBody)

	ans, err := r.DoReq(url, http.MethodPut, body)

	if strings.Contains(string(ans), "message") {
		errorNew := errors.New("Error from api rigla.")
		errorToLog := errors.New(string(ans))
		log_controller.LogErr(errorToLog)
		return []byte("Fail"), errorNew
	}

	return []byte("Success"), nil
}

func (r *Order_Rigla) MapData(order model.Order_Inter){
	r.Name = order.Name
	r.Phone = order.Phone
	r.Email = "rospharm@gmail.com"
	strIdInt, _:= strconv.Atoi(order.StoreId)
	r.StoreId = strIdInt
	r.Basket = order.Basket
	r.Comment = order.Comment
}

func (r *Order_Rigla) DoReq(url string, method string, body io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r Order_Rigla) GetInformation() string {
	str := fmt.Sprintf("Name: %s;	Phone: %s;BaseStore: %s;	Basket: %v;", r.Name, r.Phone, r.StoreId, r.Basket)
	return str
}

//todo нужен рефактор кода на мап дата
