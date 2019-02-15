package log

import (
	"github.com/sirupsen/logrus"
	
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"os"
	"time"
	"fmt"
)

/**
	使用logrus 生成的日志文件管理,默认状态下只是在控制台输出。
	参数是 日志级别，路径，最大保存时间，分割时间
	TODO 远程管理日志 如 elasticsearch

**/

var log = logrus.New()
var writer *rotatelogs.RotateLogs

//New 新建一个Log实例
func New(strLevel string, pathname string, maxAge, rotationTime time.Duration) (*logrus.Logger ,error ) {

	//默认先设置到控制台
	log.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(strLevel)

	if err != nil {
		log.WithFields(logrus.Fields{
			"strLevel": strLevel,
			"err": err.Error(),
		}).Fatal("NewLog config parse strLevel error")
	}

	log.SetLevel(level)

	if pathname == "" {
		return log,nil
	}
	now := time.Now()
	writer, err = rotatelogs.New(
		fmt.Sprintf("%s/%d%02d%02d_%02d.log",pathname, now.Year(), now.Month(), now.Day(), now.Hour()),
        rotatelogs.WithLinkName(fmt.Sprintf("%s/log.log", pathname)), // 生成软链，指向最新日志文件
        rotatelogs.WithMaxAge(maxAge),        // 文件最大保存时间
        rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
    )
    if err != nil {
		log.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("NewLog config parse rotatelogs error")
    }
	lfHook := lfshook.NewHook(lfshook.WriterMap{
        logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
        logrus.InfoLevel:  writer,
        logrus.WarnLevel:  writer,
        logrus.ErrorLevel: writer,
        logrus.FatalLevel: writer,
        logrus.PanicLevel: writer,
	})
	
	lfHook.SetFormatter(log.Formatter)
    log.AddHook(lfHook)
	return log, nil
}

//WithField 日志字段
func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

//Debug debug 日志
func Debug(format string, a ...interface{}){
	 log.Debugf(format, a...)
}

//Info Info 日志
func Info(format string, a ...interface{}){
	log.Infof(format, a...)
}

//Warn Warn 日志
func Warn(format string , a ...interface{}) {
	log.Warnf(format, a...)
}

//Error Error 日志
func Error(format string, a ...interface{}){
	log.Errorf(format, a...)
}

//Fatal Fatal 日志
func Fatal(format string, a ...interface{}){
	log.Fatalf(format, a...)
}

func Close() {
	if writer != nil {
		writer.Close()
	}
}