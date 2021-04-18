package main

import (
	"fmt"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/wallet"
)
func main() {
	svc := wallet.Service{}
	svc.RegisterAccount("+9920000001")
	file := svc.ExportToFile("massage.txt")
	fmt.Print(file)
}