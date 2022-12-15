package service

import (
	"fmt"
	"go-api/dto"
	"go-api/dto/request/userRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/entity"
	"go-api/repository"
	"strconv"

	"github.com/mashingan/smapping"
)

// UserService is a contract.....
type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userID string) entity.User
	GetDeveloper(request userRequestDTO.ListUserDeveloperRequestDTO) responseDTO.Response
	ListUserReferral(request userRequestDTO.ListUserReferralRequestDTO) responseDTO.Response
}

type userService struct {
	userRepository repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		fmt.Println("Failed map %v:", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) entity.User {
	return service.userRepository.ProfileUser(userID)
}

func (service *userService) GetDeveloper(request userRequestDTO.ListUserDeveloperRequestDTO) responseDTO.Response {
	var response responseDTO.Response
	data := service.userRepository.GetDeveloper(request)
	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "Success"
	response.ResponseData = data
	response.Summary = nil
	return response
}

func (service *userService) ListUserReferral(request userRequestDTO.ListUserReferralRequestDTO) responseDTO.Response {
	var response responseDTO.Response
	var metadataResponse responseDTO.ListUserDtoRes

	metadataResponse.Currentpage = request.Offset
	if request.Offset > 0 {
		request.Offset = request.Offset * request.Limit
	}

	var sqlStr, sqlStr2 string

	if request.StartDate != "" && request.EndDate != "" {
		sqlStr = `SELECT name,mobile_no,properti_id,tpkbs.created_at from tbl_sales ts 
	join tbl_customer tc on tc.sales_id = ts.id
	join tbl_pengajuan_kpr_by_sales tpkbs on tpkbs.customer_id = tc.id
	where ts.sales_email like '%` + request.SalesEmail + `%' and tpkbs.created_at BETWEEN '` + request.StartDate + `' and '` + request.EndDate + `'  
	limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + ``

		sqlStr2 = `SELECT count(name) from tbl_sales ts 
	join tbl_customer tc on tc.sales_id = ts.id
	join tbl_pengajuan_kpr_by_sales tpkbs on tpkbs.customer_id = tc.id
	where ts.sales_email like '%` + request.SalesEmail + `%' and tpkbs.created_at BETWEEN '` + request.StartDate + `' and '` + request.EndDate + `'  
 	`
	} else if request.EndDate == "" && request.StartDate != "" {
		sqlStr = `SELECT name,mobile_no,properti_id,tpkbs.created_at from tbl_sales ts 
		join tbl_customer tc on tc.sales_id = ts.id
		join tbl_pengajuan_kpr_by_sales tpkbs on tpkbs.customer_id = tc.id
		where ts.sales_email like '%` + request.SalesEmail + `%' and date(tpkbs.created_at) >='` + request.StartDate + `'  
		limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + ``

		sqlStr2 = `SELECT count(name) from tbl_sales ts 
		join tbl_customer tc on tc.sales_id = ts.id
		join tbl_pengajuan_kpr_by_sales tpkbs on tpkbs.customer_id = tc.id
		where ts.sales_email like '%` + request.SalesEmail + `%' and date(tpkbs.created_at) >='` + request.StartDate + `'`

	} else {
		sqlStr = `SELECT name,mobile_no,properti_id,tpkbs.created_at from tbl_sales ts 
		join tbl_customer tc on tc.sales_id = ts.id
		join tbl_pengajuan_kpr_by_sales tpkbs on tpkbs.customer_id = tc.id
		where ts.sales_email like '%` + request.SalesEmail + `%' and tc.name like '%` + request.Keyword + `%' 
		or ts.sales_email like '%` + request.SalesEmail + `%' and tc.mobile_no like '%` + request.Keyword + `%' 
		or ts.sales_email like '%` + request.SalesEmail + `%' and tpkbs.properti_id like '%` + request.Keyword + `%'
		or ts.sales_email like '%` + request.SalesEmail + `%' and tpkbs.created_at like '%` + request.Keyword + `%' 
		limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + ``

		sqlStr2 = `SELECT count(name) from tbl_sales ts 
		join tbl_customer tc on tc.sales_id = ts.id
		join tbl_pengajuan_kpr_by_sales tpkbs on tpkbs.customer_id = tc.id
		where ts.sales_email like '%` + request.SalesEmail + `%' and tc.name like '%` + request.Keyword + `%' 
		or ts.sales_email like '%` + request.SalesEmail + `%' and tc.mobile_no like '%` + request.Keyword + `%' 
		or ts.sales_email like '%` + request.SalesEmail + `%' and tpkbs.properti_id like '%` + request.Keyword + `%'
		or ts.sales_email like '%` + request.SalesEmail + `%' and tpkbs.created_at like '%` + request.Keyword + `%'`
	}

	data, totalData := service.userRepository.GetUserReferral(request, sqlStr, sqlStr2)
	metadataResponse.TotalData = totalData
	response.HttpCode = 200
	response.MetadataResponse = metadataResponse
	response.ResponseCode = "success"
	response.ResponseData = data
	response.Summary = nil

	return response
}
