package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type PlaybackCommand struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
	Status   string `json:"status"`
}

type PlaybackHandler func(filename string, status string)

func (wr *WebSocketRelay) SetPlaybackHandler(handler PlaybackHandler) {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	wr.onPlayback = handler
}

type WebSocketRelay struct {
	masterConn *websocket.Conn
	guiConn    *websocket.Conn
	mu         sync.Mutex
	upgrader   websocket.Upgrader
	onPlayback PlaybackHandler
	lastcmd    string
}

// NewWebSocketRelay creates a new relay instance
func NewWebSocketRelay() *WebSocketRelay {
	return &WebSocketRelay{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (wr *WebSocketRelay) StartServer(addr string) error {
	http.HandleFunc("/ws", wr.handleGUIConnection)
	return http.ListenAndServe(addr, nil)
}

func (wr *WebSocketRelay) handleGUIConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := wr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	wr.mu.Lock()
	// Close existing GUI connection if any
	if wr.guiConn != nil {
		wr.guiConn.Close()
	}
	wr.guiConn = conn
	wr.mu.Unlock()

	// Start forwarding messages
	go wr.forwardMessages(conn, "GUI")
}

// SetMasterConnection sets up the master node connection
func (wr *WebSocketRelay) SetMasterConnection(conn *websocket.Conn) {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if wr.masterConn != nil {
		wr.masterConn.Close()
	}
	wr.masterConn = conn

	// Start forwarding messages
	go wr.forwardMessages(conn, "Master")
}

func (wr *WebSocketRelay) Broadcast(message []byte) {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if wr.guiConn != nil {
		if err := wr.guiConn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error broadcasting to GUI: %v", err)
		}
	}
}

// forwardMessages handles message forwarding between connections
// Modify the forwardMessages function to handle playback commands
func (wr *WebSocketRelay) forwardMessages(source *websocket.Conn, sourceType string) {
	defer source.Close()

	for {
		_, message, err := source.ReadMessage()
		if err != nil {
			log.Printf("%s connection closed: %v", sourceType, err)
			return
		}

		// Try to parse the message as a playback command
		var cmd PlaybackCommand
		if err := json.Unmarshal(message, &cmd); err == nil {
			if cmd.Type == "playback_command" {
				wr.mu.Lock()
				// log.Println("Before updaing lastcmd: ", wr.lastcmd)
				wr.lastcmd = cmd.Status
				// log.Println("After updaing lastcmd: ", wr.lastcmd)
				wr.mu.Unlock()
				log.Printf("[WebSocket] Received playback command: %s - %s", cmd.Filename, cmd.Status)
				continue
			}
		}

		// Handle other message types as before
		wr.mu.Lock()
		var dest *websocket.Conn
		if sourceType == "GUI" {
			dest = wr.masterConn
		} else {
			dest = wr.guiConn
		}

		if dest != nil {
			if err := dest.WriteMessage(websocket.TextMessage, message); err != nil {
				destType := "Master"
				if sourceType == "GUI" {
					destType = "GUI"
				}
				log.Printf("Error forwarding message to %s: %v", destType, err)
			}
		}
		wr.mu.Unlock()
	}
}

func (wr *WebSocketRelay) GetLastCommand() string {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	return wr.lastcmd
}
