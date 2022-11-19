package repository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type JwtHistRepository interface {
	FindByEmail(email string) entity.JwtHistGo
	Save(entity.JwtHistGo) entity.JwtHistGo
}

type jwtHistConnection struct {
	connection *gorm.DB
}

func NewJwtHistRepository(conn *gorm.DB) JwtHistRepository {
	return &jwtHistConnection{
		connection: conn,
	}
}

func (db *jwtHistConnection) FindByEmail(email string) entity.JwtHistGo {
	var data entity.JwtHistGo
	db.connection.Debug().Raw("SELECT * FROM jwt_hist_go where email = ?", email).First(&data)
	return data
}

func (db *jwtHistConnection) Save(data entity.JwtHistGo) entity.JwtHistGo {
	var checker entity.JwtHistGo
	db.connection.Raw("SELECT * FROM jwt_hist_go where email = ?", data.Email).Take(&checker)
	if checker.Email == "" {
		db.connection.Save(&data)
	} else if checker.Email != "" {
		db.connection.Model(&entity.JwtHistGo{}).Where("email = ?", data.Email).Update("jwt", data.Jwt)
	}
	return data
}
