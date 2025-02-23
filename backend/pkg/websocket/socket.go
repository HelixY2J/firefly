package websocket

import (
    "log"
    "net/http"
    "sync"
    "encoding/json"

    "github.com/gorilla/websocket"
)

type PlaybackCommand struct {
    Type     string `json:"type"`
    Filename string `json:"filename"`
    Status   string `json:"status"`
}

// Add this callback type
type PlaybackHandler func(filename string, status string)

// Add this method to set the playback handler
func (wr *WebSocketRelay) SetPlaybackHandler(handler PlaybackHandler) {
    wr.mu.Lock()
    defer wr.mu.Unlock()
    wr.onPlayback = handler
}

// WebSocketRelay handles message forwarding between master and GUI
type WebSocketRelay struct {
    masterConn *websocket.Conn
    guiConn    *websocket.Conn
    mu         sync.Mutex
    upgrader   websocket.Upgrader
    onPlayback PlaybackHandler
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


// StartServer starts the WebSocket server for GUI connections
func (wr *WebSocketRelay) StartServer(addr string) error {
    http.HandleFunc("/ws", wr.handleGUIConnection)
    return http.ListenAndServe(addr, nil)
}

// handleGUIConnection handles incoming GUI WebSocket connections
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
                if wr.onPlayback != nil {
                    wr.onPlayback(cmd.Filename, cmd.Status)
                }
                wr.mu.Unlock()
                log.Printf("Received playback command: %s - %s", cmd.Filename, cmd.Status)
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