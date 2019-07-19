package cmq_go

import (
	"strconv"
	"fmt"
)

type Queue struct {
	queueName string
	client    *CMQClient
}

func NewQueue(queueName string, client *CMQClient) (queue *Queue) {
	return &Queue{
		queueName: queueName,
		client:    client,
	}
}

func (this *Queue) SetQueueAttributes(queueMeta QueueMeta) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["queueName"] = this.queueName

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

	_, err, code, moduleErrCode = doCall(this.client, param, "SetQueueAttributes")
	if err != nil {
		//log.Printf("client.call SetQueueAttributes failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Queue) GetQueueAttributes() (queueMeta QueueMeta, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["queueName"] = this.queueName

	resMap, err, code, moduleErrCode := doCall(this.client, param, "GetQueueAttributes")
	if err != nil {
		//log.Printf("client.call GetQueueAttributes failed: %v\n", err.Error())
		return
	}

	queueMeta.MaxMsgHeapNum = int(resMap["maxMsgHeapNum"].(float64))
	queueMeta.PollingWaitSeconds = int(resMap["pollingWaitSeconds"].(float64))
	queueMeta.VisibilityTimeout = int(resMap["visibilityTimeout"].(float64))
	queueMeta.MaxMsgSize = int(resMap["maxMsgSize"].(float64))
	queueMeta.MsgRetentionSeconds = int(resMap["msgRetentionSeconds"].(float64))
	queueMeta.CreateTime = int(resMap["createTime"].(float64))
	queueMeta.LastModifyTime = int(resMap["lastModifyTime"].(float64))
	queueMeta.ActiveMsgNum = int(resMap["activeMsgNum"].(float64))
	queueMeta.InactiveMsgNum = int(resMap["inactiveMsgNum"].(float64))
	queueMeta.RewindMsgNum = int(resMap["rewindMsgNum"].(float64))
	queueMeta.MinMsgTime = int(resMap["minMsgTime"].(float64))
	queueMeta.DelayMsgNum = int(resMap["delayMsgNum"].(float64))
	queueMeta.RewindSeconds = int(resMap["rewindSeconds"].(float64))

	return
}

func (this *Queue) SendMessage(msgBody string) (messageId string, err error, code, moduleErrCode int) {
	messageId, err, code, moduleErrCode = _sendMessage(this.client, msgBody, this.queueName, 0)
	return
}

func (this *Queue) SendDelayMessage(msgBody string, delaySeconds int) (messageId string, err error, code, moduleErrCode int) {
	messageId, err, code, moduleErrCode = _sendMessage(this.client, msgBody, this.queueName, delaySeconds)
	return
}

func _sendMessage(client *CMQClient, msgBody, queueName string, delaySeconds int) (messageId string, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["queueName"] = queueName
	param["msgBody"] = msgBody
	param["delaySeconds"] = strconv.Itoa(delaySeconds)

	resMap, err, code, moduleErrCode := doCall(client, param, "SendMessage")
	if err != nil {
		//log.Printf("client.call GetQueueAttributes failed: %v\n", err.Error())
		return
	}

	messageId = resMap["msgId"].(string)
	return
}

func (this *Queue) BatchSendMessage(msgBodys []string) (messageIds []string, err error, code, moduleErrCode int) {
	messageIds, err, code, moduleErrCode = _batchSendMessage(this.client, msgBodys, this.queueName, 0)
	return
}

func (this *Queue) BatchSendDelayMessage(msgBodys []string, delaySeconds int) (messageIds []string, err error, code, moduleErrCode int) {
	messageIds, err, code, moduleErrCode = _batchSendMessage(this.client, msgBodys, this.queueName, delaySeconds)
	return
}

func _batchSendMessage(client *CMQClient, msgBodys []string, queueName string, delaySeconds int) (messageIds []string, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	messageIds = make([]string, 0)

	if len(msgBodys) == 0 || len(msgBodys) > 16 {
		err = fmt.Errorf("message size is 0 or more than 16")
		//log.Printf("%v", err.Error())
		return
	}

	param := make(map[string]string)
	param["queueName"] = queueName
	for i, msgBody := range msgBodys {
		param["msgBody."+strconv.Itoa(i+1)] = msgBody
	}
	param["delaySeconds"] = strconv.Itoa(delaySeconds)

	resMap, err, code, moduleErrCode := doCall(client, param, "BatchSendMessage")
	if err != nil {
		//log.Printf("client.call BatchSendMessage failed: %v\n", err.Error())
		return
	}

	msgList := resMap["msgList"].([]interface{})
	for _, v := range msgList {
		msg := v.(map[string]interface{})
		messageIds = append(messageIds, msg["msgId"].(string))
	}
	return
}

