package wallet

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
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
	accountStr := ""
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	
	if path == "accounts.dump"{for _, acc := range s.accounts {
		id := strconv.Itoa(int(acc.ID))
		phone := string(acc.Phone)
		balance := strconv.Itoa(int(acc.Balance))
		accountStr +=  id + ";" + phone + ";" + balance + "\n"
		
	}
}

	if path == "payments.dump"{for _, payment := range s.payments {
		id := string(payment.ID)
		accountID := strconv.Itoa(int(payment.AccountID))
		amount := strconv.Itoa(int(payment.Amount))
		category := string(payment.Category)
		status := string(payment.Status)
		accountStr +=  id + ";" + accountID + ";" + amount + ";" + category + ";" + status + "\n"
	}
}

	
	if path == "favorites.dump"{
		for _, favorite := range s.favorites {
		id := string(favorite.ID)
		accountID := strconv.Itoa(int(favorite.AccountID))
		name := string(favorite.Name)
		amount := strconv.Itoa(int(favorite.Amount))
		category := string(favorite.Category)
		accountStr +=  id + ";" + accountID + ";" + name + ";" + amount + ";" + category + "\n"
		
		}
	}
	if path != "accounts.dump" && path != "payments.dump" && path != "favorites.dump"{
	for _, acc := range s.accounts {
		id := strconv.Itoa(int(acc.ID))
		phone := string(acc.Phone)
		balance := strconv.Itoa(int(acc.Balance))
		accountStr +=  id + ";" + phone + ";" + balance + "|"
	}
	}
	accountStr = accountStr[:len(accountStr)-1]
	_, err = file.Write([]byte(accountStr))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}


//ImportFromFile импортировать с файла
func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	content := make([]byte, 0)
	buf := make([]byte, 4)
	
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)	
			break
		}

		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}

	data := string(content)
	if path == "accounts.dump"{
		content, err := ioutil.ReadFile("accounts.dump")
		if err != nil {
			return err
		}
		rows := strings.Split(string(content), "\n")
		for _, row  := range rows {
			cols := strings.Split(row, ";")

			id, err := strconv.ParseInt(cols[0], 10, 64)
			if err != nil {
				return err
			}
			balance, err := strconv.ParseInt(cols[2], 10, 64)
			if err != nil {
				return err
			}
			flag := true
			for _, v := range s.accounts {
				if v.ID == id {
					flag = false
				}
			}
			if flag {
				account := &types.Account{
					ID:      id,
					Phone:   types.Phone(cols[1]),
					Balance: types.Money(balance),
				}
				s.accounts = append(s.accounts, account)
			}
		}
	}
	if path == "payments.dump"{
		content, err := ioutil.ReadFile("payments.dump")
		if err != nil {
			return err
		}
		rows := strings.Split(string(content), "\n")
		for _, row  := range rows {
			cols := strings.Split(row, ";")

			id := cols[0]
			if err != nil {
				return err
			}
			accountID, err := strconv.ParseInt(cols[1], 10, 64)
			if err != nil {
				return err
			}
			amount, err := strconv.ParseInt(cols[2], 10, 64)
			if err != nil {
				return err
			}
			flag := true
			for _, v := range s.payments {
				if v.ID == id {
					flag = false
				}
			}
			if flag {
				data := &types.Payment{
					ID:        id,
					AccountID: accountID,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(cols[3]),
					Status:    types.PaymentStatus(cols[4]),
				}
				s.payments = append(s.payments, data)
			}
		}
	}
	if path == "favorites.dump"{
		content, err := ioutil.ReadFile("favorites.dump")
		if err != nil {
			return err
		}
		rows := strings.Split(string(content), "\n")
		for _, row  := range rows {
			cols := strings.Split(row, ";")

			id := cols[0]
			if err != nil {
				return err
			}
			accountID, err := strconv.ParseInt(cols[1], 10, 64)
			if err != nil {
				return err
			}
			amount, err := strconv.ParseInt(cols[3], 10, 64)
			if err != nil {
				return err
			}
			flag := true
			for _, v := range s.favorites {
				if v.ID == id {
					flag = false
				}
			}
			if flag {
				data := &types.Favorite{
					ID:        id,
					AccountID: accountID,
					Name: 	   cols[2],
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(cols[4]),
				}
				s.favorites = append(s.favorites, data)
			}
		}
	}
	if path != "accounts.dump" && path != "payments.dump" && path != "favorites.dump"{
	rows := strings.Split(data, "|")
	for _, row := range rows {
		cols := strings.Split(row, ";")
		id, _ := strconv.ParseInt(cols[0],10,64)
		phone := types.Phone(cols[1])
		balance, _ := strconv.ParseInt(cols[2],10,64)
		
		account := &types.Account{
			ID:      id,
			Phone:   phone,
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, account)
	}
	}
	return nil
}

//Export экспортировать
func (s *Service) Export(dir string) error {
	_, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	} 
	if len(s.accounts) > 0 {
	s.ExportToFile("accounts.dump")
	if len(s.payments) > 0{
		s.ExportToFile("payments.dump")
	if len(s.favorites) > 0{
	s.ExportToFile("favorites.dump")
	}}}
	return nil
}

//Import импортирует данные из файла
func (s *Service) Import(dir string) error {
	_, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	} 
	if len(s.accounts) > 0 {
	s.ImportFromFile("accounts.dump")
	if len(s.payments) > 0{
		s.ImportFromFile("payments.dump")
	if len(s.favorites) > 0{
	s.ImportFromFile("favorites.dump")
	}}}
	return nil
}
