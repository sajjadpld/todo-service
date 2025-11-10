package orm

import "gorm.io/gorm"

//go:generate mockgen -source=./contract.go -destination=./mocks/orm_mock.go -package=orm_mock
type ISql interface {
	ISqlGeneric
	ISqlTx
}

type (
	ISqlGeneric interface {
		Init()
		C() *gorm.DB
		Migrate(path string)
		Seed()
		Stop()
	}

	ISqlTx interface {
		Begin()
		Commit() error
		Rollback() error
		// Resolve commit or rollback transaction by getting the error
		Resolve(dbErr error) error
	}
)
