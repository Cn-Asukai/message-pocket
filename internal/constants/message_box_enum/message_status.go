package message_box_enum

type StatusType int32

func (r StatusType) Val() int32 {
	return int32(r)
}

const (
	// Pending 发送中
	Pending StatusType = iota + 1
	// Sent 发送成功
	Sent
)
