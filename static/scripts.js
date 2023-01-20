let socket = null;

// メッセージリストを宣言
let messageList = document.getElementById("message-list")

// 最初のHTMLの読み込みと解析が完了したとき、スタイルシート、
// 画像、サブフレームの読み込みが完了するのを待たずに発生。
document.addEventListener("DOMContentLoaded", function(){

  // WebSocketオブジェクトの作成
  // https://developer.mozilla.org/ja/docs/Web/API/WebSockets_API/Writing_WebSocket_client_applications
  //　localhost:8080 にアクセスした時に一度だけ呼び出される     
  socket = new ReconnectingWebSocket(
  "ws://127.0.0.1:8080/ws",
  null,
  {
    debug: true,
    reconnectInterval: 3000 // 3s後に再接続
  })
  //　接続を確立
  // https://developer.mozilla.org/ja/docs/Web/API/WebSocket/onopen
  socket.onopen = () => {
    console.log("Successfully connected")
  }

    // 接続がCLOSEDに変わった時に呼ばれる
  // https://developer.mozilla.org/ja/docs/Web/API/WebSocket/onclose
  socket.onclose = () => {
    console.log("connection closed")
  }

  // エラーが発生した時に呼び出される
  // https://developer.mozilla.org/ja/docs/Web/API/WebSocket/onerror
  socket.onerror = error => {
    console.log("there was an error")
  }

  // サーバーからメッセージが届いたときに呼び出される
  // https://developer.mozilla.org/ja/docs/Web/API/WebSocket/onmessage
  socket.onmessage = msg => {
    let data = JSON.parse(msg.data)
    console.log({data})
    console.log("Action is", data.action)
    switch (data.action) {
      case "list_users":
        let ul = document.getElementById("online-users")
        while (ul.firstChild) ul.removeChild(ul.firstChild)

        if (data.connected_users.length > 0) {
          data.connected_users.forEach(function(item){
            let li = document.createElement("li")
            li.appendChild(document.createTextNode(item))
            ul.appendChild(li)
          })
        }
      case "broadcast":
       let message = data.message
       let username = document.getElementById("username").value
      // メッセージが自分のものかチェック
      // classをmeかotherに書き換え
      if (message.indexOf(username) > 0) {
          message = message.replace("replace", "me")
      } else {
          message = message.replace("replace", "other")
      }
      messageList.innerHTML = messageList.innerHTML + message
      break
    }
  }

  let userInput = document.getElementById("username")
  userInput.addEventListener("change", function() {
    let jsonData = {}
    // Action Nameをusernameにする
    jsonData["action"] = "username"
    jsonData["username"] = this.value;
    // user名を送信
    socket.send(JSON.stringify(jsonData))
  })

  document.getElementById("message").addEventListener("keydown", function(event) {
    if (event.code === "Enter") {
      if (!socket) {
        console.log("no connection")
        return false
      }
      // HTML要素既存の動きやイベント伝搬をキャンセル
      event.preventDefault()
      event.stopPropagation()
      sendMessage()
    }
  })

})

// ページから離脱時に発生
window.onbeforeunload = function() {
    console.log("User Leaving")
    let jsonData = {}
    jsonData["action"] = "left"
    socket.send(JSON.stringify(jsonData))
  }

function sendMessage() {
console.log("Send Message...")
let jsonData = {}
jsonData["action"] = "broadcast"
jsonData["username"] = document.getElementById("username").value
jsonData["message"] = document.getElementById("message").value
socket.send(JSON.stringify(jsonData))
document.getElementById("message").value = ""
}