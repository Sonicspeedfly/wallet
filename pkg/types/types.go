package types

//Money представляет собой денежную сумму в минимальной денежных еденицах (центы, копейки, дирамы и т.д) .
type Money int64

//PaymentCategory представляет собой категорию, в уоторой был совершён плтёж (авто, аптеки, рестораны и т.д.) .
type PaymentCategory string

//PaymentStatus представляет собой платёж
type PaymentStatus string

//Предопределённые статусы платежа
const (
	StatusOk PaymentStatus = "OK"
	StatusFail PaymentStatus = "FAIL"
	StatusInProgress PaymentStatus = "INPROGRESS"
)

//Payment представляет информацию о платеже
type Payment struct {
	ID int64
	Amount Money 
	Category PaymentCategory
	Status PaymentStatus
}

type Phone string

//Account представляет информацию о счёте пользователя
type Account struct {
	ID int64
	Phone Phone
	Balance Money
}