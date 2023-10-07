package account

import "database/sql"

type Account struct {
	Id           int            `json:"id,omitempty"`
	Uname        string         `json:"uname,omitempty"`
	Pass         string         `json:"pass,omitempty"`
	AccessToken  sql.NullString `json:"access_token,omitempty"`
	RefreshToken sql.NullString `json:"refresh_token,omitempty"`
}
