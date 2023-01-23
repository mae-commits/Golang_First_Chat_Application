# Golangを用いたチャットアプリケーション

## 概要
- Go を用いたチャットアプリケーションです。ログイン機能とチャット機能を有しています。

## 技術的にこだわった点
- GORMを用いたDB管理
    - GORMというライブラリを用いてSQLiteでユーザ情報の登録・認証・削除をGoのコード上で行っています。
- WebSocketを用いた双方向通信
    - WebSocket を用いてサーバとクライアント間の双方向通信を可能にし、複数ユーザでのチャットが同期できるようになっています。

## 操作方法
```git clone git@github.com:mae-commits/Golang_First_Chat_Application.git```

```cd Golang_First_Chat_Application```

```go mod tidy```

```go run cmd/web/*.go```

## 難しかった箇所
- ユーザログインのDB管理
    - GoからDB、DBからGoへのデータの受け渡しの仕方に関して理解するのに苦労しました。
- WebSocketを用いた双方向通信
    - WebSocketの概念の理解 (通信手段)と実際の実装方法を考えるのに苦労しました。
    - json を用いたデータの受け渡しもあまり経験がないので、理解するのに苦戦しました。

## 使用技術・言語
- Go
  - GORMを用いたデータベース管理
  - WebSocket を用いた同期管理
- JavaScript
  - WebSocket を用いた同期管理
- HTML/CSS
  - アプリ画面のデザイン・表示
- SQLite
  - ユーザ名・パスワード管理（GORM）
