package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 往文件里面写日志相关代码
var(
	//MaxSize 日志通道的大小
	MaxSize=50000
)
type FileLogger struct {
	Level       LogLevel
	FilePath    string //日志文件保存的路径
	FileName    string //日志文件保存的文件名
	maxFileSize int64  //文件大小,超过要切割
	fileObj     *os.File
	errFileObj  *os.File
	logChan     chan *logMsg
}

//构造通道结构体
type logMsg struct {
	level     LogLevel
	msg       string
	funcName  string
	line      int
	fileName  string
	timeStamp string
}

//NewFileLogger 构造函数
func NewFileLogger(levelStr, fp, fn string, maxSize int64) *FileLogger {
	logLevel, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	fl := &FileLogger{
		Level:       logLevel,
		FilePath:    fp,
		FileName:    fn,
		maxFileSize: maxSize,
		logChan:     make(chan *logMsg,MaxSize),
	}
	fl.initFile() //按照文件名称和文件路径打开文件
	if err != nil {
		panic(err)
	}
	return fl
}

//文件初始化
func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.FilePath, f.FileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed,err:%v\n", err)
		return err
	}
	errFileObj, err := os.OpenFile(fullFileName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed,err:%v\n", err)
		return err
	}
	//日志文件都已经打开了
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	//开启5个后台的Goroutine去往文件里日志
	// for i := 0; i < 5; i++ {
	// 	go f.writeLogBackgroud()
	// }
	go f.writeLogBackgroud()
	
	return nil
}

//判断是否需要记录该日志
func (f *FileLogger) enable(logLevel LogLevel) bool {
	return f.Level <= logLevel
}

//判断文件是否需要切割
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("open fileInfo failed,err:%v\n", err)
		return false
	}
	//如果当前文件大小>=日志文件大小的最大值
	return fileInfo.Size() >= f.maxFileSize
}

//SplitFile 切割文件
func (f *FileLogger) SplitFile(file *os.File) (*os.File, error) {
	nowStr := time.Now().Format("20060102150405000") //获取当前时间
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err:%v\n", err)
		return nil, err
	}
	logName := path.Join(f.FilePath, fileInfo.Name()) //拿到完整路径
	newLogName := fmt.Sprintf("%s.bak%s", logName, nowStr)
	//1. 关闭当前的日志文件
	file.Close()
	//2. 备份一下rename()  xx.log-->xx.log.back202008071709

	os.Rename(logName, newLogName) //重命名并移动到新路径
	//3.打开一个新文件
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new log file failed,err:%v\n", err)
		return nil, err
	}
	//4. 将打开的新日志文件对象赋值给fileObj
	return fileObj, nil
}
func (f *FileLogger) writeLogBackgroud() {
	for {
		if f.checkSize(f.fileObj) { //如果文件需要切割
			newFile, err := f.SplitFile(f.fileObj)
			if err != nil {
				return
			}
			f.fileObj = newFile
		}
		select {
		case logTmp := <-f.logChan:
			//把日志拼出来
			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", logTmp.timeStamp, getLogString(logTmp.level), logTmp.fileName, logTmp.funcName, logTmp.line, logTmp.msg)
			fmt.Fprintf(f.fileObj, logInfo)
			if logTmp.level >= ERROR {
				if f.checkSize(f.errFileObj) { //如果文件需要切割
					newFile, err := f.SplitFile(f.errFileObj)
					if err != nil {
						return
					}
					f.errFileObj = newFile
					//如果要记录的日志大于等于ERROR,还要再ERROR中记录一遍
				}
				fmt.Fprintf(f.fileObj, logInfo)
			}
		default:
			//取不到日志睡眠500毫秒
			time.Sleep(500 * time.Millisecond)
		}

	}
}

//记录日志的方法
func (f *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...) //对传进来的信息格式化
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		//先把日志写道通道中
		logTmp := &logMsg{
			level:     lv,
			msg:       msg,
			funcName:  funcName,
			fileName:  fileName,
			timeStamp: now.Format("2006-01-02 15::04:05"),
			line:      lineNo,
		}
		select {
		case f.logChan <- logTmp:
		default:
			//把日志丢掉，保证不被阻塞
		}
	}
}

//Debug debug级别日志
func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}

//Info ...
func (f *FileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}

// Trace ...
func (f *FileLogger) Trace(format string, a ...interface{}) {
	f.log(TRACE, format, a...)
}

//Warning ...Warning
func (f *FileLogger) Warning(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}

// Error ...
func (f *FileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR, format, a...)
}

//Fatal ...
func (f *FileLogger) Fatal(format string, a ...interface{}) {
	f.log(FATAL, format, a...)
}

//Close 关闭文件链接
func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}
