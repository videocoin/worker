module github.com/videocoin/transcode

go 1.12

require (
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/cespare/cp v1.1.1 // indirect
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/ethereum/go-ethereum v1.8.27
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/grafov/m3u8 v0.11.1
	github.com/hashicorp/consul/api v1.2.0
	github.com/huin/goupnp v1.0.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.1 // indirect
	github.com/karalabe/hid v1.0.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/common v0.7.0
	github.com/sirupsen/logrus v1.4.2
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/videocoin/cloud-api v0.0.17
	github.com/videocoin/cloud-pkg v0.0.5
	github.com/videocoin/telegraf v0.0.0-20190710212555-97388fb1c745
	go.uber.org/atomic v1.4.0 // indirect
	google.golang.org/grpc v1.23.1
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190709231704-1e4459ed25ff // indirect
)

replace github.com/videocoin/cloud-api => ../cloud-api

replace github.com/videocoin/cloud-pkg => ../cloud-pkg
