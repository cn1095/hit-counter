package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"
)

// localTimeFormatter 返回本地时区的时间字符串
func localTimeFormatter() string {
	// 以 RFC3339 格式输出本地时间
	return time.Now().Local().Format(time.RFC3339)
}

// customWriter 是一个实现了 io.Writer 的结构体，用于自定义写入行为
type customWriter struct {
	w io.Writer
}

// Write 方法会插入本地时间戳
func (cw *customWriter) Write(p []byte) (n int, err error) {
	// 插入自定义时间前缀（本地时间）
	timePrefix := fmt.Sprintf("{\"time\":\"%s\"} ", localTimeFormatter())
	full := append([]byte(timePrefix), p...)
	return cw.w.Write(full)
}

// NewLogger 创建一个支持本地时区时间戳的 Logger
func NewLogger(dir, filename string) (*log.Logger, error) {
	logger := log.New("")

	// 不使用默认 Header，因为我们自定义前缀了
	logger.SetHeader("") // 清空默认头部

	// 设置日志输出位置
	var output io.Writer
	fpath := filepath.Join(dir, filename)
	if filename == "" {
		output = os.Stdout
	} else {
		f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		output = f
	}

	// 使用 customWriter 包装输出，以加上本地时间戳
	logger.SetOutput(&customWriter{w: output})

	return logger, nil
}
