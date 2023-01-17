package main

import (
	"os"
	"time"
)

type SimpleFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

func (f SimpleFileInfo) Name() string {
	return f.name
}

func (f SimpleFileInfo) Size() int64 {
	return f.size
}

func (f SimpleFileInfo) Mode() os.FileMode {
	return f.mode
}

func (f SimpleFileInfo) ModTime() time.Time {
	return f.modTime
}

func (f SimpleFileInfo) IsDir() bool {
	return f.isDir
}

func (f SimpleFileInfo) Sys() any {
	return f.sys
}
