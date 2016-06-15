package offline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Compasses/GOProjects/apiservice/db"
	"github.com/Compasses/GOProjects/apiservice/utils"
	"github.com/julienschmidt/httprouter"
)

type offlinemiddleware struct {
	router   *httprouter.Router
	replaydb *db.ReplayDB
}

func NewMiddleware() *offlinemiddleware {
	router := httprouter.New()
	db, err := db.NewReplayDB()
	db.ReadDir("./input")
	if err != nil {
		log.Println("Open replayDB error ", err)
	}
	//db.SerilizeToFile()
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
	fmt.Println("try to get ", path[0], req.Method, string(newbody))
	res, err := middleware.replaydb.GetResponse(path[0], req.Method, string(newbody))
	if err != nil || res == nil {
		log.Println("Cannot get response from replaydb on offline mode, need hanle in offline handler ", err)
		newRq, err := http.NewRequest(req.Method, req.RequestURI, ioutil.NopCloser(bytes.NewReader(newbody)))
		if err != nil {
			log.Println("new http request failed ", err)
		}
		middleware.router.ServeHTTP(w, newRq)
	} else {
		result, _ := utils.TOJsonInterface(res)
		log.Println("Get response from replaydb on offline mode ", (result))
		resultmap := result.(map[string]interface{})
		for key, value := range resultmap {
			status, _ := strconv.Atoi(key)
			w.WriteHeader(status)
			stream := []byte("")
			if value != nil {
				stream, err = json.Marshal(value)
				if err != nil {
					log.Println("Marshal failed ", err, stream)
				}
			}
			_, err = w.Write(stream)
			if err != nil {
				log.Println("Get response from replaydb  but write error ", err)
			}
			break
		}
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
