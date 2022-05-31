package frm

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

const (
	Console = 0
	File    = 1
	Both    = 2
)

// Console = 0
// File    = 1
// Console + File = 2
func InitLogger(output int) error {
	switch output {
	case Console:
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05Z07:00"}
		Logger = zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
		return nil
	case File:
		fi := &lumberjack.Logger{
			Filename: "Logs.log",
			MaxSize:  10,
			Compress: true,
		}

		Logger = zerolog.New(fi).With().Timestamp().Caller().Logger()
		return nil

	case Both:
		fi := &lumberjack.Logger{
			Filename: "Logs.log",
			MaxSize:  10,
			Compress: true,
		}

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05Z07:00"}

		multi := zerolog.MultiLevelWriter(consoleWriter, fi)

		Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()
		return nil

	default:
		return errors.New("wrong argument")
	}

}
