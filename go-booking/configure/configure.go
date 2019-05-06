package configure

import (
	"encoding/json"
	"fmt"
	"go-booking/model/config"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type ServerConfig struct {
	GoEnv        string
	Port         string
	MockMode     bool
	RostaConf    config.RostaConfig
	PermFarmConf config.PermFarmConfig
	Apteka366    config.Apteka366
	Rigla        config.Rigla
	Code_Farm    config.Code_Farm_List
}

var Server ServerConfig

func (s *ServerConfig) LoadConfigServer() error {
	file, err := os.Open("configure/config.json")
	defer file.Close()

	if err != nil {
		return err
	}

	configByte, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(configByte, &s)
	if err != nil {
		log.Fatal(err)
	}

	s.GoEnv = os.Getenv("GO_ENV")
	if s.GoEnv == "" {
		log.Fatal("$GO_ENV must be set")
	}

	s.Port = os.Getenv("PORT")
	if s.Port == "" {
		log.Fatal("$PORT must be set")
	}

	s.MockMode, err = strconv.ParseBool(os.Getenv("MOCK_MODE"))
	if err != nil {
		log.Printf("Error converting $MOCK_MODE to an bool: %q - Using default\n", err)
		s.MockMode = false
	}

	fmt.Printf("ConfigData:\n GO_ENV: %s\n PORT: %s\n MOCK_MODE: %v\n", s.GoEnv, s.Port, s.MockMode)
	Server = *s
	fmt.Println(Server)
	return err
}

func GetURLData(code string) (string, string, string) {
	switch code {
	case Server.Code_Farm.ROSTA_CODE_FARM:
		{
			url, user, pass := Server.RostaConf.GetURLData(Server.MockMode)
			return url, user, pass
		}
	case Server.Code_Farm.PERM_CODE_FARM:
		{
			url := Server.PermFarmConf.GetUrlData(Server.MockMode)
			return url, "", ""
		}
	case Server.Code_Farm.APTEKA366_CODE_FARM:
		{
			url := Server.Apteka366.GetUrlData(Server.MockMode)
			return url, "", ""
		}
	}
	return "", "", ""
}
