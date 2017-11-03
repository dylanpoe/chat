package longpolling

import (
	"server/chatroom"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/unrolled/render"
	"strconv"
)


func Room(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")
		chatroom.Join(user)
		formatter.HTML(w,http.StatusOK,"longpolling/room", struct {
			User string
		}{user})
	}
}

func Say(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")
		message := r.FormValue("message")
		chatroom.Say(user, message)
	}
}

func Leave(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")
		chatroom.Leave(user)
		http.Redirect(w,r,"/",http.StatusMovedPermanently)
	}

}


func WaitMessages(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		lastReceived,_ := strconv.ParseInt(r.FormValue("lastReceived"),10,64)
		subscription := chatroom.Subscribe()
		defer subscription.Cancel()

		// See if anything is new in the archive.
		var events []chatroom.Event
		for _, event := range subscription.Archive {
			if event.Timestamp > lastReceived {
				events = append(events, event)
			}
		}

		// If we found one, grand.
		if len(events) > 0 {
			formatter.JSON(w,http.StatusOK,events)
			return
		}
		// Else, wait for something new.
		event := <-subscription.NewMsg

		formatter.JSON(w,http.StatusOK,[]chatroom.Event{event})

		return
	}
}

