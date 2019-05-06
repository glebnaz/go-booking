package test_version

import (
	"encoding/json"
	"fmt"
	"go-booking/configure"
	"go-booking/controllers"
	"go-booking/model"
	"io/ioutil"
	"os"
	"testing"
)

func Test_BookingOrderSucces(t *testing.T) {
	file, err := os.Open("config_test.json")
	byteConf, err := ioutil.ReadAll(file)
	var conf configure.ServerConfig
	err = json.Unmarshal(byteConf, &conf)
	if err != nil {
		fmt.Println(err)
	}
	configure.Server = conf
	fmt.Println(configure.Server.Code_Farm)
	err = Apteka366()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("Apteka366 Test: Succses")
	}
}

func Apteka366() error {
	orders, err := GetMockData("mock_data/mock_order_366.json")
	if err != nil {
		return err
	}
	fmt.Println(orders)
	for i, u := range orders {
		if u.Name!="Fail"{
			answer, _, err := controllers.BookingOrder(u, "80")
			if err != nil {
				fmt.Printf("Answer: %v,Index: %v\n", answer, i)
				return err
			}
		}else{
			answer, _, err := controllers.BookingOrder(u, "80")
			if err == nil {
				fmt.Printf("Answer: %v,Index: %v\n", answer, i)
				return err
			}
		}


	}
	return nil
}

func GetMockData(path string) ([]model.Order_Inter, error) {
	var Orders []model.Order_Inter
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	bodyMock, err := ioutil.ReadAll(file)

	err = json.Unmarshal(bodyMock, &Orders)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return Orders, err
}
