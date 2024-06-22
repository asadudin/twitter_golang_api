package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	CreatedAt     string `json:"created_at"`
	PinnedTweetID string `json:"pinned_tweet_id"`
}

type Tweet struct {
	ID        string `json:"id"`
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
	Username  string `json:"username"`
}

type UserResponse struct {
	Data     []User `json:"data"`
	Includes struct {
		Tweets []Tweet `json:"tweets"`
	} `json:"includes"`
}

func FetchUserData(username, bearerToken string) (*UserResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.twitter.com/2/users/by?usernames=%s&user.fields=created_at&expansions=pinned_tweet_id&tweet.fields=author_id,created_at", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error fetching user data: %s, Response Body: %s", resp.Status, string(body))
		return nil, fmt.Errorf("error fetching user data: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userResponse UserResponse
	if err := json.Unmarshal(body, &userResponse); err != nil {
		return nil, err
	}

	return &userResponse, nil
}

func StoreUserData(client *mongo.Client, userResponse *UserResponse) {
	db := client.Database("twitter")
	usersCollection := db.Collection("users")
	tweetsCollection := db.Collection("tweets")

	for _, user := range userResponse.Data {
		_, err := usersCollection.InsertOne(context.Background(), user)
		if err != nil {
			log.Fatal(err)
		}

		// Store tweets with username
		for _, tweet := range userResponse.Includes.Tweets {
			tweet.Username = user.Username
			_, err := tweetsCollection.InsertOne(context.Background(), tweet)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Println("User and tweet data have been stored in MongoDB!")
}
