package cmd

import (
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
	"github.com/TranTheTuan/task-service/app/infrastructure/logger"
	"github.com/TranTheTuan/task-service/wire"
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
	logger := logger.NewLogger(logrus.Fields{
		"auth_grpc_addr": viper.GetString(AuthGrpcAddr),
	})
	// defer closeFn()
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
		authUsecase, err := wire.InitAuthUsecase(viper.GetString(AuthGrpcAddr))
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_logrus.UnaryServerInterceptor(logger),
				grpc_auth.UnaryServerInterceptor(authUsecase.AuthHandler),
			)),
		)

		taskServiceServer := wire.InitTaskServiceServer(orm)

		pbTasks.RegisterTaskCreateServiceServer(grpcServer, taskServiceServer)
		pbTasks.RegisterTaskUpdateServiceServer(grpcServer, taskServiceServer)
		pbTasks.RegisterTaskDeleteServiceServer(grpcServer, taskServiceServer)
		pbTasks.RegisterTaskGetAllServiceServer(grpcServer, taskServiceServer)
		pbTasks.RegisterTaskGetByIDServiceServer(grpcServer, taskServiceServer)

		grpcAddr := viper.GetString(SystemGrpcAddr)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.WithError(err).Fatalln("listen to grpc address failed")
			panic(err)
		}
		logger.Info("serving gRPC on 0.0.0.0:9090")
		err = grpcServer.Serve(lis)
		defer func() {
			err = lis.Close()
			if err != nil {
				logger.WithError(err).Error("close grpc server failed")
			}
		}()
		if err != nil {
			panic(err)
		}
	}()
	<-c
	logger.Info("server graceful shutdown")
}
