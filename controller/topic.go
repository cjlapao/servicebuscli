package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/servicebuscli-go/entities"
	"github.com/gorilla/mux"
)

// GetTopic Get Topic by name from the service bus
func (c *Controller) GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	errorResponse := entities.ApiErrorResponse{}

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	sbTopic := sbcli.GetTopicDetails(key)
	if sbTopic == nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "Topic was not found"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	topic := entities.TopicEntity{}
	topic.FromServiceBus(sbTopic)

	json.NewEncoder(w).Encode(topic)
}

// GetTopics Gets all topics in the namespace
func (c *Controller) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	errorResponse := entities.ApiErrorResponse{}

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	sbTopic := sbcli.GetTopicDetails(key)
	if sbTopic == nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "Topic was not found"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err := sbcli.DeleteTopic(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "Topic was not found"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetArticle Gets an article by it's id from the database
func (c *Controller) GetTopicSubscriptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	errorResponse := entities.ApiErrorResponse{}

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	azTopicSubscriptions, err := sbcli.ListSubscriptions(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "Topic was not found"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	subscriptions := make([]entities.SubscriptionEntity, 0)
	for _, azsubscription := range azTopicSubscriptions {
		result := entities.SubscriptionEntity{}
		result.FromServiceBus(azsubscription)
		subscriptions = append(subscriptions, result)
	}

	json.NewEncoder(w).Encode(subscriptions)
}

// GetTopics Gets all topics in the namespace
func (c *Controller) GetTopics(w http.ResponseWriter, r *http.Request) {
	errorResponse := entities.ApiErrorResponse{}
	azTopics, err := sbcli.ListTopics()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	topics := make([]entities.TopicEntity, 0)
	for _, aztopic := range azTopics {
		topic := entities.TopicEntity{}
		topic.FromServiceBus(aztopic)
		topics = append(topics, topic)
	}

	json.NewEncoder(w).Encode(topics)
}

func (c *Controller) UpsertTopic(w http.ResponseWriter, r *http.Request) {
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

	topic := entities.TopicRequest{}
	err = json.Unmarshal(reqBody, &topic)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Body Deserialization"
		errorResponse.Message = "There was an error deserializing the body of the request"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if topic.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if r.Method == http.MethodPut {
		eTopic := sbcli.GetTopic(topic.Name)
		if eTopic == nil {
			w.WriteHeader(http.StatusNotFound)
			errorResponse.Code = http.StatusNotFound
			errorResponse.Error = "Topic Not Found"
			errorResponse.Message = "Topic name cannot be found"
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	var opts []servicebus.TopicManagementOption
	opts = make([]servicebus.TopicManagementOption, 0)
	if topic.Options != nil {
		if topic.Options.AutoDeleteOnIdle != nil {
			d, err := time.ParseDuration(*topic.Options.AutoDeleteOnIdle)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the AutoDeleteOnIdle from string"
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
			opts = append(opts, servicebus.TopicWithAutoDeleteOnIdle(&d))
		}
		if topic.Options.EnableBatchedOperation != nil {
			opts = append(opts, servicebus.TopicWithBatchedOperations())
		}
		if topic.Options.EnableDuplicateDetection != nil {
			d, err := time.ParseDuration(*topic.Options.AutoDeleteOnIdle)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the EnableDuplicateDetection from string"
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
			opts = append(opts, servicebus.TopicWithDuplicateDetection(&d))
		}
		if topic.Options.EnableExpress != nil {
			opts = append(opts, servicebus.TopicWithExpress())
		}
		if topic.Options.MaxSizeInMegabytes != nil {
			opts = append(opts, servicebus.TopicWithMaxSizeInMegabytes(*topic.Options.MaxSizeInMegabytes))
		}
		if topic.Options.DefaultMessageTimeToLive != nil {
			d, err := time.ParseDuration(*topic.Options.AutoDeleteOnIdle)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the DefaultMessageTimeToLive from string"
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
			opts = append(opts, servicebus.TopicWithMessageTimeToLive(&d))
		}
		if topic.Options.SupportOrdering != nil {
			opts = append(opts, servicebus.TopicWithOrdering())
		}
		if topic.Options.EnablePartitioning != nil {
			opts = append(opts, servicebus.TopicWithPartitioning())
		}
	}

	sbTopic, err := sbcli.CreateTopic(topic.Name, opts...)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Error Creating Topic"
		errorResponse.Message = "There was an error creating topic"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	topicE := entities.TopicEntity{}
	topicE.FromServiceBus(sbTopic)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topicE)
}
