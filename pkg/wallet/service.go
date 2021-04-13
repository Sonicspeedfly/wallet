package wallet

import (
  "errors"
  "fmt"

  "github.com/Sonicspeedfly/wallet/v1.1.0/pkg/types"
  "github.com/google/uuid"
)
//ErrPhoneRegistred ##**##
var ErrPhoneRegistred = errors.New("phone alredy registered")
//ErrAmountMustBePositive ##**##
var ErrAmountMustBePositive = errors.New("amount must be greater than 0")
//ErrAccountNotFound ##**##
var ErrAccountNotFound = errors.New("account not found")
//ErrNotEnoughtBalance ##**##
var ErrNotEnoughtBalance = errors.New("not enought balance")
//ErrPaymentNotFound ##**##
var ErrPaymentNotFound = errors.New("payment not found")

//Service ##**##
type Service struct {
  nextAccountID int64
  accounts      []*types.Account
  payments      []*types.Payment
}

//Error ##**#
type Error string

func (e Error) Error() string {
  return string(e)
}

//RegisterAccount ##**##
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
  for _, account := range s.accounts {
    if account.Phone == phone {
      return nil, ErrPhoneRegistred
    }
  }
  s.nextAccountID++
  account := &types.Account{
    ID:      s.nextAccountID,
    Phone:   phone,
    Balance: 0,
  }
  s.accounts = append(s.accounts, account)
  return account, nil
}

//Deposit ##**##
func (s *Service) Deposit(accountID int64, amount types.Money) error {
  if amount <= 0 {
    return ErrAmountMustBePositive
  }

  var account *types.Account
  for _, acc := range s.accounts {
    if acc.ID == accountID {
      account = acc
      break
    }
  }

  if account == nil {
    return ErrAccountNotFound
  }

  account.Balance += amount
  return nil
}

//Pay ##**##
func (s * Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error){
  if amount <= 0 {
    return nil, ErrAmountMustBePositive
  }

  var account *types.Account
  for _, acc := range s.accounts {
    if acc.ID == accountID {
      account = acc
      break
    }
  }
  if account == nil {
    return nil, ErrAccountNotFound
  }

  if account.Balance  <= amount {
    return nil, ErrNotEnoughtBalance
  }

  account.Balance -= amount
  paymentID := uuid.New().String()
  payment := &types.Payment{
    ID: paymentID,
    AccountID: accountID,
    Amount: amount,
    Category: category,
    Status: types.PaymentStatusInProgress,
  }
  s.payments = append(s.payments, payment)
  return payment, nil
}

//FindAccountByID ##**##
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
  var account *types.Account
  for _, acc := range s.accounts {
    if acc.ID == accountID {
      account = acc
      break
    }
  }
  if account == nil {
    return nil, ErrAccountNotFound
  }
  return account, nil
}

//FindPaymentByID ##**##
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
  for _, pay := range s.payments {
    if pay.ID == paymentID {
      return pay, nil
    }
  }
  return nil, ErrPaymentNotFound
}

//Reject ##**##
func (s *Service) Reject(paymentID string) error {
  payment, err := s.FindPaymentByID(paymentID)
  if err != nil {
    return nil
  }
  account, err := s.FindAccountByID(payment.AccountID)
  if err != nil {
    return nil
  }
  account.Balance += payment.Amount
  payment.Status = types.PaymentStatusFail
  return nil
}

//AddAccountWithBalance ##**##
func (s *Service) AddAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
  account, err := s.RegisterAccount(phone)
  if err != nil {
    return nil, fmt.Errorf("can't register account, error %v", err)
  }
  err = s.Deposit(account.ID, balance)
  if err != nil {
    return nil, fmt.Errorf("can't deposit account, error = %v", err)
  }
  return account, nil
}

//Repeat ##**##
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
  payment, err := s.FindPaymentByID(paymentID)  
  if err != nil {
    return nil, ErrAccountNotFound
  }
  paymentID = uuid.New().String()
  payment = &types.Payment{
    ID: paymentID,
    Amount: payment.Amount,
    AccountID: payment.AccountID,
    Category: payment.Category,
    Status: payment.Status,
  }
  s.payments = append(s.payments, payment)
  return payment, nil
}