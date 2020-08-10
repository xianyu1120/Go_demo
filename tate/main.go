package main

import (
	"fmt"	
	"time"
)


type Bar struct {
	percent int64  //百分比
	cur     int64  //当前进度
	total   int64  //总任务量
	rate    string //进度条
	graph   string //进度条表示符号
}

//初始化
func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start //当前进度等于初始化进度
	bar.total = total
	if bar.graph == "" {
		bar.graph = "■"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph //初始化进度条位置
	}
}
func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

//指定符号
func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}

//打印进度条
func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%% %8d/%d",bar.rate,bar.percent,bar.cur,bar.total)
}
func (bar *Bar) Finish() {
	fmt.Println()
}
func main() {
	var bar Bar
	bar.NewOption(0,100)
	// bar.NewOptionWithGraph(0,100,"#")
	for i := 0; i <= 100; i++ {
		time.Sleep(100*time.Millisecond)
		bar.Play(int64(i))
	}
	bar.Finish()
}
