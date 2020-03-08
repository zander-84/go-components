package CHelperFile

import (
	"log"
	"os"
	"strings"
	"time"
)

type File struct {
	timezZone string
}

func NewFile(timezZone string) interface{} {
	if timezZone == "" {
		timezZone = "Asia/Shanghai"
	}

	this := new(File)
	this.timezZone = timezZone
	return this
}

func (this *File) Create(path string, fileName string) (*os.File, error) {
	file, err := os.OpenFile(this.FullLogPath(path, fileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (this *File) OpenLogFile(path string, fileName string) error {
	f, err := this.Create(path, fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(f)
	log.Println("starting log...")
	return nil
}
func (this *File) FullLogPath(path string, fileName string) string {
	if strings.Count(fileName, ".") < 1 {
		fileName += ".log"
	}

	local, _ := time.LoadLocation(this.timezZone)
	return strings.TrimRight(path, "/") + "/" + time.Now().In(local).Format("2006-01-02") + "-" + fileName
}
