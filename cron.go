package main

import (
	"log"
	"time"

	"github.com/robfig/cron"

	"github.com/gin-blog/models"
)

func main() {
	log.Println("Starting...")

	//根据本地时间，创建一个新的cron job runner
	c := cron.New()

	/**
	向之前创建的cron job runner添加一个匿名函数，首先解析时间表，如果填写有问题会直接 err，无误则将 func 添加到 Schedule 队列中等待执行
	下面其实就是每秒运行一次
	时间格式：秒/分/小时/天/月/一周的哪一天
	 */
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	//启动goroutine调度程序，主体是 goroutine + for + select + timer （参考go底层实现）
	c.Start()

	//创建一个新的定时器，在指定时间后，向t1.c（定时器的chan类型的属性）中发送一个消息
	t1 := time.NewTimer(time.Second * 10)

	//for死循环，每十秒就往后顺延十秒
	for {
		//select case类似switch case，此select中只有一个case，所以如果t1.C没有数据可以被读出的话，程序会堵死
		select {
		case <-t1.C: //从chan中读取数据

			/*
				重置定时器，让定时器重新开始计时。将定时器的时间向后推10秒
				重置会将计时器更改为在持续时间d后过期。如果计时器已激活，则返回true;如果计时器已过期或已停止，则返回false
			*/
			t1.Reset(time.Second * 10)
		}
	}
}