func (this *Queue) ReceiveMessage(pollingWaitSeconds int) (msg Message, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["queueName"] = this.queueName
	if pollingWaitSeconds >= 0 {
		param["UserpollingWaitSeconds"] = strconv.Itoa(pollingWaitSeconds * 1000)
		param["pollingWaitSeconds"] = strconv.Itoa(pollingWaitSeconds)
	} else {
		param["UserpollingWaitSeconds"] = strconv.Itoa(30000)
	}

	resMap, err, code, moduleErrCode := doCall(this.client, param, "ReceiveMessage")
	if err != nil {
		//log.Printf("client.call ReceiveMessage failed: %v\n", err.Error())
		return
	}

	msg.MsgId = resMap["msgId"].(string)
	msg.ReceiptHandle = resMap["receiptHandle"].(string)
	msg.MsgBody = resMap["msgBody"].(string)
	msg.EnqueueTime = int64(resMap["enqueueTime"].(float64))
	msg.NextVisibleTime = int64(resMap["nextVisibleTime"].(float64))
	msg.FirstDequeueTime = int64(resMap["firstDequeueTime"].(float64))
	msg.DequeueCount = int(resMap["dequeueCount"].(float64))

	return
}

func (this *Queue) BatchReceiveMessage(numOfMsg, pollingWaitSeconds int) (msgs []Message, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	msgs = make([]Message, 0)
	param := make(map[string]string)
	param["queueName"] = this.queueName
	param["numOfMsg"] = strconv.Itoa(numOfMsg)
	if pollingWaitSeconds >= 0 {
		param["UserpollingWaitSeconds"] = strconv.Itoa(pollingWaitSeconds * 1000)
		param["pollingWaitSeconds"] = strconv.Itoa(pollingWaitSeconds)
	} else {
		param["UserpollingWaitSeconds"] = strconv.Itoa(30000)
	}

	resMap, err, code, moduleErrCode := doCall(this.client, param, "BatchReceiveMessage")
	if err != nil {
		//log.Printf("client.call BatchReceiveMessage failed: %v\n", err.Error())
		return
	}
	msgInfoList := resMap["msgInfoList"].([]interface{})
	for _, v := range msgInfoList {
		msgInfo := v.(map[string]interface{})
		msg := Message{}
		msg.MsgId = msgInfo["msgId"].(string)
		msg.ReceiptHandle = msgInfo["receiptHandle"].(string)
		msg.MsgBody = msgInfo["msgBody"].(string)
		msg.EnqueueTime = int64(msgInfo["enqueueTime"].(float64))
		msg.NextVisibleTime = int64(msgInfo["nextVisibleTime"].(float64))
		msg.FirstDequeueTime = int64(msgInfo["firstDequeueTime"].(float64))
		msg.DequeueCount = int(msgInfo["dequeueCount"].(float64))

		msgs = append(msgs, msg)
	}

	return
}

func (this *Queue) DeleteMessage(receiptHandle string) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["queueName"] = this.queueName
	param["receiptHandle"] = receiptHandle

	_, err, code, moduleErrCode = doCall(this.client, param, "DeleteMessage")
	if err != nil {
		//log.Printf("client.call DeleteMessage failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Queue) BatchDeleteMessage(receiptHandles []string) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	if len(receiptHandles) == 0 {
		return
	}
	param := make(map[string]string)
	param["queueName"] = this.queueName
	for i, receiptHandle := range receiptHandles {
		param["receiptHandle."+strconv.Itoa(i+1)] = receiptHandle
	}

	_, err, code, moduleErrCode = doCall(this.client, param, "BatchDeleteMessage")
	if err != nil {
		//log.Printf("client.call BatchDeleteMessage failed: %v\n", err.Error())
		return
	}
	return
}

func (this *Queue) RewindQueue(backTrackingTime int) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	if backTrackingTime <= 0 {
		return
	}
	param := make(map[string]string)
	param["queueName"] = this.queueName
	param["startConsumeTime"] = strconv.Itoa(backTrackingTime)

	_, err, code, moduleErrCode = doCall(this.client, param, "RewindQueue")
	if err != nil {
		//log.Printf("client.call RewindQueue failed: %v\n", err.Error())
		return
	}
	return
}
