package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-booking/secure"
	"io/ioutil"
	"log"
	"os"
)

type ResponseSecure struct {
	Id    string
	Token string
}

type Secure_Config struct {
	USER string
	PASS string
}

type SecretKey_Config struct {
	Secret_Key string
	id         int
}

//TODO rewrite struct
//TODO think about the secure/token.txt

var Config Secure_Config

func SignUp(c *gin.Context) (string, error) {
	var ConfigFromClient Secure_Config
	c.BindJSON(&ConfigFromClient)
	if Config.USER != ConfigFromClient.USER || Config.PASS != ConfigFromClient.PASS {
		err := errors.New("Incorrect Pass or USER\n")
		return "", err
	}
	token := new(secure.JWT)
	secretKey := token.GenerateToken("1")
	file, err := os.Create("secure/token.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	file.WriteString(secretKey)
	return token.GetToken(), nil
}

func CheckAuth(sessionToken string) bool {
	file, err := os.Open("secure/token.txt")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	Bt, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	token := new(secure.JWT)
	token.SetSignature(sessionToken)
	sk := string(Bt)
	return token.CheckToken(sk)
}

func GetConfigJWT() error {
	//file, err := os.Open("secure/secure_config.json")
	//defer file.Close()
	//if err != nil {
	//	return err
	//}
	//configByte, err := ioutil.ReadAll(file)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//err = json.Unmarshal(configByte, &Config)

	Config.USER = os.Getenv("USER")
	if Config.USER == "" {
		log.Fatal("$USER must be set")
	}

	Config.PASS = os.Getenv("PASS")
	if Config.PASS == "" {
		log.Fatal("$PASS must be set")
	}

	return nil
}
