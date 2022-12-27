package repository

import (
	"log"

	"go-api/dto/request/userRequestDTO"
	"go-api/dto/response/userResponseDTO"
	"go-api/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository is contract what userRepository can do to db
type UserRepository interface {
	CheckType(email string) string
	InsertUser(user entity.User) entity.User
	InsertUserSales(user entity.TblUser) error
	UpdateUser(user entity.User) entity.User
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entity.User
	FindByEmail2(email string) entity.TblUser
	CheckUserExist(email string) bool
	UpdateOrCreate(data entity.TblUser)
	GetLatestId() entity.TblUser
	GetDeveloper(request userRequestDTO.ListUserDeveloperRequestDTO) []userResponseDTO.UserDeveloperResponse
	GetUserReferral(request userRequestDTO.ListUserReferralRequestDTO, sqlStr string, sqlStr2 string) ([]userResponseDTO.UserReferralResponse, int)
	ProfileUser(userID string) entity.User
}

type userConnection struct {
	connection *gorm.DB
}

// NewUserRepository is creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))
	db.connection.Save(&user)
	return user
}

func (db *userConnection) InsertUserSales(user entity.TblUser) error {
	err := db.connection.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *userConnection) UpdateUser(user entity.User) entity.User {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser entity.User
		db.connection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}

	db.connection.Save(&user)
	return user
}

func (db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User
	res := db.connection.Where("email = ?", email).Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnection) FindByEmail(email string) entity.User {
	var user entity.User
	// db.connection.Where("email = ?", email).Take(&user)
	db.connection.Debug().Raw("SELECT * from users where email = ?", email).Take(&user)
	return user
}

func (db *userConnection) FindByEmail2(email string) entity.TblUser {
	var user entity.TblUser
	// db.connection.Where("email = ?", email).Take(&user)
	db.connection.Raw("SELECT *,json_extract(metadata,'$.name')as userName from tbl_user where email = ?", email).Find(&user)
	return user
}

func (db *userConnection) CheckUserExist(email string) bool {
	var user entity.TblUser
	db.connection.Debug().Raw("SELECT * from tbl_user where email = ?", email).Find(&user)
	if user.Email == "" {
		return false
	}
	return true
}

func (db *userConnection) ProfileUser(userID string) entity.User {
	var user entity.User
	db.connection.Preload("Books").Preload("Books.User").Find(&user, userID)
	return user
}

func (db *userConnection) UpdateOrCreate(data entity.TblUser) {
	var checker entity.TblUser
	db.connection.Raw("SELECT * FROM tbl_user where email = ?", data.Email).Take(&checker)
	if checker.Email == "" {
		db.connection.Save(&data)
	} else if checker.Email != "" {
		db.connection.Model(&entity.TblUser{}).Where("email = ?", data.Email).Updates(&data)
	}
}

func (db *userConnection) GetLatestId() entity.TblUser {
	var data entity.TblUser
	db.connection.Raw("SELECT id from tbl_user order by id desc limit 1").Take(&data)
	return data
}

func (db *userConnection) GetDeveloper(request userRequestDTO.ListUserDeveloperRequestDTO) []userResponseDTO.UserDeveloperResponse {
	var data []userResponseDTO.UserDeveloperResponse
	db.connection.Raw(`SELECT email,json_extract(metadata,'$.name')as name FROM tbl_user where type = 'developer' and status = 'active' and email like '%` + request.Keyword + `%' or json_extract(metadata,'$.name') like '%` + request.Keyword + `%' and type = 'developer' and status = 'active'`).Find(&data)
	return data
}

func (db *userConnection) GetUserReferral(request userRequestDTO.ListUserReferralRequestDTO, sqlStr string, sqlStr2 string) ([]userResponseDTO.UserReferralResponse, int) {
	var data []userResponseDTO.UserReferralResponse
	var totalData int
	db.connection.Raw(sqlStr).Find(&data)
	db.connection.Raw(sqlStr2).Scan(&totalData)
	return data, totalData
}

func (db *userConnection) CheckType(email string) string {
	var userType string
	db.connection.Raw(`SELECT type FROM tbl_user where email = '` + email + `'`).Scan(&userType)
	return userType
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}
