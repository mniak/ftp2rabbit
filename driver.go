package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"goftp.io/server/v2"
)

func (drv *queueDriver) Close() {
	if drv.connection != nil {
		drv.connection.Close()
	}
	if drv.channel != nil {
		drv.channel.Close()
	}
}

func NewDriver(host string, port int, username, password string) (result *queueDriver, err error) {
	result = new(queueDriver)
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	if verbose {
		fmt.Println("AMQP URL:", amqpURL)
	}
	result.connection, err = amqp.Dial(amqpURL)
	if err != nil {
		err = errors.WithMessage(err, "failed to connect to rabbit mq")
		return
	}

	result.channel, err = result.connection.Channel()
	if err != nil {
		err = errors.WithMessage(err, "failed to open a channel")
		return
	}

	result.queue, err = result.channel.QueueDeclare(
		"ftp-integration", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		err = errors.WithMessage(err, "failed to declare a queue")
		return
	}
	return
}

type queueDriver struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

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

type FileInfo struct {
	TraceID  uuid.UUID
	Contents []byte
	FileName string
}

func (drv *queueDriver) PutFile(ftpContext *server.Context, dstPath string, fileReader io.Reader, _ int64) (int64, error) {
	fmt.Println("PutFile")
	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		return 0, err
	}

	traceID, err := uuid.NewRandom()
	if err != nil {
		return 0, err
	}

	fileInfo := FileInfo{
		TraceID:  traceID,
		Contents: fileData,
		FileName: dstPath,
	}
	fileInfoBytes, err := json.Marshal(fileInfo)

	if err != nil {
		return 0, err
	}

	if verbose {
		fmt.Println("  ", string(fileInfoBytes))
	}

	err = drv.channel.PublishWithContext(context.Background(),
		"",             // exchange
		drv.queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        fileInfoBytes,
		},
	)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to publish a message")
	}

	return 0, nil
}
