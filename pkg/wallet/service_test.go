package wallet

import (
	"testing"
)

func TestService_FindAccountByID_success(t *testing.T) {
	var result int64
	s := Service{}
	account, err := s.RegisterAccount("+992936888007")
	account, err = s.FindAccountByID(1)
	result = 1
	if account.ID != result {
		t.Error(err)
	}
}

func TestService_FindAccountByID_NotFoundAccount(t *testing.T) {
	var result int64
	svc := Service{}
	account, err := svc.RegisterAccount("+992936888007")
	account, err = svc.FindAccountByID(1)
	result = 3
	result2 := ErrAccountNotFound
	if err != nil{
		t.Error(ErrPhoneRegisted)
		return
	}
	if account.ID != result {
		if result2 != ErrAccountNotFound{
			t.Error(ErrAccountNotFound)
		}
	}
}