package storage

import (
	"io"
	"os"
)

type FileInfo struct {
	Id         uint64
	ByteReader io.Reader
	Ext        Extension
	os.File
}
