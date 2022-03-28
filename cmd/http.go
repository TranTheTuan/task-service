package cmd

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "start the back order service gateway",
	Run:   runServeHTTPCmd,
}

func init() {
	serveCmd.AddCommand(httpCmd)
}

func runServeHTTPCmd(cmd *cobra.Command, args []string) {
	logger := log.Default()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	d := initDB()
	mysqlDsn := d.ToDSN()
	orm, err := gorm.Open("mysql", mysqlDsn)
	if err != nil {
		panic(err)
	}

	maxOpenConnections := viper.GetInt(MySQLMaxOpenConnections)
	maxIdleConnections := viper.GetInt(MySQLMaxIdleConnections)

	orm.DB().SetMaxOpenConns(maxOpenConnections)
	orm.DB().SetMaxIdleConns(maxIdleConnections)
	orm.DB().SetConnMaxLifetime(200 * time.Minute)

	go func() {
		// casbin.InitFromSQLLite(orm, viper.GetString(RBACFilePath))
		// pubsub.InitPubSub(orm)

		// userRepo := repository.NewUserRepository(orm)
		// userService := service.NewUserService(userRepo)
		// userUsecase := usecase.NewUserUsecase(userService)
		// authorUsecase := usecase.NewAuthorUsecase()
		// authHandler := httpHandler.NewAuthHandler(userUsecase, authorUsecase)

		router := mux.NewRouter().PathPrefix("/v1/auth/").Subrouter()
		// router.Use(pubsub.EventDispatcherMiddleware)
		// router.HandleFunc("/login", authHandler.Login).Methods("POST")
		// router.HandleFunc("/register", authHandler.Register).Methods("POST")

		// router1 := mux.NewRouter().PathPrefix("/v1/tasks/").Subrouter()
		// router1.HandleFunc("/{id:[0-9]+}", authHandler.TestAuthorize).Methods("POST")

		httpMux := http.NewServeMux()
		httpMux.Handle("/v1/tasks/", router)

		httpHandler := cors.AllowAll().Handler(httpMux)

		srv := &http.Server{
			Addr:         ":8080",
			Handler:      httpHandler,
			IdleTimeout:  60 * time.Second,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		logger.Print("server started")
		err = srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	<-c
	logger.Print("server graceful shutdown")
}
