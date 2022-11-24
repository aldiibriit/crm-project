package repository

import (
	propertiDTO "go-api/dto/properti"
	"go-api/dto/response/landingPageResponseDTO"

	"gorm.io/gorm"
)

type PropertiRepository interface {
	AdvancedFilter(sqlStr string) interface{}
	AdvancedFilter2(trxId string) []propertiDTO.TblImageProperti
	FindNearby() []landingPageResponseDTO.DetailPropertiDtoRes
}

type propertiConnection struct {
	connection *gorm.DB
}

// NewPropertiRepository creates an instance BookRepository
func NewPropertiRepository(dbConn *gorm.DB) PropertiRepository {
	return &propertiConnection{
		connection: dbConn,
	}
}

func (db *propertiConnection) AdvancedFilter(sqlStr string) interface{} {
	var result []propertiDTO.AdvancedFilter
	db.connection.Raw(sqlStr).Preload("ImageProperti").Preload("ImageProject").Find(&result)
	return result
}

func (db *propertiConnection) AdvancedFilter2(trxId string) []propertiDTO.TblImageProperti {
	var imageProject []propertiDTO.TblImageProperti
	db.connection.Raw(`SELECT * FROM tbl_image_properti WHERE trx_id = '` + trxId + `'`).Find(&imageProject)
	return imageProject
}

