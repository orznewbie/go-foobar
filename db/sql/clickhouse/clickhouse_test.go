package clickhouse

//func TestClickhouse(t *testing.T) {
//	conn := clickhouse.OpenDB(&clickhouse.Options{
//		Addr: []string{"127.0.0.1:9999"},
//		Auth: clickhouse.Auth{
//			Database: "default",
//			Username: "default",
//			Password: "",
//		},
//		TLS: &tls.Config{
//			InsecureSkipVerify: true,
//		},
//		Settings: clickhouse.Settings{
//			"max_execution_time": 60,
//		},
//		DialTimeout: 5 * time.Second,
//		Compression: &clickhouse.Compression{
//			clickhouse.CompressionLZ4,
//		},
//		Debug: true,
//	})
//	conn.SetMaxIdleConns(5)
//	conn.SetMaxOpenConns(10)
//	conn.SetConnMaxLifetime(time.Hour)
//}
