package gnetrpc

type CallType int

const (
	None CallType = iota
	Self
	Broadcast
	BroadcastExceptSelf
	BroadcastSomeone
)

type CallMode struct {
	Call CallType
	Ids  []string
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
