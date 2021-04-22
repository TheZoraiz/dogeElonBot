package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gen2brain/beeep"
)

func initTune(tuneUrl string) beep.StreamSeekCloser {
	f, err := os.Open(tuneUrl)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	// defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	return streamer
}

func playSound() {
	streamer := initTune("alert.mp3")

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func getTweets(tweets []byte) ([]TweetStruct, string) {
	var tweetStrings []TweetStruct

	// Loop works with first 5 tweets. Change the condition if max_results parameter for API endpoint is changed as well
	for i := 0; i < 5; i++ {
		str, err := jsonparser.GetString(tweets, "data", "["+strconv.Itoa(i)+"]", "text")
		if err != nil {
			return nil, "Unexpected error. Retrying in 5 seconds..."
		}

		id, err2 := jsonparser.GetString(tweets, "data", "["+strconv.Itoa(i)+"]", "id")
		if err2 != nil {
			return nil, "Unexpected error. Retrying in 5 seconds..."
		}

		tempTweet := TweetStruct{id, str}

		tweetStrings = append(tweetStrings, tempTweet)
		// fmt.Println(str + "\n\n")
	}

	return tweetStrings, ""
}

func fetchApiData(url string) ([]byte, string) {

	// Bearer token declared in a Go file in the same package

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "Unexpected error. Retrying in 5 seconds..."
	}

	bearer := bearerToken
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}

	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, "Unexpected error. Retrying in 5 seconds..."
	}
	defer resp.Body.Close()

	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return nil, "Unexpected error. Retrying in 5 seconds..."
	}
	return []byte(body), ""
}

func makeAlert(title string, body string, play bool) {
	err := beeep.Alert(title, body, "assets/warning.png")
	if err != nil {
		panic(err)
	}
	if play {
		playSound()
	}
	fmt.Println(time.Now().Format("[01-02-2006 15:04:05]") + " Notification sent")
}

func containsKeywords(tweet string) bool {
	tweet = strings.ToLower(tweet)

	// Searches these keywords in the tweet strings
	// Add more strings to this slice to check tweets for those keywords as well
	keywordList := []string{
		"doge",
		// "starship",
	}

	for _, keyword := range keywordList {
		if strings.Contains(tweet, keyword) {
			return true
		}
	}
	return false
}

func isVisited(visited []string, tweet string) bool {
	if len(visited) == 0 {
		visited = append(visited, tweet)
		return false
	}
	for _, element := range visited {
		if element == tweet {
			return true
		}
	}
	return false
}
