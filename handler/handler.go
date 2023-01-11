package handler

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("loginPage.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := html.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {

}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	Is_userName := strings.ReplaceAll(userName, " ", "")
	Is_password := strings.ReplaceAll(password, " ", "")
	if Is_userName == "" || Is_password == "" {
		// w.WriteHeader(http.StatusUnauthorized)
		// fmt.Fprintf(w, "username or password is wrong.")
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if userName == "userName" && password == "password" {
		// Database 内に該当のusername とpassword の組が存在する場合、
		// 次のページへとジャンプする
		// そうでない場合は新規登録ページへと戻す
		http.Redirect(w, r, "/chat", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login/newResistration", http.StatusFound)
	}
}

func NewResistrationHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("newResistrationPage.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := html.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

func NewResistrationPostHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	// Database にPOSTされてきたすusernameとpasswordを入力する処理
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database.")
	}
	db.AutoMigrate(&User{})
	db.Create(&User{Name: userName, Password: password})
	http.Redirect(w, r, "/login", http.StatusFound)
}
