package refresh

import (
	"server/chatroom"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/unrolled/render"
)

func Join(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")
		chatroom.Join(user)
		http.Redirect(w,r,"/refresh/room",http.StatusMovedPermanently)

	}
}

func Room(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")
		subscription := chatroom.Subscribe()
		defer subscription.Cancel()
		events := subscription.Archive
		for i, _ := range events {
			if events[i].User == user {
				events[i].User = "you"
			}
		}
		data := struct {
			User string
			Events []chatroom.Event
		}{user,events}
		formatter.HTML(w,http.StatusOK,"refresh/room",data)
	}
}

func Say(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")
		message := r.FormValue("message")
		chatroom.Say(user, message)
		http.Redirect(w,r,"/refresh/room",http.StatusMovedPermanently)

	}
}

func Leave(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user := r.FormValue("user")
		chatroom.Leave(user)
		http.Redirect(w,r,"/",http.StatusMovedPermanently)
	}
}
