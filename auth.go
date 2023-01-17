package main

import "goftp.io/server/v2"

func NewFakeAuth() server.Auth {
	return &fakeAuth{}
}

type fakeAuth struct{}

func (fa *fakeAuth) CheckPasswd(ftpContext *server.Context, username string, password string) (bool, error) {
	return true, nil
}
