package line

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

func createOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID: "1559204179",
		ClientSecret: "b45ec5f5361fdb7dc54b84e79f77eda8",
		Scopes: []string{"profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL: "https://access.line.me/oauth2/v2.1/authorize",
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
		},
		RedirectURL: "http://localhost:8080/auth/callback",
	}
}

func openServerForReceiveCallback(state string, sendCode chan<- string) {
	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query().Get("error"))  != 0 {
			fmt.Printf("cannot login to line with error: %s, please try again\n", r.URL.Query().Get("error_description"))
		}

		responseState := r.URL.Query().Get("state")

		if state != responseState {
			fmt.Printf("there is an error in line login. please try again")
			os.Exit(1)
		}

		sendCode <- r.URL.Query().Get("code")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetAccessToken() *oauth2.Token {
	ctx := context.Background()

	config := createOAuth2Config()
	// @todo need random string for state
	state := "state"
	url := config.AuthCodeURL(state)
	fmt.Println(url)

	sendCode := make(chan string)
	go openServerForReceiveCallback(state, sendCode)

	code := <-sendCode
	token, err := config.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("there is an error in line login. please try again")
		os.Exit(1)
	}

	return token
}