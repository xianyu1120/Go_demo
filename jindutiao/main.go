package main

import (
	"fmt"
	"time"
)

type Bar struct {
	percent int64  //百分比
	cur     int64  //当前进度位置
	total   int64  //总进度
	rate    string //进度条
	graph   string //显示符号
}

func main() {
	var bar Bar           //实例化结构体
	bar.NewOption(0, 100) //初始化显示图案以及总的任务量
	for i := 0; i <=100; i++ {
		time.Sleep(100 * time.Millisecond)//线程休眠,100毫秒
		bar.Play(int64(i))
	}
	bar.Finish()
}

//NewOption 默认显示图案
//start 可以不从0开始，支持断点
//total 总的任务量
func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start      //当前位置等于初始化位置
	bar.total = total    //总进度等于初始化总任务量
	if bar.graph == "" { //显示符号为空,打印"#"
		bar.graph = "#"
	}
	bar.percent = bar.getPercent() //计算百分比
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph //进度条以符号显示
	}
}

//getPercent 根据当前cur和total获取当前进度完成百分比
func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total)*100)
}

//NewOptionWithGraph 自己指定显示进度的图案
func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}

//Play 打印进度条
func (bar *Bar) Play(cur int64) {
	bar.cur = cur //进度条为指定的进度
	last := bar.percent
	bar.percent = bar.getPercent()                 //百分比为计算的百分比
	if bar.percent != last && bar.percent%2 == 0 { //百分比不等于传递或者不为偶数
		bar.rate += bar.graph //进度条添加符号
	}
	//\r 表示回车,返回到行首,\n下一行行首 
	//回车,左对齐50字符  整型长度为3不足前面补空格
	fmt.Printf("\r[%-50s] %3d %% %3d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

//Finish 结束换行
func (bar *Bar) Finish() {
	fmt.Println("....")
}
