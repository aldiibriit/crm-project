package entity

type JwtHistGo struct {
	Id    string `json:"id" gorm:"column:id"`
	Email string `json:"email" gorm:"column:email"`
	Jwt   string `json:"jwt" gorm:"column:jwt"`
}
