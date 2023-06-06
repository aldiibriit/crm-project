package main

import (
	"go-api/config"
	"go-api/controller"
	"go-api/middleware"
	externalRepository "go-api/repository/external-repo"
	internalRepository "go-api/repository/internal-repo"
	"go-api/service"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB      = config.SetupDatabaseConnection()
	dbWithBegin *gorm.DB      = config.SetupDatabaseConnection()
	minioClient *minio.Client = config.SetupMinioConnection()

	userRepository        internalRepository.UserRepository          = internalRepository.NewUserRepository(db)
	bookRepository        internalRepository.BookRepository          = internalRepository.NewBookRepository(db)
	crmRepository         internalRepository.CRMRepository           = internalRepository.NewCRMRepository(db)
	cctvRepository        internalRepository.CCTVRepository          = internalRepository.NewCCTVRepository(db)
	logActivityRepository internalRepository.LogActivityRepository   = internalRepository.NewLogActivityRepository(db)
	prestagingRepository  internalRepository.PrestagingCRMRepository = internalRepository.NewPrestagingCRMRepository(db)
	baseRepository        internalRepository.BaseRepository          = internalRepository.NewBaseRepository(db)
	stagingRepository     internalRepository.StagingCRMRepository    = internalRepository.NewStagingCRMRepository(db)
	minioRepository       externalRepository.MinioRepository         = externalRepository.NewMinioRepository(minioClient)

	jwtService           service.JWTService           = service.NewJWTService()
	userService          service.UserService          = service.NewUserService(userRepository)
	bookService          service.BookService          = service.NewBookService(bookRepository)
	authService          service.AuthService          = service.NewAuthService(userRepository, jwtService)
	crmService           service.CRMService           = service.NewCRMService(crmRepository)
	cctvService          service.CCTVService          = service.NewCCTVService(cctvRepository)
	minioService         service.MinioService         = service.NewMinioService(minioRepository)
	logActivityService   service.LogActivityService   = service.NewLogActivityService(logActivityRepository)
	baseService          service.BaseService          = service.NewBaseService(baseRepository)
	prestagingCRMService service.PrestagingCRMService = service.NewPrestagingCRMService(minioRepository, logActivityRepository, prestagingRepository, baseRepository, stagingRepository)
	qrCodeService        service.QrCodeService        = service.NewQrCodeService(minioRepository, logActivityRepository)

	authController          controller.AuthController          = controller.NewAuthController(authService, jwtService)
	userController          controller.UserController          = controller.NewUserController(userService, jwtService)
	bookController          controller.BookController          = controller.NewBookController(bookService, jwtService)
	crmController           controller.CRMController           = controller.NewCRMController(crmService, jwtService)
	cctvController          controller.CCTVController          = controller.NewCCTVController(cctvService, jwtService)
	logActivityController   controller.LogActivityController   = controller.NewLogActivityController(logActivityService)
	prestagingCRMController controller.PrestagingCRMController = controller.NewPrestagingCRMController(prestagingCRMService, jwtService)
	qrCodeController        controller.QrCodeController        = controller.NewQrCodeController(qrCodeService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user")
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

	crmRoutes := r.Group("api/crm", middleware.AuthorizeJWT(jwtService))
	{
		crmRoutes.GET("/all", crmController.GetAll)
	}

	cctvRoutes := r.Group("api/cctv", middleware.AuthorizeJWT(jwtService))
	{
		cctvRoutes.GET("/all", cctvController.GetAll)
		cctvRoutes.GET("/findBySn", cctvController.FindBySN)
	}

	logActivityRoutes := r.Group("api/log-activity", middleware.AuthorizeJWT(jwtService))
	{
		logActivityRoutes.POST("/getTimeline", logActivityController.GetTimeLine)
	}

	prestagingCRMRoutes := r.Group("api/prestaging-crm", middleware.AuthorizeJWT(jwtService))
	{
		prestagingCRMRoutes.POST("/post", prestagingCRMController.PostPrestaging)
		prestagingCRMRoutes.POST("/approve", prestagingCRMController.ApprovePrestaging)
		prestagingCRMRoutes.PUT("/reject", prestagingCRMController.RejectPrestaging)
		prestagingCRMRoutes.PUT("/reupload", prestagingCRMController.ReuploadPrestaging)
		prestagingCRMRoutes.GET("/getAllSubmittedData", prestagingCRMController.AllSubmittedData)
		prestagingCRMRoutes.POST("/getSubmittedDataBySn", prestagingCRMController.GetSubmittedDataBySn)
		prestagingCRMRoutes.POST("/getRejectedData", prestagingCRMController.GetRejectedData)
		prestagingCRMRoutes.POST("/post-prestaging-v2", prestagingCRMController.PostPrestagingV2)
	}

	qrCodeRoutes := r.Group("api/qr-code", middleware.AuthorizeJWT(jwtService))
	{
		qrCodeRoutes.POST("/generateQrCode", qrCodeController.GenerateQr)
	}

	r.Run(":7177")
}
