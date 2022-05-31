package frm

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

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

	return Proxies
}
