package propertiDTO

type AdvancedFilterDTO struct {
	Email string `json:"email"`
	Section1
	Section2
	Section3
	Section4
	Section5
	Section6
}

type Section1 struct {
	HargaMin         int      `json:"hargaMin"`
	HargaMax         int      `json:"hargaMax"`
	TipeProperti     string   `json:"tipeProperti"`
	JumlahKamarTidur []string `json:"jumlahKamarTidur"`
	JumlahKamarMandi []string `json:"jumlahKamarMandi"`
}

type Section2 struct {
	LuasTanahMinimal     int `json:"luasTanahMinimal"`
	LuasTanahMaksimal    int `json:"luasTanahMaksimal"`
	LuasBangunanMinimal  int `json:"luasBangunanMinimal"`
	LuasBangunanMaksimal int `json:"luasBangunanMaksimal"`
}

type Section3 struct {
	JenisProperti string `json:"jenisProperti"`
}

type Section4 struct {
	Dapur         bool `json:"dapur"`
	JalurTelepon  bool `json:"jalurTelepon"`
	JalurListrik  bool `json:"jalurListrik"`
	JalurPDAM     bool `json:"jalurPDAM"`
	RuangKeluarga bool `json:"ruangKeluarga"`
	RuangKerja    bool `json:"ruangKerja"`
}

type Section5 struct {
	RumahSakit  bool `json:"rumahSakit"`
	JalanTol    bool `json:"jalanTol"`
	Sekolah     bool `json:"sekolah"`
	Mall        bool `json:"mall"`
	BankATM     bool `json:"bankATM"`
	Taman       bool `json:"taman"`
	Pasar       bool `json:"pasar"`
	Farmasi     bool `json:"farmasi"`
	RumahIbadah bool `json:"rumahIbadah"`
	Restoran    bool `json:"restoran"`
	Bioskop     bool `json:"bioskop"`
	Bar         bool `json:"bar"`
	Halte       bool `json:"halte"`
	Stasiun     bool `json:"stasiun"`
	Bandara     bool `json:"bandara"`
	GerbangTol  bool `json:"gerbangTol"`
	SPBU        bool `json:"spbu"`
	Gymnasium   bool `json:"gymnasium"`
}

type Section6 struct {
	KolamRenang  bool `json:"kolamRenang"`
	TempatParkir bool `json:"tempatParkir"`
	Keamanan     bool `json:"keamanan"`
	Penghijauan  bool `json:"penghijauan"`
	Lift         bool `json:"lift"`
	ClubHouse    bool `json:"clubHouse"`
	Elevator     bool `json:"elevator"`
	Gym          bool `json:"gym"`
	Garasi       bool `json:"garasi"`
	RowJalan12   bool `json:"rowJalan12"`
}
