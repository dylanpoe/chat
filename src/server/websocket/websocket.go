package controllers

import (
	"github.com/gorilla/websocket"
	"server/chatroom"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


func Room(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")
		formatter.HTML(w,http.StatusOK,"websocket/room", struct {
			User string
		}{user})
	}
}

func WsHandler(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		chatroom.Join(user)
		defer chatroom.Leave(user)

		// Join the room.
		subscription := chatroom.Subscribe()
		defer subscription.Cancel()

		//先把历史消息推送出去
		// Send down the archive.
		for _, event := range subscription.Archive {
			if conn.WriteJSON(&event) != nil {
				// They disconnected
				return
			}
		}

		// In order to select between websocket messages and subscription events, we
		// need to stuff websocket events into a channel.
		newMessages := make(chan string)
		go func() {
			var res = struct {
				Msg string `json:"msg"`
			}{}
			for {
				err := conn.ReadJSON(&res)
				if err != nil {
					close(newMessages)
					return
				}
				newMessages <- res.Msg
			}
		}()

		// Now listen for new events from either the websocket or the chatroom.
		for {
			select {
			case event := <-subscription.NewMsg:
				if conn.WriteJSON(&event) != nil {
					// They disconnected.
					return
				}
			case msg, ok := <-newMessages:
				// If the channel is closed, they disconnected.
				if !ok {
					return
				}
				// Otherwise, say something.
				chatroom.Say(user, msg)
			}
		}

		formatter.Text(w,http.StatusOK,"ws closed")
		return

	}
}

