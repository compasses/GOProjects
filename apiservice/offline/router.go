package offline

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Compasses/GOProjects/apiservice/db"
	"github.com/julienschmidt/httprouter"
)

type offlinemiddleware struct {
	router   *httprouter.Router
	replaydb *db.ReplayDB
}

func NewMiddleware() *offlinemiddleware {
	router := httprouter.New()
	db, err := db.NewReplayDB()
	if err != nil {
		log.Println("Open replayDB error ", err)
	}

	for _, route := range routes {
		httpHandle := Logger(route.HandleFunc, route.Name)

		router.Handle(
			route.Method,
			route.Pattern,
			httpHandle,
		)
	}

	router.NotFound = LoggerNotFound(NotFoundHandler)

	return &offlinemiddleware{
		router:   router,
		replaydb: db,
	}
}

func (middleware *offlinemiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newbody := make([]byte, req.ContentLength)
	req.Body.Read(newbody)
	path := strings.Split(req.RequestURI, "?")

	res, err := middleware.replaydb.GetResponse(path[0], req.Method, string(newbody))
	if err != nil || len(res) == 0 {
		log.Println("Cannot get response from replaydb on offline mode ", err)
		newRq, err := http.NewRequest(req.Method, req.RequestURI, ioutil.NopCloser(bytes.NewReader(newbody)))
		if err != nil {
			log.Println("new http request failed ", err)
		}
		middleware.router.ServeHTTP(w, newRq)
	} else {
		log.Println("Get response from replaydb on offline mode ", string(res))
		w.Write(res)
	}
}

func ServerRouter() *httprouter.Router {
	router := httprouter.New()

	for _, route := range routes {
		httpHandle := Logger(route.HandleFunc, route.Name)

		router.Handle(
			route.Method,
			route.Pattern,
			httpHandle,
		)
	}

	router.NotFound = LoggerNotFound(NotFoundHandler)

	return router
}
