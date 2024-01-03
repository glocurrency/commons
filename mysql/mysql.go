package mysql

import (
	"fmt"
	"time"

	"github.com/glocurrency/commons/logger"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewOrm(dsn string) (*gorm.DB, error) {
	d := mysql.New(mysql.Config{DriverName: "nrmysql", DSN: dsn, DefaultStringSize: 256})
	cfg := &gorm.Config{Logger: logger.NewGormLogger(logger.Log(), 200*time.Millisecond, true)}

	orm, err := gorm.Open(d, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot open mysql session: %w", err)
	}

	return orm, nil
}

func Migrate(orm *gorm.DB, dst ...interface{}) error {
	orm = orm.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	return orm.AutoMigrate(dst...)
}

func Drop(orm *gorm.DB, dst ...interface{}) error {
	mig := orm.Migrator()
	return mig.DropTable(dst...)
}
