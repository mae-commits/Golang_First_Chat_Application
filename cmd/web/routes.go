package main

import (
	"chatapp/internal/handler"
	"net/http"

	"github.com/bmizerany/pat"
)

func routes() http.Handler {
	mux := pat.New()
	// ログインページ
	mux.Get("/login", http.HandlerFunc(handler.LoginHandler))
	// ログインページ処理のハンドラ
	mux.Post("/login/create", http.HandlerFunc(handler.CreateHandler))
	// ユーザ新規登録ページ
	mux.Get("/login/newResistration", http.HandlerFunc(handler.NewResistrationHandler))
	// ユーザ新規登録処理のハンドラ
	mux.Post("/login/newResistrationPost", http.HandlerFunc(handler.NewResistrationPostHandler))
	// ユーザ情報削除のハンドラ
	mux.Post("/login/delete", http.HandlerFunc(handler.DeleteHandler))
	// チャットページのハンドラ
	mux.Get("/chat", http.HandlerFunc(handler.ChatHandler))
	// WebSocket のエンドポイントを格納
	mux.Get("/ws", http.HandlerFunc(handler.WsEndpoint))
	// /static 配下のファイルを一括で読み込み
	fileServer := http.FileServer(http.Dir("./static/"))
	// localhost:8080/static としてファイルを配信できるようになる
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
