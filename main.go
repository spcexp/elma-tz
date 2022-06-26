package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func main() {
	countCh := make(chan int, 1)
	ch := make(chan int, 5)
	urls := []string{
		"https://golang.org",
		"https://google.com",
		"https://www.youtube.com/",
		"https://ru.wikipedia.org/",
		"https://www.coursera.org/",
		"https://www.elma-bpm.ru/",
	}
	wg := sync.WaitGroup{}
	countCh <- 0
	for _, url := range urls {
		wg.Add(1)
		ch <- 1
		go getGoEntry(url, countCh, ch, &wg)
	}
	wg.Wait()
	fmt.Println("Total: ", <-countCh)
}

func getGoEntry(url string, countCh, ch chan int, wg *sync.WaitGroup) {
	page := getPage(url)
	count := countGoEntry(page)
	fmt.Println(url+": ", count)
	currentSum := <-countCh
	currentSum += count
	go setSum(currentSum, countCh)
	<-ch
	wg.Done()
}

func getPage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
}

func countGoEntry(page string) int {
	return strings.Count(page, "Go")
}

func setSum(currentSum int, countCh chan int) {
	countCh <- currentSum
}
