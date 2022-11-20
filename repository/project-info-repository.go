package repository

import (
	propertiDTO "go-api/dto/properti"

	"gorm.io/gorm"
)

type ProjectInfoRepository interface {
	FindByID(id int) propertiDTO.Project
}

type projectInfoConnection struct {
	connection *gorm.DB
}

func NewProjectInfoRepository(conn *gorm.DB) ProjectInfoRepository {
	return &projectInfoConnection{
		connection: conn,
	}
}

func (db *projectInfoConnection) FindByID(id int) propertiDTO.Project {
	var result propertiDTO.Project
	db.connection.Raw("SELECT * from tbl_project where id", id).Take(&result)
	return result
}
