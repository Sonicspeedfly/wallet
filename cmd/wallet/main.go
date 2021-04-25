package main

import (
	"fmt"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/wallet"
)
func main() {
	svc := wallet.Service{}
	a, err := svc.RegisterAccount("+9920000001")
	if err != nil{
		fmt.Println(wallet.ErrAccountNotFound)
	}
	file := svc.Export("./information")
	fmt.Print(file)
	fmt.Println(a)
}