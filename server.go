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
// 	Password := "feryganteng@yopmail.com"
// 	encryptedPassword, _ := helper.RsaEncryptFEToBE([]byte(Password))
// 	fmt.Println("Password : ", encryptedPassword)
// 	// decodedNama, _ := base64.StdEncoding.DecodeString("D7SO96BeYYOUubuqnRj19FMmBL4N9D1wK+U2Ob5f9gC8okP6R5Z8uWZ1LJ3rEWgCmFVY8TS/iwnv2czoZe3nN+S88ATPkTFkI1LHSv+DrE711lkf9oxLnGvkmgkUDfBLXXgCAXNt+5HRwTwP2S4QKVOidj/cCTeRxsr9U0wgJMpsS2084dJ3N61ct6p4FwuZXKsR8m4mmwq20lkTEGbFPwFeAj3o2rUnAO0EOfnpOe8AznEOo6bFpVhFUyH12hLOI3l6O9ymp7X7uYQa27K+QLzPvR+Z3GY6S9GTX2in+hGjtMfRy/aK5muIAXh2PknEdGPtzusurehMOnqanD5GO6s0SJClUjE5o1LpbEYzl64vdH5YrJZ/6CMy6cny4qcQ7LToxUuQcGOEDeCXRvNVCvOgJieB5BKmwvhEDJcMCcOMG6qceAaYpZAImXS0YqBmGcO3zFHXlzXmWpbFEj8BZw6MmE8xyUT3Pfpt3/2FOPV9hKvY4aZvhYK7yhyaZWF+uQ+IGBC2IF5VWVkcRn/4geBIj79G60LvUFu8HiWeTDENCemCn3+Qi6yb9PrAANZ9ju93cuyN421tzlm4VYuds9/M82H5RbVcIo9lRCXdt80l8aliXWdFmkjOKRKUo9x8yP6Vzt7U5gkymiDh5nXrkZA8fQJoGKsy/E+o2CDNYQ4=")
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
		propertiRoutes.POST("/landingPage", propertiController.LandingPage)
	}

	testRoutes := r.Group("/test")
	{
		testRoutes.POST("", testController.TestEmail)
	}

	r.Run(":7177")
}
