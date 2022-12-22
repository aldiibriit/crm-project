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
	db *gorm.DB = config.SetupDatabaseConnection()

	userRepository         repository.UserRepository         = repository.NewUserRepository(db)
	bookRepository         repository.BookRepository         = repository.NewBookRepository(db)
	propertiRepository     repository.PropertiRepository     = repository.NewPropertiRepository(db)
	jwtHistRepository      repository.JwtHistRepository      = repository.NewJwtHistRepository(db)
	salesRepository        repository.SalesRepository        = repository.NewSalesRepository(db)
	emailAttemptRepository repository.EmailAttemptRepository = repository.NewEmailAttemptRepository(db)
	otpRepository          repository.OTPRepository          = repository.NewOTPRepository(db)
	otpAttemptRepository   repository.OTPAttemptRepository   = repository.NewOTPAttemptRepository(db)
	customerRepository     repository.CustomerRepository     = repository.NewCustomerRepository(db)
	kprRepository          repository.KPRRepository          = repository.NewKPRRepository(db)

	jwtService      service.JWTService      = service.NewJWTService(jwtHistRepository)
	userService     service.UserService     = service.NewUserService(userRepository)
	bookService     service.BookService     = service.NewBookService(bookRepository)
	baseService     service.BaseService     = service.NewBaseService(jwtHistRepository, jwtService)
	propertiService service.PropertiService = service.NewPropertiService(propertiRepository, userRepository, baseService)
	otpService      service.OTPService      = service.NewOTPService(otpRepository, emailService, otpAttemptRepository)
	authService     service.AuthService     = service.NewAuthService(userRepository, salesRepository, emailService, emailAttemptRepository, otpService, jwtService)
	emailService    service.EmailService    = service.NewEmailService(emailAttemptRepository)
	kprService      service.KPRService      = service.NewKPRService(customerRepository, kprRepository, salesRepository, emailService)
	salesService    service.SalesService    = service.NewSalesService(salesRepository, userRepository, kprRepository)

	testController     controller.TestController     = controller.NewTestController(emailService)
	authController     controller.AuthController     = controller.NewAuthController(authService, jwtService)
	userController     controller.UserController     = controller.NewUserController(userService, jwtService)
	bookController     controller.BookController     = controller.NewBookController(bookService, jwtService)
	propertiController controller.PropertiController = controller.NewPropertiController(propertiService)
	otpController      controller.OTPController      = controller.NEwOTPController(otpService)
	salesController    controller.SalesController    = controller.NewSalesController(salesService)
	kprController      controller.KPRController      = controller.NewKPRController(kprService)
)

// func main() {
// 	// helper.SignatureBRIVA()
// 	// decodedNama, _ := base64.StdEncoding.DecodeString("YfWaK9aA0mybw/heC6zYRLKlhzVngablTXNJBSHSp4yB1Pg8GcO1dmLdd0ok4Yojsz9XO7pBLpRIROjSbpyi8YvkfdiVLpbpl2sNkemoOmYVO1q62aV6u6mboJ516vkvDeE09z4LgYzAXJf0SL99EuKQDL1++qzDBbr9ubmDmgAX/C3UffgBH4yqijZfG5hGpj/UkRPcBE04g3WwtYpyDXalDQYog7gxnOmUw5h4TfJnTjCZsaNHWsAavKcE9+zbyRzMXhxcvNF6H4S8lvFAE24i3dEYzCJXVFTRSrkvPS6VrT96d0QMayJTf37LzdJRrfJzI9YojeJGw0xxnuPk5Q==")
// 	// plainNama, _ := helper.RsaDecryptFromBEInFE(decodedNama)
// 	// fmt.Println(plainNama)
// 	Password := "aldi.sptra86@gmail.com"
// 	encryptedPassword, err := helper.RsaEncryptFEToBE([]byte(Password))
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println("Password : ", encryptedPassword)
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
		authRoutes.POST("/passthroughLogin", authController.PassthroughLogin)
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
		salesRoutes.POST("/editSalesByDeveloper", salesController.EditSalesByDeveloper)
		salesRoutes.POST("/detailSalesAtDeveloper", salesController.DetailSalesByDeveloper)
		salesRoutes.POST("/deleteSalesByDeveloper", salesController.DeleteSalesByDeveloper)
		salesRoutes.POST("/draftDetail", salesController.DraftDetail)
		salesRoutes.POST("/draftDelete", salesController.DeletePengajuan)
		salesRoutes.POST("/listFinalPengajuan", salesController.ListFinalPengajuan)
	}

	userRoutes := r.Group("api/user")
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		userRoutes.POST("/userDeveloperList", userController.GetDeveloper)
		userRoutes.POST("/userReferralList", userController.GetUserReferral)
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
