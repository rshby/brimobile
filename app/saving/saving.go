package saving

type Saving struct {
	AccountNumber string `json:"account_number,omitempty"`
	AccountType   string `json:"account_type,omitempty"`
	BranchCode    string `json:"branch_code,omitempty"`
	ShortName     string `json:"short_name,omitempty"`
	Currency      string `json:"currency,omitempty"`
	Cbal          string `json:"cbal,omitempty"`
	Hold          string `json:"hold,omitempty"`
	OpeningDate   string `json:"opening_date,omitempty"`
	ProductGroup  string `json:"product_group,omitempty"`
	ProductName   string `json:"product_name,omitempty"`
	Status        string `json:"status,omitempty"`
}
