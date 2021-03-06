package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var showsuccess bool = false
var successlist []string

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	address := flag.String("url", "", "url address https://google.com")
	showsuccessresult := flag.Bool("v", false, "show success result only")
	flag.Parse()
	if *showsuccessresult {
		showsuccess = true
	}
	if *address == "" {
		println("Please Set url --url or -h for help")
		return
	}
	file, err := os.Open(exPath + "/Directory_list.txt")

	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	list := readToDisplayUsingFile1(file)
	addreswithouthttp := strings.ReplaceAll(*address, "https://", "")
	addreswithouthttp = strings.ReplaceAll(addreswithouthttp, "http://", "")
	subanddomain := strings.Split(addreswithouthttp, ".")
	for i := 0; i < len(subanddomain)-1; i++ {
		list = append(list, subanddomain[i])
	}
	newlist := append(list, addreswithouthttp)
	newlistwithcurl := pathinurl(*address)
	uniqelistwithurl := Unique(newlistwithcurl)

	for _, str := range uniqelistwithurl {
		newlist = append(newlist, str)
	}
	defer file.Close()
	for _, url := range newlist {

		if url != "" {
			if strings.HasPrefix(url, "/") {

				if strings.HasSuffix(url, "/") {
					checkurl(*address+url, url)
				} else {
					checkurl(*address+url+"/", url)
				}

			} else {
				if strings.HasSuffix(url, "/") {
					checkurl(*address+"/"+url, "/ "+url)
				} else {
					checkurl(*address+"/"+url+"/", "/"+url)
				}
			}

		}

	}
	fmt.Printf("%d %s", len(successlist), " Found")
	for _, v := range successlist {
		println(v)
	}
	//fmt.Printf(*address)
}
func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}
func pathinurl(urlrecive string) (list []string) {
	response, err := http.Get(urlrecive)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)

	re := regexp.MustCompile("href=[\"'](.*?)[\"']")
	var patharray []string
	found := re.FindAllStringSubmatch(responseString, -1)
	addresuri, _ := url.Parse(urlrecive)

	for _, fou := range found {

		perfix := ""
		if strings.Contains(strings.ToLower(fou[1]), "http") {
			perfix = fou[1]
		} else {
			if strings.HasPrefix(fou[1], "/") {
				perfix = urlrecive + fou[1]
			} else {
				perfix = urlrecive + "/" + fou[1]
			}

		}
		checkuri, err := url.Parse(perfix)

		if err != nil {
			fmt.Printf("\n url is problem  :  %s \n ", fou[1])
			continue
		}

		if len(checkuri.Path) > 3 {

			if strings.Contains(addresuri.Host, checkuri.Host) {
				//fmt.Printf("%s\n", checkuri.Path)
				pathurispli := strings.Split(checkuri.Path, "/")
				//fmt.Printf("pathurispli len   :  %d \n", len(pathurispli))
				pathcomplate := ""
				perfix := ""
				for i := 1; i < len(pathurispli); i++ {

					if pathurispli[i] == "/" || pathurispli[i] == "" {
						continue
					}
					tolowers := strings.ToLower(pathurispli[i])

					if strings.TrimSpace(tolowers) != "" {
						if strings.Contains(tolowers, ".png") || strings.Contains(tolowers, ".jpg") ||
							strings.Contains(tolowers, ".svg") || strings.Contains(tolowers, ".gif") ||
							strings.Contains(tolowers, ".css") || strings.Contains(tolowers, ".js") ||
							strings.Contains(tolowers, ".ttf") || strings.Contains(tolowers, ".ico") ||
							strings.Contains(tolowers, ".otf") || strings.Contains(tolowers, ".woff") ||
							strings.Contains(tolowers, ".woff2") || strings.Contains(tolowers, ".ico") {

						} else {
							if strings.Contains(tolowers, "?") {
								splitqus := strings.Split(pathurispli[i], "?")
								pathcomplate += perfix + splitqus[0]
							} else {
								pathcomplate += perfix + pathurispli[i]
							}

						}

					}
					perfix = "/"
				}
				//patharray = append(patharray, pathcomplate)

				patharray = append(patharray, pathcomplate) // here
				for i := len(pathurispli) - 1; 1 < i; i-- {

					//fmt.Println("beginlop" + pathurispli[i])
					lastpath := ""
					for ii := 0; ii < i; ii++ {
						if pathurispli[ii] == "" {
							continue
						}
						lastpath += pathurispli[ii] + "/"
					}
					//fmt.Println("afterlop : " + lastpath)
					patharray = append(patharray, strings.TrimRight(lastpath, "/"))
				}

			}
		}

	}
	defer response.Body.Close()
	return patharray

}
func readToDisplayUsingFile1(f *os.File) (line []string) {
	defer f.Close()
	reader := bufio.NewReader(f)
	contents, _ := ioutil.ReadAll(reader)
	lines := strings.Split(string(contents), "\n")
	return lines
}
func checkurl(url string, path string) {

	resp, err := http.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "http: server gave HTTP response to HTTPS clien") {
			os.Exit(3)
		}
		fmt.Printf("%s", err.Error())
	}
	if err == nil {
		if !showsuccess {
			fmt.Println(url + " [] " + resp.Status)
		}

		if resp.StatusCode == 200 {
			if resp.Header.Get("Content-Type") != "" {
				contenttype := resp.Header.Get("Content-Type")
				if strings.Contains(contenttype, "text/html") {
					bodyBytes, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return
					}

					bodyString := strings.ToLower(string(bodyBytes))
					if (strings.Contains(bodyString, "index of /") || strings.Contains(bodyString, "directory /")) && strings.Contains(bodyString, strings.TrimRight(path, "/")) {
						//fmt.Printf("'%s', '%s', '%s',\n", url, resp.Status, resp.Header.Get("Content-Type"))
						fmt.Println(url)
						successlist = append(successlist, url)
					}

				}
			}

		} else {
			//fmt.Printf("'%s', '%s',\n", url, resp.Status)
		}

	}

}
