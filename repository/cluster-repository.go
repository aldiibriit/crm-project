package repository

import (
	propertiDTO "go-api/dto/properti"

	"gorm.io/gorm"
)

type ClusterInfoRepository interface {
	FindByID(id int) propertiDTO.Cluster
}

type clusterInfoConnection struct {
	connection *gorm.DB
}

func NewClusterInfoRepository(conn *gorm.DB) ClusterInfoRepository {
	return &clusterInfoConnection{
		connection: conn,
	}
}

func (db *clusterInfoConnection) FindByID(id int) propertiDTO.Cluster {
	var result propertiDTO.Cluster
	db.connection.Raw(`select *,case when is_cluster = 1 then "1" when is_cluster = 0 then "0" else "NA" end as isCluster from tbl_cluster where id`, id).Take(&result)
	return result
}
