package models

type User struct{
	ID int `json:"ID"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Address  string `json:"Address"`

}
