module github.com/cat3306/gnetrpc

go 1.20

require (
	github.com/lithammer/shortuuid/v4 v4.0.0
	github.com/panjf2000/ants/v2 v2.8.2
	github.com/panjf2000/gnet/v2 v2.3.3
	github.com/valyala/bytebufferpool v1.0.0
	go.uber.org/zap v1.24.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/google/uuid v1.3.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/panjf2000/gnet/v2 => github.com/cat3306/gnet/v2 v2.3.4-0.20231219133356-c8a16922f414
