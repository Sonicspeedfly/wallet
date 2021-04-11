package wallet

import (
	"errors"

	"github.com/Sonicspeedfly/wallet/v1/pkg/types")

//ErrPhoneRegisted ##**#
var ErrPhoneRegisted = errors.New("phone already registed")
//ErrAccountNotFound ##**#
var ErrAccountNotFound = errors.New("account not found")

//Service ##**##
type Service struct{
	NextAccountID int64
	accounts []types.Account
	payments []types.Payment
}

//RegisterAccount ##**##
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegisted
		}
	}
	s.NextAccountID++
	account := &types.Account{
		ID: s.NextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, *account)
	return account, nil
}

//FindAccountByID ##**##
func (s *Service) FindAccountByID(accontID int64)  (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accontID {
			account = &acc
			break
		}
	} 
	if account == nil{
		return nil, ErrAccountNotFound
	}
	return account, nil
}