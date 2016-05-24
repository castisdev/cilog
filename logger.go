package cilog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Level :
type Level int

// Level enum
const (
	DEBUG Level = iota
	REPORT
	INFO
	SUCCESS
	WARNING
	ERROR
	FAIL
	EXCEPTION
)

func (l Level) String() string {
	m := map[Level]string{
		DEBUG:     "debug",
		REPORT:    "report",
		INFO:      "info",
		SUCCESS:   "success",
		WARNING:   "warning",
		ERROR:     "error",
		FAIL:      "fail",
		EXCEPTION: "exception",
	}
	return m[l]
}

// Logger :
type Logger struct {
	mu        sync.Mutex
	writer    io.Writer
	module    string
	moduleVer string
	minLevel  Level
}

// New :
func New(out io.Writer, module string, moduleVer string, minLevel Level) *Logger {
	return &Logger{writer: out, module: module, moduleVer: moduleVer, minLevel: minLevel}
}

// Set :
func (l *Logger) Set(out io.Writer, module string, moduleVer string, minLevel Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = out
	l.module = module
	l.moduleVer = moduleVer
	l.minLevel = minLevel
}

// SetWriter :
func (l *Logger) SetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = w
}

// SetModule :
func (l *Logger) SetModule(m string) {
	l.module = m
}

// SetModuleVer :
func (l *Logger) SetModuleVer(v string) {
	l.moduleVer = v
}

// SetMinLevel :
func (l *Logger) SetMinLevel(lvl Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = lvl
}

// GetWriter :
func (l *Logger) GetWriter() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.writer
}

// GetModule :
func (l *Logger) GetModule() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.module
}

// GetModuleVer :
func (l *Logger) GetModuleVer() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.moduleVer
}

// GetMinLevel :
func (l *Logger) GetMinLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.minLevel
}

// Log :
func (l *Logger) Log(calldepth int, lvl Level, msg string, t time.Time) {
	if l.GetMinLevel() > lvl {
		return
	}

	timeStr := t.Format("2006-01-02,15:04:05.000000")
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	m := l.GetModule() + "," + l.GetModuleVer() + "," + timeStr + "," +
		lvl.String() + "," + file + ":" + strconv.Itoa(line) + "," + msg
	if len(msg) == 0 || msg[len(msg)-1] != '\n' {
		m += "\n"
	}
	l.GetWriter().Write([]byte(m))
}

var std = New(os.Stderr, "", "", DEBUG)

// Set :
func Set(out io.Writer, module string, moduleVer string, minLevel Level) {
	std.Set(out, module, moduleVer, minLevel)
}

// SetWriter :
func SetWriter(w io.Writer) {
	std.SetWriter(w)
}

// SetModule :
func SetModule(v string) {
	std.SetModule(v)
}

// SetModuleVer :
func SetModuleVer(v string) {
	std.SetModuleVer(v)
}

// SetMinLevel :
func SetMinLevel(lvl Level) {
	std.SetMinLevel(lvl)
}

// GetWriter :
func GetWriter() io.Writer {
	return std.GetWriter()
}

// GetModule :
func GetModule() string {
	return std.GetModule()
}

// GetModuleVer :
func GetModuleVer() string {
	return std.GetModuleVer()
}

// GetMinLevel :
func GetMinLevel() Level {
	return std.GetMinLevel()
}

// Debug :
func Debug(format string, v ...interface{}) {
	std.Log(2, DEBUG, fmt.Sprintf(format, v...), time.Now())
}

// Report :
func Report(format string, v ...interface{}) {
	std.Log(2, REPORT, fmt.Sprintf(format, v...), time.Now())
}

// Info :
func Info(format string, v ...interface{}) {
	std.Log(2, INFO, fmt.Sprintf(format, v...), time.Now())
}

// Success :
func Success(format string, v ...interface{}) {
	std.Log(2, SUCCESS, fmt.Sprintf(format, v...), time.Now())
}

// Warning :
func Warning(format string, v ...interface{}) {
	std.Log(2, WARNING, fmt.Sprintf(format, v...), time.Now())
}

// Error :
func Error(format string, v ...interface{}) {
	std.Log(2, ERROR, fmt.Sprintf(format, v...), time.Now())
}

// Fail :
func Fail(format string, v ...interface{}) {
	std.Log(2, FAIL, fmt.Sprintf(format, v...), time.Now())
}

// Exception :
func Exception(format string, v ...interface{}) {
	std.Log(2, EXCEPTION, fmt.Sprintf(format, v...), time.Now())
}
