package cilog_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/castisdev/cilog"
	"github.com/google/uuid"
)

func TestLogPath_NoLog(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	expected := filepath.Join(dir, "2009-11", "2009-11-23_module.log")
	v := cilog.LogPath(dir, "module", 100, time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	if v != expected {
		t.Errorf("LogPath expected %s but %s", expected, v)
	}
}

func TestLogPath_ExistsIndex0Log(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	monthD := filepath.Join(dir, "2009-11")
	os.Mkdir(monthD, 0775)
	ioutil.WriteFile(filepath.Join(monthD, "2009-11-23_module.log"), []byte("test"), 0775)

	expected := filepath.Join(dir, "2009-11", "2009-11-23_module.log")
	v := cilog.LogPath(dir, "module", 100, time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	if v != expected {
		t.Errorf("LogPath expected %s but %s", expected, v)
	}
}

func TestLogPath_ExistsBigLog(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	monthD := filepath.Join(dir, "2009-11")
	os.Mkdir(monthD, 0775)

	const rotateSize = 4
	bigData := "12345" // length(5) is bigger than rotateSize(4)
	ioutil.WriteFile(filepath.Join(monthD, "2009-11-23_module.log"), []byte(bigData), 0775)

	expected := filepath.Join(dir, "2009-11", "2009-11-23[1]_module.log")
	v := cilog.LogPath(dir, "module", rotateSize, time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	if v != expected {
		t.Errorf("LogPath expected %s but %s", expected, v)
	}
}

func TestLogPath_ExistsIndex2Log(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	monthD := filepath.Join(dir, "2009-11")
	os.Mkdir(monthD, 0775)
	ioutil.WriteFile(filepath.Join(monthD, "2009-11-23[2]_module.log"), []byte("test"), 0775)

	expected := filepath.Join(dir, "2009-11", "2009-11-23[2]_module.log")
	v := cilog.LogPath(dir, "module", 100, time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	if v != expected {
		t.Errorf("LogPath expected %s but %s", expected, v)
	}
}

func TestLogWriter_Write(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	w := cilog.NewLogWriter(dir, "module", 5)
	w.WriteWithTime([]byte("abc"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))

	b, err := ioutil.ReadFile(filepath.Join(dir, "2009-11", "2009-11-23_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b) != "abc" {
		t.Errorf("log expected abc, but %s", string(b))
	}
}

func TestLogWriter_WriteRotateBySize(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	const rotateSize = 5
	w := cilog.NewLogWriter(dir, "module", rotateSize)
	w.WriteWithTime([]byte("abc"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	w.WriteWithTime([]byte("def"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	w.WriteWithTime([]byte("ghi"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))

	b1, err := ioutil.ReadFile(filepath.Join(dir, "2009-11", "2009-11-23_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b1) != "abcdef" {
		t.Errorf("log expected abc, but %s", string(b1))
	}

	b2, err := ioutil.ReadFile(filepath.Join(dir, "2009-11", "2009-11-23[1]_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b2) != "ghi" {
		t.Errorf("log expected abc, but %s", string(b2))
	}
}

func TestLogWriter_WriteRotateByNextDay(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	w := cilog.NewLogWriter(dir, "module", 5)
	w.WriteWithTime([]byte("abc"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	w.WriteWithTime([]byte("def"), time.Date(2009, 11, 24, 0, 0, 0, 0, time.Local))

	b1, err := ioutil.ReadFile(filepath.Join(dir, "2009-11", "2009-11-23_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b1) != "abc" {
		t.Errorf("log expected abc, but %s", string(b1))
	}

	b2, err := ioutil.ReadFile(filepath.Join(dir, "2009-11", "2009-11-24_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b2) != "def" {
		t.Errorf("log expected abc, but %s", string(b2))
	}
}

func TestLogWriter_WriteRotateByNextMonth(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	w := cilog.NewLogWriter(dir, "module", 5)
	w.WriteWithTime([]byte("abc"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	w.WriteWithTime([]byte("def"), time.Date(2009, 12, 23, 0, 0, 0, 0, time.Local))

	b1, err := ioutil.ReadFile(filepath.Join(dir, "2009-11", "2009-11-23_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b1) != "abc" {
		t.Errorf("log expected abc, but %s", string(b1))
	}

	b2, err := ioutil.ReadFile(filepath.Join(dir, "2009-12", "2009-12-23_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b2) != "def" {
		t.Errorf("log expected abc, but %s", string(b2))
	}
}

func TestLogWriter_WriteNotExistsDir(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	//os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	w := cilog.NewLogWriter(dir, "module", 5)
	w.WriteWithTime([]byte("abc"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))

	b, err := ioutil.ReadFile(filepath.Join(dir, "2009-11", "2009-11-23_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b) != "abc" {
		t.Errorf("log expected abc, but %s", string(b))
	}
}

func TestLogWriter_Write_DeleteFile_Write(t *testing.T) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	w := cilog.NewLogWriter(dir, "module", 5)
	w.WriteWithTime([]byte("abc"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))
	os.RemoveAll(dir)
	w.WriteWithTime([]byte("def"), time.Date(2009, 11, 23, 0, 0, 0, 0, time.Local))

	b, err := ioutil.ReadFile(filepath.Join(dir, "2009-11", "2009-11-23_module.log"))
	if err != nil {
		t.Error(err)
	}
	if string(b) != "def" {
		t.Errorf("log expected def, but %s", string(b))
	}
}

func BenchmarkLogWriter_Write(b *testing.B) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	w := cilog.NewLogWriter(dir, "module", 1024*1024)
	for n := 0; n < b.N; n++ {
		w.Write([]byte(fmt.Sprintf("this is log. line:%d", n)))
	}
}

func BenchmarkLogWriter_RotateBySize(b *testing.B) {
	idv4, _ := uuid.NewRandom()
	dir := path.Join("ut.dir", idv4.String())
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	w := cilog.NewLogWriter(dir, "module", 32*1024)
	for n := 0; n < b.N; n++ {
		w.Write([]byte(fmt.Sprintf("this is log. line:%d", n)))
	}
}
