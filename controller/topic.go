package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

	topic := entities.TopicRequestEntity{}
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

	isValid, validError := topic.IsValidate()

	if !isValid {
		if validError != nil {
			w.WriteHeader(int(validError.Code))
			json.NewEncoder(w).Encode(validError)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	sbTopicOptions, errResp := topic.GetOptions()

	if errResp != nil {
		w.WriteHeader(int(errResp.Code))
		json.NewEncoder(w).Encode(errResp)
		return

	}

	sbTopic, err := sbcli.CreateTopic(topic.Name, *sbTopicOptions...)

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
