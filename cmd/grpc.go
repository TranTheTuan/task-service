package cmd

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	pbTasks "github.com/TranTheTuan/pbtypes/build/go/tasks"
	"github.com/TranTheTuan/task-service/app/domain/service"
	"github.com/TranTheTuan/task-service/app/domain/usecase"
	internalGrpc "github.com/TranTheTuan/task-service/app/infrastructure/grpc"
	"github.com/TranTheTuan/task-service/app/infrastructure/grpc/client"
	"github.com/TranTheTuan/task-service/app/infrastructure/repository"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "start grpc server",
	Run:   runServeGRPCCmd,
}

func init() {
	serveCmd.AddCommand(grpcCmd)
}

func runServeGRPCCmd(cmd *cobra.Command, args []string) {
	logger := log.Default()
	lgg := logrus.WithFields(logrus.Fields{
		"auth_grpc_addr": viper.GetString(AuthGrpcAddr),
	})
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
	orm.LogMode(true)

	go func() {
		authClient, err := client.NewAuthClient(viper.GetString(AuthGrpcAddr))
		if err != nil {
			panic(err)
		}
		authUsecase := usecase.NewAuthUsecase(authClient)

		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_logrus.UnaryServerInterceptor(lgg),
				grpc_auth.UnaryServerInterceptor(authUsecase.AuthHandler),
			)),
		)

		taskRepo := repository.NewTaskRepository(orm)
		taskService := service.NewTaskService(taskRepo)
		taskUsecase := usecase.NewTaskUsecase(taskService)
		taskServiceServer := internalGrpc.NewTaskServiceServer(taskUsecase)

		pbTasks.RegisterTaskServiceServer(grpcServer, taskServiceServer)

		grpcAddr := viper.GetString(SystemGrpcAddr)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Fatalln("Failed to listen:", err)
			panic(err)
		}
		logger.Println("Serving gRPC on 0.0.0.0:9090")
		err = grpcServer.Serve(lis)
		defer func() {
			err = lis.Close()
			if err != nil {
				logger.Fatal(err)
			}
		}()
		if err != nil {
			panic(err)
		}
	}()
	<-c
	logger.Print("server graceful shutdown")
}
