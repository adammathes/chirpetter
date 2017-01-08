package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)

	// required
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")

	// optional
	lastIDFile := flags.String("last-id", "", "File to store last ID seen")
	extractLinks := flags.Bool("extract-links", false, "Pull out external links to top")

	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "CHIRPETTER")

	// ensure required flags
	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" || *lastIDFile == "" {
		log.Fatal("Required flags not set -- please include consumer keys, secret, tokens")
	}

	// figure out where we left off
	lastSeen := 0
	lf, err := os.Open(*lastIDFile)
	fmt.Fscanf(lf, "%d", &lastSeen)
	lf.Close()

	// set up the twitter client
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// get everything since last seen, or 200 if we don't know where we left off
	var homeTimelineParams *twitter.HomeTimelineParams
	if lastSeen != 0 {
		homeTimelineParams = &twitter.HomeTimelineParams{SinceID: int64(lastSeen)}
	} else {
		homeTimelineParams = &twitter.HomeTimelineParams{Count: 200}
	}
	tweets, _, _ := client.Timelines.HomeTimeline(homeTimelineParams)
	if len(tweets) == 0 {
		panic("You're all caught up! No tweets since last run.")
	}

	// save most recent for next time
	lf, err = os.Create(*lastIDFile)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(lf, "%d", tweets[0].ID)
	lf.Close()

	// take out all tweets with non-twitter links, pull them up to the top
	// this will be in random order since we fetch link titles asynchronously
	if *extractLinks {
		printLinks(tweets)
	}

	// print out each tweet, in chronological order
	// this repeats the tweets above but, i mean, whatever. it's fine.
	for i, _ := range tweets {
		tweet := tweets[len(tweets)-i-1]
		fmt.Printf("@%s\n%s\n\n", tweet.User.ScreenName, tweet.Text)
	}
}

func printLinks(tweets []twitter.Tweet) {
	linkCount := 0
	ch := make(chan string)
	for _, tweet := range tweets {
		for _, url := range tweet.Entities.Urls {
			if !blacklisted(url.ExpandedURL) {
				linkCount++
				go fetchLink(tweet, url, ch)
				break
			}
		}
	}
	for i := 0; i < linkCount; i++ {
		fmt.Println(<-ch)
	}
	return
}

func fetchLink(tweet twitter.Tweet, url twitter.URLEntity, ch chan<- string) {
	title := getTitle(url.ExpandedURL)
	ch <- fmt.Sprintf("%s\n%s\n@%s\n%s\n\n", title, url.ExpandedURL, tweet.User.ScreenName, tweet.Text)
	return
}

func blacklisted(u string) bool {
	blacklist := [...]string{"t.co", "twitter.com", "instagram.com", "facebook.com"}
	for _, b := range blacklist {
		if strings.Contains(u, b) {
			return true
		}
	}
	return false
}

func getTitle(u string) string {
	c := &http.Client{
		// give up after 5 seconds
		Timeout: 5 * time.Second,
	}
	resp, err := c.Get(u)
	if err != nil {
		return ""
	}
	root, err := html.Parse(resp.Body)
	title, ok := scrape.Find(root, scrape.ByTag(atom.Title))
	if ok {
		return scrape.Text(title)
	}
	return ""
}
