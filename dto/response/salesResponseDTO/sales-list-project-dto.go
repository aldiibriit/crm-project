package salesResponseDTO

import (
	"database/sql"
	"time"
)

type Response struct {
	HttpCode         int                     `json:"httpCode"`
	ResponseCode     string                  `json:"responseCode"`
	ResponseDesc     string                  `json:"responseDesc"`
	ResponseData     interface{}             `json:"responseData"`
	MetadataResponse MetadataResponse        `json:"metadata"`
	Summary          MetadataSummeryResponse `json:"summary"`
}

type MetadataResponse struct {
	ListUserDtoRes ListUserDtoRes `json:"listUserDtoRes"`
}

type ListUserDtoRes struct {
	Currentpage  int `json:"currentPage"`
	TotalData    int `json:"totalData"`
	TotalDataAll int `json:"totalDataAll"`
}

type MetadataSummeryResponse struct {
	ListTerdekat ListUserDtoRes `json:"listTerdekat"`
	List360      ListUserDtoRes `json:"list360"`
	ListByCity   ListUserDtoRes `json:"listByCity"`
}

type ListProject struct {
	DetailProperti  DetailProperti     `json:"-" gorm:"embedded"`
	Project         Project            `json:"project" gorm:"embedded"`
	Cluster         Cluster            `json:"-" gorm:"embedded"`
	ImagePropertiId string             `json:"-" gorm:"column:imagePropertiId"`
	ImageProjectId  string             `json:"-" gorm:"column:imageProjectId"`
	ImageProperti   []TblImageProperti `json:"-" gorm:"foreignKey:TrxId;references:ImagePropertiId"`
	ImageProject    []TblImageProperti `json:"imageProject" gorm:"foreignKey:TrxId;references:ImageProjectId"`
}

type DetailProperti struct {
	ID                    int                   `json:"-" gorm:"column:propertiId"`
	IDString              string                `json:"id" `
	GroupProperti         string                `json:"groupProperti" gorm:"column:group_properti"`
	Email                 string                `json:"email"`
	NamaProperti          string                `json:"namaProperti" gorm:"column:nama_properti"`
	DeskripsiProperti     string                `json:"deskripsiProperti" gorm:"column:deskripsi_properti"`
	HargaProperti         string                `json:"hargaProperti" gorm:"column:harga_properti"`
	JumlahLantai          string                `json:"jmlLantai" gorm:"column:jml_lantai"`
	JumlahKamarTidur      string                `json:"jmlKamarTidur" gorm:"column:jml_kmr_tidur"`
	JumlahKamarMandi      string                `json:"jmlKamarMandi" gorm:"column:jml_kmr_mandi"`
	ParkirMobilString     string                `json:"-" gorm:"parkirMobilString"`
	ParkirMobilBool       bool                  `json:"parkirMobil" `
	ProjectId             int                   `json:"-" gorm:"column:project_id"`
	ProjectIdString       string                `json:"projectId" `
	ClusterId             int                   `json:"-" gorm:"column:cluster_id"`
	ClusterIdString       string                `json:"clusterId" `
	WishlistCounterInt    int                   `json:"wishlistCounter"`
	WishlistCounterNull   sql.NullString        `json:"-" gorm:"column:wishlist_counter"`
	ViewedCounterInt      int                   `json:"viewedCounter"`
	ViewedCounterNull     sql.NullString        `json:"-" gorm:"column:viewed_vounter"`
	InformasiProperti     InformasiProperti     `json:"informasiPropert" gorm:"embedded"`
	KelengkapanProperti   KelengkapanProperti   `json:"kelengkapanProperti" gorm:"embedded"`
	SellingPropertiMethod SellingPropertiMethod `json:"sellingPropertiMethod" gorm:"embedded"`
	MediaProperti         MediaProperti         `json:"mediaProperti" gorm:"embedded"`
	Status                string                `json:"status"`
	Lt                    int                   `json:"lt"`
	Lb                    int                   `json:"lb"`
	CreatedAt             time.Time             `json:"createdAt" gorm:"column:createdAtP"`
	ModifiedAt            time.Time             `json:"modifiedAt" gorm:"column:modifiedAtP"`
}

