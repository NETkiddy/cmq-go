package cmq_go

const (
	NotifyStrategyDefault      = "BACKOFF_RETRY"
	NotifyContentFormatDefault = "JSON"
	NotifyContentFormatSimplified = "SIMPLIFIED"
)

type SubscriptionMeta struct {
	//Subscription 订阅的主题所有者的appId
	TopicOwner string
	//订阅的终端地址
	Endpoint string
	//订阅的协议
	Protocal string
	//推送消息出现错误时的重试策略
	NotifyStrategy string
	//向 Endpoint 推送的消息内容格式
	NotifyContentFormat string
	//描述了该订阅中消息过滤的标签列表（仅标签一致的消息才会被推送）
	FilterTag []string
	//Subscription 的创建时间，从 1970-1-1 00:00:00 到现在的秒值
	CreateTime int
	//修改 Subscription 属性信息最近时间，从 1970-1-1 00:00:00 到现在的秒值
	LastModifyTime int
	//该订阅待投递的消息数
	MsgCount   int
	BindingKey []string
}

func NewSubscriptionMeta() *SubscriptionMeta {
	return &SubscriptionMeta{
		TopicOwner:          "",
		Endpoint:            "",
		Protocal:            "",
		NotifyStrategy:      NotifyStrategyDefault,
		NotifyContentFormat: NotifyContentFormatDefault,
		FilterTag:           nil,
		CreateTime:          0,
		LastModifyTime:      0,
		MsgCount:            0,
		BindingKey:          nil,
	}
}
