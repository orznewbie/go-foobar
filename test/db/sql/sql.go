package sql

import (
	"github.com/jmoiron/sqlx"
	"time"
)

func NewDB(driver, dns string) *sqlx.DB {
	db, err := sqlx.Connect(driver, dns)
	if err != nil {
		panic(err)
	}
	// 最大的连接数量maxOpen
	// numOpen是已打开的连接数量，所以numOpen <= maxOpen
	db.SetMaxOpenConns(5)
	// maxIdleCount是连接池中最大的可以重复使用的连接数量，默认是2
	// maxIdleCount <= maxOpen，如果maxIdleCount设置的比maxOpen还大，会自动调整maxIdleCount = maxOpen
	// freeConn []*driverConn为可重复使用的连接，所以len(freeConn) <= maxIdleOpen
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db
}
