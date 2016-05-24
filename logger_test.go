package cilog_test

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/castisdev/cilog"
	"github.com/satori/go.uuid"
)

type stringWriter struct {
	writed string
}

func (f *stringWriter) Write(output []byte) (int, error) {
	f.writed += string(output)
	return len(output), nil
}

func TestLogger_Log(t *testing.T) {
	w := &stringWriter{}
	logger := cilog.New(w, "module", "1.0", cilog.DEBUG)
	logger.Log(1, cilog.DEBUG, "abc", time.Date(2009, 11, 23, 15, 21, 30, 123456000, time.Local))
	_, file, line, _ := runtime.Caller(0)
	file = file[strings.LastIndex(file, "/")+1:]

	expected := "module,1.0,2009-11-23,15:21:30.123456,debug," + file + ":" + strconv.Itoa(line-1) + ",abc\n"
	if w.writed != expected {
		t.Errorf("log expected %s, but %s", expected, w.writed)
	}
}

func TestLogger_SetGet(t *testing.T) {
	w := &stringWriter{}

	cilog.SetWriter(w)
	cilog.SetModule("module")
	cilog.SetModuleVer("1.0")
	cilog.SetMinLevel(cilog.DEBUG)

	cilog.GetWriter().Write([]byte("abc"))
	if w.writed != "abc" {
		t.Errorf("GetWriter error")
	}
	if cilog.GetModule() != "module" {
		t.Errorf("GetModule expected %s, but %s", "module", cilog.GetModule())
	}
	if cilog.GetModuleVer() != "1.0" {
		t.Errorf("GetModule expected %s, but %s", "1.0", cilog.GetModuleVer())
	}
	if cilog.GetMinLevel() != cilog.DEBUG {
		t.Errorf("GetModule expected %s, but %s", cilog.DEBUG, cilog.GetMinLevel())
	}
}

type dummyWriter struct{}

func (dummyWriter) Write(out []byte) (int, error) {
	return len(out), nil
}
func BenchmarkLogger_WithDummyWriter(b *testing.B) {
	cilog.Set(dummyWriter{}, "module", "1.0,", cilog.DEBUG)
	for n := 0; n < b.N; n++ {
		cilog.Info("this is log. line:%d", n)
	}
}

func BenchmarkLogger_WithLogWriter(b *testing.B) {
	dir := uuid.NewV4().String()
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	cilog.Set(cilog.NewLogWriter(dir, "module", 1024*1024), "module", "1.0,", cilog.DEBUG)
	for n := 0; n < b.N; n++ {
		cilog.Info("this is log. line:%d", n)
	}
}

func BenchmarkLogger_WithLogWriter_TwoGoroutines(b *testing.B) {
	dir := uuid.NewV4().String()
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	cilog.Set(cilog.NewLogWriter(dir, "module", 1024*1024), "module", "1.0,", cilog.DEBUG)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			cilog.Report("test1 log : %d", i)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			cilog.Report("test2 log : %d", i)
		}
		wg.Done()
	}()
	wg.Wait()
}
