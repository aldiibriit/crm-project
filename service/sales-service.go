package service

import (
	"encoding/base64"
	"fmt"
	"go-api/dto/request/salesRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/dto/response/salesResponseDTO"
	"go-api/helper"
	"go-api/repository"
	"strconv"
)

type SalesService interface {
	MISDeveloper(request salesRequestDTO.MISDeveloperRequestDTO) responseDTO.Response
	MISSuperAdmin(request salesRequestDTO.MISSuperAdminRequestDTO) responseDTO.Response
	ListProject(request salesRequestDTO.ListProjectRequest) responseDTO.Response
	EditSalesByDeveloper(request salesRequestDTO.SalesEditRequestDTO) responseDTO.Response
	DeleteSalesByDeveloper(request salesRequestDTO.SalesDeleteRequestDTO) responseDTO.Response
	DetailSalesByDeveloper(request salesRequestDTO.DetailSalesRequest) responseDTO.Response
	DraftDetail(request salesRequestDTO.DraftDetailRequest) responseDTO.Response
	DeletePengajuan(request salesRequestDTO.SalesDeleteRequestDTO) responseDTO.Response
	ListFinalPengajuan(request salesRequestDTO.FinalPengajuanRequest) responseDTO.Response
}

type salesService struct {
	salesRepository repository.SalesRepository
	userRepository  repository.UserRepository
	kprRepository   repository.KPRRepository
}

func NewSalesService(salesRepo repository.SalesRepository, userRepo repository.UserRepository, kprRepo repository.KPRRepository) SalesService {
	return &salesService{
		salesRepository: salesRepo,
		userRepository:  userRepo,
		kprRepository:   kprRepo,
	}
}

func (service *salesService) MISDeveloper(request salesRequestDTO.MISDeveloperRequestDTO) responseDTO.Response {
	var response responseDTO.Response
	var metadataResponse responseDTO.ListUserDtoRes

	var sqlStr1, sqlStrCountAll1 string

	if request.StartDate != "" && request.EndDate != "" {
		sqlStr1 = `SELECT 
	tu.id,ts.developer_email,ts.sales_email,ts.refferal_code,ts.registered_by,ts.created_at,ts.modified_at,ts.sales_name,tu.mobile_no as salesPhone,tu.status
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = '` + request.EmailDeveloper + `' and ts.created_at between '` + request.StartDate + `' and '` + request.EndDate + `'
	limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
	`

		sqlStrCountAll1 = `SELECT 
	count(tu.id)
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = '` + request.EmailDeveloper + `' and ts.created_at between '` + request.StartDate + `' and '` + request.EndDate + `'
	`
	} else if request.EndDate == "" && request.StartDate != "" {
		sqlStr1 = `SELECT 
	tu.id,ts.developer_email,ts.sales_email,ts.refferal_code,ts.registered_by,ts.created_at,ts.modified_at,ts.sales_name,tu.mobile_no as salesPhone,tu.status
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = '` + request.EmailDeveloper + `' and date(ts.created_at) >= '` + request.StartDate + `'
	limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
	`

		sqlStrCountAll1 = `SELECT 
	count(tu.id)
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = '` + request.EmailDeveloper + `' and date(ts.created_at) >= '` + request.StartDate + `'
	`
	} else {
		sqlStr1 = `SELECT 
	tu.id,ts.developer_email,ts.sales_email,ts.refferal_code,ts.registered_by,ts.created_at,ts.modified_at,ts.sales_name,tu.mobile_no as salesPhone,tu.status
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = '` + request.EmailDeveloper + `' and ts.developer_email like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.refferal_code like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.registered_by like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.created_at like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.modified_at like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.sales_name like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and tu.mobile_no like '%` + request.Keyword + `%'
	limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
	`
		sqlStrCountAll1 = `SELECT 
	count(tu.id)
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = '` + request.EmailDeveloper + `' and ts.developer_email like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.refferal_code like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.registered_by like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.created_at like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.modified_at like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.sales_name like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and tu.mobile_no like '%` + request.Keyword + `%'
	`
	}

	metadataResponse.Currentpage = request.Offset
	if request.Offset > 0 {
		request.Offset = request.Offset * request.Limit
	}
	data, totalData := service.salesRepository.FindByEmailDeveloper(sqlStr1, sqlStrCountAll1, request)
	metadataResponse.TotalData = totalData

	encryptedData := serializeMisDeveloper(data)
	response.HttpCode = 200
	response.MetadataResponse = metadataResponse
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = encryptedData
	response.Summary = nil

	return response
}

