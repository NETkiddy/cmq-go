package cmq_go

type TopicMeta struct {
	// 当前该主题的消息堆积数
	MsgCount int
	// 消息最大长度，取值范围1024-1048576 Byte（即1-1024K），默认1048576
	MaxMsgSize int
	//消息在主题中最长存活时间，从发送到该主题开始经过此参数指定的时间后，
	//不论消息是否被成功推送给用户都将被删除，单位为秒。固定为一天，该属性不能修改。
	MsgRetentionSeconds int
	//创建时间
	CreateTime int
	//修改属性信息最近时间
	LastModifyTime int
	LoggingEnabled int
	FilterType     int
}

func NewTopicMeta() *TopicMeta {
	return &TopicMeta{
		MsgCount:            0,
		MaxMsgSize:          1048576,
		MsgRetentionSeconds: 86400,
		CreateTime:          0,
		LastModifyTime:      0,
		LoggingEnabled:      0,
		FilterType:          1,
	}
}
