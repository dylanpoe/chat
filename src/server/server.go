package server

import (
	"github.com/urfave/negroni"
	"github.com/unrolled/render"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"server/refresh"
	"server/longpolling"
	ws "server/websocket"
)

func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.New()

	n.Use(negroni.NewLogger())

	n.Use(negroni.NewRecovery())

	router := httprouter.New()

	router.GET("/",IndexHandler(formatter))

	router.GET("/demo",DemoHandler(formatter))

	// Refresh
	router.GET("/refresh",refresh.Join(formatter))
	router.GET("/refresh/room",refresh.Room(formatter))
	router.POST("/refresh/room",refresh.Say(formatter))
	router.GET("/refresh/room/leave",refresh.Leave(formatter))

	// Long polling demo

	router.GET("/longpolling/room",longpolling.Room(formatter))
	router.GET("/longpolling/room/messages",longpolling.WaitMessages(formatter))
	router.POST("/longpolling/room/messages",longpolling.Say(formatter))
	router.GET("/longpolling/room/leave",longpolling.Leave(formatter))

	// WebSocket

	router.GET("/websocket/room",ws.Room(formatter))
	router.GET("/websocket/room/socket",ws.WsHandler(formatter))
	router.ServeFiles("/static/*filepath",http.Dir("static"))

	n.UseHandler(router)

	return n

}

func DemoHandler(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		demo := r.FormValue("demo")
		user := r.FormValue("user")
		switch demo {
		case "refresh":
			http.Redirect(w,r, "/refresh?user=" + user, http.StatusMovedPermanently)
		case "longpolling":
			http.Redirect(w,r, "/longpolling/room?user=" + user, http.StatusMovedPermanently)
		case "websocket":
			http.Redirect(w,r,"/websocket/room?user=" + user,http.StatusMovedPermanently)
		default:
			http.Redirect(w,r, "/websocket/room?user=" + user, http.StatusMovedPermanently)
		}

	}
}

func IndexHandler(formatter *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		formatter.HTML(w,http.StatusOK,"index", nil)
		return
	}
}