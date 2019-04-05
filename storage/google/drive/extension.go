package drive

import (
	"encoding/json"
	"fmt"
	"github.com/markdicksonjr/nibbler"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Extension struct {
	nibbler.NoOpExtension
	ConnectServiceOnInit bool
	CredentialsPath      string
	TokenFilePath        string

	Srv    *drive.Service
	Client *http.Client

	config *oauth2.Config
}

func (s *Extension) Init(app *nibbler.Application) error {
	if s.CredentialsPath == "" {
		s.CredentialsPath = app.GetConfiguration().Raw.Get("google.drive.credentials.path").String("drive-credentials.json")
	}

	if s.TokenFilePath == "" {
		s.TokenFilePath = app.GetConfiguration().Raw.Get("google.drive.tokenfile.path").String("drive-token.json")
	}

	b, err := ioutil.ReadFile(s.CredentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	s.config, err = google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	if s.ConnectServiceOnInit {
		return s.InitService()
	}

	return nil
}

func (s *Extension) InitService() error {
	client := getClient(s.config, s.TokenFilePath)

	var err error
	s.Srv, err = drive.New(client)

	return err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokenFilePath string) *http.Client {
	tok, err := tokenFromFile(tokenFilePath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFilePath, tok)
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
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
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