func (service *salesService) ListProject(request salesRequestDTO.ListProjectRequest) responseDTO.Response {
	var response responseDTO.Response
	var metadata salesResponseDTO.MetadataResponse

	var offset int
	if request.PageStart > 0 {
		offset = request.PageStart * 10
	}

	sqlStr := `SELECT 
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
	tp.modified_at as modifiedAtTP,
	(select id from tbl_sales where refferal_code = '` + request.ReferralCode + `' or sales_email = '` + request.EmailSales + `' limit 1)as salesID,
	(select refferal_code from tbl_sales where refferal_code = '` + request.ReferralCode + `' or sales_email = '` + request.EmailSales + `' limit 1)as referralCode
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
	WHERE (prop.status != 'sold') AND (prop.status != 'deleted') AND (tp.status = 'published') AND tp.email = (select developer_email from tbl_sales where refferal_code = '` + request.ReferralCode + `' or sales_email = '` + request.EmailSales + `' limit 1) limit 10 offset ` + strconv.Itoa(offset) + ``

	data, totalDataAll := service.salesRepository.ListProject(sqlStr)

	metadata.ListUserDtoRes.Currentpage = request.PageStart
	metadata.ListUserDtoRes.TotalData = len(data)
	metadata.ListUserDtoRes.TotalDataAll = int(totalDataAll)

	for i, v := range data {
		chiperImageProjectId, _ := base64.StdEncoding.DecodeString(v.ImageProjectId)
		originText, err := helper.RsaDecryptFromFEInBEJava(chiperImageProjectId)
		if err != nil {
			fmt.Println(err.Error())
		}
		data[i].ImageProjectId = originText
		listImageProject := service.salesRepository.RelationToImageProperti(originText)
		data[i].ImageProject = listImageProject
	}

	for i, v := range data {
		originIdProject := v.Project.Id
		originEmailProject := v.Project.Email
		originNoHpPic := v.Project.NoHpPic
		originIdCluster := v.Cluster.Id
		originProjectIdCluster := v.Cluster.ProjectId
		originBrosurProjectId := v.Project.BrosurProjectId

		encryptedIdProject, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originIdProject)))
		encryptedEmailProject, _ := helper.RsaEncryptBEToFE([]byte(originEmailProject))
		encryptedNoHpPic, _ := helper.RsaEncryptBEToFE([]byte(originNoHpPic))
		encryptedIdCluster, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originIdCluster)))
		encryptedProjectIdCluster, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originProjectIdCluster)))
		encryptedBrosurProjectId, _ := helper.RsaEncryptBEToFE([]byte(originBrosurProjectId))

		data[i].Project.IdString = encryptedIdProject
		data[i].Project.Email = encryptedEmailProject
		data[i].Project.NoHpPic = encryptedNoHpPic
		data[i].Cluster.IdString = encryptedIdCluster
		data[i].Cluster.ProjectIdString = encryptedProjectIdCluster
		data[i].Project.BrosurProjectId = encryptedBrosurProjectId

		for x := 0; x < len(v.ImageProject); x++ {
			originProjectTrxId := v.ImageProject[x].TrxId
			originImageName := v.ImageProject[x].ImageName
			encryptedProjectTrxId, _ := helper.RsaEncryptBEToFE([]byte(originProjectTrxId))
			encryptedImageName, _ := helper.RsaEncryptBEToFE([]byte(originImageName))
			data[i].ImageProject[x].TrxId = encryptedProjectTrxId
			data[i].ImageProject[x].ImageName = encryptedImageName
		}

		for y := 0; y < len(v.ImageProperti); y++ {
			originPropertiTrxId := v.ImageProperti[y].TrxId
			encryptedPropertiTrxId, _ := helper.RsaEncryptBEToFE([]byte(originPropertiTrxId))
			data[i].ImageProperti[y].TrxId = encryptedPropertiTrxId
		}
	}

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseDesc = "Success"
	response.ResponseData = data
	response.MetadataResponse = metadata

	return response
}

