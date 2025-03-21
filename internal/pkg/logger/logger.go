package pkg_logger

import (
	"io"
	"log"
	"os"
)

// アプリケーションロガー
type AppLogger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	WarnLog  *log.Logger
	DebugLog *log.Logger
	TestLog  *log.Logger
}

// アプリケーションロガーのインスタンス化
func NewAppLogger() *AppLogger {
	return &AppLogger{
		InfoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		WarnLog:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		DebugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		TestLog:  log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// ログ設定の初期化
func (l *AppLogger) SetUpLogger() {
	if os.Getenv("TEST_MODE") == "true" {
		l.InfoLog.SetOutput(io.Discard)
		l.ErrorLog.SetOutput(io.Discard)
		l.WarnLog.SetOutput(io.Discard)
		l.DebugLog.SetOutput(io.Discard)
	}
}
