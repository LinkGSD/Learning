package model

type User struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	Age      int64  `json:"age"`
}
