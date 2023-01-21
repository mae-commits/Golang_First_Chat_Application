package handler

import (
	"chatapp/domain"
	"chatapp/internal/account"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// コネクション情報（ユーザやPayload）の保持
var (
	// ペイロードチャネルを作成
	wsChan = make(chan domain.WsPayload)

	// コネクションマップを作成
	// keyはコネクション情報, valueにはユーザー名を入れる
	clients = make(map[domain.WebSocketConnection]string)
)

// .jet ファイルを ./html ファイルとして読み込むための処理
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WebSocketsのエンドポイント
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// HTTPサーバーコネクションをWebSocketsプロトコルにアップグレード
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// リロードした際にターミナルに表示
	log.Println("OK Client Connecting")
	// サーバからのレスポンス
	var response domain.WsJsonResponse
	response.Message = `<li>Connected to server</li>`
	// コネクション情報を格納
	conn := domain.WebSocketConnection{Conn: ws}
	// ブラウザが読み込まれた時に一度だけ呼び出されるのでユーザ名なし
	clients[conn] = ""
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
	// 非同期でWebSocketをリッスンする
	//　WebSocket のリクエストをキャッチし続ける
	go ListenForWs(&conn)

}

// WebSocket をリッスンする関数
func ListenForWs(conn *domain.WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload domain.WsPayload

	//無限ループでずっと起動させる
	// これにより、既にコネクションが確立されている時通信は、
	// ここから情報を取得することができる
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			// WebSocket のエンドポイントにアクセスした時
			// そのコネクション情報からPayloadを読み込み、
			// チャネルに格納
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// 全てのユーザーにメッセージを返す
func broadcastToAll(response domain.WsJsonResponse) {
	// clientsには全ユーザーのコネクション情報が格納されている
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websockets err")
			_ = client.Close()
			// clients Mapからclientの情報を消す
			delete(clients, client)
		}
	}
}

// wsChannel からメッセージを受け取り、
// 全ユーザにメッセージをブロードキャストするため
// Webサーバとは別のプロセスで起動する必要がある
func ListenToWsChannel() {
	var response domain.WsJsonResponse

	for {
		// メッセージが入るまで、ここでブロック
		e := <-wsChan

		switch e.Action {
		// JavaScript で設定したusername Action に紐づいている
		case "username":
			// ここで、コネクションのユーザー名を格納
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users // 後ほど構造体に追加
			// wsChan チャネルにメッセージが格納されたら、
			// 全てのユーザにメッセージを送信
			broadcastToAll(response)
			// ページ離脱時に行うハンドリング処理
		case "left":
			response.Action = "list_users"
			// clientsからユーザーを削除
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)
			// メッセージ送信時に行うハンドリング処理
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf(
				"<li class='replace'><strong>%s</strong>: %s</li>",
				e.Username,
				e.Message)
			broadcastToAll(response)
		}
	}
}

// clients 関数に格納したユーザ情報を全て取得し、ブロードキャストする
func getUserList() []string {
	var clientList []string
	for _, client := range clients {
		if client != "" {
			clientList = append(clientList, client)
		}
		// 初回コネクション時は、ユーザ名が空の状態で渡されるので、処理をスキップ
	}
	sort.Strings(clientList)
	return clientList
}

// ログインページのハンドラ
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "loginPage.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

// .jet ファイルの中身を読込、エラーハンドリング
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ログインページで入力ボタンを押した際に行われるハンドラ処理
// userName と password 入力がDB内にあるかどうか確認
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	Is_userName := strings.ReplaceAll(userName, " ", "")
	Is_password := strings.ReplaceAll(password, " ", "")
	count := account.GetUser(userName, password)
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

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "chatPage.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	count := account.DeleteUser(userName, password)
	if count != 0 {
		db.Where("name = ?", userName).Where("password = ?", password).Delete(&domain.User{})
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
func NewResistrationHandler(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "newResistrationPage.jet", nil)
	if err != nil {
		log.Println(err)
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
	db.AutoMigrate(&domain.User{})
	count := account.GetUser(userName, password)
	// 新規登録のユーザである場合は、新たにユーザ登録
	// そうでない場合は再度ブランク画面に戻る
	if count == 0 {
		db.Create(&domain.User{Name: userName, Password: password})
		// ログインページへリダイレクト
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login/newResistration", http.StatusFound)
		fmt.Println("This username has already been resistered.")
	}
}
