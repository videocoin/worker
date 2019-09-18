module github.com/videocoin/transcode

go 1.12

require (
	cloud.google.com/go v0.41.0 // indirect
	github.com/ethereum/go-ethereum v1.8.27
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gogo/protobuf v1.3.0
	github.com/google/uuid v1.0.0
	github.com/grafov/m3u8 v0.6.1
	github.com/grpc-ecosystem/grpc-gateway v1.9.3 // indirect
	github.com/hashicorp/consul/api v1.1.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/common v0.4.1
	github.com/rogpeppe/fastuuid v1.1.0 // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/videocoin/cloud-api v0.2.7
	github.com/videocoin/cloud-dispatcher v0.0.0-20190918112314-8857eafafc65
	github.com/videocoin/cloud-pkg v0.0.5
	github.com/videocoin/telegraf v0.0.0-20190710212555-97388fb1c745
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/exp v0.0.0-20190627132806-fd42eb6b336f // indirect
	golang.org/x/image v0.0.0-20190622003408-7e034cad6442 // indirect
	golang.org/x/mobile v0.0.0-20190607214518-6fa95d984e88 // indirect
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7 // indirect
	golang.org/x/sys v0.0.0-20190626221950-04f50cda93cb // indirect
	golang.org/x/tools v0.0.0-20190701194522-38ae2c8f6412 // indirect
	google.golang.org/api v0.7.0
	google.golang.org/genproto v0.0.0-20190627203621-eb59cef1c072 // indirect
	google.golang.org/grpc v1.23.0
)

replace github.com/videocoin/cloud-api => ../cloud-api
