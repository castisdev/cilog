package cilog

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type logMsg struct {
	output []byte
	t      time.Time
}

// LogWriter :
type LogWriter struct {
	lock       sync.Mutex
	dir        string
	module     string
	rotateSize int64
	curYearDay int
	fp         *os.File
	fpath      string
	queue      chan logMsg
}

// NewLogWriter :
func NewLogWriter(dir string, module string, rotateSize int64) *LogWriter {
	return &LogWriter{module: module, dir: dir, rotateSize: rotateSize}
}

func logIndex(fname string) int {
	start := strings.Index(fname, "[")
	end := strings.Index(fname, "]")
	if start == -1 || end == -1 {
		return 0
	}
	s := fname[start+1 : end]
	i, _ := strconv.ParseInt(s, 10, 32)
	return int(i)
}

// LogPath :
func LogPath(dir string, module string, maxFileSize int64, now time.Time) string {
	pre := now.Format("2006-01-02")
	monthDir := filepath.Join(dir, now.Format("2006-01"))
	// e.g. "2014-08-12[1]_example.log" or "2014-08-12_example.log"
	pattern := pre + "(\\[[0-9]+\\])?" + "_" + module + ".log"
	regx, _ := regexp.Compile(pattern)
	idx := 0
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		fname := filepath.Base(path)
		if regx.MatchString(fname) {
			curIdx := logIndex(fname)
			if idx <= curIdx {
				idx = curIdx
				if f.Size() > maxFileSize {
					idx++
				}
			}
		}
		return nil
	})

	log := pre + "_" + module + ".log"
	if idx > 0 {
		log = pre + "[" + strconv.Itoa(idx) + "]_" + module + ".log"
	}
	return filepath.Join(monthDir, log)
}

// WriteWithTime :
func (w *LogWriter) WriteWithTime(output []byte, t time.Time) (int, error) {
	if w.queue == nil {
		w.lock.Lock()
		defer w.lock.Unlock()
	}
	if w.curYearDay != t.YearDay() {
		w.closeFile()
		w.curYearDay = t.YearDay()
	}

	if _, err := os.Stat(w.fpath); os.IsNotExist(err) {
		w.closeFile()
	}

	if w.fp == nil {
		p := LogPath(w.dir, w.module, w.rotateSize, t)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			return 0, err
		}
		f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		w.fp = f
		w.fpath = p

		abspath, _ := filepath.Abs(p)
		symfilepath := filepath.Join(filepath.Dir(filepath.Dir(abspath)), w.module+".log")
		if _, err := os.Lstat(symfilepath); err == nil {
			os.Remove(symfilepath)
		}
		os.Symlink(abspath, symfilepath)
	}

	n, err := w.fp.Write(output)

	size, err := w.fp.Seek(0, os.SEEK_CUR)
	if err != nil {
		return 0, err
	}
	if size > w.rotateSize {
		w.closeFile()
	}
	return n, err
}

// Write :
func (w *LogWriter) Write(output []byte) (int, error) {
	if w.queue == nil {
		return w.WriteWithTime(output, time.Now())
	}

	w.queue <- logMsg{output, time.Now()}
	return len(output), nil
}

func (w *LogWriter) closeFile() {
	w.fp.Close()
	w.fp = nil
	w.fpath = ""
}

// Start :
func (w *LogWriter) Start() {
	w.StartWithBufferSize(1024)
}

// StartWithBufferSize :
func (w *LogWriter) StartWithBufferSize(size int) {
	w.queue = make(chan logMsg, size)
	go w.serve()
}

// Stop :
func (w *LogWriter) Stop() {
	close(w.queue)
}

func (w *LogWriter) serve() {
	for {
		select {
		case msg, ok := <-w.queue:
			if !ok {
				return
			}
			w.WriteWithTime(msg.output, msg.t)
		}
	}
}
