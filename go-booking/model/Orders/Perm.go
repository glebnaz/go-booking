package Orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-booking/configure"
	"go-booking/model"
	"io"
	"io/ioutil"
	"net/http"
)

type PermFarmOrder struct {
	Fio      string `json:"Fio"`
	Store_Id int    `json:"Idapt"`
	Duration int    //срок резервирования в часах
	Phone    string
	Basket   []GoodPermFarm `json:"Goods"`
}

type GoodPermFarm struct {
	Idpos  int     `json:"Idpos"` //Код ГЕС товара
	Amount int     `json:"Amount"`
	Cost   float64 `json:"Cost"`
}

type PermFarmBody struct {
	Ver     string
	Command Command
}

type Command struct {
	Name   string
	Params ParamsPerm
}

type ParamsPerm struct {
	Reserve []Reserve
}

type Reserve struct {
	Source   string
	Number   string
	Idapt    int
	Date     string
	Duration int
	Fio      string
	Phone    string
	Goods    []GoodPermFarm
}

func (p *PermFarmOrder) BookOrder() ([]byte, error) {
	url, _, _ := configure.GetURLData(configure.Server.Code_Farm.PERM_CODE_FARM)
	//Date:= time.Now().String()
	response, err := p.DoReq(url, http.MethodPost, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (p *PermFarmOrder) MapData(order model.Order_Inter){
	//err := c.BindJSON(&p)
	//if err != nil {
		//return err
	//}
}

func (p *PermFarmOrder) DoReq(url, method string, body io.Reader) ([]byte, error) {
	ReserveExl := Reserve{
		Source:   "test",
		Number:   "",
		Idapt:    p.Store_Id,
		Date:     "2017-01-12T10:01:01", //TODO нормальную дату
		Duration: p.Duration,
		Fio:      p.Fio,
		Phone:    p.Phone,
		Goods:    p.Basket,
	}
	var ReserveBody []Reserve
	ReserveBody = append(ReserveBody, ReserveExl)
	b := PermFarmBody{
		Ver: "1.0",
		Command: Command{
			Name: "NewReserve",
			Params: ParamsPerm{
				Reserve: ReserveBody,
			},
		},
	}
	byteBody, err := json.Marshal(b)
	body = bytes.NewReader(byteBody)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	bodyResponse := res.Body
	response, err := ioutil.ReadAll(bodyResponse)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (p *PermFarmOrder) GetInformation() string {
	var str string
	str = fmt.Sprintf("Name: %s;	Phone: %s; Basket: %v", p.Fio, p.Phone, p.Basket)
	return str
}
