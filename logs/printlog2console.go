package logs

import (
	"fmt"
)

type logErrData struct {
	err  error
	tips []string
}

var (
	logInfoCh  = make(chan string)
	logErrCh   = make(chan logErrData)
	logPanicCh = make(chan error)
)

func init() {
	go func() {
		for {
			select {
			case msg := <-logInfoCh:
				fmt.Println(msg)
			case errInfo := <-logErrCh:
				if len(errInfo.tips) > 0 {
					fmt.Println(fmt.Sprintf("%v%v", errInfo.tips, errInfo.err.Error()))
				} else {
					fmt.Println(fmt.Sprintf("%v", errInfo.err.Error()))
				}
			case panicInfo := <-logPanicCh:
				panic(panicInfo)
			default:
				break
			}
		}
	}()
}

// 打印到控制台信息
func printLogInfoToConsole(msg string) {
	if msg == "" {
		return
	}

	logInfoCh <- msg
}

// 打印到控制台错误
func printLogErrToConsole(err error, tips ...string) bool {
	if err == nil {
		return false
	}

	logErrCh <- logErrData{
		err:  err,
		tips: tips,
	}
	return true
}

// 打印到控制台Panic
func printLogPanicToConsole(err error) {
	if err == nil {
		return
	}

	logPanicCh <- err
}
