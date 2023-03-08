package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/viper"
	"goftp.io/server/v2"
)

var verbose bool

func init() {
	viper.SetDefault("Verbose", "false")

	viper.SetDefault("RabbitMQ_Host", "rabbitmq.portalbooking.svc.cluster.local")
	viper.SetDefault("RabbitMQ_Port", "5672")
	viper.SetDefault("RabbitMQ_Username", "guest")
	viper.SetDefault("RabbitMQ_Password", "guest")
	viper.SetDefault("ListenPort", "21")
	viper.SetDefault("FTP_Username", "ftp2rabbit")
	viper.SetDefault("FTP_Password", "password")
	viper.SetDefault("FTP_PassivePorts", "10000-10005")

	viper.AutomaticEnv()
	verbose = viper.GetBool("Verbose") || len(os.Args) > 1 && strings.EqualFold(os.Args[1], "--verbose")
}

func main() {
	rabbitHost := viper.GetString("RabbitMQ_Host")
	rabbitPort := viper.GetInt("RabbitMQ_Port")
	rabbitUsername := viper.GetString("RabbitMQ_Username")
	rabbitPassword := viper.GetString("RabbitMQ_Password")
	listenPort := viper.GetInt("ListenPort")
	ftpUsername := viper.GetString("FTP_Username")
	ftpPassword := viper.GetString("FTP_Password")
	ftpPassivePorts := viper.GetString("FTP_PassivePorts")

	if verbose {
		fmt.Println("RabbitHost:", rabbitHost)
		fmt.Println("RabbitPort:", rabbitPort)
		fmt.Println("RabbitUsername:", rabbitUsername)
		fmt.Println("RabbitPassword:", rabbitPassword)
		fmt.Println("ListenPort:", listenPort)
		fmt.Println("FtpUsername:", ftpUsername)
		fmt.Println("FtpPassword:", ftpPassword)
		fmt.Println("FtpPassivePorts:", ftpPassivePorts)
	}

	driver := lo.Must(NewDriver(rabbitHost, rabbitPort, rabbitUsername, rabbitPassword))
	defer driver.Close()

	serverOptions := &server.Options{
		Name:         "FTP Server",
		Driver:       driver,
		Port:         listenPort,
		Auth:         NewFakeAuth(),
		Perm:         server.NewSimplePerm(ftpUsername, ftpPassword),
		PassivePorts: ftpPassivePorts,
	}

	ftpServer := lo.Must(server.NewServer(serverOptions))
	fmt.Println("Starting...")
	lo.Must0(ftpServer.ListenAndServe())
}
