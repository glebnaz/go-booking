package log_controller

import (
	"encoding/json"
	"fmt"
	"go-booking/model"
	"time"
)

func LogOrder(Order model.Order_Inter, CodePharm string, Answer string) {
	t := time.Now()
	byteBasket, err := json.Marshal(Order.Basket)
	if err != nil {
		fmt.Println("Trouble with Log Basket")
	}
	str := fmt.Sprintf("     ORDER\nTime: %v,\nCodeFarm: %v,\nName: %v,\nPhone: %s,\nComent: %v,\nStoreId: %v,\nBasket: %v,\nAnswer: %s\n     END ORDER", t.String(), CodePharm, Order.Name, Order.Phone, Order.Comment, Order.StoreId, string(byteBasket), Answer)
	fmt.Println(str)
}