func (service *salesService) MISSuperAdmin(request salesRequestDTO.MISSuperAdminRequestDTO) responseDTO.Response {
	var response responseDTO.Response
	var metadataResponse responseDTO.ListUserDtoRes

	metadataResponse.Currentpage = request.Offset
	if request.Offset > 0 {
		request.Offset = request.Offset * request.Limit
	}

	var sqlStr1, sqlStrCountAll1 string

	if request.StartDate != "" && request.EndDate != "" {
		sqlStr1 = `SELECT ts.sales_email ,ts.sales_name,(select json_extract(metadata,'$.name') from tbl_user tu2 where tu2.email = ts.developer_email)as metadata,tp.status,tp.jenis_properti,tp.tipe_properti,ts.created_at,ts.modified_at
		FROM tbl_project tp
		JOIN tbl_sales ts on tp.email = ts.developer_email 
		JOIN tbl_user tu on tu.email = ts.sales_email 
		WHERE ts.created_at between '` + request.StartDate + `' and '` + request.EndDate + `'
		order by ts.sales_email limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
		`

		sqlStrCountAll1 = `SELECT count(ts.sales_email)
		FROM tbl_project tp
		JOIN tbl_sales ts on tp.email = ts.developer_email 
		JOIN tbl_user tu on tu.email = ts.sales_email 
		WHERE ts.created_at between '` + request.StartDate + `' and '` + request.EndDate + `'
		order by ts.sales_email
		`
	} else if request.EndDate == "" && request.StartDate != "" {
		sqlStr1 = `SELECT ts.sales_email ,ts.sales_name,(select json_extract(metadata,'$.name') from tbl_user tu2 where tu2.email = ts.developer_email)as metadata,tp.status,tp.jenis_properti,tp.tipe_properti,ts.created_at,ts.modified_at
		FROM tbl_project tp
		JOIN tbl_sales ts on tp.email = ts.developer_email 
		JOIN tbl_user tu on tu.email = ts.sales_email 
		WHERE date(ts.created_at) >= '` + request.StartDate + `'
		order by ts.sales_email limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
		`

		sqlStrCountAll1 = `SELECT count(ts.sales_email)
		FROM tbl_project tp
		JOIN tbl_sales ts on tp.email = ts.developer_email 
		JOIN tbl_user tu on tu.email = ts.sales_email 
		WHERE date(ts.created_at) >= '` + request.StartDate + `'
		order by ts.sales_email
		`
	} else {
		sqlStr1 = `SELECT ts.sales_email ,ts.sales_name,(select json_extract(metadata,'$.name') from tbl_user tu2 where tu2.email = ts.developer_email)as metadata,tp.status,tp.jenis_properti,tp.tipe_properti,ts.created_at,ts.modified_at
		FROM tbl_project tp
		JOIN tbl_sales ts on tp.email = ts.developer_email 
		JOIN tbl_user tu on tu.email = ts.sales_email 
		WHERE ts.sales_email like '%` + request.Keyword + `%' or ts.sales_name like '%` + request.Keyword + `%' or metadata like '%` + request.Keyword + `%' or tp.status like '%` + request.Keyword + `%' or tp.jenis_properti like '%` + request.Keyword + `%' or tp.tipe_properti like '%` + request.Keyword + `%'
		order by ts.sales_email limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
		`

		sqlStrCountAll1 = `SELECT count(ts.sales_email)
		FROM tbl_project tp
		JOIN tbl_sales ts on tp.email = ts.developer_email 
		JOIN tbl_user tu on tu.email = ts.sales_email 
		WHERE ts.sales_email like '%` + request.Keyword + `%' or ts.sales_name like '%` + request.Keyword + `%' or metadata like '%` + request.Keyword + `%' or tp.status like '%` + request.Keyword + `%' or tp.jenis_properti like '%` + request.Keyword + `%' or tp.tipe_properti like '%` + request.Keyword + `%'
		order by ts.sales_email 
		`
	}

	data, totalData := service.salesRepository.MISSuperAdmin(sqlStr1, sqlStrCountAll1, request)
	metadataResponse.TotalData = totalData

	encryptedData := serializeMisSuperAdmin(data)

	response.HttpCode = 200
	response.MetadataResponse = metadataResponse
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = encryptedData
	response.Summary = nil

	return response
}

