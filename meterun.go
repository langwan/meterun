package meterun

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"sync"
	"time"
)

type Request struct {
	Ok      bool
	Runtime time.Duration
}

type Report struct {
	workers      int
	works        int
	sleep        time.Duration
	maxQps       int
	maxRuntime   time.Duration
	totalRequest int
	totalBad     int
}

type PerReport struct {
	qps        int
	maxRuntime time.Duration
	minRuntime time.Duration
	p90        time.Duration
	bad        int
}

var report Report
var perRequests []Request

var gun = make(chan bool, 1)
var waitFinish sync.WaitGroup
var waitReady sync.WaitGroup
var chanFinish = make(chan bool)
var chanRequest chan *Request
var row int

var tick *time.Ticker

func Run(f func() bool, workers int, works int, sleep time.Duration, title string) {
	st := time.Now().Format("2006-01-02 15:04")
	logName := fmt.Sprintf("%s_%s.txt", title, st)
	log.SetFlags(0)
	logFile, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE, 0644)
	mw := io.MultiWriter(os.Stdout, logFile)
	if err != nil {
		panic(err)
	}
	log.SetOutput(mw)
	report.workers = workers
	report.works = works
	report.sleep = sleep

	chanRequest = make(chan *Request, workers)
	tick = time.NewTicker(time.Second)
	go monitor()

	waitFinish.Add(report.workers)
	waitReady.Add(report.workers)
	for i := 0; i < report.workers; i++ {
		go doing(f, i, &waitFinish, &waitReady)
	}

	waitReady.Wait()

	//log.Printf("wait for the workers to be ready...\n")
	for i := 0; i < report.workers; i++ {
		gun <- true
	}
	log.Printf("- START - %s - %s\n", st, title)
	waitFinish.Wait()

	close(gun)
	close(chanRequest)
	chanFinish <- true
	close(chanFinish)
	time.Sleep(time.Millisecond * 100)

	log.Printf("- END -\n")
	log.Printf("|%10s |%20s |%20s |%20s |%20s |%20s |\n", "REQ TOTAL", "REQ BAS", "WORKERS", "WORKS", "SLEEP", "")
	log.Printf("|%10d |%20d |%20d |%20d |%20s |%20s |\n", report.totalRequest, report.totalBad, report.workers, report.works, report.sleep, "")
	log.Printf("\n")
}

func monitor() {
	for {
		select {
		case request, ok := <-chanRequest:
			if ok {
				perRequests = append(perRequests, *request)
				report.totalRequest++
			}
		case <-tick.C:

			result(&perRequests)
			perRequests = make([]Request, 0)

		case _, ok := <-chanFinish:
			if ok {

				result(&perRequests)
				//requests = append(requests, perRequests)

				tick.Stop()
			}
		}
	}
}

func result(perRequest *[]Request) {

	if len(perRequests) == 0 {
		return
	}
	m := PerReport{}

	m.qps = len(perRequests)

	for _, request := range *perRequest {
		if request.Ok == false {
			m.bad++
			report.totalBad++
		}
	}

	sort.SliceStable(*perRequest, func(i, j int) bool {
		return (*perRequest)[i].Runtime < (*perRequest)[j].Runtime
	})
	m.minRuntime = perRequests[0].Runtime
	m.maxRuntime = perRequests[len(perRequests)-1].Runtime

	n90 := int(math.Floor(90.0 / 100.0 * float64(len(perRequests))))

	m.p90 = perRequests[n90].Runtime

	//table := tablewriter.NewWriter(os.Stderr)
	//table.SetHeader([]string{"流水", "QPS", "最大执行时间"})
	//log.Printf("ok")
	//data := []string{strconv.Itoa(n), strconv.Itoa(m.qps), m.maxRuntime.String()}
	//table.Append(data)
	//table.Render()
	row++

	log.Printf("|%10s |%20s |%20s |%20s |%20s |%20s |\n", "SEC", "QPS", "MAX TIME", "MIN TIME", "P90", "BAD")
	log.Printf("|%10d |%20d |%20s |%20s |%20s |%20d |\n", row, m.qps, m.maxRuntime.String(), m.minRuntime.String(), m.p90.String(), m.bad)
	log.Printf("\n")
}

func doing(f func() bool, id int, wg *sync.WaitGroup, waitReady *sync.WaitGroup) {
	waitReady.Done()
	<-gun
	//log.Printf("worker id = %d start...\n", id)
	for i := 0; i < report.works; i++ {

		start := time.Now()

		request := Request{}
		request.Ok = f()

		elapsed := time.Since(start)

		request.Ok = true
		request.Runtime = elapsed

		chanRequest <- &request

		time.Sleep(report.sleep)
	}

	wg.Done()
}
