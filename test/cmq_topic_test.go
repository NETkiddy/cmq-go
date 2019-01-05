package test

import (
	"testing"
	"github.com/NETkiddy/cmq-go"
)

var endpointTopic = "https://cmq-topic-sh.api.qcloud.com"
var endpointTopicInner = "http://cmq-topic-sh.api.tencentyun.com"
//-----------------------------------------------------------------
// 创建主题
func Test_CreateTopic(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	err, _ := account.CreateTopic("topic-test-001", 2048)
	if err != nil {
		t.Errorf("topic-test-001 created failed, %v", err.Error())
		return
	}
	t.Log("topic-test-001 created")
	err, _ = account.CreateTopic("topic-test-002", 4096)
	if err != nil {
		t.Errorf("topic-test-002 created failed, %v", err.Error())
		return
	}
	t.Log("topic-test-002 created")
}

// 删除主题
func Test_ListTopic(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	totalCount, topicList, err, _ := account.ListTopic("", -1, -1)
	if err != nil {
		t.Errorf("ListTopic failed, %v", err.Error())
		return
	}
	t.Logf("totalCount, %v", totalCount)
	t.Logf("topicList, %v", topicList)
}

// 获取主题
func Test_GetTopic(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	topic := account.GetTopic("topic-test-001")
	t.Logf("GetTopic, %v", *topic)
}

// 删除主题
func Test_DeleteTopic(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	err, _ := account.DeleteTopic("topic-test-001")
	if err != nil {
		t.Errorf("DeleteTopic failed, %v", err.Error())
		return
	}
	t.Logf("DeleteTopic, %v", "topic-test-001")
}

//-----------------------------------------------------------------
// 设置，获取主题属性
func Test_GetSetTopicAttributes(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	topic := account.GetTopic("topic-test-001")
	t.Logf("GetTopic, %v", *topic)
	// get meta
	meta, err, _ := topic.GetTopicAttributes()
	if err != nil {
		t.Errorf("GetTopicAttributes failed, %v", err.Error())
		return
	}
	t.Logf("GetTopicAttributes, before set, %v", meta)
	// set meta
	meta.MaxMsgSize = 32768
	topic.SetTopicAttributes(meta.MaxMsgSize)
	// get meta
	newMeta, err, _ := topic.GetTopicAttributes()
	if err != nil {
		t.Errorf("GetTopicAttributes failed, %v", err.Error())
		return
	}
	t.Logf("GetTopicAttributes, after set, %v", newMeta)
}

// 创建订阅者
func Test_CreateSubscribe(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	err, _ := account.CreateSubscribe("topic-test-001", "sub-test", "queue-test-001", "queue", "SIMPLIFIED")
	if err != nil {
		t.Errorf("CreateSubscribe failed, %v", err.Error())
		return
	}
	t.Logf("CreateSubscribe succeed")
}

// 获取订阅者
func Test_GetSubscribe(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	// get
	sub := account.GetSubscription("topic-test-001", "sub-test")
	t.Logf("GetSubscription succeed: %v", sub)

	// set meta
	meta, err, _ := sub.GetSubscriptionAttributes()
	if err != nil {
		t.Errorf("CreateSubscribe failed, %v", err.Error())
		return
	}
	t.Logf("GetSubscriptionAttributes succeed: %v", meta)
}

// 获取所有主题订阅者
func Test_ListSubscription(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	topic := account.GetTopic("topic-test-001")
	t.Logf("GetTopic, %v", topic)

	totalCount, subscriptionList, err, _ := topic.ListSubscription(-1, -1, "")
	if err != nil {
		t.Errorf("ListSubscription failed, %v", err.Error())
		return
	}
	t.Logf("ListSubscription totalCount, %v", totalCount)
	t.Logf("ListSubscription subscriptionList, %v", subscriptionList)
}

// 发布消息
func Test_PublishMessage(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	topic := account.GetTopic("topic-test-001")
	t.Logf("GetTopic, %v", topic)

	msgId, err, _ := topic.PublishMessage("hello world")
	if err != nil {
		t.Errorf("PublishMessage failed, %v", err.Error())
		return
	}
	t.Logf("PublishMessage msgId, %v", msgId)
}

// 批量发布消息
func Test_BatchPublishMessage(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	topic := account.GetTopic("topic-test-001")
	t.Logf("GetTopic, %v", topic)

	msgs := []string{"hello world", "foo", "bar"}
	msgIds, err, _ := topic.BatchPublishMessage(msgs)
	if err != nil {
		t.Errorf("BatchPublishMessage failed, %v", err.Error())
		return
	}
	t.Logf("BatchPublishMessage msgIds, %v", msgIds)
}

// 删除主题订阅者
func Test_DeleteSubscription(t *testing.T) {
	account := cmq_go.NewAccount(endpointTopic, secretId, secretKey)
	err, _ := account.DeleteSubscribe("topic-test-001", "sub-test")
	if err != nil {
		t.Errorf("DeleteSubscribe failed, %v", err.Error())
		return
	}
	t.Logf("DeleteSubscribe succeed")
}
