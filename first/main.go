package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var StoreName = "./counter"
var HttpMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "asura",
	Name:      "stress",
}, []string{"method"})

func main() {
	c := Counter{}
	c.loadCounter()
	go StartMetrics()

	s := gin.Default()
	s.GET("/counter", c.GetCounter)
	s.GET("/clear", c.ClearCounter)
	s.Run(":3131")
}

type Counter struct {
	counter int64
}

func StartMetrics() {
	prometheus.Register(HttpMetrics)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2121", nil)
}

func (c *Counter) GetCounter(ctx *gin.Context) {
	var start = time.Now().UnixNano()
	c.counter = c.counter + 1
	c.storeCounter()
	HttpMetrics.WithLabelValues("counter").Observe(float64(time.Now().UnixNano() - start))
	ctx.JSON(http.StatusOK, gin.H{"counter": c.counter})
}

func (c *Counter) ClearCounter(ctx *gin.Context) {
	var start = time.Now().UnixNano()
	c.counter = 0
	c.storeCounter()
	HttpMetrics.WithLabelValues("clear").Observe(float64(time.Now().UnixNano() - start))
	ctx.JSON(http.StatusOK, gin.H{"counter": c.counter})
}

func (c *Counter) storeCounter() {
	var data = []byte(strconv.FormatInt(c.counter, 10))
	ioutil.WriteFile(StoreName, data, 0666)
}

func (c *Counter) loadCounter() {
	if checkFileExist(StoreName) {
		f, err := os.Create(StoreName)
		defer f.Close()
		if err != nil {
			panic(err)
		}
		c.counter = 0
		return
	}
	data, err := ioutil.ReadFile(StoreName)
	if err != nil {
		panic(err)
	}
	c.counter, err = strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.counter = 0
	}
}

func checkFileExist(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}
