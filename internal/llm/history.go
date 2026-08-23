package llm

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	clients = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

func InitHistory() { os.MkdirAll("history", 0755) }

func SaveHistory(text string) {
	f, _ := os.OpenFile("history/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(time.Now().Format("2006-01-02 15:04:05") + " - " + text + "\n")
}

func Broadcast(msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for c := range clients {
		c.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
}
