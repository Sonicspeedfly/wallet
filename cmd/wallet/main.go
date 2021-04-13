package main

import (
	"fmt"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/wallet"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/types"
)
func main() {
	s := wallet.Service{}
	phone := types.Phone("+992936888007")
	account, err := s.RegisterAccount(phone)
	err = s.Deposit(account.ID, 10_000_00)
	if err != nil{
		return 
		}
		fmt.Println(account.Balance)
	}
