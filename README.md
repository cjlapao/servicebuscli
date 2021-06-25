# Azure Service Bus Command Line Tool

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Build](https://github.com/cjlapao/servicebuscli-go/actions/workflows/build.yml/badge.svg)](https://github.com/cjlapao/servicebuscli-go/actions/workflows/build.yml) [![Release](https://github.com/cjlapao/servicebuscli-go/actions/workflows/release.yml/badge.svg)](https://github.com/cjlapao/servicebuscli-go/actions/workflows/release.yml) [![CodeQL](https://github.com/cjlapao/servicebuscli-go/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/cjlapao/servicebuscli-go/actions/workflows/codeql-analysis.yml)

## Index

- [Azure Service Bus Command Line Tool](#azure-service-bus-command-line-tool)
  - [Index](#index)
  - [**How to Use it**](#how-to-use-it)
  - [API Mode](#api-mode)
    - [[GET] /topics](#get-topics)
    - [[POST] /topics](#post-topics)
    - [[GET] /topics/{topic_name}](#get-topicstopic_name)
    - [[DELETE] /topics/{topic_name}](#delete-topicstopic_name)
    - [[PUT] /topics/{topic_name}/send](#put-topicstopic_namesend)
    - [[PUT] /topics/{topic_name}/sendbulk](#put-topicstopic_namesendbulk)
    - [[PUT] /topics/{topic_name}/sendbulktemplate](#put-topicstopic_namesendbulktemplate)
    - [[GET] /topics/{topic_name}/subscriptions](#get-topicstopic_namesubscriptions)
    - [[POST] /topics/{topic_name}/subscriptions](#post-topicstopic_namesubscriptions)
    - [[GET] /topics/{topic_name}/{subscription_name}](#get-topicstopic_namesubscription_name)
    - [[DELETE] /topics/{topic_name}/{subscription_name}](#delete-topicstopic_namesubscription_name)
    - [[GET] /topics/{topic_name}/{subscription_name}/deadletters](#get-topicstopic_namesubscription_namedeadletters)
    - [[GET] /topics/{topic_name}/{subscription_name}/messages](#get-topicstopic_namesubscription_namemessages)
    - [[GET] /topics/{topic_name}/{subscription_name}/rules](#get-topicstopic_namesubscription_namerules)
    - [[POST] /topics/{topic_name}/{subscription_name}/rules](#post-topicstopic_namesubscription_namerules)
    - [[GET] /topics/{topic_name}/{subscription_name}/rules/{rule_name}](#get-topicstopic_namesubscription_namerulesrule_name)
    - [[DELETE] /topics/{topic_name}/{subscription_name}/rules/{rule_name}](#delete-topicstopic_namesubscription_namerulesrule_name)
    - [[GET] /queues](#get-queues)
    - [[POST] /queues](#post-queues)
    - [[GET] /queues/{queue_name}](#get-queuesqueue_name)
    - [[DELETE] /queues/{queue_name}](#delete-queuesqueue_name)
    - [[PUT] /queues/{queue_name}/send](#put-queuesqueue_namesend)
    - [[PUT] /topics/{queue_name}/sendbulk](#put-topicsqueue_namesendbulk)
    - [[PUT] /topics/{queue_name}/sendbulktemplate](#put-topicsqueue_namesendbulktemplate)
    - [[GET] /queues/{queue_name}/deadletters](#get-queuesqueue_namedeadletters)
    - [[GET] /queues/{queue_name}/messages](#get-queuesqueue_namemessages)
  - [Topics](#topics)
    - [List Topics](#list-topics)
    - [Create Topic](#create-topic)
    - [Delete Topic](#delete-topic)
    - [List Subscription for a Topic](#list-subscription-for-a-topic)
    - [Create Topic Subscription](#create-topic-subscription)
    - [Delete Topic Subscription](#delete-topic-subscription)
    - [Subscribe to a Topic Subscription](#subscribe-to-a-topic-subscription)
    - [Send a Message to a Topic](#send-a-message-to-a-topic)
  - [Queues](#queues)
    - [List Queues](#list-queues)
    - [Create Queue](#create-queue)
    - [Delete Queue](#delete-queue)
    - [Subscribe to a Queue](#subscribe-to-a-queue)
    - [Send a Message to a Queue](#send-a-message-to-a-queue)

This is a command line tool to help test service bus messages.

You will be able to do **C** *U* **RD** operations to topics/subscriptions and queues, you will also be able to send messages and subscribe to a specific queue/subscription

## **How to Use it**

Once you have it compiled you can run it with the needed options

you can also run the following command to display the help

```bash
servicebus.exe --help
```

## API Mode

The ServiceBus Client contains an API mode that gives the same functionality but using a REST api
To start the client in API mode run the following command.

```bash
servicebus.exe api
```

### [GET] /topics

Returns all the topics in the namespace

### [POST] /topics

Creates a Topic in the namespace

Example Payload:

```json
{
    "name": "example",
    "options": {
        "autoDeleteOnIdle": "24h",
        "enableBatchedOperation": true,        
        "enableDuplicateDetection": "30m",        
        "enableExpress": false,
        "maxSizeInMegabytes": 10,
        "defaultMessageTimeToLive": "1d",
        "supportOrdering": true,
        "enablePartitioning": true
    }
}
```

### [GET] /topics/{topic_name}

Returns the details of a specific topic in the namespace

### [DELETE] /topics/{topic_name}

Deletes a specific topic from the namespace, this will also delete any subscriptions and messages in the same topic

### [PUT] /topics/{topic_name}/send

Sends a message to the specific topic

Example Payload:

```json
{
    "label": "example",
    "correlationId": "test",
    "contentType": "application/json",
    "data": {
        "key": "value"
    },
    "userProperties": {
        "name": "test message"
    }
}
```

### [PUT] /topics/{topic_name}/sendbulk

Sends bulk messages to the specific topic

Example Payload:

```json
{
     "messages": [
         {
            "label": "example",
            "correlationId": "test",
            "contentType": "application/json",
            "data": {
                "key": "value1"
            },
            "userProperties": {
                "name": "test message1"
            }
         },
         {
            "label": "example",
            "correlationId": "test",
            "contentType": "application/json",
            "data": {
                "key": "value2"
            },
            "userProperties": {
                "name": "test message2"
            }
         }
     ]
}
```

### [PUT] /topics/{topic_name}/sendbulktemplate

Sends a templated bulk messages to the specific topic, it can set a wait time between batches  
You can define in how many batches you want to send the amount of message and a wait time between each batches.  

**Attention**: if the volume of messages is big, the messages will be split in batches automatically, this also happens if the batch is too small for the maximum allowed size of a payload

Example Payload:

```json
{
    "totalMessages": 50, // Number of total messages to send, if not defined it will be set to 1
    "batchOf": 5, // send the total message in how many batches, if not defined it will be set to 1
    "waitBetweenBatchesInMilli": 500, // wait between any batches, if not defined it will be 0
    "template": { // Message template
        "label": "example",
        "correlationId": "test",
        "contentType": "application/json",
        "data": {
            "key": "value2"
        },
        "userProperties": {
            "name": "test message2"
        }
    }
}
```

### [GET] /topics/{topic_name}/subscriptions

Returns all the subscriptions in the specific topic

### [POST] /topics/{topic_name}/subscriptions

Creates a subscription in the specific topic

Example Payload:

```json
{
    "name": "wiretap",
    "topicName": "example",
    "maxDeliveryCount": 5,
    "forward": {
        "to": "otherTopic",
        "in": "Topic"
    },
    "forwardDeadLetter": {
        "to": "otherQueue",
        "in": "Queue"
    },
    "rules": [
        {
            "name": "example_rule",
            "sqlFilter": "2=2",
            "sqlAction": "SET A='one'"
        }
    ],    
    "options":{
        "autoDeleteOnIdle": "24h",
        "defaultMessageTimeToLive": "1d",
        "lockDuration": "30s",
        "enableBatchedOperation": true,
        "deadLetteringOnMessageExpiration": false,
        "requireSession": false
    }
}
```

### [GET] /topics/{topic_name}/{subscription_name}

Returns a subscription detail from a specific topic

### [DELETE] /topics/{topic_name}/{subscription_name}

Deletes a subscription from a specific topic

### [GET] /topics/{topic_name}/{subscription_name}/deadletters

Gets the dead letters from a subscription in a topic

**Query Attributes**  
*qty*, *integer*: amount of messages to collect, defaults to all with a maximum of 100 messages  
*peek*, *bool*: sets the collection mode to peek, messages will remain in the subscription, defaults to false

### [GET] /topics/{topic_name}/{subscription_name}/messages

Gets the dead letters from a subscription in a topic

**Query Attributes**  
*qty*, *integer*: amount of messages to collect, defaults to all with a maximum of 100 messages  
*peek*, *bool*: sets the collection mode to peek, messages will remain in the subscription, defaults to false

### [GET] /topics/{topic_name}/{subscription_name}/rules

Gets all the rules in a subscription

### [POST] /topics/{topic_name}/{subscription_name}/rules

Creates a rule in a subscription

Example Payload:

```json
{
    "name": "example_rule",
    "sqlFilter": "2=2",
    "sqlAction": "SET A='one'"
}
```

### [GET] /topics/{topic_name}/{subscription_name}/rules/{rule_name}

Gets the details of a specific rule in a subscription

### [DELETE] /topics/{topic_name}/{subscription_name}/rules/{rule_name}

Deletes a specific rule in a subscription

### [GET] /queues

Returns all the queues in the namespace

### [POST] /queues

Creates a queue in the namespace

Example Payload:

```json
{
    "name": "example",
    "maxDeliveryCount": 5,
    "forward": {
        "to": "otherTopic",
        "in": "Topic"
    },
    "forwardDeadLetter": {
        "to": "otherQueue",
        "in": "Queue"
    },
    "options": {
        "autoDeleteOnIdle": "24h",
        "enableDuplicateDetection": "30m",        
        "maxSizeInMegabytes": 10,
        "defaultMessageTimeToLive": "1d",
        "lockDuration": "30s",
        "supportOrdering": true,
        "enablePartitioning": true,
        "requireSession": false,
        "deadLetteringOnMessageExpiration": false
    }
}
```

### [GET] /queues/{queue_name}

Returns the details of a specific queue in the namespace

### [DELETE] /queues/{queue_name}

Deletes a specific queue from the namespace, this will also delete any subscriptions and messages in the same topic

### [PUT] /queues/{queue_name}/send

Sends a message to the specific queue

Example Payload:

```json
{
    "label": "example",
    "correlationId": "test",
    "contentType": "application/json",
    "data": {
        "key": "value"
    },
    "userProperties": {
        "name": "test message"
    }
}
```

### [PUT] /topics/{queue_name}/sendbulk

Sends bulk messages to the specific queue

Example Payload:

```json
{
     "messages": [
         {
            "label": "example",
            "correlationId": "test",
            "contentType": "application/json",
            "data": {
                "key": "value1"
            },
            "userProperties": {
                "name": "test message1"
            }
         },
         {
            "label": "example",
            "correlationId": "test",
            "contentType": "application/json",
            "data": {
                "key": "value2"
            },
            "userProperties": {
                "name": "test message2"
            }
         }
     ]
}
```

### [PUT] /topics/{queue_name}/sendbulktemplate

Sends a templated bulk messages to the specific queue, it can set a wait time between batches  
You can define in how many batches you want to send the amount of message and a wait time between each batches.  

**Attention**: if the volume of messages is big, the messages will be split in batches automatically, this also happens if the batch is too small for the maximum allowed size of a payload

Example Payload:

```json
{
    "totalMessages": 50, // Number of total messages to send, if not defined it will be set to 1
    "batchOf": 5, // send the total message in how many batches, if not defined it will be set to 1
    "waitBetweenBatchesInMilli": 500, // wait between any batches, if not defined it will be 0
    "template": { // Message template
        "label": "example",
        "correlationId": "test",
        "contentType": "application/json",
        "data": {
            "key": "value2"
        },
        "userProperties": {
            "name": "test message2"
        }
    }
}
```

### [GET] /queues/{queue_name}/deadletters

Gets the dead letters from a queue

**Query Attributes**  
*qty*, *integer*: amount of messages to collect, defaults to all with a maximum of 100 messages  
*peek*, *bool*: sets the collection mode to peek, messages will remain in the subscription, defaults to false

### [GET] /queues/{queue_name}/messages

Gets the dead letters from a queue

**Query Attributes**  
*qty*, *integer*: amount of messages to collect, defaults to all with a maximum of 100 messages  
*peek*, *bool*: sets the collection mode to peek, messages will remain in the subscription, defaults to false

## Topics

### List Topics

This will list all topics in a namespace

```bash
servicebus.exe topic list
```

### Create Topic

This will create a topic in a namespace

```bash
servicebus.exe topic create --name="topic_name"
```

### Delete Topic

This will delete a topic and all subscriptions in a namespace

```bash
servicebus.exe topic delete --name="topic_name"
```

### List Subscription for a Topic

```bash
servicebus.exe topic list-subscriptions --topic="example.topic"
```

### Create Topic Subscription

This will create a subscription to a specific topic in a namespace

```bash
servicebus.exe topic create-subscription --name="topic_name" --subscription="name_of_subscription"
```

**Possible flags:**

```--forward-to``` this will create a message forwarding rule in the subscription, the format is ```topic|queue```:```[target_name]```

*Examples*:

```bash
servicebus.exe topic create-subscription --name="new.topic" --subscription="fwd-example" --forward-to="topic:example.topic"
```

in this case it will forward all messages arriving to the *topic* **new.topic** to the *topic* **example.topic**

```--forward-deadletter-to``` this will create a dead letter forwarding rule in the subscription, the format is ```topic|queue```:```[target_name]```

*Examples*:

```bash
servicebus.exe topic create-subscription --name="new.topic" --subscription="fwd-example" --forward-deadletter-to="topic:example.topic"
```

in this case it will forward all dead letters in the *topic* **new.topic** to the *topic* **example.topic**

```--with-rule``` this will create a sql *filter/action* rule in the subscription, the format is *rule_name*:*sql_filter_expression*:*sql_action_expression*

*Examples*:

```bash
servicebus.exe topic create-subscription --name="new.topic" --subscription="rule-example" --with-rule="example_rule:1=1"
```

in this example it will create a sql filter **1=1** rule named *example_rule*

```bash
servicebus.exe topic create-subscription --name="new.topic" --subscription="rule-example" --with-rule="example_rule:1=1:SET sys.label='example.com'"
```

in this example it will create a sql filter **1=1** and a action **SET sys.label='example.com'** named *example_rule*

### Delete Topic Subscription

This will delete a topic subscription for a topic in a namespace

```bash
servicebus.exe topic delete-subscription --name="topic.name" --subscription="subscription.name"
```

### Subscribe to a Topic Subscription

```bash
servicebus.exe topic subscribe --topic="topic.name" --wiretap --peek
```

**Possible flags:**

```--topic``` Name of the topic you want to subscribe, it can be repeated to get multiple subscribers
```--subscription``` Name of the subscriptio you want to subscribe, if you use the **--wiretap** this flag will not be taken into account
```--wiretap``` this will create a **wiretap** subscription in that topic as a catch all
```--peek``` this will not delete the messages from the subscription

*Examples*:

Single Subscriber creating a wiretap

```bash
servicebus.exe topic subscribe --topic="example.topic" --wiretap
```

Multiple Subscriber creating a wiretap

```bash
servicebus.exe topic subscribe --topic="example.topic1" --topic="example.topic2" --wiretap
```

### Send a Message to a Topic

```bash
servicebus.exe topic send --topic="topic.name"
```

**Possible flags:**

```--topic``` Name of the topic where to send the message

```--file``` File path with the MessageRequest entity to send, use this instead on inline --body flag

```--body``` Message body in json (please escape the json correctly as this is validated)

```--label``` Message Label

```--property``` Add a User property to the message, this flag can be repeated to add more than one property. **format:** the format will be **[key]:[value]**

*Examples*:

```bash
servicebus.exe topic send --topic="example.topic" --body='{\"example\":\"document\"}' --label="ExampleLabel"
```

## Queues

### List Queues

This will list all topics in a namespace

```bash
servicebus.exe topic list
```

### Create Queue

This will create a Queue in a Namespace

```bash
servicebus.exe queue create --name="queue.name"
```

**Possible flags:**

```--forward-to``` this will create a message forwarding rule in the queue, the format is ```topic|queue```:```[target_name]```

*Examples*:

```bash
servicebus.exe queue create --name="new.queue" --forward-to="topic:example.topic"
```

in this case it will forward all messages arriving to the *queue* **new.queue** to the *topic* **example.topic**

```--forward-deadletter-to``` this will create a dead letter forwarding rule in the queue, the format is ```topic|queue```:```[target_name]```

*Examples*:

```bash
servicebus.exe topic create-subscription --name="new.queue" --forward-deadletter-to="topic:example.topic"
```

in this case it will forward all dead letters in the *queue* **new.queue** to the *topic* **example.topic**

### Delete Queue

This will delete a topic and all subscriptions in a namespace

```bash
servicebus.exe queue delete --name="queue.name"
```

### Subscribe to a Queue

```bash
servicebus.exe queue subscribe --queue="queue.name" --wiretap --peek
```

**Possible flags:**

```--topic``` Name of the topic you want to subscribe, it can be repeated to get multiple subscribers
```--subscription``` Name of the subscriptio you want to subscribe, if you use the **--wiretap** this flag will not be taken into account
```--wiretap``` this will create a **wiretap** subscription in that topic as a catch all
```--peek``` this will not delete the messages from the subscription

*Examples*:

Single Subscriber creating a wiretap

```bash
servicebus.exe queue subscribe --queue="example.queue" --wiretap
```

Multiple Subscriber creating a wiretap

```bash
servicebus.exe queue subscribe --queue="example.queue" --queue="example.queue" --wiretap
```

### Send a Message to a Queue

```bash
servicebus.exe queue send --queue="queue.name"
```

**Possible flags:**

```--queue``` Name of the queue where to send the message

```--file``` File path with the MessageRequest entity to send, use this instead on inline --body flag

```--body``` Message body in json (please escape the json correctly as this is validated)

```--label``` Message Label

```--property``` Add a User property to the message, this flag can be repeated to add more than one property. **format:** the format will be **[key]:[value]**

*Examples*:

```bash
servicebus.exe queue send --queue="example.queue" --body='{\"example\":\"document\"}' --label=ExampleLabel
```
