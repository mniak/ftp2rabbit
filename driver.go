package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"goftp.io/server/v2"
)

func NewDriver() server.Driver {
	return &queueDriver{}
}

type queueDriver struct{}

func (drv *queueDriver) Stat(ftpContext *server.Context, path string) (os.FileInfo, error) {
	fmt.Println("Stat", ftpContext.Cmd, path)
	if strings.HasSuffix(path, "/") {
		return SimpleFileInfo{
			name:  path,
			isDir: true,
		}, nil
	}
	return nil, os.ErrNotExist
}

func (drv *queueDriver) ListDir(ftpContext *server.Context, path string, fileInfoFn func(os.FileInfo) error) error {
	fmt.Println("ListDir")
	return nil
}

func (drv *queueDriver) DeleteDir(ftpContext *server.Context, path string) error {
	fmt.Println("DeleteDir")
	return nil
}

func (drv *queueDriver) DeleteFile(ftpContext *server.Context, path string) error {
	fmt.Println("DeleteFile")
	return nil
}

func (drv *queueDriver) Rename(ftpContext *server.Context, from string, to string) error {
	fmt.Println("Rename")
	return nil
}

func (drv *queueDriver) MakeDir(ftpContext *server.Context, path string) error {
	fmt.Println("MakeDir")
	return nil
}

func (drv *queueDriver) GetFile(ftpContext *server.Context, path string, filepos int64) (int64, io.ReadCloser, error) {
	fmt.Println("GetFile")
	return 0, nil, nil
}

func (drv *queueDriver) PutFile(ftpContext *server.Context, dstPath string, fileReader io.Reader, _ int64) (int64, error) {
	fmt.Println("PutFile")
	fileData, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return 0, err
	}
	fmt.Println("  ", string(fileData))

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return 0, errors.WithMessage(err, "failed to connect to rabbit mq")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return 0, errors.WithMessage(err, "failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"ftp-integration", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to declare a queue")
	}

	err = ch.PublishWithContext(context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fileData),
		},
	)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to publish a message")
	}

	return 0, nil
}
