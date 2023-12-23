package view

import "fmt"

type Log struct {
}

func (l *Log) Debug(v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprint(v...)
}
func (l *Log) Debugf(format string, v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprintf(format, v...)
}

func (l *Log) Info(v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprint(v...)

}
func (l *Log) Infof(format string, v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprintf(format, v...)
}

func (l *Log) Warn(v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprint(v...)
}
func (l *Log) Warnf(format string, v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprintf(format, v...)
}

func (l *Log) Error(v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprint(v...)
}
func (l *Log) Errorf(format string, v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprintf(format, v...)
}

func (l *Log) Fatal(v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprint(v...)
}
func (l *Log) Fatalf(format string, v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprintf(format, v...)
}

func (l *Log) Panic(v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprint(v...)
}
func (l *Log) Panicf(format string, v ...interface{}) {
	GlobalText.msgChan <- fmt.Sprintf(format, v...)
}
