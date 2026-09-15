package handler

import (
	"fp-kpl/application/service"
	"fp-kpl/domain/order"
	"fp-kpl/domain/transaction"
	"fp-kpl/infrastructure/adapter/payment_gateway"
	"fp-kpl/infrastructure/database/config"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/repository"
	"fp-kpl/presentation/controller"
	"fp-kpl/presentation/middleware"
	"fp-kpl/presentation/route"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	app  *gin.Engine
	once sync.Once
)

func initApp() {
	once.Do(func() {
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

		app = gin.Default()
		app.Use(middleware.CORSMiddleware())

		app.Static("/assets", "./assets")

		if os.Getenv("IS_LOGGER") == "true" {
			route.LoggerRoute(app)
		}

		route.UserRoute(app, userController, jwtService)
		route.TableRoute(app, tableController, jwtService)
		route.CategoryRoute(app, categoryController, jwtService)
		route.MenuRoute(app, menuController, jwtService, userService)
		route.TransactionRoute(app, transactionController, jwtService, userService)
		route.OrderRoute(app, orderController, jwtService)
	})
}

// Handler harus ber-package handler dan mengekspor fungsi Handler
func Handler(w http.ResponseWriter, r *http.Request) {
	initApp()
	app.ServeHTTP(w, r)
}