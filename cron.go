package main

import (
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	log.Println("Starting...")

	//会根据本地时间创建一个新（空白）的 Cron job runner
	c := cron.New()
	/*
		func New() *Cron {
		    return NewWithLocation(time.Now().Location())
		}

		// NewWithLocation returns a new Cron job runner.
		func NewWithLocation(location *time.Location) *Cron {
		    return &Cron{
		        entries:  nil,
		        add:      make(chan *Entry),
		        stop:     make(chan struct{}),
		        snapshot: make(chan []*Entry),
		        running:  false,
		        ErrorLog: nil,
		        location: location,
		    }
		}
	*/
	//AddFunc 会向 Cron job runner 添加一个 func ，以按给定的时间表运行
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
	})
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
	})
	/*
				func (c *Cron) AddJob(spec string, cmd Job) error {
				    schedule, err := Parse(spec)
				    if err != nil {
				        return err
				    }
				    c.Schedule(schedule, cmd)
				    return nil
				}
			会首先解析时间表，如果填写有问题会直接 err，无误则将 func 添加到 Schedule 队列中等待执行
		func (c *Cron) Schedule(schedule Schedule, cmd Job) {
		    entry := &Entry{
		        Schedule: schedule,
		        Job:      cmd,
		    }
		    if !c.running {
		        c.entries = append(c.entries, entry)
		        return
		    }

		    c.add <- entry
		}
	*/

	//在当前执行的程序中启动 Cron 调度程序。其实这里的主体是 goroutine + for + select + timer 的调度控制哦
	c.Start()
	/*
		func (c *Cron) Run() {
		    if c.running {
		        return
		    }
		    c.running = true
		    c.run()
		}
	*/

	t1 := time.NewTimer(time.Second * 10)
	//time.NewTimer 会创建一个新的定时器，持续你设定的时间 d 后发送一个 channel 消息
	//for + select 阻塞 select 等待 channel
	//t1.Reset 会重置定时器，让它重新开始计时 （注意，本文适用于 “t.C已经取走，可直接使用 Reset”）
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
