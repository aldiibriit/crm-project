package CCTVRequestDTO

type FindBySNRequest struct {
	SnCctv string `json:"snCctv" binding:"required"`
}
