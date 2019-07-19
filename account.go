package cmq_go

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
)

const (
	DEFAULT_ERROR_CODE = -1
)

type Account struct {
	client *CMQClient
}

func NewAccount(endpoint, secretId, secretKey string) *Account {
	return &Account{
		client: NewCMQClient(endpoint, "/v2/index.php", secretId, secretKey, "POST"),
	}
}

func (this *Account) SetProxy(proxyUrl string) {
	this.client.setProxy(proxyUrl)
}

func (this *Account) UnsetProxy() {
	this.client.unsetProxy()
}

func (this *Account) setSignMethod(method string) (err error) {
	return this.client.setSignMethod(method)
}

func (this *Account) CreateQueue(queueName string, queueMeta QueueMeta) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	if queueName == "" {
		err = fmt.Errorf("createQueue failed: queueName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["queueName"] = queueName
	if queueMeta.MaxMsgHeapNum > 0 {
		param["maxMsgHeapNum"] = strconv.Itoa(queueMeta.MaxMsgHeapNum)
	}
	if queueMeta.PollingWaitSeconds > 0 {
		param["pollingWaitSeconds"] = strconv.Itoa(queueMeta.PollingWaitSeconds)
	}
	if queueMeta.VisibilityTimeout > 0 {
		param["visibilityTimeout"] = strconv.Itoa(queueMeta.VisibilityTimeout)
	}
	if queueMeta.MaxMsgSize > 0 {
		param["maxMsgSize"] = strconv.Itoa(queueMeta.MaxMsgSize)
	}
	if queueMeta.MsgRetentionSeconds > 0 {
		param["msgRetentionSeconds"] = strconv.Itoa(queueMeta.MsgRetentionSeconds)
	}
	if queueMeta.RewindSeconds > 0 {
		param["rewindSeconds"] = strconv.Itoa(queueMeta.RewindSeconds)
	}

	_, err, code, moduleErrCode = doCall(this.client, param, "CreateQueue")
	if err != nil {
		//log.Printf("client.call CreateQueue failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Account) DeleteQueue(queueName string) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	if queueName == "" {
		err = fmt.Errorf("deleteQueue failed: queueName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["queueName"] = queueName

	_, err, code, moduleErrCode = doCall(this.client, param, "DeleteQueue")
	if err != nil {
		//log.Printf("client.call DeleteQueue failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Account) ListQueue(searchWord string, offset, limit int) (
	totalCount int, queueList []string, err error, code int) {
	code = DEFAULT_ERROR_CODE
	queueList = make([]string, 0)
	param := make(map[string]string)
	if searchWord != "" {
		param["searchWord"] = searchWord
	}
	if offset >= 0 {
		param["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		param["limit"] = strconv.Itoa(limit)
	}

	resMap, err, code, _ := doCall(this.client, param, "ListQueue")
	if err != nil {
		//log.Printf("client.call ListQueue failed: %v\n", err.Error())
		return
	}
	totalCount = int(resMap["totalCount"].(float64))
	resQueueList := resMap["queueList"].([]interface{})
	for _, v := range resQueueList {
		queue := v.(map[string]interface{})
		queueList = append(queueList, queue["queueName"].(string))
	}
	return
}

func (this *Account) GetQueue(queueName string) (queue *Queue) {
	return NewQueue(queueName, this.client)
}

func (this *Account) GetTopic(topicName string) (topic *Topic) {
	return NewTopic(topicName, this.client)
}

func (this *Account) CreateTopic(topicName string, maxMsgSize int) (err error, code, moduleErrCode int) {
	err, code, moduleErrCode = _createTopic(this.client, topicName, maxMsgSize, 1)
	return
}

func _createTopic(client *CMQClient, topicName string, maxMsgSize, filterType int) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	if topicName == "" {
		err = fmt.Errorf("createTopic failed: topicName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["topicName"] = topicName
	param["filterType"] = strconv.Itoa(filterType)
	if maxMsgSize < 1024 || maxMsgSize > 1048576 {
		err = fmt.Errorf("createTopic failed: Invalid parameter: maxMsgSize > 1024KB or maxMsgSize < 1KB")
		//log.Printf("%v", err.Error())
		return
	}
	param["maxMsgSize"] = strconv.Itoa(maxMsgSize)

	_, err, code, moduleErrCode = doCall(client, param, "CreateTopic")
	if err != nil {
		//log.Printf("client.call CreateTopic failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Account) DeleteTopic(topicName string) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	if topicName == "" {
		err = fmt.Errorf("deleteTopic failed: topicName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["topicName"] = topicName

	_, err, code, moduleErrCode = doCall(this.client, param, "DeleteTopic")
	if err != nil {
		//log.Printf("client.call DeleteTopic failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Account) ListTopic(searchWord string, offset, limit int) (
	totalCount int, topicList []string, err error, code int) {
	code = DEFAULT_ERROR_CODE
	topicList = make([]string, 0)
	param := make(map[string]string)
	if searchWord != "" {
		param["searchWord"] = searchWord
	}
	if offset > 0 {
		param["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		param["limit"] = strconv.Itoa(limit)
	}

	resMap, err, code, _ := doCall(this.client, param, "ListTopic")
	if err != nil {
		//log.Printf("client.call ListTopic failed: %v\n", err.Error())
		return
	}
	totalCount = int(resMap["totalCount"].(float64))
	resTopicList := resMap["topicList"].([]interface{})
	for _, v := range resTopicList {
		topic := v.(map[string]interface{})
		topicList = append(topicList, topic["topicName"].(string))
	}
	return
}

func (this *Account) CreateSubscribe(topicName, subscriptionName, endpoint, protocol, notifyContentFormat string) (
	err error, code, moduleErrCode int) {
	err, code, moduleErrCode = _createSubscribe(
		this.client, topicName, subscriptionName, endpoint, protocol, nil, nil,
		NotifyStrategyDefault, notifyContentFormat)
	return
}

func _createSubscribe(client *CMQClient, topicName, subscriptionName, endpoint, protocol string, filterTag []string,
	bindingKey []string, notifyStrategy, notifyContentFormat string) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	if topicName == "" {
		err = fmt.Errorf("createSubscribe failed: topicName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["topicName"] = topicName

	if subscriptionName == "" {
		err = fmt.Errorf("createSubscribe failed: subscriptionName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["subscriptionName"] = subscriptionName

	if endpoint == "" {
		err = fmt.Errorf("createSubscribe failed: endpoint is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["endpoint"] = endpoint

	if protocol == "" {
		err = fmt.Errorf("createSubscribe failed: protocal is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["protocol"] = protocol

	if notifyStrategy == "" {
		err = fmt.Errorf("createSubscribe failed: notifyStrategy is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["notifyStrategy"] = notifyStrategy

	if notifyContentFormat == "" {
		err = fmt.Errorf("createSubscribe failed: notifyContentFormat is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["notifyContentFormat"] = notifyContentFormat

	if filterTag != nil {
		for i, tag := range filterTag {
			param["filterTag."+strconv.Itoa(i+1)] = tag
		}
	}

	if bindingKey != nil {
		for i, key := range bindingKey {
			param["bindingKey."+strconv.Itoa(i+1)] = key
		}
	}

	_, err, code, moduleErrCode = doCall(client, param, "Subscribe")
	if err != nil {
		//log.Printf("client.call Subscribe failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Account) DeleteSubscribe(topicName, subscriptionName string) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	if topicName == "" {
		err = fmt.Errorf("createSubscribe failed: topicName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["topicName"] = topicName

	if subscriptionName == "" {
		err = fmt.Errorf("createSubscribe failed: subscriptionName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["subscriptionName"] = subscriptionName

	_, err, code, moduleErrCode = doCall(this.client, param, "Unsubscribe")
	if err != nil {
		//log.Printf("client.call Unsubscribe failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Account) GetSubscription(topicName, subscriptionName string) *Subscription {
	return NewSubscription(topicName, subscriptionName, this.client)
}

func getErrFromMessage(s string) int {
	//the error message is like this 
	// {
  //   "code": "5100",
  //   "message": "(100004)projectId不正确"
	// }
	if strings.ContainsAny(s, "(") && strings.ContainsAny(s, ")") {
		start := strings.IndexByte(s, '(')
		end := strings.IndexByte(s, ')')
		errCode := s[start+1:end]
		result, err := strconv.Atoi(errCode)
		if err == nil {
			return result
		}
	}
	return 0
}

//According to the latest cmq api document, more error info about cmq is in the "message" of resp
func doCall(client *CMQClient, param map[string]string, opration string) (resMap map[string]interface{}, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	res, err := client.call(opration, param)
	if err != nil {
		//log.Printf("client.call %v failed: %v\n", opration, err.Error())
		return
	}
	//log.Printf("res: %v", res)

	resMap = make(map[string]interface{}, 0)
	err = json.Unmarshal([]byte(res), &resMap)
	if err != nil {
		//log.Printf(err.Error())
		return
	}
	code = int(resMap["code"].(float64))
	if code != 0 {
		moduleErrCode = getErrFromMessage(resMap["message"].(string))
		err = fmt.Errorf("%v failed: code: %v, message: %v, requestId: %v", opration, code, resMap["message"], resMap["requestId"])
		//log.Printf(err.Error())
		return
	}

	return
}
