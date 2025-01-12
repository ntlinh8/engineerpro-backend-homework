package models

type User struct{
	ID int `json:"ID"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password []byte `json:"Password"`
	Address  string `json:"Address"`

}
