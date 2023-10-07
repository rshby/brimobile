// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AccountResponse struct {
	ID    int    `json:"id"`
	Uname string `json:"uname"`
	Pass  string `json:"pass"`
}

type CreateAccountResponse struct {
	ID    int    `json:"id"`
	Uname string `json:"uname"`
	Pass  string `json:"pass"`
}

type InqAccountSaving struct {
	AccountNumber    string `json:"accountNumber"`
	AvailableBalance string `json:"availableBalance"`
	AccountType      string `json:"accountType"`
	BranchCode       string `json:"branchCode"`
	Currency         string `json:"currency"`
	OpeningDate      string `json:"openingDate"`
	ProductGroup     string `json:"productGroup"`
	ProductName      string `json:"productName"`
	Status           string `json:"status"`
	CurrentBalance   string `json:"currentBalance"`
	ShortName        string `json:"shortName"`
}

type InsertSavingRequest struct {
	AccountNumber string `json:"accountNumber"`
	ShortName     string `json:"shortName"`
	Cbal          string `json:"cbal"`
}

type LoginResponse struct {
	LoginAt      string `json:"loginAt"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
