package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	HOST     = "localhost"
	PORT     = "3306"
	DB_NAME  = "go_db_2"
	PASSWORD = "password"
	DRIVER   = "mysql"
	USER     = "root"
)

var srcConn string = fmt.Sprintf("%s:%s@/%s", USER, PASSWORD, DB_NAME)

// Creates a new connection to the database using MySQL driver
func ConnectDB() (*DBConn, error) {
	conn, err := sql.Open(DRIVER, srcConn)
	if err != nil {
		return nil, fmt.Errorf("Error occured opening connection:\n %w", err)
	}

	return &DBConn{conn: conn}, nil
}

func castToUserRes(u *UserReq) *UserRes {
	var r UserRes
	r.Email = u.Email
	r.FirstName = u.FirstName
	r.LastName = u.LastName
	r.UserName = u.UserName

	fmt.Printf("request from User -> %+v\n", u)
	fmt.Printf("casted response -> %+v\n", r)

	return &r
}
func (c *DBConn) CreateUser(u *UserReq) (*UserRes, error) {
	fmt.Println("pinging...")
	if err := c.conn.Ping(); err != nil {
		return nil, fmt.Errorf("Database did not return ping ->\n:%w", err)
	}

	q := `
	INSERT INTO users (user_name, first_name, last_name, email)
	VALUES (?, ?, ?, ?);
	`

	res, err := c.conn.Exec(q, &u.UserName, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		return nil, fmt.Errorf("Error occured trying to insert user:\n %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Error occured trying to get insertId for user:\n %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		log.Println("WARNING: Could not get number of roles affected ->", err)
	}
	if n >= 2 {
		log.Println("WARNING: Rows affected for insert ->", n)
	}

	ur := castToUserRes(u)
	ur.Id = id

	fmt.Println("User Created")
	fmt.Println(ur)

	return ur, nil
}
