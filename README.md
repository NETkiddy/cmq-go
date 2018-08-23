# cmq-go
Tencent CMQ Golang SDK

与腾讯云官方SDK接口保持一致

## Getting Started

所有的API中都有[test case](https://github.com/NETkiddy/cmq-go/tree/master/test)，这里使用[CreateQueue](https://github.com/NETkiddy/cmq-go/blob/master/test/cmq_queue_test.go#L15)举例。
```
var secretId = "YourTencentSecretId"
var secretKey = "YourTencentSecretKey"
var endpointQueue = "https://cmq-queue-sh.api.qcloud.com"
var endpointQueueInner = "http://cmq-queue-sh.api.tencentyun.com"

// 创建队列
func Test_CreateQueue(t *testing.T) {
    //创建账户
	account := cmq_go.NewAccount(endpointQueue, secretId, secretKey)
	
	//设置队列metadata
	meta := cmq_go.QueueMeta{}
	meta.PollingWaitSeconds = 10
	meta.VisibilityTimeout = 10
	meta.MaxMsgSize = 1048576
	meta.MsgRetentionSeconds = 345600

    //创建队列queue-test-001
	err := account.CreateQueue("queue-test-001", meta)
	if err != nil {
		t.Errorf("queue-test-001 created failed, %v", err.Error())
		return
	}
	t.Log("queue-test-001 created")

    //创建队列queue-test-002
	err = account.CreateQueue("queue-test-002", meta)
	if err != nil {
		t.Errorf("queue-test-002 created failed, %v", err.Error())
		return
	}
	t.Log("queue-test-002 created")
}
```

## Test Case

```
测试单个方法
go test -v -test.run Test_CreateQueue
```


## API Status
### 队列模型
#### 队列相关接口
- [x] CreateQueue		
- [x] ListQueue 	
- [x] GetQueueAttributes 		
- [x] SetQueueAttributes	
- [x] DeleteQueue

#### 消息相关接口
- [x] SendMessage	
- [x] BatchSendMessage	
- [x] ReceiveMessage	
- [x] BatchReceiveMessage	
- [x] DeleteMessage	
- [x] BatchDeleteMessage

### 主题模型
#### 主题相关接口
- [x] 	CreateTopic
- [x] 	SetTopicAttributes	
- [x] 	ListTopic
- [x] 	GetTopicAttributes
- [x] 	DeleteTopic

#### 消息相关接口
- [x] PublishMessage	
- [x] BatchPublishMessage

#### 订阅相关接口
- [x] Subscribe	
- [x] ListSubscriptionByTopic	
- [x] SetSubscriptionAttributes	
- [x] GetSubscriptionAttributes	
- [x] Unsubscribe

