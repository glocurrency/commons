package mysql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glocurrency/commons/mysql" // Adjust to your actual import path
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Dummy model for testing
type User struct {
	ID   uint
	Name string
}

// setupMockDB creates a GORM instance backed by go-sqlmock
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := gormmysql.New(gormmysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true, // Crucial for mocking MySQL
	})

	orm, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return orm, mock
}

func TestNewOrm(t *testing.T) {
	t.Run("fails with invalid DSN", func(t *testing.T) {
		// DSN format is strict; an obviously broken one should fail to connect/ping
		orm, err := mysql.NewOrm("invalid-dsn-format")

		require.Error(t, err)
		assert.Nil(t, orm)
		assert.Contains(t, err.Error(), "cannot open mysql session")
	})
}

func TestMigrate(t *testing.T) {
	orm, mock := setupMockDB(t)

	t.Run("successfully runs automigrate", func(t *testing.T) {
		// 1. GORM checks the database schema existence first
		mock.ExpectQuery("^SELECT SCHEMA_NAME from Information_schema.SCHEMATA").
			WillReturnRows(sqlmock.NewRows([]string{"SCHEMA_NAME"}).AddRow("test_db"))

		// 2. GORM checks if the table exists (it might skip this depending on GORM version,
		// but if your logs ask for it, keep it)
		mock.ExpectQuery("^SELECT count").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// 3. GORM creates the table
		mock.ExpectExec("^CREATE TABLE `users`").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := mysql.Migrate(orm, &User{})

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDrop(t *testing.T) {
	orm, mock := setupMockDB(t)

	t.Run("successfully drops table", func(t *testing.T) {
		// 1. GORM disables foreign key checks to prevent dependency errors
		mock.ExpectExec("^SET FOREIGN_KEY_CHECKS = 0").
			WillReturnResult(sqlmock.NewResult(0, 0))

		// 2. GORM drops the table (using regex to catch the optional CASCADE)
		mock.ExpectExec("^DROP TABLE IF EXISTS `users`").
			WillReturnResult(sqlmock.NewResult(0, 0))

		// 3. GORM re-enables foreign key checks
		mock.ExpectExec("^SET FOREIGN_KEY_CHECKS = 1").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := mysql.Drop(orm, &User{})

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
