package mysql

import (
	"database/sql"

	"gorm.io/gorm"
)

type Tx interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
}

// FakeTx is a test helper that implements Tx interface.
// It executes SQL statements without wrapping them in transaction.
type FakeTx struct{}

func (*FakeTx) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return fc(&gorm.DB{})
}
