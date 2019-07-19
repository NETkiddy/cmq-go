package cmq_go

import (
	"strconv"
	"fmt"
)

type Topic struct {
	topicName string
	client    *CMQClient
}

func NewTopic(topicName string, client *CMQClient) (queue *Topic) {
	return &Topic{
		topicName: topicName,
		client:    client,
	}
}

func (this *Topic) SetTopicAttributes(maxMsgSize int) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	if maxMsgSize < 1024 || maxMsgSize > 1048576 {
		err = fmt.Errorf("Invalid parameter maxMsgSize < 1KB or maxMsgSize > 1024KB")
		//log.Printf("%v", err.Error())
		return
	}
	param := make(map[string]string)
	param["topicName"] = this.topicName
	if maxMsgSize > 0 {
		param["maxMsgSize"] = strconv.Itoa(maxMsgSize)
	}

	_, err, code, moduleErrCode = doCall(this.client, param, "SetTopicAttributes")
	if err != nil {
		//log.Printf("client.call SetTopicAttributes failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Topic) GetTopicAttributes() (meta TopicMeta, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["topicName"] = this.topicName

	resMap, err, code, moduleErrCode := doCall(this.client, param, "GetTopicAttributes")
	if err != nil {
		//log.Printf("client.call GetTopicAttributes failed: %v\n", err.Error())
		return
	}
	pmeta := NewTopicMeta()
	pmeta.MsgCount = int(resMap["msgCount"].(float64))
	pmeta.MaxMsgSize = int(resMap["maxMsgSize"].(float64))
	pmeta.MsgRetentionSeconds = int(resMap["msgRetentionSeconds"].(float64))
	pmeta.CreateTime = int(resMap["createTime"].(float64))
	pmeta.LastModifyTime = int(resMap["lastModifyTime"].(float64))

	meta = *pmeta
	return
}

func (this *Topic) PublishMessage(message string) (msgId string, err error, code, moduleErrCode int) {
	msgId, err, code, moduleErrCode = _publishMessage(this.client, this.topicName, message, nil, "")
	return
}

func _publishMessage(client *CMQClient, topicName, msg string, tagList []string, routingKey string) (
	msgId string, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["topicName"] = topicName
	param["msgBody"] = msg
	if routingKey != "" {
		param["routingKey"] = routingKey
	}
	if tagList != nil {
		for i, tag := range tagList {
			param["msgTag."+strconv.Itoa(i+1)] = tag
		}
	}
	resMap, err, code, moduleErrCode := doCall(client, param, "PublishMessage")
	if err != nil {
		//log.Printf("client.call GetTopicAttributes failed: %v\n", err.Error())
		return
	}
	msgId = resMap["msgId"].(string)

	return
}

func (this *Topic) BatchPublishMessage(msgList []string) (msgIds []string, err error, code, moduleErrCode int) {
	msgIds, err, code, moduleErrCode = _batchPublishMessage(this.client, this.topicName, msgList, nil, "")
	return
}

func _batchPublishMessage(client *CMQClient, topicName string, msgList, tagList []string, routingKey string) (
	msgIds []string, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["topicName"] = topicName
	if routingKey != "" {
		param["routingKey"] = routingKey
	}
	if msgList != nil {
		for i, msg := range msgList {
			param["msgBody."+strconv.Itoa(i+1)] = msg
		}
	}
	if tagList != nil {
		for i, tag := range tagList {
			param["msgTag."+strconv.Itoa(i+1)] = tag
		}
	}

	resMap, err, code, moduleErrCode := doCall(client, param, "BatchPublishMessage")
	if err != nil {
		//log.Printf("client.call BatchPublishMessage failed: %v\n", err.Error())
		return
	}
	resMsgList := resMap["msgList"].([]interface{})
	for _, v := range resMsgList {
		msg := v.(map[string]interface{})
		msgIds = append(msgIds, msg["msgId"].(string))
	}

	return
}

func (this *Topic) ListSubscription(offset, limit int, searchWord string) (totalCount int, subscriptionList []string, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["topicName"] = this.topicName
	if searchWord != "" {
		param["searchWord "] = searchWord
	}
	if offset >= 0 {
		param["offset "] = strconv.Itoa(offset)
	}
	if limit > 0 {
		param["limit "] = strconv.Itoa(limit)
	}

	resMap, err, code, moduleErrCode := doCall(this.client, param, "ListSubscriptionByTopic")
	if err != nil {
		//log.Printf("client.call ListSubscriptionByTopic failed: %v\n", err.Error())
		return
	}

	totalCount = int(resMap["totalCount"].(float64))
	resSubscriptionList := resMap["subscriptionList"].([]interface{})
	for _, v := range resSubscriptionList {
		subscribe := v.(map[string]interface{})
		subscriptionList = append(subscriptionList, subscribe["subscriptionName"].(string))
	}

	return
}
