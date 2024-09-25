package user_alert

import "fmt"

func AlertOnEmail(msg string) error {

	fmt.Println("*====WARNING====*")
	fmt.Println(msg)

	return nil
}
