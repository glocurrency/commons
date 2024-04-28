package mysql

import (
	"fmt"
	"time"

	"github.com/glocurrency/commons/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewOrm(dsn string) (*gorm.DB, error) {
	dialector := mysql.New(mysql.Config{DSN: dsn, DefaultStringSize: 256})

	cfg := &gorm.Config{Logger: logger.NewGormLogger(logger.Log(), 200*time.Millisecond, true)}

	orm, err := gorm.Open(dialector, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot open mysql session: %w", err)
	}

	return orm, nil
}

func Migrate(orm *gorm.DB, dst ...interface{}) error {
	return orm.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").
		AutoMigrate(dst...)
}

func Drop(orm *gorm.DB, dst ...interface{}) error {
	mig := orm.Migrator()
	return mig.DropTable(dst...)
}
