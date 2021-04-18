package wallet

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/Sonicspeedfly/wallet/v1.1.0/pkg/types"
	"github.com/google/uuid"
)

//ErrPhoneRegistered - телефон уже регитрирован
var ErrPhoneRegistered = errors.New("phone already registred")

//ErrAmountMustBePositive - счёт не может быть пустым
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")

//ErrAccountNotFound - пользователь не найден
var ErrAccountNotFound = errors.New("account not found")

//ErrNotEnoughtBalance - на счете недостаточно средств
var ErrNotEnoughtBalance = errors.New("account not enough balance")

//ErrPaymentNotFound - платеж не найден
var ErrPaymentNotFound = errors.New("payment not found")

// ErrFavoriteNotFound - Избранное не найдено
var ErrFavoriteNotFound = errors.New("favorite not found")

// Service представляет информацию о пользователе.
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

// RegisterAccount - метод для регистрация нового прользователя.
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
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

//Pay метод для регистрации платижа
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
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

	if account.Balance < amount {
		return nil, ErrNotEnoughtBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

// FindAccountByID ищем пользователя по ID
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

// FindPaymentByID ищем платёж по ID
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment

	for _, pay := range s.payments {
		if pay.ID == paymentID {
			payment = pay
		}
	}

	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}

// FindFavoriteByID ищем платёж по ID в Избранное
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

//Deposit method
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}

	account.Balance += amount
	return nil
}

// Reject метод для отмены покупок
func (s *Service) Reject(paymentID string) error {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	acc, err := s.FindAccountByID(pay.AccountID)
	if err != nil {
		return err
	}

	pay.Status = types.PaymentStatusFail
	acc.Balance += pay.Amount

	return nil
}

// Repeat позволāет по идентификатору повторитþ платёж
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(pay.AccountID, pay.Amount, pay.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// FavoritePayment добавления новых Избранных
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favoriteID := uuid.New().String()
	newFavorite := &types.Favorite{
		ID:        favoriteID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, newFavorite)
	return newFavorite, nil
}

//PayFromFavorite для совершения платежа в Избранное
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

//ExportToFile экспортирует аккаунт в файл
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func()  {
		if err := file.Close(); err != nil {
			log.Print(err)
		}
	}()
	var id int64
	var phone string
	var balance int64
	for _, account := range s.accounts {
		id = account.ID
		phone = string(account.Phone)
		balance = int64(account.Balance)
	_, err = file.Write([]byte(strconv.FormatInt(int64(id),10)+(";")+(phone)+(";")+(strconv.FormatInt(int64(balance),10))+("|")))
	if err != nil {
		return err
	}
}

	return nil
}