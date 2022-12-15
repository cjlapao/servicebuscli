package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/servicebuscli/entities"
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

	topics := make([]entities.TopicResponseEntity, 0)
	for _, aztopic := range azTopics {
		topic := entities.TopicResponseEntity{}
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

	topic := entities.TopicResponseEntity{}
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

	isValid, validError := topic.IsValid()
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

	topicE := entities.TopicResponseEntity{}
	topicE.FromServiceBus(sbTopic)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topicE)
}

func (c *Controller) SendTopicMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
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

	message := entities.MessageRequest{}
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

func (c *Controller) SendBulkTopicMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
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

	bulk := entities.BulkMessageRequest{}
	err = json.Unmarshal(reqBody, &bulk)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Body Deserialization"
		errorResponse.Message = "There was an error deserializing the body of the request"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = sbcli.SendBulkTopicMessage(topicName, bulk.Messages...)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Error Sending Topic Message"
		errorResponse.Message = "There was an error sending bulk messages to topic " + topicName
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := entities.ApiSuccessResponse{
		Message: "Sent " + fmt.Sprint(len(bulk.Messages)) + " Messages successfully to " + topicName + " topic",
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) SendBulkTemplateTopicMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	reqBody, err := ioutil.ReadAll(r.Body)
	errorResponse := entities.ApiErrorResponse{}
	maxSizeOfPayload := 262144

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

	bulk := entities.BulkTemplateMessageRequest{}
	err = json.Unmarshal(reqBody, &bulk)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Body Deserialization"
		errorResponse.Message = "There was an error deserializing the body of the request"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if bulk.BatchOf > bulk.TotalMessages {
		bulk.BatchOf = 1
	}

	totalMessageSent := 0

	messageTemplate, err := json.Marshal(bulk.Template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = err.Error()
		errorResponse.Message = "There was an error checking the templated message payload size"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	messageTemplateSize := len(messageTemplate)
	totalMessagePayloadSize := messageTemplateSize * bulk.TotalMessages
	minimumBatchSize := totalMessagePayloadSize / maxSizeOfPayload
	if minimumBatchSize > bulk.BatchOf {
		bulk.BatchOf = minimumBatchSize + 5
		logger.Info("The total payload was too big for the given batch size, increasing it to the minimum of " + fmt.Sprint(bulk.BatchOf))
	}

	batchSize := bulk.TotalMessages / bulk.BatchOf
	messages := make([]entities.MessageRequest, 0)

	for i := 0; i < batchSize; i++ {
		messages = append(messages, bulk.Template)
	}

	if bulk.WaitBetweenBatchesInMilli > 0 {
		for i := 0; i < bulk.BatchOf; i++ {
			if i > 0 && bulk.WaitBetweenBatchesInMilli > 0 {
				logger.Info("Waiting " + fmt.Sprint(bulk.WaitBetweenBatchesInMilli) + "ms for next batch")
				time.Sleep(time.Duration(bulk.WaitBetweenBatchesInMilli) * time.Millisecond)
			}

			err = sbcli.SendBulkTopicMessage(topicName, messages...)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = err.Error()
				errorResponse.Message = "There was an error sending bulk messages to topic " + topicName
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
			totalMessageSent += len(messages)
		}
	} else {
		logger.Info("Sending all batches without waiting")
		var waitFor sync.WaitGroup
		waitFor.Add(bulk.BatchOf)
		for i := 0; i < bulk.BatchOf; i++ {
			go sbcli.SendParallelBulkTopicMessage(&waitFor, topicName, messages...)
			totalMessageSent += len(messages)
		}

		waitFor.Wait()
		logger.Success("Finished sending all the batches to service bus")
	}

	if totalMessageSent < bulk.TotalMessages {
		missingMessageCount := bulk.TotalMessages - totalMessageSent
		logger.Info("Sending remaining " + fmt.Sprint(missingMessageCount) + " messages in the payload")
		messages := make([]entities.MessageRequest, 0)
		for x := 0; x < missingMessageCount; x++ {
			messages = append(messages, bulk.Template)
		}

		err = sbcli.SendBulkTopicMessage(topicName, messages...)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse.Code = http.StatusBadRequest
			errorResponse.Error = err.Error()
			errorResponse.Message = "There was an error sending bulk messages to topic " + topicName
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		bulk.BatchOf += 1
	}

	response := entities.ApiSuccessResponse{
		Message: "Sent " + fmt.Sprint(bulk.TotalMessages) + " Messages in " + fmt.Sprint(bulk.BatchOf) + " batches successfully to " + topicName + " topic",
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}
