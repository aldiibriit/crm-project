package InternalRepository

import "gorm.io/gorm"

type BaseRepository interface {
	StartTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB)
	RollbackTransaction(tx *gorm.DB)
}

type baseConnection struct {
	connection *gorm.DB
}

func NewBaseRepository(conn *gorm.DB) BaseRepository {
	return &baseConnection{
		connection: conn,
	}
}

func (db *baseConnection) StartTransaction() *gorm.DB {
	return db.connection.Begin()
}

func (db *baseConnection) CommitTransaction(tx *gorm.DB) {
	tx.Debug().Commit()
}

func (db *baseConnection) RollbackTransaction(tx *gorm.DB) {
	tx.Debug().Rollback()
}
