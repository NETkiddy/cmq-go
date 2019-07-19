package cmq_go

import (
	"strconv"
)

type Subscription struct {
	topicName        string
	subscriptionName string
	client           *CMQClient
}

func NewSubscription(topicName, subscriptionName string, client *CMQClient) *Subscription {
	return &Subscription{
		topicName:        topicName,
		subscriptionName: subscriptionName,
		client:           client,
	}
}

func (this *Subscription) ClearFilterTags() (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["topicName"] = this.topicName
	param["subscriptionName "] = this.subscriptionName

	_, err, code, moduleErrCode = doCall(this.client, param, "ClearSubscriptionFilterTags")
	if err != nil {
		//log.Printf("client.call ClearSubscriptionFilterTags failed: %v\n", err.Error())
		return
	}

	return
}

func (this *Subscription) SetSubscriptionAttributes(meta SubscriptionMeta) (err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["topicName"] = this.topicName
	param["subscriptionName "] = this.subscriptionName
	if meta.NotifyStrategy != "" {
		param["notifyStrategy"] = meta.NotifyStrategy
	}
	if meta.NotifyContentFormat != "" {
		param["notifyContentFormat"] = meta.NotifyContentFormat
	}
	if meta.FilterTag != nil {
		for i, flag := range meta.FilterTag {
			param["filterTag."+strconv.Itoa(i+1)] = flag
		}
	}
	if meta.BindingKey != nil {
		for i, binding := range meta.BindingKey {
			param["bindingKey."+strconv.Itoa(i+1)] = binding
		}
	}

	_, err, code, moduleErrCode = doCall(this.client, param, "SetSubscriptionAttributes")
	if err != nil {
		//log.Printf("client.call SetSubscriptionAttributes failed: %v\n", err.Error())
		return
	}

	return
}

func (this *Subscription) GetSubscriptionAttributes() (meta *SubscriptionMeta, err error, code, moduleErrCode int) {
	code = DEFAULT_ERROR_CODE
	param := make(map[string]string)
	param["topicName"] = this.topicName
	param["subscriptionName"] = this.subscriptionName

	resMap, err, code, moduleErrCode := doCall(this.client, param, "GetSubscriptionAttributes")
	if err != nil {
		//log.Printf("client.call GetSubscriptionAttributes failed: %v\n", err.Error())
		return
	}

	meta = NewSubscriptionMeta()
	meta.FilterTag = make([]string, 0)
	meta.BindingKey = make([]string, 0)
	meta.TopicOwner = resMap["topicOwner"].(string)
	meta.Endpoint = resMap["endpoint"].(string)
	meta.Protocal = resMap["protocol"].(string)
	meta.NotifyStrategy = resMap["notifyStrategy"].(string)
	meta.NotifyContentFormat = resMap["notifyContentFormat"].(string)
	meta.CreateTime = int(resMap["createTime"].(float64))
	meta.LastModifyTime = int(resMap["lastModifyTime"].(float64))
	meta.MsgCount = int(resMap["msgCount"].(float64))
	if filterTag, found := resMap["filterTag"]; found {
		for _, v := range filterTag.([]interface{}) {
			filter := v.(string)
			meta.FilterTag = append(meta.FilterTag, filter)
		}
	}
	if bindingKey, found := resMap["bindingKey"]; found {
		for _, v := range bindingKey.([]interface{}) {
			binding := v.(string)
			meta.BindingKey = append(meta.BindingKey, binding)
		}
	}

	return
}
