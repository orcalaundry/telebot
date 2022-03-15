package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TELEBOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	s := newServer(os.Getenv("API_KEY"))
	// The server.registerRoutes() method handles all the logic for registering
	// handlers.
	s.registerRoutes(b)

	b.Start()
}

type house string

const (
	leo house = "14"
)

// A map of routes to houses. Each key should be prefixed with a "/", since it
// should serve as a route for the bot.
var houses map[string]house = map[string]house{
	"/leo": leo,
}

// server represents a worker for this bot. Although we are not writing a
// server.
type server struct {
	client *http.Client
	apiKey string
}

func newServer(apiKey string) *server {
	return &server{
		client: &http.Client{Timeout: 5 * time.Second},
		apiKey: apiKey,
	}
}

// registerRoutes contains all the logic needed to attach handlers to the
// tele.Bot instance.
func (s *server) registerRoutes(b *tele.Bot) {
	for route, house := range houses {
		b.Handle(route, s.newHouseHandler(house))
	}
}

// newHouseHandler returns a handler for the default house routes, e.g. "/leo".
// It is a closure on both the house and the GET request to make.
func (s *server) newHouseHandler(h house) tele.HandlerFunc {
	// Create a GET request, and re-use it for any future invocations of this
	// handler.
	r, _ := http.NewRequest("GET", os.Getenv("API_URL")+"/machine?floor="+string(h), nil)
	r.Header.Add("x-api-key", s.apiKey)

	return tele.HandlerFunc(func(c tele.Context) error {
		resp, err := s.client.Do(r)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var ms []*machine
		if err := json.NewDecoder(resp.Body).Decode(&ms); err != nil {
			return err
		}
		if err := c.Send(machines(ms).String()); err != nil {
			return err
		}
		return nil
	})
}
