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

// LogWriter :
type LogWriter struct {
	lock       sync.Mutex
	dir        string
	module     string
	rotateSize int64
	curDate    string
	fp         *os.File
	fpath      string
}

// NewLogWriter :
func NewLogWriter(dir string, module string, rotateSize int64) *LogWriter {
	return &LogWriter{module: module, dir: dir, rotateSize: rotateSize}
}

func logIndex(fname string) int {
	if matches, _ := regexp.MatchString("\\[[0-9]+\\]", fname); matches {
		s := fname[strings.Index(fname, "[")+1 : strings.Index(fname, "]")]
		i, _ := strconv.ParseInt(s, 10, 32)
		return int(i)
	}
	return 0
}

// LogPath :
func LogPath(dir string, module string, maxFileSize int64, now time.Time) string {
	pre := now.Format("2006-01-02")
	monthDir := filepath.Join(dir, now.Format("2006-01"))
	// e.g. "2014-08-12[1]_example.log" or "2014-08-12_example.log"
	pattern := pre + "(\\[[0-9]+\\])?" + "_" + module + ".log"
	idx := 0
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		fname := filepath.Base(path)
		if matches, _ := regexp.MatchString(pattern, fname); matches {
			if idx <= logIndex(fname) {
				idx = logIndex(fname)
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
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.curDate != t.Format("2006-01-02") {
		w.closeFile()
		w.curDate = t.Format("2006-01-02")
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
	return w.WriteWithTime(output, time.Now())
}

func (w *LogWriter) closeFile() {
	w.fp.Close()
	w.fp = nil
	w.fpath = ""
}