func (service *salesService) EditSalesByDeveloper(request salesRequestDTO.SalesEditRequestDTO) responseDTO.Response {
	var response responseDTO.Response

	user := service.userRepository.FindByEmail2(request.Email)

	if len(user.Email) > 0 {
		response.HttpCode = 500
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = "Email already exist"
		response.Summary = nil
		return response
	}

	err := service.salesRepository.EditSalesByDeveloper(request)
	if err != nil {
		response.HttpCode = 500
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = err.Error()
		response.Summary = nil
		return response
	}

	encryptedData := serializeUpdatedSales(request)

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = encryptedData
	response.ResponseDesc = "Success"
	response.Summary = nil

	return response
}

func (service *salesService) DetailSalesByDeveloper(request salesRequestDTO.DetailSalesRequest) responseDTO.Response {
	var response responseDTO.Response

	data := service.salesRepository.DetailSalesByDeveloper(request)
	if data.EmailDeveloper == "" && data.EmailSales == "" {
		response.HttpCode = 404
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = "Data not found"
		response.Summary = nil
		return response
	}

	encryptedData := serializeDetailSales(data)

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = encryptedData
	response.ResponseDesc = "Success"
	response.Summary = nil

	return response
}

func (service *salesService) DeleteSalesByDeveloper(request salesRequestDTO.SalesDeleteRequestDTO) responseDTO.Response {
	var response responseDTO.Response

	err := service.salesRepository.DeleteSalesByDeveloper(request)
	if err != nil {
		response.HttpCode = 500
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = err.Error()
		response.Summary = nil
		return response
	}

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = nil
	response.ResponseDesc = "Success"
	response.Summary = nil

	return response
}

func (service *salesService) DraftDetail(request salesRequestDTO.DraftDetailRequest) responseDTO.Response {
	var response responseDTO.Response
	data := service.salesRepository.DraftDetail(request)
	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = data
	response.ResponseDesc = "Success"
	response.Summary = nil
	return response
}

func (service *salesService) DeletePengajuan(request salesRequestDTO.SalesDeleteRequestDTO) responseDTO.Response {
	var response responseDTO.Response

	data := service.kprRepository.Delete(request.ID)

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = data
	response.ResponseDesc = "Success"
	response.Summary = nil

	return response
}

