module github.com/celo-org/celo-blockchain

go 1.13

require (
	github.com/Azure/azure-storage-blob-go v0.7.0
	github.com/VictoriaMetrics/fastcache v1.6.0
	github.com/aristanetworks/goarista v0.0.0-20170210015632-ea17b1a17847
	github.com/aws/aws-sdk-go v1.25.48
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/buraksezer/consistent v0.0.0-20191006190839-693edf70fd72
	github.com/celo-org/celo-bls-go v0.2.4
	github.com/cespare/cp v0.1.0
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/cloudflare/cloudflare-go v0.14.0
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea
	github.com/docker/docker v1.4.2-0.20180625184442-8e610b2b55bf
	github.com/dop251/goja v0.0.0-20200721192441-a695b0cdd498
	github.com/ethereum/go-ethereum v1.10.8
	github.com/fatih/color v1.7.0
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5
	github.com/go-stack/stack v1.8.0
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.3
	github.com/gorilla/websocket v1.4.2
	github.com/graph-gophers/graphql-go v0.0.0-20201113091052-beb923fada29
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/hdevalence/ed25519consensus v0.0.0-20201207055737-7fde80a9d5ff
	github.com/holiman/uint256 v1.2.0
	github.com/huin/goupnp v1.0.2
	github.com/influxdata/influxdb v1.8.3
	github.com/jackpal/go-nat-pmp v1.0.2-0.20160603034137-1fa385a6f458
	github.com/julienschmidt/httprouter v1.2.0
	github.com/karalabe/usb v0.0.0-20190919080040-51dc0efba356
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mattn/go-colorable v0.1.8
	github.com/mattn/go-isatty v0.0.12
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/gomega v1.10.1
	github.com/pborman/uuid v0.0.0-20170112150404-1b00554d8222
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7
	github.com/prometheus/tsdb v0.7.1
	github.com/rjeczalik/notify v0.9.1
	github.com/rs/cors v1.7.0
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible
	github.com/shopspring/decimal v1.2.0
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200619000410-60c24ae608a6
	gopkg.in/urfave/cli.v1 v1.20.0
)

// Use our fork which disables bitcode
replace golang.org/x/mobile => github.com/celo-org/mobile v0.0.0-20210324213558-66ac87d7fb95
