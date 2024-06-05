package entity

type CreateDto struct {
	Data []byte `json:"data"`
	Meta string `json:"meta"`
}
