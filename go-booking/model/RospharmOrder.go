package model

import "errors"

type Order_Inter struct {
	StoreId string  `json:"store_id"`
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Comment string  `json:"comment"`
	Basket  []Pharm `json:"basket"`
}

type Pharm struct {
	Product_id string `json:"product_id"`
	Qty        int    `json:"qty"`
}

func (o *Order_Inter) CheckData() error {
	if len(o.StoreId) == 0 {
		err := errors.New("StoreId is not define")
		return err
	}
	if len(o.Name) == 0 {
		err := errors.New("Name is not define")
		return err
	}
	if len(o.Basket) == 0 {
		err := errors.New("Basket is not define")
		return err
	}
	return nil
}
