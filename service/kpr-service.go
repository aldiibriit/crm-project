package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-api/dto/request/KPRRequestDTO"
	"go-api/dto/request/brispotRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/entity"
	"go-api/helper"
	"go-api/repository"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mashingan/smapping"
)

type KPRService interface {
	PengajuanKPR(request KPRRequestDTO.PengajuanKPRRequest) responseDTO.Response
	ListPengajuan(request KPRRequestDTO.ListPengajuanKPR) responseDTO.Response
}

type kprService struct {
	customerRepository repository.CustomerRepository
	kprRepository      repository.KPRRepository
	salesRepository    repository.SalesRepository
	emailService       EmailService
	userRepository     repository.UserRepository
}

func NewKPRService(customerRepo repository.CustomerRepository, kprRepo repository.KPRRepository, salesRepo repository.SalesRepository, emailServ EmailService, userRepo repository.UserRepository) KPRService {
	return &kprService{
		customerRepository: customerRepo,
		kprRepository:      kprRepo,
		salesRepository:    salesRepo,
		emailService:       emailServ,
		userRepository:     userRepo,
	}
}

func (service *kprService) PengajuanKPR(request KPRRequestDTO.PengajuanKPRRequest) responseDTO.Response {
	var response responseDTO.Response
	// var emailRequest emailRequestDTO.EmailRequestDTO

	customer := entity.TblCustomer{}
	pengajuanKPR := entity.TblPengajuanKprBySales{}
	customer.CreatedAt = time.Now()
	customer.ModifiedAt = time.Now()
	pengajuanKPR.CreatedAt = time.Now()
	pengajuanKPR.ModifiedAt = time.Now()
	pengajuanKPR.Status = "on_reviewed"
	salesID := service.salesRepository.GetIDByRefCode(request.ReferralCode)
	request.SalesID = salesID

	err := smapping.FillStruct(&customer, smapping.MapFields(&request))
	if err != nil {
		response.HttpCode = 500
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = "failed map " + err.Error()
		response.Summary = nil

		return response
	}

	err = smapping.FillStruct(&pengajuanKPR, smapping.MapFields(&request))
	if err != nil {
		response.HttpCode = 500
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = "failed map " + err.Error()
		response.Summary = nil

		return response
	}

	// data := service.customerRepository.Insert(customer)
	// pengajuanKPR.CustomerID = data.ID
	// service.kprRepository.Insert(pengajuanKPR)

	var alamatKTP map[string]interface{}

	json.Unmarshal([]byte(request.AlamatKTP), &alamatKTP)

	fmt.Println("Alamat KTP : ", alamatKTP)

	// requestData := brispotRequestDTO.RequestData{

	// }
	brispotSubmitRequest := brispotRequestDTO.BrispotSubmitRequestDTO{
		RequestUser:  "00132152",
		Branch:       "509",
		Nik:          request.NIK,
		Nama:         request.Name + " TEST ",
		JenisKelamin: "L",
		Alamat:       fmt.Sprintf("%v", alamatKTP["alamat"]),
		Rt:           fmt.Sprintf("%v", alamatKTP["rt"]),
		Rw:           fmt.Sprintf("%v", alamatKTP["rw"]),
		Provinsi:     fmt.Sprintf("%v", alamatKTP["provinsi"]),
		// Kota:             fmt.Sprintf("%v", alamatKTP["kota"]),
		Remark:           "testing golang aldi",
		Kota:             "Jakarta Barat",
		Kecamatan:        fmt.Sprintf("%v", alamatKTP["kecamatan"]),
		Kelurahan:        fmt.Sprintf("%v", alamatKTP["kelurahan"]),
		KodePos:          fmt.Sprintf("%v", alamatKTP["kodePos"]),
		TempatLahir:      request.Pob,
		TanggalLahir:     request.TanggalLahir,
		StatusPernikahan: "1",
		Amount:           request.JumlahPinjaman,
		Periode:          request.Periode,
		NomorHandphone:   request.MobileNo,
		Email:            request.Email,
	}

	jsonByte, err := json.Marshal(&brispotSubmitRequest)
	if err != nil {
		response.HttpCode = 500
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = err.Error()
		response.Summary = nil

		return response
	}

	if !passToBrispot(string(jsonByte)) {
		response.HttpCode = 500
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = "error while pass to brispot"
		response.Summary = nil

		return response
	}

	// emailRequest.Action = 6
	// emailRequest.EmailBody = request.Name + ` Berminat untuk membeli salah satu properti anda. Mohon cek MIS`
	// emailRequest.Name = request.Name
	// emailRequest.Subject = "Pengajuan KPR"
	// emailRequest.ToAddres = request.SalesEmail

	// if !service.emailService.SendMessage(emailRequest) {
	// 	response.HttpCode = 422
	// 	response.MetadataResponse = nil
	// 	response.ResponseCode = "99"
	// 	response.ResponseData = nil
	// 	response.ResponseDesc = "error while send email to sales"
	// 	response.Summary = nil
	// }

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = brispotSubmitRequest
	response.ResponseDesc = "Success"
	response.Summary = nil

	return response
}