type InformasiProperti struct {
	Sertifikat      string    `json:"sertifikat" gorm:"column:sertifikat"`
	JumlahLantai    int       `json:"jmlLantai" gorm:"column:jml_lantai"`
	HadapRumah      string    `json:"hadapRumah" gorm:"column:hadap_rumah"`
	KamarPembantu   int       `json:"kamarPembantu" gorm:"column:kamar_pembantu"`
	KondisiProperti string    `json:"kondisiProperti" gorm:"column:kondisi_properti"`
	DayaListrik     string    `json:"dayaListrik" gorm:"column:daya_listrik"`
	TahunBangun     string    `json:"tahunBangun" gorm:"column:tahun_bangun"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:createdAtIP"`
	ModifiedAt      time.Time `json:"modifiedAt" gorm:"column:modifiedAtIP"`
}

type KelengkapanProperti struct {
	Dapur         bool      `json:"dapur" gorm:"column:dapur"`
	JalurListrik  bool      `json:"jalurListrik" gorm:"column:jalur_listrik"`
	JalurPDAM     bool      `json:"jalurPDAM" gorm:"column:jalur_pdam"`
	JalurTelepone bool      `json:"jalurTelepone" gorm:"column:jalur_telepone"`
	RuangKeluarga bool      `json:"ruangKeluarga" gorm:"column:ruang_keluarga"`
	RuangKerja    bool      `json:"ruangKerja" gorm:"column:ruang_kerja"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:createdAtKP"`
	ModifiedAt    time.Time `json:"modifiedAt" gorm:"column:modifiedAtKP"`
}

type SellingPropertiMethod struct {
	Method     string    `json:"method"`
	Type       string    `json:"type"`
	Duration   int       `json:"duration"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:createdAtSPM"`
	ModifiedAt time.Time `json:"modifiedAt" gorm:"column:modifiedAtSPM"`
}

type MediaProperti struct {
	ClusterId       int       `json:"clusterId" gorm:"column:cluster_id"`
	ImagePropertiId string    `json:"imagePropertiId" gorm:"column:image_properti_id"`
	YoutubeUrl      string    `json:"youtubeUrl" gorm:"column:youtube_url"`
	Virtual360Url   string    `json:"virtual360Url" gorm:"column:virtual360url"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:createdAtMP"`
	ModifiedAt      time.Time `json:"modifiedAt" gorm:"column:modifiedAtMP"`
}

type Project struct {
	Id                int               `json:"-" gorm:"column:projectId"`
	IdString          string            `json:"id" `
	Email             string            `json:"email" gorm:"column:emailProject"`
	NamaProyek        string            `json:"namaProyek" gorm:"column:nama_proyek"`
	KisaranHarga      string            `json:"kisaranHarga" gorm:"column:kisaran_harga"`
	PicProyek         string            `json:"picProyek" gorm:"column:pic_proyek"`
	NoHpPic           string            `json:"noHpPic" gorm:"column:no_hp_pic"`
	BrosurProjectId   string            `json:"brosurProjectId" gorm:"column:brosur_project_id"`
	TipeProperti      string            `json:"tipeProperti" gorm:"column:tipe_properti"`
	JenisProperti     string            `json:"jenisProperti" gorm:"column:jenis_properti"`
	DeskripsiProyek   string            `json:"deskirpsiProyek" gorm:"deskripsi_proyek"`
	StockUnits        int               `json:"stockUnits" gorm:"column:stock_units"`
	AlamatProperti    AlamatProperti    `json:"alamatProperti" gorm:"embedded"`
	FasilitasProperti FasilitasProperti `json:"fasilitasProperti" gorm:"embedded"`
	AksesProperti     AksesProperti     `json:"aksesProperti" gorm:"embedded"`
	MediaProject      MediaProject      `json:"mediaProject" gorm:"embedded"`
	CreatedAt         time.Time         `json:"createdAt" gorm:"column:createdAtTP"`
	ModifiedAt        time.Time         `json:"modifiedAt" gorm:"column:modifiedAtTP"`
	Status            string            `json:"status" gorm:"column:statusProject"`
}

type AlamatProperti struct {
	TipeProperti  string    `json:"tipeProperti" gorm:"column:tipePropertiAP"`
	JenisProperti string    `json:"jenisProperti" gorm:"column:jenisPropertiAP"`
	Alamat        string    `json:"alamat"`
	Provinsi      string    `json:"provinsi"`
	Longitude     float64   `json:"longitude"`
	Latitude      float64   `json:"latitude"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:createdAtAP"`
	ModifiedAt    time.Time `json:"modifiedAt" gorm:"column:modifiedAtAP"`
}

