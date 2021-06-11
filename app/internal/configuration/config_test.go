package configuration

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

const (
	testAfsFolder = "test"
	testAfsFile   = testAfsFolder + "/config.json"
)

type osFSMock struct{}

var (
	appFsMem = afero.NewMemMapFs()
	afsMem   = &afero.Afero{Fs: appFsMem}
	config   = Config{
		Server: Server{
			Address: "127.0.0.1",
			Port:    "8080",
		},
		RedisDB: RedisDB{
			Address: "127.0.0.1",
			Port:    "6379",
		},
	}
	errPrinter = func(err error) {
		if err != nil {
			log.Println("Caught error: ", err.Error())
		}
	}
	fsTest fileSystem = osFSMock{}
)

func (osFSMock) Open(name string) (ifile, error) {
	return afsMem.Open(name)
}
func (osFSMock) Stat(name string) (os.FileInfo, error) {
	return afsMem.Stat(name)
}

func (osFSMock) ReadFile(name string) ([]byte, error) {
	return afsMem.ReadFile(name)
}

func init() {
	err := afsMem.MkdirAll(testAfsFolder, 0755)
	errPrinter(err)

	file, err := afsMem.Create(testAfsFile)
	defer func() {
		err := file.Close()
		errPrinter(err)
	}()
	errPrinter(err)

	err = json.NewEncoder(file).Encode(config)
	errPrinter(err)
}

func TestReadConfig(t *testing.T) {
	b, err := afsMem.ReadFile(testAfsFile)
	errPrinter(err)
	var config Config

	err = json.Unmarshal(b, &config)
	errPrinter(err)
	readConfig, err := ReadConfig(testAfsFile, fsTest)
	errPrinter(err)
	if !reflect.DeepEqual(config, *readConfig) {
		t.Errorf("JSON files are not equal. Have: %v, want: %v", *readConfig, config)
	}
}
