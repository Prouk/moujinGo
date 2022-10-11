package bin

import (
	"os"
	"time"
)

type Logger struct {
	ToConsole bool
	FilePath  string
	Level     int
	File      *os.File
}

var Reset = "\033[0m"

func (Logger) InitLogger(toConsole bool, path string, level int) (Logger, error) {
	var logger Logger
	var err error
	logger.FilePath = path
	logger.Level = level
	if toConsole {
		logger.ToConsole = true
	} else {
		logger.File, err = os.Create(path)
		defer logger.File.Close()
	}
	return logger, err
}

func (l *Logger) PassLog(msg string, level int) {
	if level > l.Level {
		return
	}
	currentTime := time.Now()
	if l.ToConsole {
		println("\033[32m" + currentTime.Format("2006-01-02 15:04:05") + ` - OK : ` + msg + Reset)
	} else {
		l.File.WriteString(currentTime.Format("2006-01-02 15:04:05") + ` - OK : ` + msg + `\n`)
	}
}

func (l *Logger) WarningLog(msg string, level int) {
	if level > l.Level {
		return
	}
	currentTime := time.Now()
	if l.ToConsole {
		println("\033[33m" + currentTime.Format("2006-01-02 15:04:05") + ` - WARNING : ` + msg + Reset)
	} else {
		l.File.WriteString(currentTime.Format("2006-01-02 15:04:05") + ` - WARNING : ` + msg + `\n`)
	}
}

func (l *Logger) ErrorLog(msg string, level int) {
	if level > l.Level {
		return
	}
	currentTime := time.Now()
	if l.ToConsole {
		println("\033[31m" + currentTime.Format("2006-01-02 15:04:05") + ` - ERROR : ` + msg + Reset)
	} else {
		l.File.WriteString(currentTime.Format("2006-01-02 15:04:05") + ` - ERROR : ` + msg + `\n`)
	}
}

func (l *Logger) CommentLog(msg string, level int) {
	if level > l.Level {
		return
	}
	currentTime := time.Now()
	if l.ToConsole {
		println("\033[37m" + currentTime.Format("2006-01-02 15:04:05") + ` - COMMENT : ` + msg + Reset)
	} else {
		l.File.WriteString(currentTime.Format("2006-01-02 15:04:05") + ` - COMMENT : ` + msg + `\n`)
	}
}
