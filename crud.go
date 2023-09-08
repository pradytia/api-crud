package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
)

const (
	dbDriver   = "mysql"
	dbUsername = "root"
	dbPassword = "pass"
	dbName     = "my_db"
)

type UserModels struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func dbConnect() (db *sql.DB) {

	dataSourceName := fmt.Sprintf("%s:%s@/%s", dbUsername, dbPassword, dbName)
	fmt.Println(dataSourceName)
	db, err := sql.Open(dbDriver, dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetListUser(res http.ResponseWriter, req *http.Request) {

	db := dbConnect()

	defer db.Close()

	var sqlRow *sql.Rows
	var errSql error

	UserID := req.URL.Query().Get("id")

	if UserID == "" {
		sqlRow, errSql = db.Query("SELECT id, email FROM users ORDER BY id ASC LIMIT 10")
		if errSql != nil {
			panic(errSql.Error())
		}
	} else {
		sqlRow, errSql = db.Query("SELECT id, email from users where id = ?", UserID)
		if errSql != nil {
			panic(errSql.Error())
		}
	}

	var userList []UserModels
	user := UserModels{}

	for sqlRow.Next() {

		if err := sqlRow.Scan(&user.ID, &user.Email); err != nil {
			panic(err.Error())
		}

		userList = append(userList, user)

		if err := json.NewEncoder(res).Encode(userList); err != nil {
			panic(err.Error())
		}
	}
	fmt.Println("user list ", userList)
}

func GetListUserByID(res http.ResponseWriter, req *http.Request) {

	db := dbConnect()

	defer db.Close()

	userID := req.URL.Query().Get("id")

	fmt.Println("userID", reflect.TypeOf(userID))

	user := UserModels{}

	err := db.QueryRow("SELECT id, email from users where id = ?", userID).Scan(&user.ID, &user.Email)

	if err != nil {
		fmt.Println("error query ", err.Error())
		panic(err.Error())
	}

	fmt.Println("result get by id ", user)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", GetListUser).Methods("GET")
	router.HandleFunc("/test", GetListUserByID).Methods("GET")

	fmt.Println("server running on port 8080")

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
