package main

import (
	"fp-kpl/application/service"
	"fp-kpl/command"
	"fp-kpl/domain/order"
	"fp-kpl/domain/transaction"
	"fp-kpl/infrastructure/adapter/payment_gateway"
	"fp-kpl/infrastructure/database/config"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/repository"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"
	"fp-kpl/presentation/route"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func args(db *gorm.DB) bool {
	if len(os.Args) > 1 {
		flag := command.Commands(db)
		return flag
	}
	return true
}

func run(server *gin.Engine) {
	// Hanya load static jika folder assets ada (Mencegah error di Vercel read-only filesystem)
	if _, err := os.Stat("./assets"); err == nil {
		server.Static("/assets", "./assets")
	}

	if os.Getenv("IS_LOGGER") == "true" {
		route.LoggerRoute(server)
	}

	// Deteksi port Vercel (PORT) atau fallback ke GOLANG_PORT / 8888
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("GOLANG_PORT")
		if port == "" {
			port = "8888"
		}
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "0.0.0.0:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func main() {
	db := config.SetUpDatabaseConnection()

	jwtService := service.NewJWTService()
	dbTransactionRepository := db_transaction.NewRepository(db)

	userRepository := repository.NewUserRepository(dbTransactionRepository)
	tableRepository := repository.NewTableRepository(dbTransactionRepository)
	categoryRepository := repository.NewCategoryRepository(dbTransactionRepository)
	menuRepository := repository.NewMenuRepository(dbTransactionRepository)
	orderRepository := repository.NewOrderRepository(dbTransactionRepository)
	transactionRepository := repository.NewTransactionRepository(dbTransactionRepository)

	transactionDomainService := transaction.NewService(transactionRepository)
	orderDomainService := order.NewService()

	paymentGateway := payment_gateway.NewMidtransAdapter(db, transactionDomainService)

	userService := service.NewUserService(userRepository, jwtService, dbTransactionRepository)
	tableService := service.NewTableService(tableRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	menuService := service.NewMenuService(menuRepository, categoryRepository)
	orderService := service.NewOrderService(orderRepository, menuRepository, orderDomainService)
	transactionService := service.NewTransactionService(transactionRepository, userRepository, tableRepository, orderRepository, menuRepository, transactionDomainService, paymentGateway, dbTransactionRepository, orderService)

	userController := controller.NewUserController(userService)
	tableController := controller.NewTableController(tableService)
	categoryController := controller.NewCategoryController(categoryService)
	menuController := controller.NewMenuController(menuService)
	transactionController := controller.NewTransactionController(transactionService)
	orderController := controller.NewOrderController(orderService)

	defer config.CloseDatabaseConnection(db)

	if !args(db) {
		return
	}

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// Endpoint root untuk cek status di Vercel
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Backend Kantara Go is running on Vercel!",
		})
	})

	route.UserRoute(server, userController, jwtService)
	route.TableRoute(server, tableController, jwtService)
	route.CategoryRoute(server, categoryController, jwtService)
	route.MenuRoute(server, menuController, jwtService, userService)
	route.TransactionRoute(server, transactionController, jwtService, userService)
	route.OrderRoute(server, orderController, jwtService)

	run(server)
}