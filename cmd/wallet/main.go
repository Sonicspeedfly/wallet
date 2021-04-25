package main

import (
	"fmt"

	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/types"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/wallet"
)
func main() {
	svc := wallet.Service{}
	payments := []types.Payment{
		{
			ID: "1",
			AccountID: 1,
			Amount: 10000,
			Category: "a",
			Status: "INPROGRESS",
		},
		{
			ID: "1",
			AccountID: 1,
			Amount: 10000,
			Category: "a",
			Status: "INPROGRESS",
		},
	}


	file := svc.HistoryToFiles(payments, "./information", 1)
	fmt.Print(file, payments)

}