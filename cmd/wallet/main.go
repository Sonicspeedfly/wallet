package main

import (
	"fmt"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/wallet"
)
func main() {
	svc := wallet.Service{}
	file := svc.Import("./information")
	fmt.Print(file)

}