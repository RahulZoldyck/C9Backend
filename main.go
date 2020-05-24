package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"io/ioutil"
	"log"
	"net/http"
)

const sqlServerUrl = "sqlserver://username:password@localhost:1433?database=C9Unity"

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

type User struct {
	UserName     string `json:"user_name" gorm:"column:cAccId"`
	PasswordHash string `json:"password_hash" gorm:"column:cPassword"`
	AuthLevel    string `gorm:"column:cAuthLevel"`
	Email        string `json:"email" gorm:"-"`
	Mode         int    `gorm:"-"`
}

func (u User) TableName() string {
	if u.Mode == 1 {
		return "Auth.TblAccount"
	} else {
		return "dbo.WebAccount"
	}
}

func (u *User) register() {
	u.Mode = 0
	db.FirstOrCreate(&u, u)
	u.Mode = 1
	db.FirstOrCreate(&u, u)
}

func (u *User) login() bool {
	var user User
	db.Where("cAccId = ?", u.UserName).Find(&user)
	return user.PasswordHash == u.PasswordHash
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &user)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &user)
	result := user.login()
	outcome := ""
	if result {
		outcome = "success"
	} else {
		outcome = "failed"
	}
	fmt.Fprintf(w, "{outcome : "+outcome+"}")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/", registerHandler)
	http.HandleFunc("/", loginHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("mssql", sqlServerUrl)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	handleRequests()

}
