package configuration

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type Config struct {
	Server  Server        `json:"server"`
	RedisDB RedisDB       `json:"redisdb"`
	Timeout time.Duration `json:"shu_timeout"`
}

type Server struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type RedisDB struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type fileSystem interface {
	Open(name string) (ifile, error)
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
}

type ifile interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

type OsFileSystem struct{}

func (OsFileSystem) Open(name string) (ifile, error) {
	return os.Open(name)
}
func (OsFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
func (OsFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func ReadConfig(path string, fs fileSystem) (*Config, error) {
	b, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
