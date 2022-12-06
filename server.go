package main

import (
	"go-api/config"
	"go-api/controller"
	"go-api/middleware"
	"go-api/repository"
	"go-api/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db                     *gorm.DB                          = config.SetupDatabaseConnection()
	userRepository         repository.UserRepository         = repository.NewUserRepository(db)
	bookRepository         repository.BookRepository         = repository.NewBookRepository(db)
	propertiRepository     repository.PropertiRepository     = repository.NewPropertiRepository(db)
	jwtHistRepository      repository.JwtHistRepository      = repository.NewJwtHistRepository(db)
	salesRepository        repository.SalesRepository        = repository.NewSalesRepository(db)
	emailAttemptRepository repository.EmailAttemptRepository = repository.NewEmailAttemptRepository(db)
	otpRepository          repository.OTPRepository          = repository.NewOTPRepository(db)
	otpAttemptRepository   repository.OTPAttemptRepository   = repository.NewOTPAttemptRepository(db)
	customerRepository     repository.CustomerRepository     = repository.NewCustomerRepository(db)
	jwtService             service.JWTService                = service.NewJWTService(jwtHistRepository)
	userService            service.UserService               = service.NewUserService(userRepository)
	bookService            service.BookService               = service.NewBookService(bookRepository)
	baseService            service.BaseService               = service.NewBaseService(jwtHistRepository, jwtService)
	propertiService        service.PropertiService           = service.NewPropertiService(propertiRepository, userRepository, baseService)
	otpService             service.OTPService                = service.NewOTPService(otpRepository, emailService, otpAttemptRepository)
	authService            service.AuthService               = service.NewAuthService(userRepository, salesRepository, emailService, emailAttemptRepository, otpService)
	emailService           service.EmailService              = service.NewEmailService(emailAttemptRepository)
	kprService             service.KPRService                = service.NewKPRService(customerRepository)
	salesService           service.SalesService              = service.NewSalesService(salesRepository, userRepository)
	testController         controller.TestController         = controller.NewTestController(emailService)
	authController         controller.AuthController         = controller.NewAuthController(authService, jwtService)
	userController         controller.UserController         = controller.NewUserController(userService, jwtService)
	bookController         controller.BookController         = controller.NewBookController(bookService, jwtService)
	propertiController     controller.PropertiController     = controller.NewPropertiController(propertiService)
	otpController          controller.OTPController          = controller.NEwOTPController(otpService)
	salesController        controller.SalesController        = controller.NewSalesController(salesService)
	kprController          controller.KPRController          = controller.NewKPRController(kprService)
)

// func main() {
// 	Password := "3853"
// 	encryptedPassword, err := helper.RsaEncryptFEToBE([]byte(Password))
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println("Password : ", encryptedPassword)
// 	// decodedNama, _ := base64.StdEncoding.DecodeString("fmDO1p4jofEp1vXNKJHJinv/UwLNT7JCrfMWFVMkWYiGhScWJdHLD6LwTkLQiYIJj13dQKAoiWwhL6XLvBzO38rkKmYy5LiGMtLezkAlmFL0YADqEsjY6xAzOZw8j58jnTsrN6ZSTaUjc9jJLVxOj3yHBKORtP1k2A7R2x46J22LqBhDsFBwijb/m5iCBCgOWyVMLHHhdqJlMTCo4cVduL6los6T1Elfmqew3ko8USPKfB+C9DSVzmtdBaJy+FyLnwf0cp9y57mgTcHjBatsCCX9/uYtcZAB3hRzML7d3jM4aVpTvJttPLE37cFq+Kl/gnbH55HSIDSp+GTdMA/u7IEvXWo7CBst1ciDfRxaCh/Rz/ax0PLSx1LxeBhyYiMPjOQ6Erj8ZZ4Gzf8APVJjsjqz4cdte5hFZ8q04Az6I760LyOHN6dgXRy7mE9GOl+gKWRN9pWRb4c+1T2ZX54D+gwITORF3rzDhOGcac9o75nYKFbcNI+uLepDDqYt8/2tnlKh2Qt77beq06+YbuygjIOmw54jkkpxgCiiwFdlg8tnkrudnO+UyXhQOF615i41XhdsCs7hl9aufcqGkzyVVuIsyP69bhYifweL1udUMrIEsHmxhAugIK94tH1kcdEgI2LFYRqSKA+uyWVV65r1B3aca2ds9SHPt/InmIuYnMM=")
// 	// plainNama, _ := helper.RsaDecryptFromBEInFE(decodedNama)
// 	// fmt.Println(plainNama)
// 	// id := uuid.New()
// 	// fmt.Println(reflect.TypeOf(id))
// 	// fmt.Println(reflect.TypeOf(id.String()))
// 	// 	helper.GenRsaKeyForBE(1024)
// 	// 	helper.GenRsaKeyForFE(1024)
// }

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/createToken", authController.CreateToken)
		authRoutes.POST("/registerUserSales", authController.RegisterSales)
		authRoutes.POST("/activateUser", authController.ActivateUser)
		authRoutes.POST("/passwordConfirmation", authController.PasswordConfirmation)
	}

	otpRoutes := r.Group("api/otp")
	{
		otpRoutes.POST("/validateOTP", otpController.ValidateOTP)
	}

	salesRoutes := r.Group("api/salesMIS")
	{
		salesRoutes.POST("/atDeveloper", salesController.MISDeveloper)
		salesRoutes.POST("/atSuperAdmin", salesController.MISSuperAdmin)
		salesRoutes.POST("/listProject", salesController.ListProject)
		salesRoutes.PUT("/editSalesByDeveloper", salesController.EditSalesByDeveloper)
		salesRoutes.POST("/detailSalesAtDeveloper", salesController.DetailSalesByDeveloper)
		salesRoutes.POST("/deleteSalesByDeveloper", salesController.DeleteSalesByDeveloper)
	}

	userRoutes := r.Group("api/user")
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		userRoutes.POST("/userDeveloperList", userController.GetDeveloper)
	}

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.All)
		bookRoutes.POST("/", bookController.Insert)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
	}

	propertiRoutes := r.Group("api/properti")
	{
		propertiRoutes.POST("/advancedFilter", propertiController.AdvancedFilter)
		propertiRoutes.POST("/landingPage", propertiController.LandingPage)
	}

	testRoutes := r.Group("/test")
	{
		testRoutes.POST("", testController.TestEmail)
	}

	kprRoutes := r.Group("api/kpr")
	{
		kprRoutes.POST("/pengajuanKPR", kprController.PengajuanKPR)
	}

	r.Run(":7177")
}
