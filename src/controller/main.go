package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"net/http"

	cjlog "github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/servicebuscli/entities"
	"github.com/cjlapao/servicebuscli/servicebus"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"

	"github.com/gorilla/mux"
)

var connStr = os.Getenv("SERVICEBUS_CONNECTION_STRING")
var port = os.Getenv("SERVICEBUS_CLI_HTTP_PORT")
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
	if port == "" {
		port = "10000"
	}

	logger.Info("Api Server starting on port " + port + ".")
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.HandleFunc("/", homePage)
	_ = NewAPIController(router)
	logger.Success("API Server ready on port " + port + ".")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.FencedCode
	parser := parser.NewWithExtensions(extensions)
	md, err := ioutil.ReadFile("README.md")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	// unsanitizedHtml := markdown.ToHTML(md, parser, nil)
	htmlHeader := []byte("<html><head><link rel=\"stylesheet\" href=\"https://cdnjs.cloudflare.com/ajax/libs/prism/1.5.0/themes/prism.min.css\"</head>><body><script src=\"https://cdnjs.cloudflare.com/ajax/libs/prism/1.5.0/prism.min.js\"></script>")
	htmlBody := markdown.ToHTML(md, parser, nil)
	htmlFooter := []byte("</body></html>")
	unsanitizedHtml := make([]byte, 0)
	unsanitizedHtml = append(unsanitizedHtml, htmlHeader...)
	unsanitizedHtml = append(unsanitizedHtml, htmlBody...)
	unsanitizedHtml = append(unsanitizedHtml, htmlFooter...)

	// html := bluemonday.UGCPolicy().SanitizeBytes(unsanitizedHtml)
	w.Write(unsanitizedHtml)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.Header().Add("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

func ServiceBusConnectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("[%v] %v route requested by %v.", r.Method, r.URL.Path, r.RemoteAddr)
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
	controller.Router.HandleFunc("/config", controller.SetConnectionString).Methods("POST")
	controller.Router.HandleFunc("/topics", controller.GetTopics).Methods("GET")
	controller.Router.HandleFunc("/topics", controller.CreateTopic).Methods("POST")
	controller.Router.HandleFunc("/topics", controller.CreateTopic).Methods("PUT")
	controller.Router.HandleFunc("/topics/{topicName}", controller.GetTopic).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}", controller.DeleteTopic).Methods("DELETE")
	controller.Router.HandleFunc("/topics/{topicName}/send", controller.SendTopicMessage).Methods("PUT")
	controller.Router.HandleFunc("/topics/{topicName}/sendbulk", controller.SendBulkTopicMessage).Methods("PUT")
	controller.Router.HandleFunc("/topics/{topicName}/sendbulktemplate", controller.SendBulkTemplateTopicMessage).Methods("PUT")
	// Subscriptions Controllers
	controller.Router.HandleFunc("/topics/{topicName}/subscriptions", controller.GetTopicSubscriptions).Methods("GET")
	controller.Router.HandleFunc("/topics/{topicName}/subscriptions", controller.UpsertTopicSubscription).Methods("POST")
	controller.Router.HandleFunc("/topics/{topicName}/subscriptions", controller.UpsertTopicSubscription).Methods("PUT")
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
	controller.Router.HandleFunc("/queues", controller.UpsertQueue).Methods("POST")
	controller.Router.HandleFunc("/queues", controller.UpsertQueue).Methods("PUT")
	controller.Router.HandleFunc("/queues/{queueName}", controller.GetQueue).Methods("GET")
	controller.Router.HandleFunc("/queues/{queueName}", controller.DeleteQueue).Methods("DELETE")
	controller.Router.HandleFunc("/queues/{queueName}/send", controller.SendQueueMessage).Methods("PUT")
	controller.Router.HandleFunc("/queues/{queueName}/sendbulk", controller.SendBulkQueueMessage).Methods("PUT")
	controller.Router.HandleFunc("/queues/{queueName}/sendbulktemplate", controller.SendBulkTemplateQueueMessage).Methods("PUT")
	controller.Router.HandleFunc("/queues/{queueName}/deadletters", controller.GetQueueDeadLetterMessages).Methods("GET")
	controller.Router.HandleFunc("/queues/{queueName}/messages", controller.GetQueueMessages).Methods("GET")

	return controller
}

// GetTopics Gets all topics in the namespace
func (c *Controller) SetConnectionString(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	errorResponse := entities.ApiErrorResponse{}

	// Body cannot be null error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Empty Body"
		errorResponse.Message = "The body of the request is null or empty"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	connection := entities.ConfigRequest{}
	err = json.Unmarshal(reqBody, &connection)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Body Deserialization"
		errorResponse.Message = "There was an error deserializing the body of the request"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	if connection.ConnectionString == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Empty Connection String"
		errorResponse.Message = "Connection string cannot be empty"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	sbcli = servicebus.NewCli(connection.ConnectionString)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Error Creating Topic"
		errorResponse.Message = "There was an error creating topic"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	os.Setenv("SERVICEBUS_CONNECTION_STRING", connection.ConnectionString)
	w.WriteHeader(http.StatusAccepted)
}