type FasilitasProperti struct {
	KolamRenang   bool      `json:"kolamRenang" gorm:"column:kolamRenang"`
	TempatParkir  bool      `json:"tempatParkir" gorm:"column:tempatParkir"`
	Keamanan24Jam bool      `json:"keamanan24Jam" gorm:"column:keamanan24Jam"`
	Penghijauan   bool      `json:"penghijauan" gorm:"column:penghijauan"`
	RumahSakit    bool      `json:"rumahSakit" gorm:"column:rumahSakit"`
	Lift          bool      `json:"lift" gorm:"column:lift"`
	ClubHouse     bool      `json:"clubHouse" gorm:"column:clubHouse"`
	Elevator      bool      `json:"elevator" gorm:"column:elevator"`
	Gym           bool      `json:"gym" gorm:"column:gym"`
	JogingTrack   bool      `json:"jogingTrack" gorm:"column:jogingTrack"`
	Garasi        bool      `json:"garasi" gorm:"column:garasi"`
	RowJalan12    bool      `json:"rowJalan12" gorm:"column:rowJalan12"`
	Cctv          bool      `json:"cctv" gorm:"column:cctv"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:createdAtFP"`
	ModifiedAtFP  time.Time `json:"ModifiedAtFP" gorm:"column:ModifiedAtFP"`
}

type AksesProperti struct {
	RumahSakit  bool      `json:"rumahSakit" gorm:"column:rumahSakitAP"`
	JalanTol    bool      `json:"jalanTol" gorm:"column:jalanTol"`
	Sekolah     bool      `json:"sekolah" gorm:"column:sekolah"`
	Mall        bool      `json:"mall" gorm:"column:mall"`
	BankAtm     bool      `json:"bankAtm" gorm:"column:bankAtm"`
	Pasar       bool      `json:"pasar" gorm:"column:pasar"`
	Farmasi     bool      `json:"farmasi" gorm:"column:farmasi"`
	RumahIbadah bool      `json:"rumahIbadah" gorm:"column:rumahIbadah"`
	Restoran    bool      `json:"restoran" gorm:"column:restoran"`
	Taman       bool      `json:"taman" gorm:"column:taman"`
	Bioskop     bool      `json:"bioskop" gorm:"column:bioskop"`
	Bar         bool      `json:"bar" gorm:"column:bar"`
	Halte       bool      `json:"halte" gorm:"column:halte"`
	Stasiun     bool      `json:"stasiun" gorm:"column:stasiun"`
	Bandara     bool      `json:"bandara" gorm:"column:bandara"`
	GerbangTol  bool      `json:"gerbangTol" gorm:"column:gerbangTol"`
	SPBU        bool      `json:"spbu" gorm:"column:spbu"`
	Gymnasium   bool      `json:"gymnasium" gorm:"column:gymnasium"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:createdAtAP2"`
	ModifiedAt  time.Time `json:"modifiedAt" gorm:"column:modifiedAtAP2"`
}

type MediaProject struct {
	YoutubeUrl    string    `json:"youtubeUrl" gorm:"column:youtubeUrlMediaProject"`
	Virtual360Url string    `json:"virtual360Url" gorm:"virtual360UrlMediaProject"`
	CreatedAt     time.Time `json:"createdAt" gorm:"createdAtMediaProject"`
	ModifiedAt    time.Time `json:"modifiedAt" gorm:"modifiedAtMediaProject"`
}

type Cluster struct {
	Id              int       `json:"-" gorm:"column:clusterId"`
	IdString        string    `json:"id"`
	ProjectId       int       `json:"-" gorm:"column:projectIdTc"`
	ProjectIdString string    `json:"projectId" `
	Name            string    `json:"name" gorm:"column:deskripsiCluster"`
	Deskripsi       string    `json:"deskripsi" gorm:"column:deskripsiCluster"`
	StockUnits      string    `json:"stockUnits" gorm:"column:stock_units"`
	IsCluster       bool      `json:"isCluster" gorm:"column:isCluster"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:createdAtTC"`
	ModifiedAt      time.Time `json:"modifiedAt" gorm:"column:modifiedAtTC"`
	Status          string    `json:"status" gorm:"column:statusCluster"`
}

type TblImageProperti struct {
	ImageName  string    `json:"imageName" gorm:"column:image_name"`
	TrxId      string    `json:"trxId" gorm:"column:trx_id"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at"`
	ModifiedAt time.Time `json:"modifiedAt" gorm:"column:modified_at"`
}
