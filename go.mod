module github.com/videocoin/transcode

go 1.12

require (
	github.com/ethereum/go-ethereum v1.8.27
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gogo/protobuf v1.2.1
	github.com/google/uuid v1.1.1
	github.com/grafov/m3u8 v0.11.1
	github.com/hashicorp/consul/api v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/common v0.7.0
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/videocoin/cloud-api v0.0.17
	github.com/videocoin/cloud-pkg v0.0.5
	github.com/videocoin/telegraf v0.0.0-20190710212555-97388fb1c745
	google.golang.org/api v0.3.1
	google.golang.org/grpc v1.23.1
)

replace github.com/videocoin/cloud-api => ../cloud-api
