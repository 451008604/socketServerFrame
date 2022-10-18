package logs

import "fmt"

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
					fmt.Println(errInfo.tips, errInfo.err.Error())
				} else {
					fmt.Println(errInfo.err.Error())
				}
			case panicInfo := <-logPanicCh:
				panic(panicInfo)
			default:
				break
			}
		}
	}()
}

// PrintLogInfoToConsole 打印到控制台信息
func PrintLogInfoToConsole(msg string) {
	if msg == "" {
		return
	}
	logInfoCh <- msg
}

// PrintLogErrToConsole 打印到控制台错误
func PrintLogErrToConsole(err error, tips ...string) bool {
	if err == nil {
		return false
	}

	logErrCh <- logErrData{
		err:  err,
		tips: tips,
	}
	return true
}

// PrintLogPanicToConsole 打印到控制台Panic
func PrintLogPanicToConsole(err error) {
	if err == nil {
		return
	}
	logPanicCh <- err
}
