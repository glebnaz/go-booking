package controllers

import (
	"errors"
	"go-booking/model"
	"go-booking/model/Orders"
	"go-booking/model/config"
)

func IdentificationCompany(pharmId string, code_pharm_list config.Code_Farm_List) (model.Order, error) {
	switch pharmId {
	case code_pharm_list.APTEKA366_CODE_FARM:
		{
			return &Orders.SixAndSixOrder{}, nil
		}
	case code_pharm_list.PERM_CODE_FARM:
		{
			return &Orders.PermFarmOrder{}, nil
		}
	case code_pharm_list.ROSTA_CODE_FARM:
		{
			return &Orders.RostaOrder{}, nil
		}
	case code_pharm_list.RIGLA_CODE_FARM:
		{
			return &Orders.Order_Rigla{}, nil
		}
	}

	return nil, errors.New("Code_Farm is incorrect")
}
