package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	MySQLUserName           = "system-mysql-username"
	MySQLPassword           = "system-mysql-password"
	MySQLDatabase           = "system-mysql-database"
	MySQLHost               = "system-mysql-host"
	MySQLPort               = "system-mysql-port"
	MySQLCharset            = "system-mysql-charset"
	MySQLLoc                = "system-mysql-loc"
	MySQLMaxOpenConnections = "system-mysql-max-open-conns"
	MySQLMaxIdleConnections = "system-mysql-max-idle-conns"

	RBACFilePath          = "system-rbac-file-path"
	AuthGrpcAddr          = "system-auth-grpc-addr"
	SystemGrpcAddr        = "system-grpc-addr"
	SystemGRPCGatewayAddr = "system-grpc-gw-addr"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "authen-go",
	Short: "authentication root cmd",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bm.yaml)")
	initConfiguration()

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".task")
		viper.SetEnvPrefix("task")
	}

	viper.SetEnvPrefix("task")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initConfiguration() {
	rootCmd.PersistentFlags().String("system-mode", "root", "mysql username")
	rootCmd.PersistentFlags().String("system-gorm-log-mode", "root", "mysql username")

	rootCmd.PersistentFlags().String(MySQLUserName, "root", "mysql username")
	rootCmd.PersistentFlags().String(MySQLPassword, "", "mysql password")
	rootCmd.PersistentFlags().String(MySQLDatabase, "task_service", "mysql database name")
	rootCmd.PersistentFlags().String(MySQLHost, "127.0.0.1", "mysql host")
	rootCmd.PersistentFlags().String(MySQLPort, "3306", "mysql port")
	rootCmd.PersistentFlags().String(MySQLCharset, "utf8mb4", "mysql default database character set. Recommend utf8mb4 for better Unicode support")
	rootCmd.PersistentFlags().String(MySQLLoc, "Local", "mysql password")
	rootCmd.PersistentFlags().String(MySQLMaxOpenConnections, "20", "mysql SetMaxOpenConns")
	rootCmd.PersistentFlags().String(MySQLMaxIdleConnections, "2", "mysql SetMaxIdleConns")

	rootCmd.PersistentFlags().String(RBACFilePath, "cmd/config/rbac.conf", "rbac config file path")
	rootCmd.PersistentFlags().String(AuthGrpcAddr, ":8080", "address of auth grpc server")
	rootCmd.PersistentFlags().String(SystemGrpcAddr, ":9090", "address of auth grpc server")
	rootCmd.PersistentFlags().String(SystemGRPCGatewayAddr, ":9091", "address of auth grpc server")

	//Bind flags to viper
	_ = viper.BindPFlag("system-mode", rootCmd.PersistentFlags().Lookup("system-mode"))
	_ = viper.BindPFlag("system-gorm-log-mode", rootCmd.PersistentFlags().Lookup("system-gorm-log-mode"))

	_ = viper.BindPFlag(MySQLUserName, rootCmd.PersistentFlags().Lookup(MySQLUserName))
	_ = viper.BindPFlag(MySQLPassword, rootCmd.PersistentFlags().Lookup(MySQLPassword))
	_ = viper.BindPFlag(MySQLDatabase, rootCmd.PersistentFlags().Lookup(MySQLDatabase))
	_ = viper.BindPFlag(MySQLHost, rootCmd.PersistentFlags().Lookup(MySQLHost))
	_ = viper.BindPFlag(MySQLPort, rootCmd.PersistentFlags().Lookup(MySQLPort))
	_ = viper.BindPFlag(MySQLCharset, rootCmd.PersistentFlags().Lookup(MySQLCharset))
	_ = viper.BindPFlag(MySQLLoc, rootCmd.PersistentFlags().Lookup(MySQLLoc))
	_ = viper.BindPFlag(MySQLMaxOpenConnections, rootCmd.PersistentFlags().Lookup(MySQLMaxOpenConnections))
	_ = viper.BindPFlag(MySQLMaxIdleConnections, rootCmd.PersistentFlags().Lookup(MySQLMaxIdleConnections))

	_ = viper.BindPFlag(RBACFilePath, rootCmd.PersistentFlags().Lookup(RBACFilePath))
	_ = viper.BindPFlag(AuthGrpcAddr, rootCmd.PersistentFlags().Lookup(AuthGrpcAddr))
	_ = viper.BindPFlag(SystemGrpcAddr, rootCmd.PersistentFlags().Lookup(SystemGrpcAddr))
	_ = viper.BindPFlag(SystemGRPCGatewayAddr, rootCmd.PersistentFlags().Lookup(SystemGRPCGatewayAddr))
}
