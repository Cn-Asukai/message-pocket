package message_box_enum

type SourceType int32

func (r SourceType) Val() int32 {
	return int32(r)
}

const (
	// SourceTypeEO EdgeOne
	SourceTypeEO SourceType = iota + 1
)
