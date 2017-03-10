package libs

import (
	"log"
	"io"
)

var (
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
	Trace *log.Logger
)

func NewLog(traceHandle io.Writer, infoHandle io.Writer,warningHandle io.Writer,errorHandle io.Writer) {
	Info = log.New(infoHandle,"INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle,"WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle,"ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Trace = log.New(traceHandle,"TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
}