package InternalRepository

import (
	CCTVResponseDTO "go-api/dto/response/cctv"

	"gorm.io/gorm"
)

type CCTVRepository interface {
	All() []CCTVResponseDTO.Cctv
	FindBySN(sn string) CCTVResponseDTO.Cctv
}

type cctvConnection struct {
	connection *gorm.DB
}

func NewCCTVRepository(conn *gorm.DB) CCTVRepository {
	return &cctvConnection{
		connection: conn,
	}
}

func (db *cctvConnection) All() []CCTVResponseDTO.Cctv {
	var result []CCTVResponseDTO.Cctv
	db.connection.Raw("SELECT * FROM tb_cctv").Find(&result)
	return result
}

func (db *cctvConnection) FindBySN(sn string) CCTVResponseDTO.Cctv {
	var result CCTVResponseDTO.Cctv
	db.connection.Raw(`SELECT * FROM tb_cctv where sn_cctv = '` + sn + `'`).Find(&result)
	return result
}
