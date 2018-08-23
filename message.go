package cmq_go

type Message struct {
	/** 服务器返回的消息ID */
	MsgId string
	/** 每次消费唯一的消息句柄，用于删除等操作 */
	ReceiptHandle string
	/** 消息体 */
	MsgBody string
	/** 消息发送到队列的时间，从 1970年1月1日 00:00:00 000 开始的毫秒数 */
	EnqueueTime int64
	/** 消息下次可见的时间，从 1970年1月1日 00:00:00 000 开始的毫秒数 */
	NextVisibleTime int64
	/** 消息第一次出队列的时间，从 1970年1月1日 00:00:00 000 开始的毫秒数 */
	FirstDequeueTime int64
	/** 出队列次数 */
	DequeueCount int
	MsgTag       []string
}
