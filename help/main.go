package help

import (
	"runtime"
	"strings"

	"github.com/cjlapao/common-go/log"
	"github.com/fatih/color"
)

var logger = log.Get()

// PrintMainCommandHelper Prints specific Help
func PrintMainCommandHelper() {
	logger.Info("")
	logger.Info("Please choose a command:")
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus [command]")
	logger.Info("")
	logger.Info("Available Commands:")
	logger.Info("  topic         Service bus topic command")
	logger.Info("  queue         Service bus topic command")
}

func PrintMissingServiceBusConnectionHelper() {
	logger.Error("Service bus connection string was not found")
	logger.Info("")
	logger.Info("Please add the SERVICEBUS_CONNECTION_STRING to your environment and try again")
	logger.Info("")
	logger.Info("Example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		logger.Info("  export SERVICEBUS_CONNECTION_STRING=\"{your connection string}\"")
	case "windows":
		logger.Info("  $env:SERVICEBUS_CONNECTION_STRING=\"{your connection string}\"")
	}
}

func PrintTopicMainCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus topic [subcommand]")
	logger.Info("")
	logger.Info("Available Sub-Commands:")
	logger.Info("  list                 Lists all Topics in a Namespace")
	logger.Info("  create               Creates a Topic in a Namespace")
	logger.Info("  delete               Deletes a Topic in a Namespace")
	logger.Info("  send                 Sends a Json Message to a specific Topic in a Namespace")
	logger.Info("  list-subscriptions   List all Subscriptions on a Topic in a Namespace")
	logger.Info("  create-subscription  Creates a Subscription on a specific Topic in a Namespace")
	logger.Info("  delete-subscription  Deletes a Subscription from a specific Topic in a Namespace")
	logger.Info("  subscribe            Subscribe to a Subscription and prints the message")
}

func PrintTopicListSubscriptionsCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus topic list-subscriptions [Options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --name               Topic name to list subscriptions")
	logger.Info("")
	logger.Info("Example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v topic list-subscriptions %v", color.HiYellowString("servicebus"), color.HiBlackString("--name=example.topic"))
	case "windows":
		color.White("%v topic list-subscriptions %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--name=example.topic"))
	}
}

func PrintTopicDeleteTopicCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus topic delete [Options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --name               Topic name to delete")
	logger.Info("")
	logger.Info("Example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v topic delete %v", color.HiYellowString("servicebus"), color.HiBlackString("--name=example.topic"))
	case "windows":
		color.White("%v topic delete %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--name=example.topic"))
	}
}

func PrintTopicCreateSubscriptionCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus topic create-subscription [Options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --name                   Topic name to create subscription on")
	logger.Info("  --subscription           Subscription name to create")
	logger.Info("  --forward-to             Creates a forward to rule in the subscription")
	logger.Info("                           the format will be topic|queue:[name_of_the_target]")
	logger.Info("                           example --forward-to=topic:example.topic")
	logger.Info("  --forward-deadletter-to  Creates a forward to rule for dead letters in the subscription")
	logger.Info("                           the format will be topic|queue:[name_of_the_target]")
	logger.Info("                           example: --forward-deadletter-to=topic:example.topic")
	logger.Info("  --with-rule              Creates a sql filter/action rule for the subscription")
	logger.Info("                           the format will be [rule_name]:[sql_filter_expression]:[sql_action_expression]")
	logger.Info("                           example with only filter: --with-rule=example:2=2")
	logger.Info("                           example with filter and action: --with-rule=example:2=2:SET sys.label='example'")
	logger.Info("")
	logger.Info("Example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v topic create-subscription %v", color.HiYellowString("servicebus"), color.HiBlackString("--name=example.topic --subscription=example.subscription --forward-to=queue:example.queue --with-rule=example:1=1"))
	case "windows":
		color.White("%v topic create-subscription %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--name=example.topic --subscription=example.subscription --forward-to=queue:example.queue --with-rule=example:1=1"))
	}
}

func PrintTopicCreateTopicCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus topic create [Options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --name               Topic name to create")
	logger.Info("")
	logger.Info("Example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v topic create %v", color.HiYellowString("servicebus"), color.HiBlackString("--name=example.topic"))
	case "windows":
		color.White("%v topic create %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--name=example.topic"))
	}
}

func PrintTopicDeleteSubscriptionCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus topic create-subscription [Options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --name                   Topic name to delete subscription on")
	logger.Info("  --subscription           Subscription name to delete")
	logger.Info("")
	logger.Info("Example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v topic delete-subscription %v", color.HiYellowString("servicebus"), color.HiBlackString("--name=example.topic --subscription=example.subscription"))
	case "windows":
		color.White("%v topic delete-subscription %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--name=example.topic  --subscription=example.subscription"))
	}
}

func PrintTopicSubscribeCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus topic subscribe [options]")
	logger.Info("")
	logger.Info("Please choose a option:")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  %v=string         Name of the topic to listen to (mandatory)", "--topic")
	logger.Info("                         this flag can be repeated to listen to several topics")
	logger.Info("  %v=string  Name of the subscription to listen to (mandatory)", "--subscription")
	logger.Info("  %v              connects to a wiretap in the topic, if this subscription", "--wiretap")
	logger.Info("                         does not exist it will be created and deleted on exit")
	logger.Info("                         this will also override the %v flag", "--subscription")
	logger.Info("  %v                 peeks into the subscription leaving the messages there", "--peek")
	logger.Info("")
	logger.Info("example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		logger.Info("Single topic subscriber:")
		color.White("%v topic subscribe %v", color.HiYellowString("servicebus"), color.HiBlackString("--topic=example.topic --wiretap"))
		logger.Info("")
		logger.Info("Multiple topics subscriber")
		color.White("%v topic subscribe %v", color.HiYellowString("servicebus"), color.HiBlackString("--topic=example.topic --topic=example.topic2 --wiretap"))
	case "windows":
		logger.Info("Single topic subscriber:")
		color.White("%v topic subscribe %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--topic=example.topic --wiretap"))
		logger.Info("")
		logger.Info("Multiple topics subscriber")
		color.White("%v topic subscribe %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--topic=example.topic --topic=example.topic2 --wiretap"))
	}
}

func PrintTopicSendCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus topic send [options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --topic    string     Name of the topic where to send the message")
	logger.Info("  --tenant   guid       Id of the tenant")
	logger.Info("  --body     json       Message body in json (please escape the json correctly as this is validated)")
	logger.Info("  --domain   string     Forwarding topology Message Domain")
	logger.Info("  --name     string     Forwarding topology Message Name")
	logger.Info("  --version  string     Forwarding topology Version")
	logger.Info("  --sender   string     Forwarding topology Sender")
	logger.Info("  --label    string     Message Label")
	logger.Info("  --property key:value  Add a User property to the message")
	logger.Info("                        This option can be repeated to add more than one property")
	logger.Info("                        the format will be [key]:[value]")
	logger.Info("                        example: X-Sender:example")
	logger.Info("  --default             With this flag the tool will generate a default TimeService sample")
	logger.Info("                        using the forwarding topology format")
	logger.Info("  --uno                 With this flag the tool will convert the default TimeService sample")
	logger.Info("                        message to Uno format")
	logger.Info("")
	logger.Info("example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v topic send %v", color.HiYellowString("servicebus"), color.HiBlackString("--topic=example.topic --body='{\\\"example\\\":\\\"document\\\"}' --domain=ExampleService --name=Example --version=\"2.1\" --sender=ExampleSender --label=ExampleLabel"))
	case "windows":
		color.White("%v topic send %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--topic=example.topic --body='{\\\"example\\\":\\\"document\\\"}' --domain=ExampleService --name=Example --version=\"2.1\" --sender=ExampleSender --label=ExampleLabel"))
	}
}

func PrintQueueMainCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus queue [subcommand]")
	logger.Info("")
	logger.Info("Available Sub-Commands:")
	logger.Info("  list                 Lists all Queues in a Namespace")
	logger.Info("  create               Creates a Queues in a Namespace")
	logger.Info("  delete               Deletes a Queues in a Namespace")
	logger.Info("  send                 Sends a Json Message to a specific Queue in a Namespace")
	logger.Info("  subscribe            Subscribe to a Queue and prints the messages")
}

func PrintQueueDeleteCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus queue delete [Options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --name               Queue name to delete")
	logger.Info("")
	logger.Info("Example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v queue delete %v", color.HiYellowString("servicebus"), color.HiBlackString("--name=example.queue"))
	case "windows":
		color.White("%v queue delete %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--name=example.queue"))
	}
}

func PrintQueueCreateCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus queue create-subscription [Options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --name                   Queue name to create subscription on")
	logger.Info("  --forward-to             Creates a forward to rule in the subscription")
	logger.Info("                           the format will be topic|queue:[name_of_the_target]")
	logger.Info("                           example --forward-to=topic:example.topic")
	logger.Info("  --forward-deadletter-to  Creates a forward to rule for dead letters in the subscription")
	logger.Info("                           the format will be topic|queue:[name_of_the_target]")
	logger.Info("                           example: --forward-deadletter-to=topic:example.topic")
	logger.Info("")
	logger.Info("Example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v queue create-subscription %v", color.HiYellowString("servicebus"), color.HiBlackString("--name=example.queue --forward-to=topic:example.topic --with-rule=example:1=1"))
	case "windows":
		color.White("%v queue create-subscription %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--name=example.queue --forward-to=topic:example.topic --with-rule=example:1=1"))
	}
}

func PrintQueueSendCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus queue send [options]")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  --queue    string     Name of the queue where to send the message")
	logger.Info("  --tenant   guid       Id of the tenant")
	logger.Info("  --body     json       Message body in json (please escape the json correctly as this is validated)")
	logger.Info("  --domain   string     Forwarding topology Message Domain")
	logger.Info("  --name     string     Forwarding topology Message Name")
	logger.Info("  --version  string     Forwarding topology Version")
	logger.Info("  --sender   string     Forwarding topology Sender")
	logger.Info("  --label    string     Message Label")
	logger.Info("  --property key:value  Add a User property to the message")
	logger.Info("                        This option can be repeated to add more than one property")
	logger.Info("                        the format will be [key]:[value]")
	logger.Info("                        example: X-Sender:example")
	logger.Info("  --default             With this flag the tool will generate a default TimeService sample")
	logger.Info("                        using the forwarding topology format")
	logger.Info("  --uno                 With this flag the tool will convert the default TimeService sample")
	logger.Info("                        message to Uno format")
	logger.Info("")
	logger.Info("example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		color.White("%v queue send %v", color.HiYellowString("servicebus"), color.HiBlackString("--queue=example.queue --body='{\\\"example\\\":\\\"document\\\"}' --domain=ExampleService --name=Example --version=\"2.1\" --sender=ExampleSender --label=ExampleLabel"))
	case "windows":
		color.White("%v queue send %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--queue=example.queue --body='{\\\"example\\\":\\\"document\\\"}' --domain=ExampleService --name=Example --version=\"2.1\" --sender=ExampleSender --label=ExampleLabel"))
	}
}

func PrintQueueSubscribeCommandHelper() {
	logger.Info("")
	logger.Info("Usage:")
	logger.Info("  servicebus queue subscribe [options]")
	logger.Info("")
	logger.Info("Please choose a option:")
	logger.Info("")
	logger.Info("Available Options:")
	logger.Info("  %v=string         Name of the queue to listen to (mandatory)", "--queue")
	logger.Info("                         this flag can be repeated to listen to several queues")
	logger.Info("  %v                 peeks into the subscription leaving the messages there", "--peek")
	logger.Info("")
	logger.Info("example:")
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		logger.Info("Single topic subscriber:")
		color.White("%v queue subscribe %v", color.HiYellowString("servicebus"), color.HiBlackString("--queue=example.queue"))
		logger.Info("")
		logger.Info("Multiple topics subscriber")
		color.White("%v queue subscribe %v", color.HiYellowString("servicebus"), color.HiBlackString("--queue=example.queue --queue=example.queue2"))
	case "windows":
		logger.Info("Single topic subscriber:")
		color.White("%v queue subscribe %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--queue=example.queue"))
		logger.Info("")
		logger.Info("Multiple topics subscriber")
		color.White("%v queue subscribe %v", color.HiYellowString("servicebus.exe"), color.HiBlackString("--queue=example.queue --topic=example.queue2"))
	}
}
