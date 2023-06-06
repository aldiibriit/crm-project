package LogActivityResponseDTO

type Response struct {
	HttpCode        int      `json:"-"`
	ResponseMessage string   `json:"responseMessage"`
	ResponseCode    string   `json:"responseCode"`
	Timeline        Timeline `json:"timeline"`
}

type Timeline struct {
	Prestaging  DetailTimeline `json:"preStaging"`
	Staging     DetailTimeline `json:"staging"`
	StagingLive DetailTimeline `json:"stagingLive"`
}

type DetailTimeline struct {
	Sn           string `json:"sn" gorm:"column:sn"`
	Category     string `json:"category" gorm:"column:category"`
	SubmittedAt  string `json:"submittedAt" gorm:"column:submittedAt"`
	RejectedAt   string `json:"rejectedAt" gorm:"column:rejectedAt"`
	ReuploadedAt string `json:"reuploadedAt" gorm:"column:reuploadedAt"`
	ApprovedAt   string `json:"approvedAt" gorm:"column:approvedAt"`
}
