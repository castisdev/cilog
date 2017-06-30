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
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

	expected := "module,1.0,2009-11-23,15:21:30.123456,Debug,cilog_test::" +
		file + ":" + strconv.Itoa(line-1) + ",,abc\n"
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

func TestLogger_PackageBase(t *testing.T) {
	v1 := cilog.PackageBase("github.com/castisdev/cdn/cache.(*Server).Serve")
	if v1 != "cache" {
		t.Errorf("package expected %s, but %s", "cache", v1)
	}
	v2 := cilog.PackageBase("github.com/castisdev/cdn/cache/filecache.(*Server).readOne.func1")
	if v2 != "filecache" {
		t.Errorf("package expected %s, but %s", "filecache", v2)
	}
	v3 := cilog.PackageBase("github.com/castisdev/cdn/cache/filecache.NewServer")
	if v3 != "filecache" {
		t.Errorf("package expected %s, but %s", "filecache", v3)
	}
	v4 := cilog.PackageBase("main.main")
	if v4 != "main" {
		t.Errorf("package expected %s, but %s", "main", v4)
	}
}

type dummyWriter struct{}

func (dummyWriter) Write(out []byte) (int, error) {
	return len(out), nil
}
func BenchmarkLogger_WithDummyWriter(b *testing.B) {
	cilog.Set(dummyWriter{}, "module", "1.0,", cilog.DEBUG)
	for n := 0; n < b.N; n++ {
		cilog.Infof("this is log. line:%d", n)
	}
}

func BenchmarkLogger_WithLogWriter(b *testing.B) {
	dir := uuid.NewV4().String()
	os.Mkdir(dir, 0775)
	defer os.RemoveAll(dir)

	cilog.Set(cilog.NewLogWriter(dir, "module", 1024*1024), "module", "1.0,", cilog.DEBUG)
	for n := 0; n < b.N; n++ {
		cilog.Infof("this is log. line:%d", n)
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
			cilog.Reportf("test1 log : %d", i)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			cilog.Reportf("test2 log : %d", i)
		}
		wg.Done()
	}()
	wg.Wait()
}

func TestLevel_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		level   cilog.Level
		wantErr bool
	}{
		{
			name:    "unmarshal debug",
			data:    "level: debug",
			level:   cilog.DEBUG,
			wantErr: false,
		},
		{
			name:    "unmarshal report",
			data:    "level: report",
			level:   cilog.REPORT,
			wantErr: false,
		},
		{
			name:    "unmarshal info",
			data:    "level: info",
			level:   cilog.INFO,
			wantErr: false,
		},
		{
			name:    "unmarshal success",
			data:    "level: success",
			level:   cilog.SUCCESS,
			wantErr: false,
		},
		{
			name:    "unmarshal warning",
			data:    "level: warning",
			level:   cilog.WARNING,
			wantErr: false,
		},
		{
			name:    "unmarshal error",
			data:    "level: error",
			level:   cilog.ERROR,
			wantErr: false,
		},
		{
			name:    "unmarshal fail",
			data:    "level: fail",
			level:   cilog.FAIL,
			wantErr: false,
		},
		{
			name:    "unmarshal exception",
			data:    "level: exception",
			level:   cilog.EXCEPTION,
			wantErr: false,
		},
		{
			name:    "unmarshal critical",
			data:    "level: critical",
			level:   cilog.CRITICAL,
			wantErr: false,
		},
		{
			name:    "unmarshal invalid error",
			data:    "level: unknown",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ret := struct {
			Level cilog.Level `yaml:"level"`
		}{}
		var err error
		if err = yaml.Unmarshal([]byte(tt.data), &ret); (err != nil) != tt.wantErr {
			t.Errorf("Level.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
		}
		if err == nil {
			assert.Equal(t, tt.level, ret.Level, "%s", tt.name)
		}
	}
}

func TestLevel_MarshalYAML(t *testing.T) {
	type testLevel struct {
		Level cilog.Level `yaml:"level"`
	}
	tests := []struct {
		name    string
		data    testLevel
		yamlStr string
		wantErr bool
	}{
		{
			name:    "unmarshal debug",
			data:    testLevel{cilog.DEBUG},
			yamlStr: "level: debug\n",
			wantErr: false,
		},
		{
			name:    "unmarshal report",
			data:    testLevel{cilog.REPORT},
			yamlStr: "level: report\n",
			wantErr: false,
		},
		{
			name:    "unmarshal info",
			data:    testLevel{cilog.INFO},
			yamlStr: "level: info\n",
			wantErr: false,
		},
		{
			name:    "unmarshal success",
			data:    testLevel{cilog.SUCCESS},
			yamlStr: "level: success\n",
			wantErr: false,
		},
		{
			name:    "unmarshal warning",
			data:    testLevel{cilog.WARNING},
			yamlStr: "level: warning\n",
			wantErr: false,
		},
		{
			name:    "unmarshal error",
			data:    testLevel{cilog.ERROR},
			yamlStr: "level: error\n",
			wantErr: false,
		},
		{
			name:    "unmarshal fail",
			data:    testLevel{cilog.FAIL},
			yamlStr: "level: fail\n",
			wantErr: false,
		},
		{
			name:    "unmarshal exception",
			data:    testLevel{cilog.EXCEPTION},
			yamlStr: "level: exception\n",
			wantErr: false,
		},
		{
			name:    "unmarshal critical",
			data:    testLevel{cilog.CRITICAL},
			yamlStr: "level: critical\n",
			wantErr: false,
		},
		{
			name:    "marshal invalid error",
			data:    testLevel{-1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		var out []byte
		var err error
		if out, err = yaml.Marshal(tt.data); (err != nil) != tt.wantErr {
			t.Errorf("Level.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
		}
		if err == nil {
			assert.Equal(t, tt.yamlStr, string(out), "%s", tt.name)
		}
	}
}
