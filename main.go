package main

import (
	"fmt"

	"github.com/samber/lo"
	"goftp.io/server/v2"
)

func main() {
	serverOptions := &server.Options{
		Name:   "Custom FTP Collector",
		Driver: NewDriver(),
		Port:   10021,
		Auth:   NewFakeAuth(),
		Perm:   server.NewSimplePerm("theuser", "thegroup"),
	}

	ftpServer := lo.Must(server.NewServer(serverOptions))
	fmt.Println("Starting...")
	lo.Must0(ftpServer.ListenAndServe())
}
