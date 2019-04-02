package line

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pkg/browser"
	"github.com/qapquiz/go-healthcheck/random"
	"golang.org/x/oauth2"
)

type lineLoginHandler struct {
	state string
	sendCode chan<- string
}

func (lineLoginHandler *lineLoginHandler) authLineCallback(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("error"))  != 0 {
		fmt.Printf("Cannot login to line with error: %s. Please try again\n", r.URL.Query().Get("error_description"))
	}

	responseState := r.URL.Query().Get("state")

	if lineLoginHandler.state != responseState {
		fmt.Printf("There is an error in Line login. please try again")
		os.Exit(1)
	}

	lineLoginHandler.sendCode <- r.URL.Query().Get("code")
}

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
	handler := lineLoginHandler{
		state,
		sendCode,
	}

	http.HandleFunc("/auth/callback", handler.authLineCallback)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetAccessToken() (*oauth2.Token, error) {
	ctx := context.Background()

	config := createOAuth2Config()

	state := random.RandStringBytesRemainder(7)
	url := config.AuthCodeURL(state)

	sendCode := make(chan string)
	go openServerForReceiveCallback(state, sendCode)
	err := browser.OpenURL(url)
	if err != nil {
		fmt.Println("Cannot open the url. You have to open it manually")
		fmt.Println("Login with Line URL: ", url)
	}

	code := <-sendCode
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}