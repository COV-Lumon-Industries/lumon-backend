package schemas

type TopUpWalletDetails struct {
	Amount int64 `gorm:"not null" json:"amount"`
}
