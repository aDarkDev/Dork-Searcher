package main

/*
* Dork Seacrher Project
* By https://github.com/ConfusedCharacter/
* 
*/

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	showBanner()
	var dork, engine string
	var save bool
	flag.StringVar(&dork, "dork", "", "")
	flag.StringVar(&dork, "d", "", "")

	flag.StringVar(&engine, "engine", "", "")
	flag.StringVar(&engine, "e", "", "")

	flag.BoolVar(&save, "save", false, "")
	flag.BoolVar(&save, "s", false, "")

	flag.Usage = func() {
		fmt.Println("\033[92m" + usage + "\033[0m")
	}
	flag.Parse()
	engine = strings.ToLower(engine)
	Validengine := engine == "google" || engine == "ask" || engine == "duck"
	ValidSave := save == true || save == false
	if Validengine && ValidSave {
		search(engine, dork, save)
	} else {
		fmt.Println(ERR, " Parameters Invalid")
		fmt.Println("\033[92m" + usage + "\033[0m")
	}

}

func search(engine string, dork string, save bool) {
	var baseURL, regex, method, dataSend string
	fmt.Scan(dataSend)
	switch engine {
	case "google":
		method = "get"
		regex = `"><a href="\/url\?q=(.*?)&amp;sa=U&amp;`
		baseURL = "https://www.google.com/search?q=" + dork + "&gws_rd=cr,ssl&client=ubuntu&ie=UTF-8&start=1"
	case "ask":
		method = "get"
		regex = `target=\"_blank\" href='(.*?)' data-unified=`
		baseURL = "https://www.ask.com/web?q=" + dork
	case "duck":
		method = "post"
		regex = `<a class="result__url" href="(.*?)"`
		baseURL = "https://html.duckduckgo.com/html/"
		dataSend = "q=" + dork + "&b="
	}

	fmt.Println(WRN, "if You are using your own ip, it may be block by provider.")
	fmt.Println(ENG, "Engine: \t\t", "["+strings.Title(engine)+"]")
	fmt.Println(INF, "Dork:\t\t", "["+dork+"]")
	var body string
	if method == "get" {
		resp, err := http.Get(baseURL)
		handelError(err)
		fmt.Println(INF, "Status Code:\t", "["+strconv.Itoa(resp.StatusCode)+"]")
		data, err := ioutil.ReadAll((resp.Body))
		handelError(err)
		body = string(data)
		defer resp.Body.Close()
	} else {
		BodyToSend := []byte(dataSend)
		bodyReader := bytes.NewReader(BodyToSend)
		req, err := http.NewRequest("POST", baseURL, bodyReader)
		req.Header.Set("accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`)
		req.Header.Set("accept-language", `en-US,en;q=0.9`)
		req.Header.Set("content-type", `application/x-www-form-urlencoded`)
		req.Header.Set("origin", `https://html.duckduckgo.com`)
		req.Header.Set("referer", `https://html.duckduckgo.com/`)
		req.Header.Set("user-agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36`)
		resp, err := http.DefaultClient.Do(req)
		handelError(err)
		fmt.Println(INF, "Status Code:\t", "["+strconv.Itoa(resp.StatusCode)+"]")
		data, err := ioutil.ReadAll((resp.Body))
		handelError(err)
		body = string(data)
		defer resp.Body.Close()
	}
	fmt.Println(INF, "Body Recived:\t", SUCCESS)
	fmt.Println(INF, "Scalping Links...\t")
	var compileRegex = regexp.MustCompile(regex)
	found := compileRegex.FindAllStringSubmatch(body, -1)
	all_urls := []string{}
	for i := range found {
		url_rs := found[i][1]
		all_urls = append(all_urls, url_rs)
	}
	fmt.Print(URL + " ")
	fmt.Println(strings.Join(all_urls, "\n"+URL+" "))
	if save {
		f, err := os.OpenFile("DorkSearch-Result.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		handelError(err)
		f.Write([]byte(strings.Join(all_urls, "\n")))
		f.Close()
		fmt.Println(INF, "Results Saved in DorkSearch-Result.txt")
	}

}
func handelError(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

func showBanner() {
	author := "@ConfusedCharacter"
	version := "v1.1"
	banner := `
	 ___           _      __                     _               
	/   \___  _ __| | __ / _\ ___  __ _ _ __ ___| |__   ___ _ __ 
       / /\ / _ \| '__| |/ / \ \ / _ \/ _' | '__/ __| '_ \ / _ \ '__|
      / /_// (_) | |  |   <  _\ \  __/ (_| | | | (__| | | |  __/ |   
     /___,' \___/|_|  |_|\_\ \__/\___|\__,_|_|  \___|_| |_|\___|_|   
	
	 ` + version + "  -  " + author

	fmt.Print("\033[96m" + banner + "\033[0m\n\n")
}

type Color string

const (
	ENG     Color = "[\033[96mENG\033[0m]"
	INF           = "[\033[0;94mINF\033[0m]"
	URL           = "[\033[92mURL\033[0m]"
	SUCCESS       = "[\033[92mSUCCESS\033[0m]"
	WRN           = "[\033[93mWRN\033[0m]"
	ERR           = "[\033[91mERR\033[0m]"
	usage         = `Options:
	-d, --dork   <dork>         Search dork
	-e, --engine <engine>       Search engine (default: Google.com)
								(engines: Google , Ask , Duck(DuckDuckGo))	
	-s, --save   <true/false>     Save Dork Result in txt file (default: false)`
)
