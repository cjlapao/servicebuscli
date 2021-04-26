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
	// Topics Controllers
	controller.Router.HandleFunc("/topics", controller.GetTopics).Methods("GET")
	controller.Router.HandleFunc("/topics", controller.CreateTopic).Methods("POST")
	controller.Router.HandleFunc("/topics", controller.CreateTopic).Methods("PUT")
	controller.Router.HandleFunc("/topics/{topicName}", controller.GetTopic).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}", controller.DeleteTopic).Methods("DELETE")
	controller.Router.HandleFunc("/topics/{topicName}/send", controller.SendTopicMessage).Methods("PUT")
	// Subscriptions Controllers
	controller.Router.HandleFunc("/topics/{topicName}/subscriptions", controller.GetTopicSubscriptions).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}/subscriptions", controller.CreateTopicSubscription).Methods("POST")
	controller.Router.HandleFunc("/topics/{topicName}/subscriptions", controller.CreateTopicSubscription).Methods("PUT")
	controller.Router.HandleFunc("/topics/{topicName}/{subscriptionName}", controller.GetTopicSubscription).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}/{subscriptionName}", controller.DeleteTopicSubscription).Methods("DELETE")
	controller.Router.HandleFunc("/topics/{topicName}/{subscriptionName}/deadletters", controller.GetSubscriptionDeadLetterMessages).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}/{subscriptionName}/messages", controller.GetSubscriptionMessages).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}/{subscriptionName}/rules", controller.GetSubscriptionRules).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}/{subscriptionName}/rules", controller.CreateSubscriptionRule).Methods("POST")
	controller.Router.HandleFunc("/topics/{topicName}/{subscriptionName}/rules/{ruleName}", controller.GetSubscriptionRule).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}/{subscriptionName}/rules/{ruleName}", controller.DeleteSubscriptionRule).Methods("DELETE")
	// Queues Controllers
	controller.Router.HandleFunc("/queues", controller.GetQueues).Methods("GET")
	controller.Router.HandleFunc("/queues/{queueName}/send", controller.SendQueueMessage).Methods("PUT")

	return controller
}
