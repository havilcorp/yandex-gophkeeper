package entity

type SaveDto struct {
	Data []byte `json:"data"`
	Meta string `json:"meta"`
}
