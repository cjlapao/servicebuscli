package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/servicebuscli-go/entities"
	"github.com/gorilla/mux"
)

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

// GetTopic Get Topic by name from the namespace
func (c *Controller) GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	errorResponse := entities.ApiErrorResponse{}

	if topicName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	sbTopic := sbcli.GetTopicDetails(topicName)
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

// DeleteTopic Deletes a topic in the namespace
func (c *Controller) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	errorResponse := entities.ApiErrorResponse{}

	if topicName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	sbTopic := sbcli.GetTopicDetails(topicName)
	if sbTopic == nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "Topic was not found"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err := sbcli.DeleteTopic(topicName)
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

func (c *Controller) CreateTopic(w http.ResponseWriter, r *http.Request) {
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
	var sbTopic *servicebus.TopicEntity
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

	topicExists := sbcli.GetTopic(topic.Name)

	if topicExists != nil {
		w.WriteHeader(http.StatusBadRequest)
		found := entities.ApiSuccessResponse{
			Message: "The topic " + topic.Name + " already exists, ignoring",
		}
		json.NewEncoder(w).Encode(found)
		return
	}

	sbTopic, err = sbcli.CreateTopic(topic.Name, *sbTopicOptions...)

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

func (c *Controller) SendTopicMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["name"]
	reqBody, err := ioutil.ReadAll(r.Body)
	errorResponse := entities.ApiErrorResponse{}

	// Topic Name cannot be nil
	if topicName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Body cannot be nil error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Empty Body"
		errorResponse.Message = "The body of the request is null or empty"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	message := entities.ServiceBusMessage{}
	err = json.Unmarshal(reqBody, &message)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Body Deserialization"
		errorResponse.Message = "There was an error deserializing the body of the request"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	sbMessage, err := message.ToServiceBus()

	// Convert to ServiceBus Message error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Conversion"
		errorResponse.Message = "There was an error converting the request to a service bus message"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = sbcli.SendTopicServiceBusMessage(topicName, sbMessage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Error Sending Topic Message"
		errorResponse.Message = "There was an error sending message to topic " + topicName
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := entities.ApiSuccessResponse{
		Message: "Message " + message.Label + " was sent successfully to " + topicName + " topic",
		Data:    message.Data,
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}
