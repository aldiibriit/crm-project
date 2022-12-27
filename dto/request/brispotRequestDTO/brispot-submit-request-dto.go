package brispotRequestDTO

type BrispotSubmitRequestDTO struct {
	// RequestMethod string      `json:"requestMethod"`
	// RequestUser   string      `json:"requestUser"`
	// RequestData   RequestData `json:"requestData"`
	Branch           string `json:"branch"`
	Nik              string `json:"nik"`
	Nama             string `json:"nama"`
	JenisKelamin     string `json:"jenisKelamin"`
	Alamat           string `json:"alamat"`
	Rt               string `json:"rt"`
	Rw               string `json:"rw"`
	Provinsi         string `json:"provinsi"`
	Kota             string `json:"kota"`
	Kecamatan        string `json:"kecamatan"`
	Kelurahan        string `json:"kelurahan"`
	TempatLahir      string `json:"tempatLahir"`
	TanggalLahir     string `json:"tanggalLahir"`
	StatusPernikahan string `json:"statusPernikahan"`
	Amount           string `json:"amount"`
	Periode          string `json:"periode"`
	NomorHandphone   string `json:"nomorHandphone"`
	Email            string `json:"email"`
	KodePos          string `json:"kodepos"`
	RequestUser      string `json:"requestUser"`
	Remark           string `json:"remark"`
}

type RequestData struct {
	Branch           string `json:"branch"`
	Nik              string `json:"nik"`
	Nama             string `json:"nama"`
	JenisKelamin     string `json:"jenis_kelamin"`
	Alamat           string `json:"alamat"`
	Rt               string `json:"rt"`
	Rw               string `json:"rw"`
	Provinsi         string `json:"string"`
	Kota             string `json:"kota"`
	Kecamatan        string `json:"kecamatan"`
	Kelurahan        string `json:"kelurahan"`
	TempatLahir      string `json:"tempat_lahir"`
	TanggalLahir     string `json:"tanggal_lahir"`
	StatusPernikahan string `json:"status_pernikahan"`
	Amount           string `json:"amount"`
	Periode          string `json:"periode"`
	NomorHandphone   string `json:"nomor_handphone"`
	Email            string `json:"email"`
	KodePos          string `json:"kode_pos"`
}
