package service

import (
	"encoding/base64"
	"fmt"
	propertiDTO "go-api/dto/properti"
	"go-api/helper"
	"go-api/repository"
	"strconv"
	"strings"
)

type PropertiService interface {
	AdvancedFilter(b propertiDTO.AdvancedFilterDTO, bearerToken string) propertiDTO.Response
}

type propertiService struct {
	baseService        BaseService
	propertiRepository repository.PropertiRepository
	userRepository     repository.UserRepository
}

// NewBookService .....
func NewPropertiService(propertiRepo repository.PropertiRepository, userRepo repository.UserRepository, baseServ BaseService) PropertiService {
	return &propertiService{
		propertiRepository: propertiRepo,
		userRepository:     userRepo,
		baseService:        baseServ,
	}
}

func (service *propertiService) AdvancedFilter(request propertiDTO.AdvancedFilterDTO, bearerToken string) propertiDTO.Response {
	var response propertiDTO.Response
	var listProperti []propertiDTO.AdvancedFilter
	cipherTextEmail, _ := base64.StdEncoding.DecodeString(request.Email)
	decryptedEmail, _ := helper.RsaDecryptFromFEInBE(cipherTextEmail)
	request.Email = decryptedEmail
	isVerified, errString := service.baseService.ValidateToken(bearerToken, request.Email)
	if !isVerified && errString != "" {
		listProperti = make([]propertiDTO.AdvancedFilter, 0)
		response.HttpCode = 500
		response.ResponseCode = "500"
		response.ResponseDesc = errString
		response.ResponseData = make([]interface{}, 0)
		return response
	} else {

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
	WHERE (prop.status != 'sold') AND (prop.status != 'deleted') AND (tp.status = 'published')`

		if request.HargaMin != 0 {
			strHargaMin := strconv.Itoa(request.HargaMin)
			sqlStr += ` AND CAST(prop.harga_properti AS DECIMAL(20,2)) >= CAST(` + strHargaMin + ` AS DECIMAL(20,2))`
		}

		if request.HargaMax != 0 {
			strHargaMax := strconv.Itoa(request.HargaMax)
			sqlStr += ` AND CAST(prop.harga_properti AS DECIMAL(20,2)) <= CAST(` + strHargaMax + ` AS DECIMAL(20,2))`
		}

		if request.TipeProperti != "" {
			sqlStr += ` AND tp.tipe_properti = ` + request.TipeProperti + ``
		}

		if len(request.JumlahKamarTidur) != 0 {
			strCondition := strings.Join(request.JumlahKamarTidur, ",")
			sqlStr += ` AND prop.jml_kmr_tidur in (` + strCondition + `)`
		}

		if len(request.JumlahKamarMandi) != 0 {
			strCondition := strings.Join(request.JumlahKamarMandi, ",")
			sqlStr += ` AND prop.jml_kmr_mandi in (` + strCondition + `)`
		}

		if request.LuasTanahMinimal != 0 {
			sqlStr += ` AND prop.lt > ` + strconv.Itoa(request.LuasTanahMinimal) + ``
		}

		if request.LuasTanahMaksimal != 0 {
			sqlStr += ` AND prop.lt < ` + strconv.Itoa(request.LuasTanahMaksimal) + ``
		}

		if request.LuasBangunanMinimal != 0 {
			sqlStr += ` AND prop.lt > ` + strconv.Itoa(request.LuasBangunanMinimal) + ``
		}

		if request.LuasBangunanMaksimal != 0 {
			sqlStr += ` AND prop.lt < ` + strconv.Itoa(request.LuasBangunanMaksimal) + ``
		}

		if request.JenisProperti != "" {
			sqlStr += ` AND tp.jenis_properti = '` + request.JenisProperti + `'`
		}

		if request.Dapur {
			sqlStr += ` AND (tkp.dapur = b'1') = true`
		}

		if request.JalurTelepon {
			sqlStr += ` AND (tkp.jalur_telepone = b'1') = true`
		}

		if request.JalurListrik {
			sqlStr += ` AND (tkp.jalur_listrik = b'1') = true`
		}

		if request.JalurPDAM {
			sqlStr += ` AND (tkp.jalurpdam = b'1') = true`
		}

		if request.RuangKeluarga {
			sqlStr += ` AND (tkp.ruang_keluarga = b'1') = true`
		}

		if request.RuangKerja {
			sqlStr += ` AND (tkp.ruang_kerja = b'1') = true`
		}

		if request.RumahSakit {
			sqlStr += ` AND (tap.rumah_sakit = b'1') = true`
		}

		if request.JalanTol {
			sqlStr += ` AND (tap.jalan_tol = b'1') = true`
		}

		if request.Sekolah {
			sqlStr += ` AND (tap.sekolah = b'1') = true`
		}

		if request.Mall {
			sqlStr += ` AND (tap.mall = b'1') = true`
		}

		if request.BankATM {
			sqlStr += ` AND (tap.bank_atm = b'1') = true`
		}

		if request.Taman {
			sqlStr += ` AND (tap.taman = b'1') = true`
		}

		if request.Pasar {
			sqlStr += ` AND (tap.pasar = b'1') = true`
		}

		if request.Farmasi {
			sqlStr += ` AND (tap.farmasi = b'1') = true`
		}

		if request.RumahIbadah {
			sqlStr += ` AND (tap.rumah_ibadah = b'1') = true`
		}

		if request.Restoran {
			sqlStr += ` AND (tap.restoran = b'1') = true`
		}

		if request.Bioskop {
			sqlStr += ` AND (tap.bioskop = b'1') = true`
		}

		if request.Bar {
			sqlStr += ` AND (tap.bar = b'1') = true`
		}

		if request.Halte {
			sqlStr += ` AND (tap.halte = b'1') = true`
		}

		if request.Stasiun {
			sqlStr += ` AND (tap.stasiun = b'1') = true`
		}

		if request.Bandara {
			sqlStr += ` AND (tap.bandara = b'1') = true`
		}

		if request.GerbangTol {
			sqlStr += ` AND (tap.gerbang_tol = b'1') = true`
		}

		if request.SPBU {
			sqlStr += ` AND (tap.spbu = b'1') = true`
		}

		if request.Gymnasium {
			sqlStr += ` AND (tap.gymnasium = b'1') = true`
		}

		if request.KolamRenang {
			sqlStr += ` AND (fp.kolam_renang = b'1') = true`
		}

		if request.TempatParkir {
			sqlStr += ` AND (fp.tempat_parkir = b'1') = true`
		}

		if request.Keamanan {
			sqlStr += ` AND (fp.keamanan = b'1') = true`
		}

		if request.Penghijauan {
			sqlStr += ` AND (fp.penghijauan = b'1') = true`
		}

		if request.Lift {
			sqlStr += ` AND (fp.lift = b'1') = true`
		}

		if request.ClubHouse {
			sqlStr += ` AND (fp.club_house = b'1') = true`
		}

		if request.Elevator {
			sqlStr += ` AND (fp.elevator = b'1') = true`
		}

		if request.Gym {
			sqlStr += ` AND (fp.gym = b'1') = true`
		}

		if request.Garasi {
			sqlStr += ` AND (fp.garasi = b'1') = true`
		}

		if request.RowJalan12 {
			sqlStr += ` AND (fp.row_jalan12 = b'1') = true`
		}

		sqlStr += ` limit 10`
		data := service.propertiRepository.AdvancedFilter(sqlStr)

		listProperti = data.([]propertiDTO.AdvancedFilter)
		for i, v := range listProperti {
			chiperImageProjectId, _ := base64.StdEncoding.DecodeString(v.ImageProjectId)
			originText, err := helper.RsaDecryptFromFEInBEJava(chiperImageProjectId)
			if err != nil {
				fmt.Println(err.Error())
			}
			listProperti[i].ImageProjectId = originText
			listImageProject := service.propertiRepository.AdvancedFilter2(originText)
			listProperti[i].ImageProject = listImageProject
		}

		for i, v := range listProperti {
			originIdProperti := v.DetailProperti.ID
			originGroupProperti := v.DetailProperti.GroupProperti
			originEmail := v.DetailProperti.Email
			originClusterId := v.DetailProperti.ClusterId
			originProjectId := v.DetailProperti.ProjectId
			originIdProject := v.Project.Id
			originEmailProject := v.Project.Email
			originNoHpPic := v.Project.NoHpPic
			originIdCluster := v.Cluster.Id
			originProjectIdCluster := v.Cluster.ProjectId

			encryptedIdProperti, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originIdProperti)))
			encryptedGroupProperti, _ := helper.RsaEncryptBEToFE([]byte(originGroupProperti))
			encryptedEmail, _ := helper.RsaEncryptBEToFE([]byte(originEmail))
			encryptedClusterId, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originClusterId)))
			encryptedProjectId, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originProjectId)))
			encryptedIdProject, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originIdProject)))
			encryptedEmailProject, _ := helper.RsaEncryptBEToFE([]byte(originEmailProject))
			encryptedNoHpPic, _ := helper.RsaEncryptBEToFE([]byte(originNoHpPic))
			encryptedIdCluster, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originIdCluster)))
			encryptedProjectIdCluster, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(originProjectIdCluster)))

			listProperti[i].DetailProperti.IDString = encryptedIdProperti
			listProperti[i].DetailProperti.GroupProperti = encryptedGroupProperti
			listProperti[i].DetailProperti.Email = encryptedEmail
			listProperti[i].DetailProperti.ClusterIdString = encryptedClusterId
			listProperti[i].DetailProperti.ProjectIdString = encryptedProjectId
			listProperti[i].Project.IdString = encryptedIdProject
			listProperti[i].Project.Email = encryptedEmailProject
			listProperti[i].Project.NoHpPic = encryptedNoHpPic
			listProperti[i].Cluster.IdString = encryptedIdCluster
			listProperti[i].Cluster.ProjectIdString = encryptedProjectIdCluster

			for x := 0; x < len(v.ImageProject); x++ {
				originProjectTrxId := v.ImageProject[x].TrxId
				encryptedProjectTrxId, _ := helper.RsaEncryptBEToFE([]byte(originProjectTrxId))
				listProperti[i].ImageProject[x].TrxId = encryptedProjectTrxId
			}

			for y := 0; y < len(v.ImageProperti); y++ {
				originPropertiTrxId := v.ImageProperti[y].TrxId
				encryptedPropertiTrxId, _ := helper.RsaEncryptBEToFE([]byte(originPropertiTrxId))
				listProperti[i].ImageProperti[y].TrxId = encryptedPropertiTrxId
			}
		}

		response.ResponseData = listProperti
	}

	return response
}
