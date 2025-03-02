package logger

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
)

var logFile *os.File

func GetFile() *os.File {
	if logFile == nil {
		os.MkdirAll("./tmp", os.ModePerm)

		var err error
		logFile, err = os.OpenFile("./tmp/app.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Panic(err)
		}
	}

	return logFile
}

func GetHandlers() slog.Handler {
	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: "15:04:05",
		NoColor:    false,
	})

	fileHandler := slog.NewTextHandler(GetFile(), nil)

	return slogmulti.Fanout(consoleHandler, fileHandler)
}

func GetMultiWriter() io.Writer {
	return io.MultiWriter(os.Stdout, GetFile())
}

func SetUp() {
	logger := slog.New(GetHandlers())
	slog.SetDefault(logger)
}

func CleanUp() {
	if logFile != nil {
		logFile.Close()
	}
}
