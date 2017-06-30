package cilog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
	CRITICAL
)

// Output :
func (l Level) Output() string {
	m := map[Level]string{
		DEBUG:     "Debug",
		REPORT:    "Report",
		INFO:      "Information",
		SUCCESS:   "Success",
		WARNING:   "Warning",
		ERROR:     "Error",
		FAIL:      "Fail",
		EXCEPTION: "Exception",
		CRITICAL:  "Critical",
	}
	return m[l]
}

// String :
func (l Level) String() string {
	return l.Output()
}

// UnmarshalYAML :
func (l *Level) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	var err error
	if err = unmarshal(&s); err != nil {
		return err
	}
	if *l, err = LevelFromString(s); err != nil {
		return err
	}
	return nil
}

// MarshalYAML :
func (l Level) MarshalYAML() (interface{}, error) {
	switch l {
	case DEBUG:
		return "debug", nil
	case REPORT:
		return "report", nil
	case INFO:
		return "info", nil
	case SUCCESS:
		return "success", nil
	case WARNING:
		return "warning", nil
	case ERROR:
		return "error", nil
	case FAIL:
		return "fail", nil
	case EXCEPTION:
		return "exception", nil
	case CRITICAL:
		return "critical", nil
	}
	return "", errors.New("invalid log level")
}

// LevelFromString :
func LevelFromString(s string) (Level, error) {
	m := map[string]Level{
		"debug":     DEBUG,
		"report":    REPORT,
		"info":      INFO,
		"success":   SUCCESS,
		"warning":   WARNING,
		"error":     ERROR,
		"fail":      FAIL,
		"exception": EXCEPTION,
		"critical":  CRITICAL,
	}
	v, ok := m[s]
	if !ok {
		return DEBUG, fmt.Errorf("invalid level string [%s]", v.Output())
	}
	return v, nil
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

// Log : "module,1.0,2009-11-23,15:21:30.123456,Debug,package1::src.go:56,,this is a example"
func (l *Logger) Log(calldepth int, lvl Level, msg string, t time.Time) {
	if l.GetMinLevel() > lvl {
		return
	}
	timeStr := t.Format("2006-01-02,15:04:05.000000")
	var file string
	var line int
	var ok bool
	var pc uintptr
	var pkg string
	pc, file, line, ok = runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
		pkg = "???"
	} else {
		file = filepath.Base(file)
		pkg = PackageBase(runtime.FuncForPC(pc).Name())
	}

	m := l.GetModule() + "," + l.GetModuleVer() + "," + timeStr + "," +
		lvl.Output() + "," + pkg + "::" + file + ":" + strconv.Itoa(line) + ",," + msg
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

// StdLogger :
func StdLogger() *Logger {
	return std
}

// Debugf :
func Debugf(format string, v ...interface{}) {
	std.Log(2, DEBUG, fmt.Sprintf(format, v...), time.Now())
}

// Reportf :
func Reportf(format string, v ...interface{}) {
	std.Log(2, REPORT, fmt.Sprintf(format, v...), time.Now())
}

// Infof :
func Infof(format string, v ...interface{}) {
	std.Log(2, INFO, fmt.Sprintf(format, v...), time.Now())
}

// Successf :
func Successf(format string, v ...interface{}) {
	std.Log(2, SUCCESS, fmt.Sprintf(format, v...), time.Now())
}

// Warningf :
func Warningf(format string, v ...interface{}) {
	std.Log(2, WARNING, fmt.Sprintf(format, v...), time.Now())
}

// Errorf :
func Errorf(format string, v ...interface{}) {
	std.Log(2, ERROR, fmt.Sprintf(format, v...), time.Now())
}

// Failf :
func Failf(format string, v ...interface{}) {
	std.Log(2, FAIL, fmt.Sprintf(format, v...), time.Now())
}

// Exceptionf :
func Exceptionf(format string, v ...interface{}) {
	std.Log(2, EXCEPTION, fmt.Sprintf(format, v...), time.Now())
}

// Criticalf :
func Criticalf(format string, v ...interface{}) {
	std.Log(2, CRITICAL, fmt.Sprintf(format, v...), time.Now())
}

// PackageBase : funcName string format : runtime.FuncForPC(pc).Name()
func PackageBase(funcName string) string {
	pkgStart := strings.LastIndex(funcName, "/") + 1
	return funcName[pkgStart : strings.Index(funcName[pkgStart:], ".")+pkgStart]
}
