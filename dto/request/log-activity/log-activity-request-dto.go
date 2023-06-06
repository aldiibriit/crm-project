package LogActivityRequest

type InsertRequest struct {
	Category          string `json:"category"`
	Sn                string `json:"sn"`
	StatusDescription string `json:"statusDescription"`
}

type GetTimeLine struct {
	Sn string `json:"sn" binding:"required"`
}
