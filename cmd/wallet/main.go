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
	file := svc.ExportToFile("massage.txt")
	read := svc.ImportFromFile("massage.txt")
	fmt.Print(file,read)
	fmt.Println(a)
}