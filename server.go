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
	testController         controller.TestController         = controller.NewTestController(emailService)
	authController         controller.AuthController         = controller.NewAuthController(authService, jwtService)
	userController         controller.UserController         = controller.NewUserController(userService, jwtService)
	bookController         controller.BookController         = controller.NewBookController(bookService, jwtService)
	propertiController     controller.PropertiController     = controller.NewPropertiController(propertiService)
	otpController          controller.OTPController          = controller.NEwOTPController(otpService)
)

// func main() {
// 	nama := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njg1ODkxOTgsImlhdCI6MTY2ODU4ODg5OCwic3ViIjoiMSJ9.THit_ZcNbO6h9tLHZbXYRh-m0H5EV3LE-8-n4nCxueM"
// 	encrypt, _ := helper.RsaEncryptFEToBE([]byte(nama))
// 	fmt.Println(encrypt)
// 	// // decodedNama, _ := base64.StdEncoding.DecodeString("XJMxFQg1jlCIzG5VXOcZqOWlCxkfcfIt7GgZNG19BYFHbIxxkM1XjSqpGgCjWxdb0fjtMzKtlct9/L9fLKB/xJo8MZGWK+uLcppVyB626K0X4SGml3iFS5YL/p/EpHFTu6h5YvpIBxBgZx9F1r6LMl+UAiaJiE28xPT4hy+oAj6ZcSM8P/0cUUyCOvxPHBeSXh7RGvKLYebTMdglgCMVACe28pkFKc4vIEd7f8q9cvQkJ52n6/TG9cN1zjBMVS4ZwUbSI/1slzgwbNrUF21ZnSg7nJkIS+ndUWZ4t0H0ynSRe1Pnt45uAyofq0qV4atZcNKyZgRp9zTNyI2cCx1l8DonmuYrOlQA367Rsu6wEJIsjrmXM47T2LO7s0h6QMJR/wScSY1RRPxLnmb13XnLp244MJZBJ66ZDGsASa/XLqcn6WL7FbTazVt94ixqLAUjGalkZTTE/i5RAwVP36X1LxKUFmAhz5OzIj6MhspT3sbhGVN1NamVoE1Phll3e6lQI7VRfaSSvHzEZeR33LXmCeFxas+/V4eeJYuagrNwgd+Sym2RJk8/4vnk1E5pdkPDUlCOz2JCGUM0rlOBDClMshZqvQuKTTaGaXHgNJ1nO61bJ+4zy4PpuyZYPOXoBonw2IfyhduKoasxDYvsYn2CWExHcpnwFjAmeyApwFPz9NQ=")
// 	// // plainNama, _ := helper.RsaDecryptFromBEInFE(decodedNama)
// 	// // fmt.Println(plainNama)
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
