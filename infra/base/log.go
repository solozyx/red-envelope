package base

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	// 使用 log 包别名 覆盖标准log日志库
	log "github.com/sirupsen/logrus"
	// 第三方 logrus 日志格式 兼容 logrus 格式
	"github.com/x-cray/logrus-prefixed-formatter"
)

// 日志输出贯穿应用整个生命周期
func init() {
	// 定义日志格式
	// logrus原生日志格式
	// formatter := &log.TextFormatter{}
	// logrus 自定义日志格式
	formatter := &prefixed.TextFormatter{}
	// 开启时间戳 设置时间戳输出格式
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	// prefixed 强制日志格式化
	formatter.ForceFormatting = true
	// 控制台高亮显示
	formatter.ForceColors = true
	formatter.DisableColors = false
	// prefixed 自定义高亮显示颜色
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
	})
	log.SetFormatter(formatter)

	// 日志级别 默认 info
	// 通过系统环境变量控制日志级别
	level := os.Getenv("log.debug")
	if level == "true" {
		log.SetLevel(log.DebugLevel)
	}

	// 显示文件名代码行数
	log.SetReportCaller(true)

	// 日志文件 和 滚动存储
	logFileSettings()

	// log.Info("测试")
	// log.Debug("debug")
}

// logrus 默认没有日志文件输出功能 通过三方hook支持
// github.com/lestrrat/go-file-rotatelogs
func logFileSettings() {
	// 配置日志输入目录
	logPath, _ := filepath.Abs("./logs")
	log.Infof("log dir: %s", logPath)
	logFileName := "resk"
	// 日志文件最大保存时间 24h
	maxAge := 24 * time.Hour
	// 日志切割时间间隔 1h
	rotationTime := time.Hour * 1

	os.MkdirAll(logPath, os.ModePerm)

	baseLogPath := path.Join(logPath, logFileName)
	// 设置滚动日志输出
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime)) // 日志切割时间间隔

	if err != nil {
		log.Errorf("config local file system logger error = %+v", err)
	}

	writers := io.MultiWriter(writer, os.Stdout)
	log.SetOutput(writers)
}
