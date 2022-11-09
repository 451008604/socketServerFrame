package logs

import (
	"log"
	"os"
	"runtime"
	"strconv"
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
	fileInfoCh  = make(chan fileInfoData, 3000)
	fileErrCh   = make(chan fileErrData, 3000)
	filePanicCh = make(chan filePanicData, 3000)
)

var logPath = "./logs/"
var todayFlag = time.Time{}
var currentFileName = ""
var sliceFlag = 0
var myLog = log.Default()

func init() {
	go func() {
		for {
			select {
			// 日志信息
			case msgInfo := <-fileInfoCh:
				setLogFile()
				myLog.Println(msgInfo.stack, msgInfo.prefix, msgInfo.info)

			// 错误信息
			case errInfo := <-fileErrCh:
				setLogFile()
				if len(errInfo.tips) > 0 {
					myLog.Println(errInfo.stack, errInfo.prefix, errInfo.tips, errInfo.err.Error())
				} else {
					myLog.Println(errInfo.stack, errInfo.prefix, errInfo.err.Error())
				}

			// panic信息
			case panicInfo := <-filePanicCh:
				setLogFile()
				myLog.Panicln(panicInfo.stack, panicInfo.prefix, panicInfo.err.Error())

			default:
				break
			}
		}
	}()
}

func setLogFile() {
	today := time.Now()
	// 每天重置分段标记
	if todayFlag.Day() != today.Day() {
		sliceFlag = 0
	}
	todayFlag = today
	timeStamp := today.Format("0102")

	// 要写入的日志文件名称
	fileName := logPath + "log-" + timeStamp + "-" + strconv.Itoa(sliceFlag) + ".log"
	fileInfo, err := os.Stat(fileName)

	// 文件存在
	if !os.IsNotExist(err) {
		// 体积超过限制则建立新的日志文件(1024*1024*50=50M)
		if fileInfo.Size() >= 1024*1024*50 {
			sliceFlag++
			setLogFile()
			return

			// 服务重启时开启新的文件片段
		} else if currentFileName == "" {
			sliceFlag++
			setLogFile()
			return
		}
	}
	// 设置保存日志的文件
	if currentFileName == fileName {
		return
	}
	currentFileName = fileName
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	myLog.SetFlags(log.LstdFlags)
	myLog.SetOutput(file)
}

// 打印信息到日志文件
func printLogInfoToFile(msg string) {
	if msg == "" {
		return
	}

	fileInfoCh <- fileInfoData{
		prefix: "[info]\t",
		info:   msg,
		stack:  getCallerStack(),
	}
}

// 打印错误到日志文件
func printLogErrToFile(err error, tips ...string) bool {
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

// 打印Panic到日志文件
func printLogPanicToFile(err error) {
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
	_, file, line, _ := runtime.Caller(3)
	s := file[strings.LastIndex(file, "/")+1:]
	return s + ":" + strconv.Itoa(line) + "\t"
}
