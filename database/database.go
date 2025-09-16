package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	DRIVER = "mysql"
)

func loadEnv() error {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Could not load env for database")
		return err
	}

	return nil
}

func ConnectDB() (*DBConn, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}

	HOST := os.Getenv("SQL_HOST")
	PASSWORD := os.Getenv("SQL_PASSWORD")
	USER := os.Getenv("SQL_USER")
	DB_NAME := os.Getenv("DB_NAME")
	PORT, err := strconv.ParseInt(os.Getenv("SQL_PORT"), 10, 64)

	if HOST == "" || PASSWORD == "" || USER == "" || DB_NAME == "" || PORT == 0 {
		return nil, fmt.Errorf("Check env!, some values have not been set properly\n")
	}

	var connStr string = fmt.Sprintf("%s:%s@/%s", USER, PASSWORD, DB_NAME)
	conn, err := sql.Open(DRIVER, connStr)
	if err != nil {
		return nil, fmt.Errorf("Error occured opening connection:\n %w", err)
	}

	return &DBConn{Conn: conn}, nil
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
	if err := c.Conn.Ping(); err != nil {
		return nil, fmt.Errorf("Database did not return ping ->\n:%w", err)
	}

	q := `
	INSERT INTO users (user_name, first_name, last_name, email)
	VALUES (?, ?, ?, ?);
	`

	res, err := c.Conn.Exec(q, &u.UserName, &u.FirstName, &u.LastName, &u.Email)
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