func (service *salesService) ListFinalPengajuan(request salesRequestDTO.FinalPengajuanRequest) responseDTO.Response {
	var response responseDTO.Response
	var metadataResponse responseDTO.ListUserDtoRes
	var sqlStr, sqlStrCount string

	metadataResponse.Currentpage = request.Offset
	if request.Offset > 0 {
		request.Offset = request.Offset * request.Limit
	}

	if request.UserType == "sales" {
		sqlStr = `
		SELECT sales_name,json_extract(metadata,'$.name')as developerName,tc.name as customerName,tpkbs.created_at,jenis_properti,tipe_properti from tbl_pengajuan_kpr_by_sales tpkbs
		join tbl_customer tc on tc.id = tpkbs.customer_id 
		join tbl_sales ts on ts.id = tc.sales_id
		join tbl_project tp on tp.email = ts.developer_email
		join tbl_user tu on tu.email = ts.developer_email 
		where ts.sales_email = '` + request.Email + `' limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
		`
		sqlStrCount = `
		SELECT count(sales_name) from tbl_pengajuan_kpr_by_sales tpkbs
		join tbl_customer tc on tc.id = tpkbs.customer_id 
		join tbl_sales ts on ts.id = tc.sales_id
		join tbl_project tp on tp.email = ts.developer_email
		join tbl_user tu on tu.email = ts.developer_email 
		where ts.sales_email = '` + request.Email + `'
		`
	} else if request.UserType == "developer" {
		sqlStr = `
		SELECT sales_name,json_extract(metadata,'$.name')as developerName,tc.name as customerName,tpkbs.created_at,jenis_properti,tipe_properti from tbl_pengajuan_kpr_by_sales tpkbs
		join tbl_customer tc on tc.id = tpkbs.customer_id 
		join tbl_sales ts on ts.id = tc.sales_id
		join tbl_project tp on tp.email = ts.developer_email
		join tbl_user tu on tu.email = ts.developer_email 
		where ts.developer_email = '` + request.Email + `' limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
		`
		sqlStrCount = `
		SELECT count(sales_name) from tbl_pengajuan_kpr_by_sales tpkbs
		join tbl_customer tc on tc.id = tpkbs.customer_id 
		join tbl_sales ts on ts.id = tc.sales_id
		join tbl_project tp on tp.email = ts.developer_email
		join tbl_user tu on tu.email = ts.developer_email 
		where ts.developer_email = '` + request.Email + `'
		`
	} else {
		sqlStr = `
		SELECT sales_name,json_extract(metadata,'$.name')as developerName,tc.name as customerName,tpkbs.created_at,jenis_properti,tipe_properti from tbl_pengajuan_kpr_by_sales tpkbs
		join tbl_customer tc on tc.id = tpkbs.customer_id 
		join tbl_sales ts on ts.id = tc.sales_id
		join tbl_project tp on tp.email = ts.developer_email
		join tbl_user tu on tu.email = ts.developer_email 
		limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
		`
		sqlStrCount = `
		SELECT count(sales_name) from tbl_pengajuan_kpr_by_sales tpkbs
		join tbl_customer tc on tc.id = tpkbs.customer_id 
		join tbl_sales ts on ts.id = tc.sales_id
		join tbl_project tp on tp.email = ts.developer_email
		join tbl_user tu on tu.email = ts.developer_email 
		`
	}

	data, totalData := service.salesRepository.ListFinalPengajuan(sqlStr, sqlStrCount)
	metadataResponse.TotalData = int(totalData)

	response.HttpCode = 200
	response.MetadataResponse = metadataResponse
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = data
	response.Summary = nil

	return response
}

func serializeMisDeveloper(request interface{}) []salesResponseDTO.MISDeveloper {
	data := request.([]salesResponseDTO.MISDeveloper)
	result := make([]salesResponseDTO.MISDeveloper, len(data))
	for i, v := range data {
		encryptedIdRes, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(v.ID)))
		encryptedEmaiDeveloper, _ := helper.RsaEncryptBEToFE([]byte(v.EmailDeveloper))
		encryptedEmaiSales, _ := helper.RsaEncryptBEToFE([]byte(v.EmailSales))
		encryptedRefferalCode, _ := helper.RsaEncryptBEToFE([]byte(v.RefferalCode))
		encryptedRegisteredBy, _ := helper.RsaEncryptBEToFE([]byte(v.RegisteredBy))
		encryptedSalesName, _ := helper.RsaEncryptBEToFE([]byte(v.SalesName))
		encryptedSalesPhone, _ := helper.RsaEncryptBEToFE([]byte(v.SalesPhone))
		// encryptedCreatedAt, _ := helper.RsaEncryptBEToFE([]byte(v.CreatedAt.String()))
		// encryptedModifiedAt, _ := helper.RsaEncryptBEToFE([]byte(v.ModifiedAt.String()))

		result[i].IDResponse = encryptedIdRes
		result[i].EmailSales = encryptedEmaiSales
		result[i].EmailDeveloper = encryptedEmaiDeveloper
		result[i].RefferalCode = encryptedRefferalCode
		result[i].RegisteredBy = encryptedRegisteredBy
		result[i].CreatedAtRes = v.CreatedAt.String()
		result[i].ModifiedAtRes = v.ModifiedAt.String()
		result[i].SalesName = encryptedSalesName
		result[i].SalesPhone = encryptedSalesPhone
		result[i].Status = v.Status
	}

	return result
}

