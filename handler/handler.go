package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 接続されるクライアント
var clients = make(map[*websocket.Conn]bool)

// メッセージブロードキャストチャネル
var broadcast = make(chan Message)

// アップグレーダ
var upgrader = websocket.Upgrader{}

type User struct {
	gorm.Model
	Name     string `json:"username"`
	Password string
}

// メッセージ用構造体
type Message struct {
	User
	Message string `json:"message"`
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

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("chatPage.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := html.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
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
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	count := deleteUser(userName, password)
	if count != 0 {
		db.Where("name = ?", userName).Where("password = ?", password).Delete(&User{})
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
	count := getUser(userName, password)
	// 新規登録のユーザである場合は、新たにユーザ登録
	// そうでない場合は再度ブランク画面に戻る
	if count == 0 {
		db.Create(&User{Name: userName, Password: password})
		// ログインページへリダイレクト
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login/newResistration", http.StatusFound)
		fmt.Println("This username has already been resistered.")
	}
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

func deleteUser(userName string, password string) (count int64) {
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database.")
	}
	db.Model(&User{}).Where("name = ?", userName).Where("password = ?", password).Count(&count)
	return count
}

// チャットページのハンドラ
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	// クライアントを登録
	clients[ws] = true

	for {
		var message Message
		// 新しいメッセージをJSONとして読み込み、Message構造体にマッピング
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// 受け取ったメッセージをbroadcast チャネルに送る
		broadcast <- message
	}
}

func HandleMessages() {
	for {
		// broadcast チャネルからメッセージを受け取る
		message := <-broadcast
		// 接続中の全クライアントにメッセージを送る
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
