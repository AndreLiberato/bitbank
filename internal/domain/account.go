package domain

const (
	AccountTypeSimple  = "simples"
	AccountTypeBonus   = "bonus"
	AccountTypeSavings = "poupanca"
)

type Account struct {
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
	Type    string  `json:"type"`
	Points  int     `json:"points"`
}

func (a Account) IsBonus() bool {
	return a.Type == AccountTypeBonus
}

func (a Account) IsSavings() bool {
	return a.Type == AccountTypeSavings
}
