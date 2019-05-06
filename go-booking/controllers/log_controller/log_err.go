package log_controller

import (
	"fmt"
	"time"
)

func LogErr(errorFromFunc error) {
	//file, err := os.OpenFile("log/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//defer file.Close()
	//if err != nil {
	//	fmt.Printf("ERR when log err : %v\n", err)
	//	file, _ = os.Create("log/log.txt")
	//}

	fmt.Println("ERROR!")
	if errorFromFunc == nil {
		t := time.Now()
		msg := t.String() + ";" + "err nil(mistake in code)" + ";\n"
		//file.WriteString(msg)
		fmt.Println(msg)

	} else {
		t := time.Now()
		msg := t.String() + ";" + errorFromFunc.Error() + ";\n"
		//file.WriteString(msg)
		fmt.Println(msg)
	}
	fmt.Println("END ERROR!")
}
