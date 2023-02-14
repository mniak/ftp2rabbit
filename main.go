package main

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/viper"
	"goftp.io/server/v2"
)

func main() {
	viper.AutomaticEnv()

	viper.SetDefault("RabbitMQ_Host", "localhost")
	viper.SetDefault("RabbitMQ_Port", "5672")
	viper.SetDefault("RabbitMQ_Username", "guest")
	viper.SetDefault("RabbitMQ_Password", "guest")
	viper.SetDefault("ListenPort", "10021")
	viper.SetDefault("FTP_Username", "ftp2rabbit")
	viper.SetDefault("FTP_Password", "password")

	rabbitHost := viper.GetString("RabbitMQ_Host")
	rabbitPort := viper.GetInt("RabbitMQ_Port")
	rabbitUsername := viper.GetString("RabbitMQ_Username")
	rabbitPassword := viper.GetString("RabbitMQ_Password")
	listenPort := viper.GetInt("ListenPort")
	ftpUsername := viper.GetString("FTP_Username")
	ftpPassword := viper.GetString("FTP_Password")

	driver := lo.Must(NewDriver(rabbitHost, rabbitPort, rabbitUsername, rabbitPassword))
	defer driver.Close()

	serverOptions := &server.Options{
		Name:   "FTP Server",
		Driver: driver,
		Port:   listenPort,
		Auth:   NewFakeAuth(),
		Perm:   server.NewSimplePerm(ftpUsername, ftpPassword),
	}

	ftpServer := lo.Must(server.NewServer(serverOptions))
	fmt.Println("Starting...")
	lo.Must0(ftpServer.ListenAndServe())
}
