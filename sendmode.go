package gnetrpc

type CallType int

const (
	None CallType = iota
	Self
	Broadcast
	BroadcastExceptSelf
)

type CallMode struct {
	Call CallType
	FDs  []int
}
