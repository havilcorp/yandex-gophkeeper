// Package модель данных
package entity

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
