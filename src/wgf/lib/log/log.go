//Copyright 2014 walu
//Use of this source code is governed by a MIT
//license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	LEVEL_DEBUG int	= 1
	LEVEL_INFO int	= 2
	LEVEL_WARNING int = 3
	LEVEL_FATAL int	= 4
)

var (
	ConfTimeLayout string = time.RFC3339
	ConfTimeLocationName string = "Asia/Shanghai"
	ConfLogWriter io.Writer = os.Stdout
	ConfMinLogLevel = 2
)

type Logger struct {
	timeLayout string
	timeLocation *time.Location
	logWriter io.Writer
	minLogLevel int
}

func NewLogger() *Logger {
	l := new(Logger)
	l.SetMinLogLevel(ConfMinLogLevel)
	l.SetLogWriter(ConfLogWriter)
	l.SetTimeLocation(ConfTimeLocationName)
	l.SetTimeLayout(ConfTimeLayout)
	return l
}

func (l *Logger) SetMinLogLevel(level int) {
	l.minLogLevel = level
}

func (l *Logger) SetLogWriter(logWriter io.Writer) {
	l.logWriter = logWriter
}

//设置时区，参数错误则使用Asia/Shanghai
func (l *Logger) SetTimeLocation(name string) {
	var err error
	if l.timeLocation, err = time.LoadLocation(name); nil!=err {
		l.timeLocation, _ = time.LoadLocation("Asia/Shanghai")
		l.Warningf("time_location_error %s, use Asia/Shanghai instead", err.Error())
	}
}

func (l *Logger) SetTimeLayout(layout string) {
	l.timeLayout = layout
}

func (l Logger) LogWriter() io.Writer {
	return l.LogWriter
}

func (l Logger) Debug(v ...interface{}) {
	if l.minLogLevel > LEVEL_DEBUG {
		return
	}

	l.output("DEBUG", v...)
}

func (l Logger) Debugf(format string, v ...interface{}) {
	if l.minLogLevel > LEVEL_DEBUG {
		return
	}

	l.outputf("DEBUG", format, v...)
}

func (l Logger) Info(v ...interface{}) {
	if l.minLogLevel > LEVEL_INFO {
		return
	}

	l.output("INFO", v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	if l.minLogLevel > LEVEL_INFO {
		return
	}

	l.outputf("INFO", format, v...)
}

func (l Logger) Warning(v ...interface{}) {
	if l.minLogLevel > LEVEL_WARNING {
		return
	}

	l.output("WARNING", v...)
}

func (l Logger) Warningf(format string, v ...interface{}) {
	if l.minLogLevel > LEVEL_WARNING {
		return
	}

	l.outputf("WARNING", format, v...)
}

//call os.Exit(1) after log was writted
func (l Logger) Fatal(v ...interface{}) {
	if l.minLogLevel > LEVEL_FATAL {
		return
	}

	l.output("FATAL", v...)
	os.Exit(1)
}

//call os.Exit(1) after log was writted
func (l Logger) Fatalf(format string, v ...interface{}) {
	if l.minLogLevel > LEVEL_FATAL {
		return
	}

	l.outputf("FATAL", format, v...)
	os.Exit(1)
}

func (l Logger) Print(v ...interface{}) {
	fmt.Fprintln(l.logWriter, v...)
}

func (l Logger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(l.logWriter, format+"\n", v...)
}

func (l Logger) output(level string, log ...interface{}) {
	fmt.Fprintf(l.logWriter, "%s %s %s\n", time.Now().In(l.timeLocation).Format(l.timeLayout), level, fmt.Sprint(log...))
}

func (l Logger) outputf(level string, format string, v ...interface{}) {
	fmt.Fprintf(l.logWriter, "%s %s %s\n", time.Now().In(l.timeLocation).Format(l.timeLayout), level, fmt.Sprintf(format, v...))
}

func init() {
	ConfLogWriter = os.Stdout
}

