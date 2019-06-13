package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/jasonlvhit/gocron"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Configure struct {
	Urls         []URL  `json:"urls"`
	Bot          Bot    `json:"bot"`
	Minutes2Call uint64 `json:"minutes2call"`
	Times2OK     int    `json:"times2ok"`
	Debug        bool   `json:"debug"`
}
type Bot struct {
	Token     string `json:"token"`
	ChannelOk int64  `json:"channelok"`
	ChannelKo int64  `json:"channelko"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type URL struct {
	Url    string     `json:"url"`
	Header []KeyValue `json:"header"`
}

var conf = Configure{}

var counter int

type result struct {
	url     string
	elapsed time.Duration
	status  string
	size    int
	err     error
}

func getURLs() {
	ch := make(chan result)

	for _, url := range conf.Urls {
		go get(ch, url)
	}

	var results []result
	for range conf.Urls {
		res := <-ch
		results = append(results, res)
	}

	sort.Slice(results, func(i, j int) bool {
		res1, res2 := results[i], results[j]
		if (res1.err == nil) == (res2.err == nil) {
			return results[i].elapsed < results[j].elapsed
		}
		return res1.err == nil
	})

	var str string
	var errors int = 0
	for _, res := range results {
		if res.err != nil {
			SendMessage(conf.Bot.ChannelKo, fmt.Sprintf("%s %v %s %v\n", res.url, res.elapsed, res.status, res.err))
			errors++
		} else {
			if res.status != "200 OK" {
				str = fmt.Sprintf("%s %v %s %d\n", res.url, res.elapsed, res.status, res.size)
				SendMessage(conf.Bot.ChannelKo, str)
				errors++
			} else {
				if conf.Debug {
					SendMessage(conf.Bot.ChannelOk, fmt.Sprintf("%s %v %s %d\n", res.url, res.elapsed, res.status, res.size))
				}
			}
		}
	}
	if errors == 0 {
		if counter%conf.Times2OK == 0 {
			SendMessage(conf.Bot.ChannelOk, "Status: Ok")
		}
		counter++
	} else {
		counter = 0
	}
}

func get(ch chan result, url URL) {
	start := time.Now()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url.Url, nil)
	if url.Header != nil {
		for _, kv := range url.Header {
			req.Header.Set(kv.Key, kv.Value)
		}
	}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode >= 400 {
		ch <- result{
			url:     url.Url,
			elapsed: time.Since(start),
			err:     fmt.Errorf("http status code %d", resp.StatusCode),
			status:  resp.Status,
		}
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- result{
			url:     url.Url,
			elapsed: time.Since(start),
			err:     err,
			status:  resp.Status,
		}
		return
	}

	ch <- result{
		url:     url.Url,
		elapsed: time.Since(start),
		status:  resp.Status,
		size:    len(body),
	}
}

func SendMessage(channel int64, msg string) {

	b, err := tb.NewBot(tb.Settings{
		Token:  conf.Bot.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	c := &tb.Chat{}
	c.ID = channel

	b.Send(c, msg)

}

func main() {
	var confPath string = "monit-telegram.json"

	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Println("Some problem: ", err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		fmt.Println("Some problem: ", err.Error())
		os.Exit(1)
	}

	counter = 0
	s := gocron.NewScheduler()
	s.Every(conf.Minutes2Call).Minute().Do(getURLs)
	<-s.Start()
}
