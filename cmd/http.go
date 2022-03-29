package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pbTasks "github.com/TranTheTuan/pbtypes/build/go/tasks"
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
		gwMux := runtime.NewServeMux(
			runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
				md := make(map[string]string)
				if method, ok := runtime.RPCMethod(ctx); ok {
					md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
				}
				if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
					md["pattern"] = pattern // /v1/example/login
				}
				return metadata.New(md)
			}),
		)
		opts := []grpc.DialOption{grpc.WithInsecure()}
		grpcAddr := viper.GetString(SystemGrpcAddr)

		err = pbTasks.RegisterTaskServiceHandlerFromEndpoint(context.Background(), gwMux, grpcAddr, opts)

		httpMux := http.NewServeMux()
		httpMux.Handle("/", gwMux)
		gwAddr := viper.GetString(SystemGRPCGatewayAddr)
		httpHandler := cors.AllowAll().Handler(httpMux)

		srv := &http.Server{
			Addr:         gwAddr,
			Handler:      httpHandler,
			IdleTimeout:  60 * time.Second,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		logger.Println("Serving gRPC gateway on 0.0.0.0:9091")
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	<-c
	logger.Print("server graceful shutdown")
}
