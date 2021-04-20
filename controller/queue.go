package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cjlapao/servicebuscli-go/entities"
	"github.com/gorilla/mux"
)

// GetQueues Gets all queues in the namespace
func (c *Controller) GetQueues(w http.ResponseWriter, r *http.Request) {
	errorResponse := entities.ApiErrorResponse{}
	azQueues, err := sbcli.ListQueues()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Error Query"
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	queues := make([]entities.QueueEntity, 0)
	for _, azQueue := range azQueues {
		queue := entities.QueueEntity{}
		queue.FromServiceBus(azQueue)
		queues = append(queues, queue)
	}

	json.NewEncoder(w).Encode(queues)
}

// UpsertQueue Update or Insert a Queue in the current namespace
// func (c *Controller) UpsertQueue(w http.ResponseWriter, r *http.Request) {
// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	errorResponse := entities.ApiErrorResponse{}

// 	// Body cannot be null error
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		errorResponse.Code = http.StatusBadRequest
// 		errorResponse.Error = "Empty Body"
// 		errorResponse.Message = "The body of the request is null or empty"
// 		json.NewEncoder(w).Encode(errorResponse)
// 		return
// 	}

// 	topic := entities.QueueEntity{}
// 	err = json.Unmarshal(reqBody, &topic)

// 	// Body deserialization error
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		errorResponse.Code = http.StatusBadRequest
// 		errorResponse.Error = "Failed Body Deserialization"
// 		errorResponse.Message = "There was an error deserializing the body of the request"
// 		json.NewEncoder(w).Encode(errorResponse)
// 		return
// 	}

// 	isValid, validError := topic.IsValidate()

// 	if !isValid {
// 		if validError != nil {
// 			w.WriteHeader(int(validError.Code))
// 			json.NewEncoder(w).Encode(validError)
// 			return
// 		} else {
// 			w.WriteHeader(http.StatusBadRequest)
// 		}
// 	}

// 	sbTopicOptions, errResp := topic.GetOptions()

// 	if errResp != nil {
// 		w.WriteHeader(int(errResp.Code))
// 		json.NewEncoder(w).Encode(errResp)
// 		return

// 	}

// 	sbTopic, err := sbcli.CreateTopic(topic.Name, *sbTopicOptions...)

// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		errorResponse.Code = http.StatusBadRequest
// 		errorResponse.Error = "Error Creating Topic"
// 		errorResponse.Message = "There was an error creating topic"
// 		json.NewEncoder(w).Encode(errorResponse)
// 		return
// 	}

// 	topicE := entities.TopicEntity{}
// 	topicE.FromServiceBus(sbTopic)
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(topicE)
// }

func (c *Controller) SendQueueMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queueName := vars["name"]
	reqBody, err := ioutil.ReadAll(r.Body)
	errorResponse := entities.ApiErrorResponse{}

	// Topic Name cannot be nil
	if queueName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Queue name is null"
		errorResponse.Message = "Queue name cannot be null"
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

	err = sbcli.SendQueueServiceBusMessage(queueName, sbMessage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Error Sending Topic Message"
		errorResponse.Message = "There was an error sending message to queue " + queueName
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := entities.ApiSuccessResponse{
		Message: "Message " + message.Label + " was sent successfully to " + queueName + " queue",
		Data:    message.Data,
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}
