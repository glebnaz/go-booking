package routing

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-booking/controllers"
	"go-booking/model"
	"net/http"
	"strings"
)

func ConfigureRouters() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Status": "Ok"})
	})

	//orderCreate
	router.POST("booking/order/:pharmacy_id", func(c *gin.Context) {
		bareir := c.GetHeader("Authorization")
		bareirArr := strings.Split(bareir, " ")
		if len(bareirArr) < 2 {
			err := errors.New("Header is not right")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}else{
			token := bareirArr[1]
			if !controllers.CheckAuth(token) {
				err := errors.New("TOKEN UNDERFIND")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}else{
				var proxy_order model.Order_Inter
				err:=c.BindJSON(&proxy_order)
				if err!=nil{
					err := errors.New("Can't parse Order" )
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}else{
					id := c.Param("pharmacy_id")
					answer,status,err:=controllers.BookingOrder(proxy_order,id)
					c.JSON(status, gin.H{"Answer":answer,"error": err.Error()})
				}
			}
		}
	})

	router.POST("sign_in", func(c *gin.Context) {
		token, err := controllers.SignUp(c)
		if err != nil {
			c.JSON(http.StatusNotExtended, gin.H{"err": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	})

	return router
}
