package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func QueryEndpoint(address string, threadPrefix string, log chan string) {
	client := &http.Client{}
	i := 0
	for true {
		i += 1
		req, err := http.NewRequest("GET", address, nil)
		if err != nil {
			fmt.Println("Error creating request")
			break
		}
		correlationId := threadPrefix + "-" + strconv.Itoa(i)
		req.Header.Add("X-CorrelationId", correlationId)
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			log <- "Failed after " + fmt.Sprint(time.Since(start).Seconds()) + " seconds: " + err.Error()
			continue
		}
		log <- "Response " + resp.Status + " after " + fmt.Sprint(time.Since(start).Seconds()) + " seconds. (CorrelationId: " + correlationId + ")"
	}
}

func LogFromChannel(c chan string) {
	i := 0
	for {
		i += 1
		log := <-c
		fmt.Println(strconv.Itoa(i)+":", log)
	}
}

func CreateLoggingPipeline() chan string {
	c := make(chan string)
	go LogFromChannel(c)
	return c
}

func main() {
	endpoint := flag.String("e", "http://google.com", "The endpoint to hit")
	parallelism := flag.Int("P", runtime.NumCPU(), "The number of parallel requests")
	flag.Parse()
	log := CreateLoggingPipeline()
	for i := 0; i < *parallelism-1; i++ {
		go QueryEndpoint(*endpoint, "jseely-load"+strconv.Itoa(i), log)
	}
	QueryEndpoint(*endpoint, "jseely-load"+strconv.Itoa(*parallelism-1), log)
}
