package logs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type fileInfoData struct {
	info  string
	stack string
}

type fileErrData struct {
	err   error
	tips  []string
	stack string
}

type filePanicData struct {
	err   error
	stack string
}

var (
	fileInfoCh  = make(chan fileInfoData)
	fileErrCh   = make(chan fileErrData)
	filePanicCh = make(chan filePanicData)
)

var todayFlag = 0

func init() {
	go func() {
		for {
			select {
			// 日志信息
			case msgInfo := <-fileInfoCh:
				setLogFile()
				log.Println(msgInfo.stack, msgInfo.info)

			// 错误信息
			case errInfo := <-fileErrCh:
				setLogFile()
				if len(errInfo.tips) > 0 {
					log.Println(errInfo.stack, errInfo.tips, errInfo.err.Error())
				} else {
					log.Println(errInfo.stack, errInfo.err.Error())
				}

			// panic信息
			case panicInfo := <-filePanicCh:
				setLogFile()
				log.Println(panicInfo.stack, panicInfo.err.Error())
				panic(panicInfo)

			default:
				break
			}
		}
	}()
}

func setLogFile() {
	today := time.Now()
	if todayFlag == today.Day() {
		return
	}
	todayFlag = today.Day()
	file, err := os.OpenFile("./logs/log-"+today.Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	log.SetFlags(log.LstdFlags)
	log.SetOutput(file)
}

// PrintLogInfoToFile 打印信息到日志文件
func PrintLogInfoToFile(msg string) {
	if msg == "" {
		return
	}
	// 获取堆栈信息
	_, file, line, _ := runtime.Caller(1)
	s := file[strings.LastIndex(file, "/")+1:]

	fileInfoCh <- fileInfoData{
		info:  msg,
		stack: fmt.Sprintf("%s:%d\t", s, line),
	}
}

// PrintLogErrToFile 打印错误到日志文件
func PrintLogErrToFile(err error, tips ...string) bool {
	if err == nil {
		return false
	}
	// 获取堆栈信息
	_, file, line, _ := runtime.Caller(1)
	s := file[strings.LastIndex(file, "/")+1:]

	fileErrCh <- fileErrData{
		err:   err,
		tips:  tips,
		stack: fmt.Sprintf("%s:%d\t", s, line),
	}
	return true
}

// PrintLogPanicToFile 打印Panic到日志文件
func PrintLogPanicToFile(err error) {
	if err == nil {
		return
	}
	// 获取堆栈信息
	_, file, line, _ := runtime.Caller(1)
	s := file[strings.LastIndex(file, "/")+1:]

	filePanicCh <- filePanicData{
		err:   err,
		stack: fmt.Sprintf("%s:%d\t", s, line),
	}
}
