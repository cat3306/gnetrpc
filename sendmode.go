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

func CallNone() *CallMode {
	return &CallMode{
		Call: None,
	}
}
func CallSelf() *CallMode {
	return &CallMode{
		Call: Self,
	}
}

func CallBroadcast() *CallMode {
	return &CallMode{
		Call: Broadcast,
	}
}

func CallBroadcastExceptSelf() *CallMode {
	return &CallMode{
		Call: BroadcastExceptSelf,
	}
}
