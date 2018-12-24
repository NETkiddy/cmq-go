package test

import (
	"testing"
	"github.com/NETkiddy/cmq-go"
)

var secretId = "YourTencentSecretId"
var secretKey = "YourTencentSecretKey"
var endpointQueue = "https://cmq-queue-sh.api.qcloud.com"
var endpointQueueInner = "http://cmq-queue-sh.api.tencentyun.com"

//-----------------------------------------------------------------
// 创建队列
func Test_CreateQueue(t *testing.T) {
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	meta := cmq_go.QueueMeta{}
	meta.PollingWaitSeconds = 10
	meta.VisibilityTimeout = 10
	meta.MaxMsgSize = 1048576
	meta.MsgRetentionSeconds = 345600

	err := account.CreateQueue("queue-test-001", meta)
	if err != nil {
		t.Errorf("queue-test-001 created failed, %v", err.Error())
		return
	}
	t.Log("queue-test-001 created")

	err = account.CreateQueue("queue-test-002", meta)
	if err != nil {
		t.Errorf("queue-test-002 created failed, %v", err.Error())
		return
	}
	t.Log("queue-test-002 created")
}

// 列出当前帐号下所有队列名字
func Test_ListQueue(t *testing.T) {
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	totalCount, queueList, err := account.ListQueue("", -1, -1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("totalCount: %v", totalCount)
	t.Logf("queueList: %v", queueList)
}

// 删除队列
func Test_DeleteQueue(t *testing.T) {
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	err := account.DeleteQueue("queue-test-002")
	if err != nil {
		t.Error(err)
		return
	}
}

// 获得队列实例
func Test_GetQueue(t *testing.T) {
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	queue := account.GetQueue("queue-test-001")
	t.Logf("GetQueue: %v", queue)
}

//-----------------------------------------------------------------
// 设置队列属性
func Test_SetQueueAttributes(t *testing.T) {
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	meta := cmq_go.QueueMeta{}
	meta.PollingWaitSeconds = 10
	meta.VisibilityTimeout = 10
	meta.MaxMsgSize = 1048576
	meta.MsgRetentionSeconds = 345600

	queue := account.GetQueue("queue-test-001")

	err := queue.SetQueueAttributes(meta)
	if err != nil {
		t.Error(err)
		return
	}
}

//-----------------------------------------------------------------
// 获取队列属性
func Test_GetQueueAttributes(t *testing.T) {
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	queue := account.GetQueue("queue-test-001")
	newMeta, err := queue.GetQueueAttributes()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("GetQueueAttributes: %v", newMeta)
}

// 发送，接受，删除消息
func Test_SendReceiveDeleteMessage(t *testing.T) {
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	queue := account.GetQueue("queue-test-001")
	// send
	msgId, err := queue.SendMessage("hello world")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("SendMessage msgId: %v", msgId)
	// receive
	msg, err := queue.ReceiveMessage(10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ReceiveMessage msgId: %v", msg.MsgId)
	// delete
	err = queue.DeleteMessage(msg.ReceiptHandle)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("DeleteMessage msgId: %v", msg.MsgId)
}

// 批量发送，接收，删除消息
func Test_BatchSendReceiveDeleteMessage(t *testing.T) {
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	queue := account.GetQueue("queue-test-001")
	// batch send
	msgBodys := []string{"hello world", "foo", "bar"}
	msgIds, err := queue.BatchSendMessage(msgBodys)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("BatchSendMessage msgId: %v", msgIds)
	// batch receive
	msgs, err := queue.BatchReceiveMessage(10, 10)
	if err != nil {
		t.Error(err)
		return
	}
	handlers := make([]string, 0)
	msgIds = msgIds[0:0]
	for _, msg := range msgs {
		handlers = append(handlers, msg.ReceiptHandle)
		msgIds = append(msgIds, msg.MsgId)
	}
	t.Logf("BatchReceiveMessage msgIds: %v", msgIds)
	// batch delete
	err = queue.BatchDeleteMessage(handlers)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("BatchDeleteMessage msgId: %v", msgIds)
}
