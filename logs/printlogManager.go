package logs

var debugConfig = false

func SetDebugConfig(v bool) {
	debugConfig = v
}

// PrintLogInfoToConsole 打印到控制台信息
func PrintLogInfoToConsole(msg string) {
	if msg == "" {
		return
	}
	if !debugConfig {
		printLogInfoToFile(msg)
		return
	}
	printLogInfoToConsole(msg)
}

// PrintLogErrToConsole 打印到控制台错误
func PrintLogErrToConsole(err error, tips ...string) bool {
	if err == nil {
		return false
	}
	if !debugConfig {
		printLogErrToFile(err, tips...)
		return false
	}

	return printLogErrToConsole(err, tips...)
}

// PrintLogPanicToConsole 打印到控制台Panic
func PrintLogPanicToConsole(err error) {
	if err == nil {
		return
	}
	if !debugConfig {
		printLogPanicToFile(err)
		return
	}

	printLogPanicToConsole(err)
}

// PrintLogInfoToFile 打印信息到日志文件
func PrintLogInfoToFile(msg string) {
	printLogInfoToFile(msg)
}

// PrintLogErrToFile 打印错误到日志文件
func PrintLogErrToFile(err error, tips ...string) bool {
	return printLogErrToFile(err, tips...)
}

// PrintLogPanicToFile 打印Panic到日志文件
func PrintLogPanicToFile(err error) {
	printLogPanicToFile(err)
}