func (db *propertiConnection) FindNearby() []landingPageResponseDTO.DetailPropertiDtoRes {
	var data []landingPageResponseDTO.DetailPropertiDtoRes
	db.connection.Raw(`SELECT 
	prop.*,
	tip.*,
	(tkp.dapur = b'1')as dapur,
	(tkp.jalur_listrik = b'1')as jalur_listrik,
	(tkp.jalurpdam = b'1')as jalur_pdam,
	(tkp.jalur_telepone = b'1')as jalur_telepone,
	(tkp.ruang_keluarga = b'1')as ruang_keluarga,
	(tkp.ruang_kerja = b'1')as ruang_kerja,
	(fp.kolam_renang  = b'1')as kolamRenang,
	(fp.tempat_parkir = b'1')as tempatParkir,
	(fp.keamanan24jam = b'1')as keamanan24Jam,
	(fp.penghijauan = b'1')as penghijauan,
	(fp.rumah_sakit = b'1')as rumahSakit,
	(fp.lift = b'1') as lift,
	(fp.club_house = b'1')as clubHouse,
	(fp.elevator = b'1')as elevator,
	(fp.gym = b'1')as gym,
	(fp.joging_track = b'1')as jogingTrack,
	(fp.garasi = b'1')as garasi,
	(fp.row_jalan12 = b'1')as rowJalan12,
	(fp.cctv = b'1')as cctv,
	(tap.rumah_sakit = b'1')as rumahSakitAP,
	(tap.jalan_tol = b'1')as jalanTol,
	(tap.sekolah = b'1')as sekolah,
	(tap.mall = b'1') as mall,
	(tap.bank_atm = b'1')as bankAtm,
	(tap.pasar = b'1')as pasar,
	(tap.farmasi = b'1')as farmasi,
	(tap.rumah_ibadah = b'1')as rumahIbadah,
	(tap.restoran = b'1')as restoran,
	(tap.taman = b'1')as taman,
	(tap.bioskop = b'1') as bioskop,
	(tap.bar = b'1')as bar,
	(tap.halte = b'1')as halte,
	(tap.stasiun = b'1')as stasiun,
	(tap.bandara = b'1')as bandara,
	(tap.gerbang_tol = b'1')as gerbangTol,
	(tap.spbu = b'1')as spbu,
	(tap.gymnasium = b'1') as gymnasium,
	tspm.*,
	tmp.*,
	tp.*,
	ap.*,
	prop.id as propertiId,
	tp.id as projectId,
	tp.jenis_properti,
	tp.tipe_properti,
	ap.jenis_properti as jenisPropertiAP,
	ap.tipe_properti as tipePropertiAP,
	tp.email as emailProject,
	tp.status as statusProject,
	tmp2.youtube_url as youtubeUrlMediaProject,
	tmp2.virtual360url as virtual360UrlMediaProject,
	tmp.youtube_url as youtubeUrlMediaProperti,
	tmp.virtual360url as virtual360UrlMediaProperti,
	tip2.image_name as imageNameIP,
	tip2.trx_id as trx_id,
	tmp.image_properti_id as imagePropertiId,
	tmp2.image_id as imageProjectId,
	(select case when prop.cluster_id = '-1' then null else tc.id end as x from tbl_cluster tc where tc.id = prop.cluster_id)as clusterId,
	(select case when prop.cluster_id = '-1' then null else tc.project_id end as x from tbl_cluster tc where tc.id = prop.cluster_id)as projectIdTC,
	(select case when prop.cluster_id = '-1' then null else tc.deskripsi end as x from tbl_cluster tc where tc.id = prop.cluster_id)as deskripsiCluster,
	(select case when prop.cluster_id = '-1' then null else tc.stock_units end as x from tbl_cluster tc where tc.id = prop.cluster_id)as stockUnits,
	(select case when prop.cluster_id = '-1' then null when tc.is_cluster = 0 then false else true end as x from tbl_cluster tc where tc.id = prop.cluster_id)as isCluster,
	(select case when prop.cluster_id = '-1' then null else tc.created_at end as x from tbl_cluster tc where tc.id = prop.cluster_id)as createdAtTC,
	(select case when prop.cluster_id = '-1' then null else tc.modified_at end as x from tbl_cluster tc where tc.id = prop.cluster_id)as modifiedAtTC,
	(select case when prop.cluster_id = '-1' then null else tc.status end as x from tbl_cluster tc where tc.id = prop.cluster_id)as statusCluster,
	tkp.created_at as createdAtKP,
	tkp.modified_at as modifiedAtKP,
	tip.created_at as createdAtIP,
	tip.modified_at as modifiedAtIP,
	tspm.created_at as createdAtSPM,
	tspm.modified_at as modifiedAtSPM,
	prop.created_at as createdAtP,
	prop.modified_at as modifiedAtP,
	tmp.created_at as createdAtMP,
	tmp.modified_at as modifiedAtMP,
	ap.created_at as createdAtAP,
	ap.modified_at as modifiedAtAP,
	fp.created_at as createdAtFP,
	fp.modified_at as modifiedAtFP,
	tap.created_at as createdAtAP2,
	tap.modified_at as modifiedAtAP2,
	tmp2.created_at as createdAtMediaProject,
	tmp2.modified_at as modifiedAtMediaProject,
	tp.created_at as createdAtTP,
	tp.modified_at as modifiedAtTP
	FROM tbl_project tp 
	INNER JOIN tbl_alamat_properti ap ON tp.alamat_properti_id = ap.id
	INNER JOIN tbl_fas_properti fp ON tp.fase_properti_id = fp.id
	INNER JOIN tbl_properti prop ON prop.project_id = tp.id
	INNER JOIN tbl_informasi_properti tip on tip.id = prop.info_properti_id
	INNER JOIN tbl_kelengkapan_properti tkp on tkp.id = prop.kelengkapan_properti_id 
	INNER JOIN tbl_selling_properti_method tspm on tspm.id = prop.selling_properti_method_id
	INNER JOIN tbl_media_properti tmp on tmp.id = prop.media_properti_id
	INNER JOIN tbl_akses_properti tap on tap.id = tp.akses_properti_id 
	INNER JOIN tbl_media_project tmp2 on tmp2.id = tp.media_project_id
	INNER JOIN tbl_image_properti tip2 on tip2.id = tmp.id
	WHERE (prop.status != 'sold') AND (prop.status != 'deleted') AND (tp.status = 'published') and length(COALESCE(tmp2.virtual360url,"")) > 10 limit 10 `).Find(&data)

	return data
}
