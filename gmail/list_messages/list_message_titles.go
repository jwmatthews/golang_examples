/**
Originally from:  https://github.com/gsuitedevs/go-samples/blob/master/gmail/quickstart/quickstart.go


 * @license
 * Copyright Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
*/
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	"github.com/mvdan/xurls"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "../token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getMessages(req *gmail.UsersMessagesListCall, nextToken ...string) (*gmail.ListMessagesResponse, string) {
	ntoken := ""
	if len(nextToken) > 0 {
		ntoken = nextToken[0]
	}
	r, err := req.PageToken(ntoken).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}
	fmt.Println("nextToken = ", nextToken, "   r.NextPageToken = ", r.NextPageToken)
	return r, r.NextPageToken
}

func getRFC282Headers(headers []*gmail.MessagePartHeader) (from string, subject string, err error) {
	for _, h := range headers {
		switch h.Name {
		case "Subject":
			subject = h.Value
		case "From":
			from = h.Value
		}
	}
	return
}

func getSubject(headers []*gmail.MessagePartHeader) (subject string) {
	_, subject, err := getRFC282Headers(headers)
	if err != nil {
		log.Fatalf("Unable to parse headers: %v", err)
	}
	return
}

func getMessageContent(payload *gmail.MessagePart) (content string) {
	if len(payload.Parts) > 0 {
		for _, part := range payload.Parts {
			if part.MimeType == "text/html" {
				data, err := base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					log.Fatalf("Unable to decode message: %v", err)
				}
				content = string(data)
			}
		}
	} else {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			log.Fatalf("Unable to decode message: %v", err)
		}
		content = string(data)
	}
	return
}

func getAllURLs(message string) []string {
	return xurls.Strict().FindAllString(message, -1)
}

func main() {
	b, err := ioutil.ReadFile("../credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	svc, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	req := svc.Users.Messages.List("me").Q("label:Recipes").
		MaxResults(25)

	nextToken := ""
	for ok := true; ok; ok = nextToken != "" {
		var r *gmail.ListMessagesResponse
		r, nextToken = getMessages(req, nextToken)
		fmt.Println(len(r.Messages), " messages found")

		for _, m := range r.Messages {
			msg, _ := svc.Users.Messages.Get("me", m.Id).Format("full").Do()
			subject := getSubject(msg.Payload.Headers)
			fmt.Println("\n---", subject)
			content := getMessageContent(msg.Payload)
			urls := getAllURLs(content)
			if len(urls) > 0 {
				fmt.Println("\t", urls[0])
			} else {
				fmt.Println("\tN/A")
			}
		}
	}
}
