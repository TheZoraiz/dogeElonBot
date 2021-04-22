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

type TweetStruct struct {
	id   string
	text string
}

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

func main() {
	fmt.Println("This script checks Elon Musk's twitter account every 3 minutes and notifies you via desktop notifications whether or not he made a tweet with the keyword \"doge\".")
	fmt.Println("If NO tweets with the \"doge\" keyword were used, desktop notifications will show every 30 minutes.")
	fmt.Println("Minimize this prompt. Closing it will shut down the script.")

	visitedList := []string{}
	errorCount := 0
	counter := 0

	fmt.Println("\nLive logs:")
	for true {
		if errorCount == 10 {
			fmt.Println("Too many errors. It's possible that the API is not responding as expected or there's something wrong with your internet connection. Please restart or try again later.")
			break
		}

		// Elon's last 5 tweets
		data, err := fetchApiData("https://api.twitter.com/2/users/44196397/tweets?max_results=5")

		tweetsList, err2 := getTweets(data)

		if err != "" || err2 != "" {
			errorCount++
			fmt.Println(err2)
			time.Sleep(5 * time.Second)
			continue
		}

		foundTweets := false
		for _, tweet := range tweetsList {
			if containsKeywords(tweet.text) {
				if !isVisited(visitedList, tweet.text) {
					if counter == 0 {
						title := "Elon made a Dogecoin tweet recently in his last 5 tweets"
						body := "\"" + tweet.text + "\"\n\nhttps://twitter.com/elonmusk/status/" + tweet.id
						makeAlert(title, body, true)

						fmt.Println(time.Now().Format("[01-02-2006 15:04:05]") + " DOGE TWEET!!!")
					} else {
						title := "😱 ELON MADE A DOGECOIN TWEET 😱"
						body := "\"" + tweet.text + "\"\n\nhttps://twitter.com/elonmusk/status/" + tweet.id
						makeAlert(title, body, true)

						fmt.Println(time.Now().Format("[01-02-2006 15:04:05]") + " DOGE TWEET!!!")
					}

					visitedList = append(visitedList, tweet.text)
					foundTweets = true
				}

			}
		}
		if foundTweets {
			counter++
			time.Sleep(180 * time.Second)
			continue
		}

		fmt.Println(time.Now().Format("[01-02-2006 15:04:05]") + " No tweet :(")

		// Following code shows alerts on 30 minute intervals

		if counter == 0 {
			// If first lookup returned no tweets
			title := "No recent Elon tweets about Dogecoin"
			body := "Better luck next time..."
			makeAlert(title, body, false)

		} else if counter%10 == 0 {
			// Shows "no tweet" alert every 30 mins (unless a tweet was found at exactly 30 minutes interval)
			title := "No recent Elon tweets about Dogecoin"
			body := "Better luck next time..."
			makeAlert(title, body, false)

		}
		counter++
		time.Sleep(180 * time.Second)
	}
}
