package logs

type errData struct {
	err  error
	tips string
}

var infoCh = make(chan string)
var errCh = make(chan errData)
var panicCh = make(chan error)

func init() {
	go func() {
		for {
			select {
			case msg := <-infoCh:
				println(msg)
			case errInfo := <-errCh:
				println(errInfo.tips, errInfo.err.Error())
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

	str := ""
	for i := 0; i < len(tips); i++ {
		str += tips[i]
	}
	errCh <- errData{
		err:  err,
		tips: str,
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