func serializeUpdatedSales(request interface{}) salesRequestDTO.SalesEditRequestDTO {
	data := request.(salesRequestDTO.SalesEditRequestDTO)
	var result salesRequestDTO.SalesEditRequestDTO

	encryptedEmail, _ := helper.RsaEncryptBEToFE([]byte(data.Email))
	encryptedID, _ := helper.RsaEncryptBEToFE([]byte(data.ID))
	encryptedSalesPhone, _ := helper.RsaEncryptBEToFE([]byte(data.SalesPhone))
	encryptedSalesName, _ := helper.RsaEncryptBEToFE([]byte(data.SalesName))

	result.Email = encryptedEmail
	result.ID = encryptedID
	result.SalesName = encryptedSalesName
	result.SalesPhone = encryptedSalesPhone

	return result
}

func serializeDetailSales(request interface{}) salesResponseDTO.MISDeveloper {
	data := request.(salesResponseDTO.MISDeveloper)
	var result salesResponseDTO.MISDeveloper

	encryptedIdRes, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(data.ID)))
	encryptedEmaiDeveloper, _ := helper.RsaEncryptBEToFE([]byte(data.EmailDeveloper))
	encryptedEmaiSales, _ := helper.RsaEncryptBEToFE([]byte(data.EmailSales))
	encryptedRefferalCode, _ := helper.RsaEncryptBEToFE([]byte(data.RefferalCode))
	encryptedRegisteredBy, _ := helper.RsaEncryptBEToFE([]byte(data.RegisteredBy))
	encryptedSalesName, _ := helper.RsaEncryptBEToFE([]byte(data.SalesName))
	encryptedSalesPhone, _ := helper.RsaEncryptBEToFE([]byte(data.SalesPhone))
	encryptedCreatedAt, _ := helper.RsaEncryptBEToFE([]byte(data.CreatedAt.String()))
	encryptedModifiedAt, _ := helper.RsaEncryptBEToFE([]byte(data.ModifiedAt.String()))

	result.IDResponse = encryptedIdRes
	result.EmailSales = encryptedEmaiSales
	result.EmailDeveloper = encryptedEmaiDeveloper
	result.RefferalCode = encryptedRefferalCode
	result.RegisteredBy = encryptedRegisteredBy
	result.CreatedAtRes = encryptedCreatedAt
	result.ModifiedAtRes = encryptedModifiedAt
	result.SalesName = encryptedSalesName
	result.SalesPhone = encryptedSalesPhone
	result.Status = data.Status

	return result
}

func serializeMisSuperAdmin(request interface{}) []salesResponseDTO.MISSuperAdmin {
	data := request.([]salesResponseDTO.MISSuperAdmin)
	result := make([]salesResponseDTO.MISSuperAdmin, len(data))
	for i, v := range data {
		encryptedSalesName, _ := helper.RsaEncryptBEToFE([]byte(v.SalesName))
		encryptedMetadata, _ := helper.RsaEncryptBEToFE([]byte(v.Metadata))
		encryptedStatus, _ := helper.RsaEncryptBEToFE([]byte(v.Status))
		encryptedJenisProperti, _ := helper.RsaEncryptBEToFE([]byte(v.JenisProperti))
		encryptedTipeProperti, _ := helper.RsaEncryptBEToFE([]byte(v.TipeProperti))

		result[i].SalesName = encryptedSalesName
		result[i].Metadata = encryptedMetadata
		result[i].Status = encryptedStatus
		result[i].JenisProperti = encryptedJenisProperti
		result[i].TipeProperti = encryptedTipeProperti
		result[i].CreatedAtRes = v.CreatedAt.String()
		result[i].ModifiedAtRes = v.ModifiedAt.String()
	}

	return result
}
