package db

import "database/sql"

type DBConn struct {
	Conn *sql.DB
}

// this is going to be the schema for our request
type UserReq struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserRes struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
