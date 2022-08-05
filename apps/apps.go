package apps

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

type MockLoggerFuncs struct {
	Trace func(a ...any)
	Debug func(a ...any)
	Info  func(a ...any)
	Warn  func(a ...any)
	Error func(a ...any)
	Fatal func(a ...any)
}

var (
	TraceLogger *log.Logger
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	FatalLogger *log.Logger
	Logs        MockLoggerFuncs
)

func logInit(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	TraceLogger = log.New(traceHandle, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	DebugLogger = log.New(infoHandle, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	InfoLogger = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	WarnLogger = log.New(warningHandle, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	ErrorLogger = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	FatalLogger = log.New(errorHandle, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	Logs.Trace = TraceLogger.Println
	Logs.Debug = DebugLogger.Println
	Logs.Info = InfoLogger.Println
	Logs.Warn = WarnLogger.Println
	Logs.Error = ErrorLogger.Println
	Logs.Fatal = FatalLogger.Println
}

func init() {
	logInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}
