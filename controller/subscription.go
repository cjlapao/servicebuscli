package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cjlapao/servicebuscli-go/entities"
	"github.com/gorilla/mux"
)

// GetTopicSubscriptions Gets All the subscriptions in a topic
func (c *Controller) GetTopicSubscriptions(w http.ResponseWriter, r *http.Request) {
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

	azTopic := sbcli.GetTopicDetails(topicName)
	if azTopic == nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "The Topic " + topicName + " was not found in the service bus " + sbcli.Namespace.Name
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	azTopicSubscriptions, err := sbcli.ListSubscriptions(topicName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Error Getting Subscriptions"
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	subscriptions := make([]entities.SubscriptionResponseEntity, 0)
	for _, azSubscription := range azTopicSubscriptions {
		result := entities.SubscriptionResponseEntity{}
		result.FromServiceBus(azSubscription)
		subscriptions = append(subscriptions, result)
	}

	json.NewEncoder(w).Encode(subscriptions)
}

// GetTopicSubscription Gets a subscription from a topic in the namespace
func (c *Controller) GetTopicSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	subscriptionName := vars["subscriptionName"]
	errorResponse := entities.ApiErrorResponse{}

	// Checking for null parameters
	if topicName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if subscriptionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Subscription Name is null"
		errorResponse.Message = "Subscription Name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Checks if the topic exists, if not issuing an error
	azTopic := sbcli.GetTopicDetails(topicName)
	if azTopic == nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "The Topic " + topicName + " was not found in the service bus " + sbcli.Namespace.Name
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	azSubscription, err := sbcli.GetSubscription(topicName, subscriptionName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Subscription not found"
		errorResponse.Message = "The Subscription " + subscriptionName + " was not found on " + topicName + " topic in the service bus " + sbcli.Namespace.Name
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	apiSubscription := entities.SubscriptionResponseEntity{}
	apiSubscription.FromServiceBus(azSubscription)
	json.NewEncoder(w).Encode(apiSubscription)
}

// UpsertTopicSubscription Creates a subscription in a topic
func (c *Controller) UpsertTopicSubscription(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	errorResponse := entities.ApiErrorResponse{}
	upsert := false
	if r.Method == "PUT" {
		upsert = true
	}

	// Checking for null parameters
	if topicName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Body cannot be null error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Empty Body"
		errorResponse.Message = "The body of the request is null or empty"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	bodySubscription := entities.SubscriptionRequestEntity{}
	err = json.Unmarshal(reqBody, &bodySubscription)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Body Deserialization"
		errorResponse.Message = "There was an error deserializing the body of the request"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	bodySubscription.TopicName = topicName
	isValid, validError := bodySubscription.ValidateSubscriptionRequest()
	if !isValid {
		if validError != nil {
			w.WriteHeader(int(validError.Code))
			json.NewEncoder(w).Encode(validError)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	_, errResp := bodySubscription.GetOptions()

	if errResp != nil {
		w.WriteHeader(int(errResp.Code))
		json.NewEncoder(w).Encode(errResp)
		return

	}

	if !upsert {
		subscriptionExists, _ := sbcli.GetSubscription(bodySubscription.TopicName, bodySubscription.Name)

		if subscriptionExists != nil {
			w.WriteHeader(http.StatusBadRequest)
			found := entities.ApiSuccessResponse{
				Message: "The Subscription " + bodySubscription.Name + " already exists in topic " + bodySubscription.TopicName + ", ignoring",
			}
			json.NewEncoder(w).Encode(found)
			return
		}

	}

	err = sbcli.CreateSubscription(bodySubscription, upsert)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Error Creating Subscription"
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	subscriptionE := entities.SubscriptionResponseEntity{}
	createdSubscription, err := sbcli.GetSubscription(bodySubscription.TopicName, bodySubscription.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Error Creating Subscription"
		errorResponse.Message = "There was an error creating subscription"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	subscriptionE.FromServiceBus(createdSubscription)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscriptionE)
}

// DeleteTopicSubscription Deletes subscription from a topic in the namespace
func (c *Controller) DeleteTopicSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	subscriptionName := vars["subscriptionName"]
	errorResponse := entities.ApiErrorResponse{}

	if topicName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if subscriptionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Subscription Name is null"
		errorResponse.Message = "Subscription Name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	sbTopic := sbcli.GetTopicDetails(topicName)
	if sbTopic == nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "The Topic " + topicName + " was not found in the service bus " + sbcli.Namespace.Name
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	sbSubscription, err := sbcli.GetSubscription(topicName, subscriptionName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Topic not found"
		errorResponse.Message = "The Subscription " + subscriptionName + " was not found on topic " + topicName + " in the service bus " + sbcli.Namespace.Name
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = sbcli.DeleteSubscription(topicName, subscriptionName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Error Deleting Subscription"
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := entities.SubscriptionResponseEntity{}
	response.FromServiceBus(sbSubscription)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

// GetSubscriptionMessages Gets messages from a topic subscription
func (c *Controller) GetSubscriptionMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	subscriptionName := vars["subscriptionName"]
	queryValues := r.URL.Query()
	qtyValue := queryValues.Get("qty")
	peekValue := queryValues.Get("peek")
	if qtyValue == "" {
		qtyValue = "0"
	}

	qty, qtyErr := strconv.Atoi(qtyValue)
	if qtyErr != nil {
		qty = 0
	}

	peek := false
	if peekValue == "true" {
		peek = true
	}

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

	// Subscription Name cannot be nil
	if subscriptionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	result, err := sbcli.GetSubscriptionActiveMessages(topicName, subscriptionName, qty, peek)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Message Data Deserialization"
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		response := make([]entities.ServiceBusMessageRequest, 0)
		for _, msg := range result {
			entityMsg := entities.ServiceBusMessageRequest{}
			entityMsg.FromServiceBus(&msg)
			response = append(response, entityMsg)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
}

// GetSubscriptionMessages Gets messages from a topic subscription
func (c *Controller) GetSubscriptionDeadLetterMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	subscriptionName := vars["subscriptionName"]
	queryValues := r.URL.Query()
	qtyValue := queryValues.Get("qty")
	peekValue := queryValues.Get("peek")
	if qtyValue == "" {
		qtyValue = "0"
	}

	qty, qtyErr := strconv.Atoi(qtyValue)
	if qtyErr != nil {
		qty = 0
	}

	peek := false
	if peekValue == "true" {
		peek = true
	}

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

	// Subscription Name cannot be nil
	if subscriptionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	result, err := sbcli.GetSubscriptionDeadLetterMessages(topicName, subscriptionName, qty, peek)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Message Data Deserialization"
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		response := make([]entities.ServiceBusMessageRequest, 0)
		for _, msg := range result {
			entityMsg := entities.ServiceBusMessageRequest{}
			entityMsg.FromServiceBus(&msg)
			response = append(response, entityMsg)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
}

// GetSubscriptionMessages Gets messages from a topic subscription
func (c *Controller) GetSubscriptionRules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	subscriptionName := vars["subscriptionName"]
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

	// Subscription Name cannot be nil
	if subscriptionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	result, err := sbcli.GetSubscriptionRules(topicName, subscriptionName)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed to GET Subscription Rules"
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		response := make([]entities.RuleResponseEntity, 0)
		for _, msg := range result {
			entityMsg := entities.RuleResponseEntity{}
			entityMsg.FromServiceBus(msg)
			response = append(response, entityMsg)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
}

// GetSubscriptionMessages Gets messages from a topic subscription
func (c *Controller) GetSubscriptionRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	subscriptionName := vars["subscriptionName"]
	ruleName := vars["ruleName"]
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

	// Subscription Name cannot be nil
	if subscriptionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Subscription Name cannot be nil
	if ruleName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Rule Name is Null"
		errorResponse.Message = "Rule Name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	result, err := sbcli.GetSubscriptionRule(topicName, subscriptionName, ruleName)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Failed to GET Subscription Rule " + ruleName
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := entities.RuleResponseEntity{}
	response.FromServiceBus(result)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetSubscriptionMessages Gets messages from a topic subscription
func (c *Controller) CreateSubscriptionRule(w http.ResponseWriter, r *http.Request) {
	errorResponse := entities.ApiErrorResponse{}
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	subscriptionName := vars["subscriptionName"]

	reqBody, err := ioutil.ReadAll(r.Body)
	// Body cannot be null error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Empty Body"
		errorResponse.Message = "The body of the request is null or empty"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ruleRequest := entities.RuleRequestEntity{}
	err = json.Unmarshal(reqBody, &ruleRequest)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Failed Body Deserialization"
		errorResponse.Message = "There was an error deserializing the body of the request"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Topic Name cannot be nil
	if topicName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Subscription Name cannot be nil
	if subscriptionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	_, err = sbcli.GetSubscription(topicName, subscriptionName)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Failed to GET Subscription " + subscriptionName
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	subscription := entities.SubscriptionRequestEntity{}
	subscription.TopicName = topicName
	subscription.Name = subscriptionName

	err = sbcli.CreateSubscriptionRule(subscription, ruleRequest)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Failed to CREATE Subscription Rule" + ruleRequest.Name
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	rule, err := sbcli.GetSubscriptionRule(topicName, subscriptionName, ruleRequest.Name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Failed to CREATE Subscription Rule" + ruleRequest.Name
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := entities.RuleResponseEntity{}
	response.FromServiceBus(rule)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetSubscriptionMessages Gets messages from a topic subscription
func (c *Controller) DeleteSubscriptionRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]
	subscriptionName := vars["subscriptionName"]
	ruleName := vars["ruleName"]
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

	// Subscription Name cannot be nil
	if subscriptionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Subscription Name cannot be nil
	if ruleName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Rule Name is Null"
		errorResponse.Message = "Rule Name cannot be null"
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	result, err := sbcli.DeleteSubscriptionRule(topicName, subscriptionName, ruleName)

	// Body deserialization error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse.Code = http.StatusNotFound
		errorResponse.Error = "Failed to DELETE Subscription Rule " + ruleName
		errorResponse.Message = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := entities.RuleResponseEntity{}
	response.FromServiceBus(result)
	w.WriteHeader(http.StatusAccepted)
}
