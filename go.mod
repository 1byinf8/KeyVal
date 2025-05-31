module badies

go 1.23

toolchain go1.24.3

require (
	google.golang.org/grpc v1.72.2
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)

replace github.com/1byinf8/KeyVal => ../KeyVal
