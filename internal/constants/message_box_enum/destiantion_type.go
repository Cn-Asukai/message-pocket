package message_box_enum

type DestinationType int32

func (r DestinationType) Val() int32 {
	return int32(r)
}

const (
	// DestinationQQGroup QQ群
	DestinationQQGroup DestinationType = iota + 1
)

