package controller

import (
	"fmt"
	"log"
	"os"

	cjlog "github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/servicebuscli-go/servicebus"

	"net/http"

	"github.com/gorilla/mux"
)

var connStr = os.Getenv("SERVICEBUS_CONNECTION_STRING")
var router mux.Router
var logger = cjlog.Get()
var ver = version.Get()
var sbcli = servicebus.NewCli(connStr)

// Controllers Controller structure
type Controller struct {
	Router *mux.Router
}

func RestApiModuleProcessor() {
	logger.Notice("Starting Service Bus Client API module v%v", ver.String())
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.HandleFunc("/", homePage)
	_ = NewAPIController(router)
	logger.Success("Finished Init")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Homepage!")
	fmt.Println("endpoint Hit: homepage")
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func ServiceBusConnectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("[%v] %v route...", r.Method, r.URL.Path)
		if connStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No Azure Service Bus Connection String defined"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NewAPIController  Creates a new controller
func NewAPIController(router *mux.Router) Controller {
	controller := Controller{
		Router: router,
	}

	controller.Router.Use(ServiceBusConnectionMiddleware)
	controller.Router.Use(commonMiddleware)
	controller.Router.HandleFunc("/topics", controller.GetTopics).Methods("GET")
	controller.Router.HandleFunc("/topics/{name}", controller.GetTopic).Methods("GET")
	controller.Router.HandleFunc("/topics/{name}", controller.DeleteTopic).Methods("DELETE")
	controller.Router.HandleFunc("/topics/{name}/subscriptions", controller.GetTopicSubscriptions).Methods("GET")
	controller.Router.HandleFunc("/topics", controller.UpsertTopic).Methods("POST")

	controller.Router.HandleFunc("/queues", controller.GetQueues).Methods("GET")

	return controller
}
