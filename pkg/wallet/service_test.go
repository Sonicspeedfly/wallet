package wallet

import (
	"reflect"
	"testing"
	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/types"
)

type testService struct {
	Service
  }
  func newTestService() *testService {
	return &testService{Service: Service{}}
  }
  func TestFinfAccountByID_empty(t *testing.T) {
	svc := Service{}
	result, err := svc.FindAccountByID(1)
	if err != ErrAccountNotFound || result != nil {
	  t.Error("Ошибка 1")
	}
  }
  
  func TestFinfAccountByID_notEmpty(t *testing.T) {
	svc := Service{}
	result, err := svc.RegisterAccount("+992000000001")
	result, err = svc.FindAccountByID(3)
	if err != ErrAccountNotFound || result != nil {
	  t.Error("Ошибка 2")
	}
  }
  
  func TestDeposit(t *testing.T) {
	//
	s := Service{}
	//
	phone := types.Phone("+992000000001")
	account, err := s.RegisterAccount(phone)
	if err != nil {
	  t.Errorf("Reject(): can't register account, error = %v", err)
	  return
	}
	//
	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
	  t.Errorf("Reject(): can't deposit account, error = %v", err)
	  return
	}
  }
  
  func TestReject_succecs(t *testing.T) {
	//
	s := Service{}
	//
	phone := types.Phone("+992000000001")
	account, err := s.RegisterAccount(phone)
	if err != nil {
	  t.Errorf("Reject(): can't register account, error = %v", err)
	  return
	}
	//
	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
	  t.Errorf("Reject(): can't deposit account, error = %v", err)
	  return
	}
	//
	payment, err := s.Pay(account.ID, 1_000_00, "auto")
	if err != nil {
	  t.Errorf("Reject(): can't make payment, error = %v", err)
	  return
	}
	//
	err = s.Reject(payment.ID)
	if err != nil {
	  t.Errorf("Reject(): can't reject payment, error = %v", err)
	  return
	}
  }
  
  func TestServise_FindPaymentByID_success(t *testing.T) {
	//
	s := newTestService()
	account, err := s.AddAccountWithBalance ("+992000000001", 10_000_00)
	if err != nil {
	  t.Error(err)
	  return
	}
	//
	payment, err := s.Pay(account.ID, 1_000_00, "auto")
	if err != nil {
	  t.Errorf("FindPaymentByID(): can't create payment, error = %v", err)
	  return
	}
	//
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
	  t.Errorf("FindPaymentByID(): can't find payment, error = %v", err)
	  return
	}
	//
	if !reflect.DeepEqual(payment,got) {
	  t.Errorf("FindPaymentByID(): wrong payment returned, error = %v", err)
	  return
	}
  }