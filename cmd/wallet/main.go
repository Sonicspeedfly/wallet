package main

import (
	"fmt"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/wallet"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/types"
)
func main() {
	//
	s := wallet.Service{}
	phone := types.Phone("+992000000001")
	account, err := s.RegisterAccount(phone)
	if err != nil{
		return
	}
	err = s.Deposit(account.ID, 10_000_00)
	payment, err := s.Pay(account.ID, 1_000_00, "auto")
	fmt.Println(payment)
	payment,err = s.Repeat(payment.ID)
	if err != nil{
		return
	}
	fmt.Println(payment)
}