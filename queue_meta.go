package cmq_go

const (
	/** 缺省消息接收长轮询等待时间 */
	DEFAULT_POLLING_WAIT_SECONDS = 0
	/** 缺省消息可见性超时 */
	DEFAULT_VISIBILITY_TIMEOUT = 30
	/** 缺省消息最大长度，单位字节 */
	DEFAULT_MAX_MSG_SIZE = 1048576
	/** 缺省消息保留周期，单位秒 */
	DEFAULT_MSG_RETENTION_SECONDS = 345600
)

type QueueMeta struct {
	/** 最大堆积消息数 */
	MaxMsgHeapNum int
	/** 消息接收长轮询等待时间 */
	PollingWaitSeconds int
	/** 消息可见性超时 */
	VisibilityTimeout int
	/** 消息最大长度 */
	MaxMsgSize int
	/** 消息保留周期 */
	MsgRetentionSeconds int
	/** 队列创建时间 */
	CreateTime int
	/** 队列属性最后修改时间 */
	LastModifyTime int
	/** 队列处于Active状态的消息总数 */
	ActiveMsgNum int
	/** 队列处于Inactive状态的消息总数 */
	InactiveMsgNum int
	/** 已删除的消息，但还在回溯保留时间内的消息数量 */
	RewindMsgNum int
	/** 消息最小未消费时间 */
	MinMsgTime int
	/** 延时消息数量 */
	DelayMsgNum int
	/** 回溯时间 */
	RewindSeconds int
}
