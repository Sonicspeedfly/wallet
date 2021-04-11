package main

import (
	"fmt"
	"github.com/Sonicspeedfly/wallet/v1/pkg/wallet"
)
func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992936888007")
	account, err = svc.FindAccountByID(1)
	if err != nil {
		fmt.Println("Аккаунт пользователя не найден")
		return
	}
	fmt.Println(account.ID)
}