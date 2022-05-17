package frm

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/BanAnna0970/frm"
)

var pMan manager

type manager struct {
	total     int
	goodProxy int
	badProxy  int
	ping      time.Duration
}

type Proxy struct {
	NetFmt       string
	FastFmt      string
	BansQuantity int
}

func GetProxies(forSite string) []*Proxy {
	f, err := os.OpenFile("proxies.txt", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	var Proxies []*Proxy
	reg := regexp.MustCompile(`:`)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "@") {
			Proxies = append(Proxies, &Proxy{NetFmt: "http://" + text, FastFmt: text})
		} else if strings.Count(text, ":") == 3 {
			userPass := reg.FindAllStringIndex(text, -1)
			fmted := string([]byte(text)[userPass[1][1]:]) + "@" + string([]byte(text)[:userPass[1][0]])
			Proxies = append(Proxies, &Proxy{NetFmt: "http://" + fmted, FastFmt: fmted})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	pMan.total = len(Proxies)
	for _, p := range Proxies {
		go ping(p, forSite)
	}

	time.Sleep(5 * time.Second)

	if float32(pMan.goodProxy) < float32(pMan.total)*0.9 {
		log.Fatal(`< 90% of good proxies`)
	}

	fmt.Printf("Average ping: %v\n\nDo you want to continue?\nY/n\n", pMan.ping/time.Duration(pMan.goodProxy))

	var shouldContinue string

	fmt.Scanln(&shouldContinue)

	if strings.ToUpper(shouldContinue) != "Y" {
		log.Fatal("Bad ping")
	}

	return Proxies
}

func ping(proxy *Proxy, urll string) {
	transport := &http.Transport{}

	var client = &http.Client{
		Timeout:   2 * time.Second,
		Transport: transport,
	}

	req, err := http.NewRequest("HEAD", urll, nil)
	if err != nil {
		return
	}
	now := time.Now()
	proxyUrl, err := url.Parse(proxy.NetFmt)
	if err != nil {
		logger.Logger.Error().Err(err)
		return
	}
	transport.Proxy = http.ProxyURL(proxyUrl)
	resp, err := client.Do(req)
	if err != nil {
		pMan.badProxy++
		fmt.Printf("[%v] - timeout\n", proxy.NetFmt)
		return
	}
	resp.Body.Close()
	pMan.goodProxy++
	pMan.ping = +time.Since(now)
	fmt.Printf("[%v] - %v\n", proxy.NetFmt, time.Since(now))
}
