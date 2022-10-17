package logs

import "fmt"

type errData struct {
	err  error
	tips []string
}

var infoCh = make(chan string)
var errCh = make(chan errData)
var panicCh = make(chan error)

func init() {
	go func() {
		for {
			select {
			case msg := <-infoCh:
				fmt.Println(msg)
			case errInfo := <-errCh:
				if len(errInfo.tips) > 0 {
					fmt.Println(errInfo.tips, errInfo.err.Error())
				} else {
					fmt.Println(errInfo.err.Error())
				}
			case panicInfo := <-panicCh:
				panic(panicInfo)
			default:
				break
			}
		}
	}()
}

// PrintToConsoleInfo 打印到控制台信息
func PrintToConsoleInfo(msg string) {
	if msg == "" {
		return
	}
	infoCh <- msg
}

// PrintToConsoleErr 打印到控制台错误
func PrintToConsoleErr(err error, tips ...string) bool {
	if err == nil {
		return false
	}

	errCh <- errData{
		err:  err,
		tips: tips,
	}
	return true
}

// PrintToConsolePanic 打印到控制台Panic
func PrintToConsolePanic(err error) {
	if err == nil {
		return
	}
	panicCh <- err
}
