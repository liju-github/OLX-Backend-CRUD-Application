package model

type UserWithProducts struct {
	User     User      `json:"user"`
	Products []Product `json:"products"`
}
