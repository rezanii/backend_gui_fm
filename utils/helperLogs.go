package utils

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)


/*
	initHandlers is a function to load config env,call function log and load arguments
    authored by Irma P 10/09/21
*/
func InitHandlers() {
	WriteLogFpmBackend()
	log.Info("Prepare Config.....")

	if len(os.Args) < 2 {
		log.Error( "empty argument")
	}
	env_file := os.Args[1]

	configEnv := Configuration{
		EnvName: env_file,
	}
	Param = configEnv
	log.Info("env file : ", configEnv.EnvName)
}
/*
	WriteLogEodApi is a function config log and create file log
    authored by Irma P 10/09/21
*/
func WriteLogFpmBackend() {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
	)

	fileLogs := "logs/" + "fpm-backend" + (strings.Replace(formatted, "-", "", -1)) + ".log"
	file, err := os.OpenFile(fileLogs, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)
	}
	log.SetReportCaller(true)

	log.SetFormatter(&Formatter{
		NoFieldsSpace: true,
		TrimMessages: true,
		ShowFullLevel: true,
		NoColors:      false,
		CallerFirst:   true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			filename := path.Base(f.File)
			return fmt.Sprintf("[%s:%d]", filename, f.Line)
		},
		TimestampFormat: "[2006-01-02 15:04:05]",

	})

	log.SetOutput(file)
}
/*
	ReadBody is a function read body request from user
    authored by Irma P 10/09/21
*/
func ReadBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}

type BodyLogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}
/*
	Write is a function implement from fuction BodyLogWriter to write respon
    authored by Irma P 10/09/21
*/
func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}


