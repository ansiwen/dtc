module github.com/niclabs/dtc/v3

go 1.12

require (
	github.com/google/uuid v1.1.1
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/miekg/pkcs11 v1.0.2
	github.com/niclabs/dtcconfig v1.0.5 // indirect
	github.com/niclabs/dtcnode/v3 v3.0.0
	github.com/niclabs/tcecdsa v0.0.4
	github.com/niclabs/tcrsa v0.0.4
	github.com/pebbe/zmq4 v1.0.0
	github.com/spf13/viper v1.4.0
)

replace (
	github.com/niclabs/dtcnode/v3 v3.0.0 => /mnt/data/code/go/dtcnode
	github.com/niclabs/tcecdsa v0.0.3 => /mnt/data/code/go/tcecdsa
	github.com/niclabs/tcpaillier v0.0.5 => /mnt/data/code/go/tcpaillier

)
