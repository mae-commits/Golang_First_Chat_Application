package domain

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// ユーザの構造体
type User struct {
	gorm.Model
	Name     string
	Password string
}

// WebSocket からの返却用データの構造体
type WsJsonResponse struct {
	// WebSocket へ送る命令を格納
	Action string `json:"action"`
	// ブラウザから受け取ったメッセージを格納
	Message string `json:"message"`
	// ユーザリスト情報
	ConnectedUsers []string `json:"connected_users"`
}

// WebSocketコネクション情報を格納
type WebSocketConnection struct {
	*websocket.Conn
}

// WebSocket送信データを格納
// ブラウザから送信されたデータを格納
type WsPayload struct {
	Action   string              `json:"action"`
	Message  string              `json:"message"`
	Username string              `json:"username"`
	Conn     WebSocketConnection `json:"-"`
}
