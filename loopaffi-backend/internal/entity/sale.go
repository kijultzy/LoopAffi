package entity

type Sale struct {
	ID          string  `json:"id" db:"id"`
	Date        string  `json:"date" db:"date"`
	Amount      float64 `json:"amount" db:"amount"`
	AffiliateID string  `json:"affiliate_id" db:"affiliate_id"`
	Status      string  `json:"status" db:"status"`
}
