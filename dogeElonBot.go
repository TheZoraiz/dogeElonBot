package main

import (
	"fmt"
	"time"
)

type TweetStruct struct {
	id   string
	text string
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

			title := "dogeElonBot crashed. Too many errors."
			body := "Please restart it and check your internet connection"
			makeAlert(title, body, true)

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
						body := "\"" + tweet.text + "\"\n\nhttps://twitter.com/elonmusk/status/" + tweet.id + "\n(This only shows up as the first notification)"
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
