package rpclog

import "testing"

func TestLog(t *testing.T) {
	Error("err msg")
	Errorf("this is:%s", "err msg")
	Info("info msg")
	Errorf("this is:%s", "info msg")

	Warn("warn msg")
	Warnf("this is:%s", "warn msg")

	Debug("debug msg")
	Debugf("this is:%s", "debug msg")
}
