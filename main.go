package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
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
	"time"
)

var showsuccess bool = false
var successlist []string
var outputFormat string

type Result struct {
	URL         string `json:"url"`
	Status      string `json:"status"`
	ContentType string `json:"content_type"`
}

func main() {
	http.DefaultClient.Timeout = 2 * time.Second
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	address := flag.String("url", "", "URL address (e.g., https://google.com)")
	showsuccessresult := flag.Bool("v", false, "Show success results only")
	outputPath := flag.String("o", "", "Output file path")
	format := flag.String("f", "text", "Output format (text, json, csv)")
	flag.Parse()

	if *showsuccessresult {
		showsuccess = true
	}
	if *address == "" {
		log.Fatal("Please set URL with --url flag or use -h for help")
	}

	outputFormat = *format
	if *outputPath != "" && outputFormat != "text" && outputFormat != "json" && outputFormat != "csv" {
		log.Fatal("Invalid format. Supported formats: text, json, csv")
	}

	if *outputPath != "" {
		outDir := filepath.Dir(*outputPath)
		if err := os.MkdirAll(outDir, 0755); err != nil {
			log.Fatal("Error creating output directory:", err)
		}
	}
	appPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get application path: %v\n", err)
		return
	}
	filepathS := "Directory_list.txt"
	appDir := filepath.Dir(appPath)
	defaultLocalPath := filepath.Join(appDir, filepathS)
	defaultGlobalPath := "/usr/local/bin/" + filepathS
	fmt.Printf("Checking for  in %s\n and %s\n", appDir, defaultGlobalPath)
	// Check if the file exists in the application's directory
	configfilepath := defaultLocalPath
	if _, err := os.Stat(configfilepath); os.IsNotExist(err) {
		// If not found in the app directory, fall back to /usr/local/bin
		fmt.Printf(" not found in %s, trying %s\n", appDir, defaultGlobalPath)
		configfilepath = defaultGlobalPath
	}
	file, err := os.Open(configfilepath)
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

	newlist = append(newlist, uniqelistwithurl...)
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
	fmt.Printf("\nFound %d directory listings\n", len(successlist))
	if *outputPath != "" {
		saveResults(*outputPath, successlist)
	} else {
		for _, v := range successlist {
			fmt.Println(v)
		}
	}
}

func Unique(slice []string) []string {
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

func pathinurl(urlrecive string) (list []string) {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	response, err := client.Get(urlrecive)
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
				pathurispli := strings.Split(checkuri.Path, "/")
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
				patharray = append(patharray, pathcomplate)
				for i := len(pathurispli) - 1; 1 < i; i-- {
					lastpath := ""
					for ii := 0; ii < i; ii++ {
						if pathurispli[ii] == "" {
							continue
						}
						lastpath += pathurispli[ii] + "/"
					}
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
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url)
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
						fmt.Println(url)
						successlist = append(successlist, url)
					}
				}
			}
		} else {
		}
	}
}

func saveResults(outputPath string, results []string) {
	var output string
	switch outputFormat {
	case "json":
		jsonResults := make([]Result, len(results))
		for i, url := range results {
			jsonResults[i] = Result{URL: url, Status: "200", ContentType: "text/html"}
		}
		jsonData, err := json.MarshalIndent(jsonResults, "", "  ")
		if err != nil {
			log.Fatal("Error creating JSON:", err)
		}
		output = string(jsonData)
	case "csv":
		var csvData strings.Builder
		csvData.WriteString("URL,Status,Content-Type\n")
		for _, url := range results {
			csvData.WriteString(fmt.Sprintf("%s,200,text/html\n", url))
		}
		output = csvData.String()
	default:
		output = strings.Join(results, "\n")
	}

	if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
		log.Fatal("Error saving results:", err)
	}
	fmt.Printf("Results saved to: %s\n", outputPath)
}
