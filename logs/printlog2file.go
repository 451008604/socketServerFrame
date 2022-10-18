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
	prefix string // 标识
	info   string // 信息
	stack  string // 堆栈
}

type fileErrData struct {
	prefix string   // 标识
	err    error    // 错误信息
	tips   []string // 自定义提示
	stack  string   // 堆栈
}

type filePanicData struct {
	prefix string // 标识
	err    error  // 错误信息
	stack  string // 堆栈
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
				log.Println(msgInfo.stack, msgInfo.prefix, msgInfo.info)

			// 错误信息
			case errInfo := <-fileErrCh:
				setLogFile()
				if len(errInfo.tips) > 0 {
					log.Println(errInfo.stack, errInfo.prefix, errInfo.tips, errInfo.err.Error())
				} else {
					log.Println(errInfo.stack, errInfo.prefix, errInfo.err.Error())
				}

			// panic信息
			case panicInfo := <-filePanicCh:
				setLogFile()
				log.Println(panicInfo.stack, panicInfo.prefix, panicInfo.err.Error())
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

	fileInfoCh <- fileInfoData{
		prefix: "[info]\t",
		info:   msg,
		stack:  getCallerStack(),
	}
}

// PrintLogErrToFile 打印错误到日志文件
func PrintLogErrToFile(err error, tips ...string) bool {
	if err == nil {
		return false
	}

	fileErrCh <- fileErrData{
		prefix: "[err]\t",
		err:    err,
		tips:   tips,
		stack:  getCallerStack(),
	}
	return true
}

// PrintLogPanicToFile 打印Panic到日志文件
func PrintLogPanicToFile(err error) {
	if err == nil {
		return
	}

	filePanicCh <- filePanicData{
		prefix: "[panic]\t",
		err:    err,
		stack:  getCallerStack(),
	}
}

// 获取堆栈信息
func getCallerStack() string {
	_, file, line, _ := runtime.Caller(2)
	s := file[strings.LastIndex(file, "/")+1:]
	return fmt.Sprintf("%s:%d\t", s, line)
}
