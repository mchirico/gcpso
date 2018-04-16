package main

import (
	"fmt"
	"time"
	"sync"
)

type Command interface {
	Execute() []byte
	ID() []rune
	Done(Reporter)
}

type Report interface {
	ID([]rune) []rune
	getID() []rune
	Result() []byte
	Start() time.Time
	Stop() time.Time
	Delta() time.Duration
	Summary() string
}

var command = make(chan Command)
var reports = make(chan Report)

type Reporter struct {
	id     []rune
	result []byte
	start  time.Time
	stop   time.Time
	delta  time.Duration
	mu    sync.Mutex
}

func (r *Reporter) Result() []byte {
	return r.result
}
func (r *Reporter) Start() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.start = time.Now()
	return r.start

}
func (r *Reporter) Stop() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.delta = time.Since(r.start)
	r.stop = time.Now()
	return r.stop
}
func (r *Reporter) Delta() time.Duration {
	return r.delta
}
func (r *Reporter) ID(id []rune) []rune {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id = id
	return id
}

func (r *Reporter) getID() []rune {
	return r.id
}

func (r Reporter) Summary() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	msgFmt := "start: %s, stop: %s, delta: %s"
	msg := fmt.Sprintf(msgFmt,r.start,r.stop,r.delta)
	return msg
}


type PingCommand struct {
	result []byte
	id     chan []rune
}

func (p *PingCommand) Execute() []byte {

	p.result = []byte("react pings")

	return []byte("react pings")
}
func (p *PingCommand) ID() []rune {
	return []rune("end")
}
func (p *PingCommand) Done(r Reporter) {

	for i:=0; i < 5; i++ {
		fmt.Printf("report id: %v\n", string(r.id))
		time.Sleep(time.Duration(1) *
			1000 * time.Millisecond)

	}

	p.id <- []rune("All done")
}


type PingCommandSlow struct {
	result []byte
	id     chan []rune
}

func (p *PingCommandSlow) Execute() []byte {

	p.result = []byte("react pings")

	return []byte("react pings")
}
func (p *PingCommandSlow) ID() []rune {
	return []rune("end Slow")
}
func (p *PingCommandSlow) Done(r Reporter) {

	for i:=0; i < 3; i++ {
		fmt.Printf("Ping Slow report id: %v\n", string(r.id))
		time.Sleep(time.Duration(3) *
			1000 * time.Millisecond)

	}

	p.id <- []rune("All done")
}









func Worker() {
	var report = &Reporter{}
	for {
		select {
		case cmd := <-command:
			report.Start()
			report.result = cmd.Execute()
			report.id = cmd.ID()
			report.Stop()

			// Do not want to block
			go cmd.Done(*report)

		case reports <- report:

		}
	}
}

func init() {
	go Worker()
}

func Ping() {
	p := &PingCommand{}
	p.id = make(chan []rune)
	command <- p

	<-p.id
}
func PingSlow() {
	p := &PingCommandSlow{}
	p.id = make(chan []rune)
	command <- p

	<-p.id
}


func PingReport() {
	r := <-reports
	fmt.Printf("\n\nreports: %v id: %v\n",
		r.Summary(),
		string(r.getID()))

}
func main() {

	go PingSlow()
	go Ping()

	for i:=0; i < 10; i++ {
		PingReport()
		time.Sleep(time.Duration(1) *
			1000 * time.Millisecond)
	}

}
