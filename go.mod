module github.com/orznewbie/gotmpl

go 1.16

require (
	cloud.google.com/go v0.103.0
	github.com/99designs/gqlgen v0.17.9
	github.com/beeker1121/goque v2.1.0+incompatible
	github.com/dgraph-io/dgo/v210 v210.0.0-20211129111319-4c8247ebe697
	github.com/go-co-op/gocron v1.15.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gomodule/redigo v1.8.9
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3
	github.com/jmoiron/sqlx v1.3.5
	github.com/json-iterator/go v1.1.11
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lib/pq v1.2.0
	github.com/mailru/easyjson v0.7.7
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/redislabs/redisgraph-go v2.0.2+incompatible
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.5
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/vektah/gqlparser/v2 v2.4.4
	go.einride.tech/aip v0.54.1
	go.etcd.io/etcd/api/v3 v3.5.2
	go.etcd.io/etcd/client/v3 v3.5.2
	go.uber.org/zap v1.21.0
	google.golang.org/genproto v0.0.0-20220628213854-d9e0b6570c03
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gotest.tools/v3 v3.2.0 // indirect
)

replace github.com/vektah/gqlparser/v2 v2.4.4 => github.com/vektah/gqlparser/v2 v2.3.0