func (service *kprService) ListPengajuan(request KPRRequestDTO.ListPengajuanKPR) responseDTO.Response {
	var response responseDTO.Response

	// check userType
	if service.userRepository.CheckType(request.Email) != "superadmin" {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = "Permission denied"
		response.Summary = nil

		return response
	}

	data := service.kprRepository.KPRAllUser()

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = data
	response.ResponseDesc = "Success"
	response.Summary = nil

	return response
}

func passToBrispot(jsonString string) bool {
	url := "https://api.privateopen.sandbox.bri.co.id/gateway/Oauth2/1.0/accessToken"
	method := "POST"

	payload := strings.NewReader("client_id=b78f2b47-b258-49cb-9a0d-cab6f4d5f80a&client_secret=eff57260-b4c2-4b5d-b04a-b59b2a9f51d1&grant_type=client_credentials&scope=apiKPRBrispot")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "TS015ba640=01d4c679fed33239583546a1c96b3f8fa7de4d7e56256a5d3c734980c774b4ebceb059cdd06f47912fa72d0836bf7361dcde55bea8")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var responseCreateToken map[string]interface{}
	// accessToken := fmt.Sprintf("%v", responseCreateToken["access_token"])
	json.Unmarshal([]byte(string(body)), &responseCreateToken)

	ts := time.Now().UTC()
	timestampFormatted := ts.Format(time.RFC3339Nano)
	timestamp := timestampFormatted[:23] + "Z"

	signature := helper.HMACBuilder("/gateway/apiKPRBrispot/1.0/insertPrakarsaMortgage", "POST", "Bearer "+responseCreateToken["access_token"].(string), jsonString, "eff57260-b4c2-4b5d-b04a-b59b2a9f51d1", timestamp)

	if responseCreateToken["access_token"].(string) != "" {
		fmt.Println("bisa create token", responseCreateToken["access_token"])
		fmt.Println(signature)

		fmt.Println("signature", signature)
		fmt.Println("json string", jsonString)
		fmt.Println("ts", timestamp)
		fmt.Println("bearertoken", responseCreateToken["access_token"])
		url := "https://api.privateopen.sandbox.bri.co.id/gateway/apiKPRBrispot/1.0/insertPrakarsaMortgage"
		method := "POST"

		// client := &http.Client{}
		// stringRequest := strings.NewReader(jsonString)
		payload := []byte(jsonString)
		fmt.Println(jsonString)
		req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("error 1")
			return false
		}
		req.Header["BRI-Signature"] = []string{signature}
		req.Header["BRI-Timestamp"] = []string{timestamp}
		req.Header["Authorization"] = []string{`Bearer ` + responseCreateToken["access_token"].(string) + ``}
		req.Header["Content-Type"] = []string{"application/json"}
		// req.Header.Add("Cookie", "TS015ba640=01d4c679fe9ad1be95c03962152f0d53883dc4fc57a78611ca05cd70308628fa76f8371a7a354b9c602209f6112dcdf6f236480c1a")
		req.Header.Add("Cookie", "PHPSESSID=j4qqnnupjdffjh7iohlgbt1v37; TS015ba640=01d4c679fe8c6caeea34b579ed9c7351250a19d2d2e61ea941bf1b1f99bdadcb494e1740c364f958893a6018a2bc982e3d64661d7dbdd7bcfa3e08efc1d35a0956882a1e549544d6f9dd256b23011c8ad137cbc201; ci_session=a%3A5%3A%7Bs%3A10%3A%22session_id%22%3Bs%3A32%3A%22d48fb6d6e298068483aac5cdbaf5f978%22%3Bs%3A10%3A%22ip_address%22%3Bs%3A12%3A%2210.35.65.158%22%3Bs%3A10%3A%22user_agent%22%3Bs%3A21%3A%22PostmanRuntime%2F7.29.2%22%3Bs%3A13%3A%22last_activity%22%3Bi%3A1672029326%3Bs%3A9%3A%22user_data%22%3Bs%3A0%3A%22%22%3B%7Dd6698797d098931b1eb594a99ab5fa4f56a1a0a5")

		// reqDump, err := httputil.DumpRequestOut(req, true)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Printf("REQUEST:\n%s", string(reqDump))
		// res, err := client.Do(req)
		// if err != nil {
		// 	fmt.Println("error 2")
		// 	return false
		// }
		// defer res.Body.Close()

		// body, err := ioutil.ReadAll(res.Body)
		// if err != nil {
		// 	fmt.Println("error 3")
		// 	return false
		// }
		// var responseInsertPrakarsaMortgage map[string]interface{}
		// // accessToken := fmt.Sprintf("%v", responseInsertPrakarsaMortgage["access_token"])
		// json.Unmarshal([]byte(string(body)), &responseInsertPrakarsaMortgage)
		// fmt.Println(responseInsertPrakarsaMortgage)
	} else {
		fmt.Println("ga bisa create token")
		return false
	}

	return true
}
