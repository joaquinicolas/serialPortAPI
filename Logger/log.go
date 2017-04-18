package logger

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Error   *log.Logger
	Warning *log.Logger
	Info    *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	file, err := os.OpenFile("logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)

	}

	fileStdout := io.MultiWriter(file, os.Stdout)
	fileStdErr := io.MultiWriter(file, os.Stderr)

	Trace = log.New(fileStdout,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(fileStdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(fileStdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(fileStdErr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
