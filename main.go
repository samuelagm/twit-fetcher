package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"net/http"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/resty.v1"
)

type report struct {
	Title       string   `json:"title"`
	FromTwitter bool     `json:"from_twitter"`
	Urls        []string `json:"urls"`
	Image       string   `json:"img"`
	Anonymous   bool     `json:"anonymous"`
	Long        float32  `json:"long"`
	Lat         float32  `json:"lat"`
	Loc         string   `json:"loc"`
	Body        string   `json:"body"`
	Featured    bool     `json:"featured"`
	Author      string   `json:"author"`
	Time        string   `json:"time"`
	Approved    bool     `json:"approved"`
	Upvotes     int      `json:"upvotes"`
	Downvotes   int      `json:"downvotes"`
	IsVideo     bool     `json:"isVideo"`
}

// Map a Go map function
func Map(vs []twitter.URLEntity, f func(twitter.URLEntity) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func convertTweetImageToBase64(tweet *twitter.Tweet, reqClient *resty.Client) string {
	var image string = ""
	url := strings.Replace(tweet.User.ProfileImageURL, "_normal", "", -1)
	if resp, err := reqClient.R().
		SetDoNotParseResponse(true).
		Get(url); err == nil {
		if body, err := ioutil.ReadAll(resp.RawBody()); err == nil {
			image = base64.RawStdEncoding.EncodeToString(body)
		}
	}
	return image
}

func sendPost(tweet *twitter.Tweet) {
	apiURL := "http://www.uprightapi.cloud"
	post := &report{
		Title:       tweet.Text,
		FromTwitter: true,
		Urls:        getEntityURLs(tweet),
		Image:       strings.Replace(tweet.User.ProfileImageURL, "_normal", "", -1),
		Anonymous:   true,
		Long:        0.0,
		Lat:         0.0,
		Loc:         "",
		Body:        getTweetText(tweet),
		Featured:    false,
		Author:      tweet.User.ScreenName,
		Approved:    true,
		Time:        tweet.CreatedAt,
		Upvotes:     0,
		Downvotes:   0,
		IsVideo:     false,
	}

	if _, err := resty.New().R().SetBody(post).Post(apiURL + "/post/createpost"); err == nil {
		log.Println("Post sent")
	} else {
		log.Println(err)
	}
}

func getTweetText(tweet *twitter.Tweet) string {
	var text string
	if tweet.ExtendedTweet != nil {
		text = tweet.ExtendedTweet.FullText
	} else {
		text = tweet.Text
	}
	return text
}

func getEntityURLs(tweet *twitter.Tweet) []string {
	var urls []string = []string{}
	if tweet.Entities.Urls != nil {
		urls = Map(tweet.Entities.Urls, func(entity twitter.URLEntity) string {
			return entity.URL
		})
	}
	return urls
}


func startServer(){
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Tweet Fetcher is working...")
	})
	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
	
}

func main() {

	config := oauth1.NewConfig("vizBLoVyy7jCO6YodnZfPQ9uw", "cGqyg85zFBJNsQzSNPq1gRKGWoiF0tswk7cZVIYcx0QCK8hw6v")
	token := oauth1.NewToken("145848213-IC11OPslrRVXRXQCE9v0YLJW1Yle5oLzlFANA5T6", "ztMup6vGkUg7rqtwOejl9zKq2uKwMJ9QabwqMheTDsz5j")
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		go sendPost(tweet)
	}
	demux.Event = func(event *twitter.Event) {
		log.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")
	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"#upright4nigeria", "#Upright4Nigeria"},
		Language:      []string{"en"},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go startServer()
	go demux.HandleChan(stream.Messages)
	fmt.Println("Processing Stream...")
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
