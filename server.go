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
// 	// Password := "feryganteng@yopmail.com"
// 	// encryptedPassword, _ := helper.RsaEncryptFEToBE([]byte(Password))
// 	// fmt.Println("Password : ", encryptedPassword)
// 	// decodedNama, _ := base64.StdEncoding.DecodeString("HZ6wNUg5TFd0E3hsq5BWPtPBGE3jlWsWWyfXlTuatglziKMhLj5yArbvuMlZrjw53PXH7+j0vI6m36jKzES3/wMC9C48Co7+aN4N/PBNU3eRQ8HQMuSkVdhbm17IrHTfL30iJcJxkkmwFIReHyROLlZD4M4f4H3KwTDnZ3v36h9qT5KBgmcwciRMJGhE2KjJqmahoW/5/4ElipOfImUycg0C4vngz40ZR1qHbt+84tgdhlhSXr6c4Txa+Xcz1fB/vRfv0bInFr+KLqJvFpD9ueelMKrhojT5zBjPp3bmeduwuja4R1PraAXFkKdCChzEkf0trOhrez/NHxlgpExrPjLihy+L8bJ1ixFNFzi3ZEwhzeyLToFZMExvtmi8o/E30OWVu7Mb19mZYKiThq8BNqwqaVYPT5RceIcdqbfNNVB3dr4XxxJ9zJHEshk27+6IGUO6PLFFka5/vLa4zAHZnDsw3BrVevqks/2EmI6Q6QxuQtDiRDy6xMVcQG8hbOqlzU9sliU5sXBU2lA6lnqQ6j0x2Lj3EQlufcSDmfxGqgRhBXZjgtolLokZ9xo+MmlW80B2YDwvs98Ob1hHZ26uhAZUJYMIwdx4aqEiJdFVxKBYQsyLptkoiyL6Si96fC9R7Uy13aZFvE79PLKHLtMkaWOSQWNo3PZWm6/dedo+jIw=")
// 	// plainNama, _ := helper.RsaDecryptFromBEInFE(decodedNama)
// 	// fmt.Println(plainNama)
// 	// id := uuid.New()
// 	// fmt.Println(reflect.TypeOf(id))
// 	// fmt.Println(reflect.TypeOf(id.String()))
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
		propertiRoutes.POST("/landingPage", propertiController.LandingPage)
	}

	testRoutes := r.Group("/test")
	{
		testRoutes.POST("", testController.TestEmail)
	}

	r.Run(":7177")
}
