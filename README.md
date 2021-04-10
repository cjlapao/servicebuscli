# Azure Service Bus Command Line Tool

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Build](https://github.com/cjlapao/servicebuscli-go/actions/workflows/build.yml/badge.svg)](https://github.com/cjlapao/servicebuscli-go/actions/workflows/build.yml) [![Release](https://github.com/cjlapao/servicebuscli-go/actions/workflows/release.yml/badge.svg)](https://github.com/cjlapao/servicebuscli-go/actions/workflows/release.yml) [![CodeQL](https://github.com/cjlapao/servicebuscli-go/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/cjlapao/servicebuscli-go/actions/workflows/codeql-analysis.yml)

## Index

- [Azure Service Bus Command Line Tool](#azure-service-bus-command-line-tool)
  - [Index](#index)
  - [**How to Use it**](#how-to-use-it)
  - [**Topics**](#topics)
    - [**List Topics**](#list-topics)
    - [**Create Topic**](#create-topic)
    - [**Delete Topic**](#delete-topic)
    - [**List Subscription for a Topic**](#list-subscription-for-a-topic)
    - [**Create Topic Subscription**](#create-topic-subscription)
    - [**Delete Topic Subscription**](#delete-topic-subscription)
    - [Subscribe to a Topic Subscription](#subscribe-to-a-topic-subscription)
    - [**Send a Message to a Topic**](#send-a-message-to-a-topic)
  - [**Queues**](#queues)
    - [**List Queues**](#list-queues)
    - [**Create Queue**](#create-queue)
    - [**Delete Queue**](#delete-queue)
    - [**Subscribe to a Queue**](#subscribe-to-a-queue)
    - [**Send a Message to a Queue**](#send-a-message-to-a-queue)

This is a command line tool to help test service bus messages.

You will be able to do **C** *U* **RD** operations to topics/subscriptions and queues, you will also be able to send messages and subscribe to a specific queue/subscription

## **How to Use it**

Once you have it compiled you can run it with the needed options

you can also run the following command to display the help

```bash
servicebus.exe --help
```

## **Topics**

### **List Topics**

This will list all topics in a namespace

```bash
servicebus.exe topic list
```

### **Create Topic**

This will create a topic in a namespace

```bash
servicebus.exe topic create --name="topic_name"
```

### **Delete Topic**

This will delete a topic and all subscriptions in a namespace

```bash
servicebus.exe topic delete --name="topic_name"
```

### **List Subscription for a Topic**

```bash
servicebus.exe topic list-subscriptions --topic="example.topic"
```

### **Create Topic Subscription**

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

### **Delete Topic Subscription**

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

### **Send a Message to a Topic**

```bash
servicebus.exe topic send --topic="topic.name"
```

**Possible flags:**

```--topic``` Name of the topic where to send the message

```--tenant``` Id of the tenant

```--body``` Message body in json (please escape the json correctly as this is validated)

```--domain``` Forwarding topology Message Domain

```--name``` Forwarding topology Message Name

```--version``` Forwarding topology Version

```--sender``` Forwarding topology Sender

```--label``` Message Label

```--property``` Add a User property to the message, this flag can be repeated to add more than one property. **format:** the format will be **[key]:[value]**

```--default``` With this flag the tool will generate a default TimeService sample using the forwarding topology format

```--uno``` With this flag the tool will convert the default TimeService sample message to Uno format

*Examples*:

```bash
servicebus.exe topic send --topic="example.topic" --body='{\"example\":\"document\"}' --domain="ExampleService" --name="Example" --version=\"2.1\" --sender="ExampleSender" --label="ExampleLabel"
```

## **Queues**

### **List Queues**

This will list all topics in a namespace

```bash
servicebus.exe topic list
```

### **Create Queue**

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

### **Delete Queue**

This will delete a topic and all subscriptions in a namespace

```bash
servicebus.exe queue delete --name="queue.name"
```

### **Subscribe to a Queue**

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

### **Send a Message to a Queue**

```bash
servicebus.exe queue send --queue="queue.name"
```

**Possible flags:**

```--queue``` Name of the queue where to send the message

```--tenant``` Id of the tenant

```--body``` Message body in json (please escape the json correctly as this is validated)

```--domain``` Forwarding topology Message Domain

```--name``` Forwarding topology Message Name

```--version``` Forwarding topology Version

```--sender``` Forwarding topology Sender

```--label``` Message Label

```--property``` Add a User property to the message, this flag can be repeated to add more than one property. **format:** the format will be **[key]:[value]**

```--default``` With this flag the tool will generate a default TimeService sample using the forwarding topology format

```--uno``` With this flag the tool will convert the default TimeService sample message to Uno format

*Examples*:

```bash
servicebus.exe queue send --queue="example.queue" --body='{\"example\":\"document\"}' --domain=ExampleService --name=Example --version=\"2.1\" --sender=ExampleSender --label=ExampleLabel
```
