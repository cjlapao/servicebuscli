package controller

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/servicebuscli-go/entities"
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
