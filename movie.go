package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"io"
	"time"

	"./pup"
	"github.com/fatih/color"
)

type Yts struct {
	SearchResults struct {
		Status string `json:"status"`
		Movies []struct {
			Url string `json:"url"`
			Img string `json:"img"`
			Title string `json:"title"`
			Year string `json:"year"`
		} `json:"data"`
	}

	DownloadResults []struct {
		Href string `json:"href"`
		Text string `json:"text"`
		Title string `json:"Title"`
	}
}

var (
	timeout = time.Duration(5 * time.Second)
	HttpClient = http.Client{
		Timeout: timeout,
	}
	YTS = Yts{}
	MOVIES_VERSION string = "1.0.0"
)

func UrlEncoded(str string) string {
    u, err := url.Parse(str)
    if err != nil {
        return ""
    }
    return u.String()
}

func main(){
	search := ParseFlages()
	searchMovie(search)
	if len(YTS.SearchResults.Movies) <= 0 {
		color.Set(color.FgMagenta, color.Bold)
		fmt.Println("No movies found.")
		color.Unset()
		os.Exit(0)
	}
	fmt.Printf("%s",listMovies())
	downloadIndex := getConformation()
	torrentHtml := torrentPage(YTS.SearchResults.Movies[downloadIndex].Url)
	os.Args = []string{
		"anbuksv",
		"#movie-info > p a json{}",
	}
	nodes := pup.Run(torrentHtml)
	json.Unmarshal(nodes,&YTS.DownloadResults)
	color.Set(color.FgYellow, color.Bold)
	fmt.Println(YTS.DownloadResults[len(YTS.DownloadResults)-1].Href) //By default high resolution torrent link will be printed
	color.Unset()
}

func ParseFlages() string {
	cmds := os.Args[1:]
	nonFlagCmds := make([]string, len(cmds))
	n := 0
	for i := 0; i < len(cmds); i++ {
		cmd := cmds[i]
		switch cmd {
		case "-h","--help","help":
			PrintMoviesHelp(os.Stdout, 0)
		case "--version":
			fmt.Println(MOVIES_VERSION)
			os.Exit(0)
		default:
			nonFlagCmds[n] = cmds[i]
			n++
		}
	}
	return strings.Join(nonFlagCmds," ")
}

func PrintMoviesHelp(w io.Writer, exitCode int) {
	helpString := `Usage
    movies [movie name] [flags]
Version
    %s
Flags
    -h --help          display this help
    --version          display version
`
	fmt.Fprintf(w, helpString, MOVIES_VERSION)
	os.Exit(exitCode)
}

func getConformation() int {
	var downloadIndex int
	color.Set(color.FgWhite, color.Bold)
	fmt.Print("awesome-movie> ")
	color.Unset()
	fmt.Scan(&downloadIndex)
	downloadIndex = downloadIndex - 1
	if downloadIndex < 0 || downloadIndex >= len(YTS.SearchResults.Movies){
		fmt.Println("awesome-movie> Please enter valid number.")
		return getConformation()
	}
	return downloadIndex
}

func onHttpError(err error){
	color.Set(color.FgRed, color.Bold)
	fmt.Println(err)
	color.Unset()
	os.Exit(0)
}

func searchMovie(query string) {
	resp,err := HttpClient.Get("https://yts.am/ajax/search?query="+UrlEncoded(query))
	if err != nil{
		onHttpError(err)		
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body),&YTS.SearchResults)
}

func torrentPage(url string) io.ReadCloser {
	resp,err := HttpClient.Get(url)
	if err != nil {
		onHttpError(err)
	}
	return resp.Body
}

func listMovies() string {
	var _movies string = ""
	for index,movie := range(YTS.SearchResults.Movies){
		_movies = _movies + strconv.Itoa(index+1) + ". " + movie.Title + " (" + movie.Year + ")\n"
	}
	return _movies
}