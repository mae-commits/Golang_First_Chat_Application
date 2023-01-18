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

// ログインページのハンドラ
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("loginPage.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := html.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

// チャットページのハンドラ
func ChatHandler(w http.ResponseWriter, r *http.Request) {

}

// ログインページで入力ボタンを押した際に行われるハンドラ処理
// userName と password 入力がDB内にあるかどうか確認
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	Is_userName := strings.ReplaceAll(userName, " ", "")
	Is_password := strings.ReplaceAll(password, " ", "")
	count := getUser(userName, password)
	if Is_userName == "" || Is_password == "" {
		// w.WriteHeader(http.StatusUnauthorized)
		// fmt.Fprintf(w, "username or password is wrong.")
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if count != 0 {
		// Database 内に該当のusername とpassword の組が存在する場合、
		// 次のページへとジャンプする
		// そうでない場合は新規登録ページへと戻す
		http.Redirect(w, r, "/chat", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login/newResistration", http.StatusFound)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("deletePage.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := html.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var count int64
	db.Model(&User{}).Where("name = ?", userName).Where("password = ?", password).Count(&count)
	if count == 0 {
	} else {
		db.Delete(&User{Name: userName, Password: password})
	}
	http.Redirect(w, r, "/login", http.StatusFound)
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
	// Database にPOSTされてきたusernameとpasswordを入力する処理
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database.")
	}
	db.AutoMigrate(&User{})
	db.Create(&User{Name: userName, Password: password})
	// ログインページへリダイレクト
	http.Redirect(w, r, "/login", http.StatusFound)
}

// ユーザ認証の確認
// DB 中に該当のusername とpassword の組がない場合
// エラーを返す
func getUser(userName string, password string) (count int64) {
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database.")
	}
	db.Model(&User{}).Where("name = ?", userName).Where("password = ?", password).Count(&count)
	return count
}
