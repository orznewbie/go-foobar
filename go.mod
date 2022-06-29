module github.com/orznewbie/gotmpl

go 1.16

require (
	github.com/99designs/gqlgen v0.17.9
	github.com/RichardKnop/machinery/v2 v2.0.11
	github.com/beeker1121/goque v2.1.0+incompatible
	github.com/dgraph-io/dgo/v210 v210.0.0-20211129111319-4c8247ebe697
	github.com/go-redis/redis_rate/v9 v9.1.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/json-iterator/go v1.1.11
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lib/pq v1.2.0
	github.com/mailru/easyjson v0.7.7
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.1
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/urfave/cli/v2 v2.8.1
	github.com/vektah/gqlparser/v2 v2.4.4
	github.com/vmihailenco/taskq/v3 v3.2.8
	go.etcd.io/etcd/api/v3 v3.5.2
	go.etcd.io/etcd/client/v3 v3.5.2
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/vektah/gqlparser/v2 v2.4.4 => github.com/vektah/gqlparser/v2 v2.3.0
