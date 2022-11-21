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
	jwtService             service.JWTService                = service.NewJWTService(jwtHistRepository)
	userService            service.UserService               = service.NewUserService(userRepository)
	bookService            service.BookService               = service.NewBookService(bookRepository)
	baseService            service.BaseService               = service.NewBaseService(jwtHistRepository, jwtService)
	propertiService        service.PropertiService           = service.NewPropertiService(propertiRepository, userRepository, baseService)
	otpService             service.OTPService                = service.NewOTPService(otpRepository, emailService, otpAttemptRepository)
	authService            service.AuthService               = service.NewAuthService(userRepository, salesRepository, emailService, emailAttemptRepository, otpService)
	emailService           service.EmailService              = service.NewEmailService(emailAttemptRepository)
	salesService           service.SalesService              = service.NewSalesService(salesRepository)
	testController         controller.TestController         = controller.NewTestController(emailService)
	authController         controller.AuthController         = controller.NewAuthController(authService, jwtService)
	userController         controller.UserController         = controller.NewUserController(userService, jwtService)
	bookController         controller.BookController         = controller.NewBookController(bookService, jwtService)
	propertiController     controller.PropertiController     = controller.NewPropertiController(propertiService)
	otpController          controller.OTPController          = controller.NEwOTPController(otpService)
	salesController        controller.SalesController        = controller.NewSalesController(salesService)
)

// func main() {
// 	Password := "3"
// 	encryptedPassword, _ := helper.RsaEncryptFEToBE([]byte(Password))
// 	fmt.Println("Password : ", encryptedPassword)
// 	// decodedNama, _ := base64.StdEncoding.DecodeString("bGkE14BjXuuP/Sn/0Rd9DO5FupG2IOGbbSZzYZwT0UUZjsGTD2I4c4X5X6TR0gCuT40cOk16+bSg5UAzY9V560GyY79cq/7DjwnMiGd6X9tzZmrp1tGxujZcf5hnbLqZs/tblK/l/8g/Bk67Wx6ASB7CNC2GbPELspc5r3io0tQ7CTXzDgbmGPtex1gqQamPI07chewqtiqyRsGgJJ3g5KZf5GHN3iCMKCLrLs2HTfBPOdWqPy2e5a29aqPdTlmpVaQmgZSC+kbM7IgCn9OHZNdEd2ZJU0z38PSlWrz3TDwWV7lst8mil3z5BUPDjF2DVRExOkRqPvPZckiLP3Xk6W/ugV9pHkKrxRsice0CAIMFmc3vc0lUWCIMkX2Oc+zHbtaMdz6RcrNAlqB/WLIbpXrzr61u3k0PVFShns074bbwrygsl/M24g/2aLotZWIvUoCNDNhTOjOVXOKUQD+ZzCgGNWPvHWwwq8PphLR1WoUdffqX2ZC+uSrawhiOKoOEytuBN27hiQ4J8WqNYINb8zWsug59TjcKW6mX9Gy47NJ3CpZf+H8KQiJZkZaEfyFPsDOJanXUWu8VlyVGr38SFfHPx/1CHebKWwMV8YjKelkszccNDY9SPlMIqAktt2G63r6atEpuNxLeFhNL0tgKct3LdcG6FXezR3X02m72Tt8=")
// 	// plainNama, _ := helper.RsaDecryptFromBEInFE(decodedNama)
// 	// fmt.Println(plainNama)
// 	// id := uuid.New()
// 	// fmt.Println(reflect.TypeOf(id))
// 	// fmt.Println(reflect.TypeOf(id.String()))
// }

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

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
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
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
	}

	testRoutes := r.Group("/test")
	{
		testRoutes.POST("", testController.TestEmail)
	}

	r.Run(":7177")
}
