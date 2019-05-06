package controllers

import (
	"go-booking/configure"
	"go-booking/controllers/log_controller"
	"go-booking/model"
	"net/http"
)

func BookingOrder(proxy_order model.Order_Inter, id string) (string, int, error) {
	order, err := IdentificationCompany(id, configure.Server.Code_Farm)
	if err != nil {
		log_controller.LogErr(err)
		return "Fail", http.StatusBadRequest, err
	}
	err = proxy_order.CheckData()
	if err != nil {
		log_controller.LogErr(err)
		return "Fail", http.StatusPaymentRequired, err
	}
	order.MapData(proxy_order)
	answer, err := order.BookOrder()
	if err != nil {
		log_controller.LogErr(err)
		return "Fail", http.StatusBadRequest, err
	}
	log_controller.LogOrder(proxy_order, id, string(answer))
	return string(answer), http.StatusOK, nil
}
