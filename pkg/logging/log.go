package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	// Logging to a file.
	filePath := getLogFilePath()
	fileName := getLogFileName()
	f, err := openLogFile(fileName, filePath)
	if err != nil {
		log.Fatalln(err)
	}
	//打印gin 框架默认的输出内容到日志,如果需要同时将日志写入文件和控制台，请使用以下代码。
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	//如果不需要日志打印到控制台，请使用以下代码
	//gin.DefaultWriter = io.MultiWriter(f)

	//使用自定义日志格式输出
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

/*func Setup() {
	filePath := getLogFileFullPath()
	F = openLogFile(filePath) //E:\jianyu\gogin\runtime\logs
	//log.New：创建一个新的日志记录器。out定义要写入日志数据的IO句柄。prefix定义每个生成的日志行的开头。flag定义了日志记录属性
	//log.New接受三个参数：
	//io.Writer：日志都会写到这个Writer中；可以使用io.MultiWriter实现多目的地输出
	//prefix：前缀，也可以后面调用logger.SetPrefix设置；
	//flag：选项，也可以后面调用logger.SetFlag设置。
	//上面代码将日志输出到一个bytes.Buffer，然后将这个buf打印到标准输出。

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
	//log.LstdFlags：日志记录的格式属性之一，其余的选项如下 flag可选值 在log包里首先定义了一些常量，它们是日志输出前缀的标识:
	//const (
	//	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	//	Ltime                         // the time in the local time zone: 01:23:23
	//	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	//	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	//	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	//	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	//	LstdFlags     = Ldate | Ltime // initial values for the standard logger
	//	Ldate：输出当地时区的日期，如2020/02/07；
	//	Ltime：输出当地时区的时间，如11:45:45；
	//	Lmicroseconds：输出的时间精确到微秒，设置了该选项就不用设置Ltime了。如11:45:45.123123；
	//	Llongfile：输出长文件名+行号，含包名，如github.com/darjun/go-daily-lib/log/flag/main.go:50；
	//	Lshortfile：输出短文件名+行号，不含包名，如main.go:50；
	//	LUTC：如果设置了Ldate或Ltime，将输出 UTC 时间，而非当地时区。
	//)
}*/

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
