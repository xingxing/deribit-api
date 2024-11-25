package models

import "encoding/json"

var EmptyParams = json.RawMessage("{}")

// privateParams is interface for methods require access_token
type PrivateParams interface {
	SetToken(token string)
}

// Token is used to embedded in params for private methods
type Token struct {
	AccessToken string `json:"access_token"`
}

func (t *Token) SetToken(token string) {
	t.AccessToken = token
}
