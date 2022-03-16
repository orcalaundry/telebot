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

	w := newWorker(os.Getenv("API_KEY"), b)
	// The worker.registerRoutes() method handles all the logic for registering
	// handlers.
	w.registerRoutes()

	w.start()
}

type house string

const (
	leo house = "14"
)

// A map of routes to houses. Each key should be prefixed with a "/", since it
// must serve as a route for the bot.
var houses map[string]house = map[string]house{
	"/leo": leo,
}

// worker is a wrapper around telebot.Bot.
type worker struct {
	client *http.Client
	apiKey string
	b      *tele.Bot
}

func newWorker(apiKey string, bot *tele.Bot) *worker {
	return &worker{
		client: &http.Client{Timeout: 5 * time.Second},
		apiKey: apiKey,
		b:      bot,
	}
}

func (w *worker) start() {
	w.b.Start()
}

// registerRoutes contains all the logic needed to attach handlers to the
// tele.Bot instance.
func (w *worker) registerRoutes() {
	// Register all the default routes, like "/leo", "/ursa", etc.
	for route, house := range houses {
		w.b.Handle(route, w.newHouseHandler(house))
	}
}

// newHouseHandler returns a handler for the default house routes, e.g. "/leo".
// It is a closure on both the house and the GET request to make.
func (w *worker) newHouseHandler(h house) tele.HandlerFunc {
	// Create a GET request, and re-use it for any future invocations of this
	// handler. This is safe to use since our request body is nil. See
	// https://github.com/golang/go/issues/19653#issuecomment-341539160.
	r, _ := http.NewRequest("GET", os.Getenv("API_URL")+"/machine?floor="+string(h), nil)
	// We don't actually need the X-API-KEY header for this request, but wouldn't
	// hurt to add it here anyway. Note that the key name is case-insensitive.
	r.Header.Add("x-api-key", w.apiKey)

	return tele.HandlerFunc(func(c tele.Context) error {
		resp, err := w.client.Do(r)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		// We rely on the json package to decode most fields, including an RFC3339
		// format time string for the machines' last started at time. The other
		// fields are patched in later.
		var ms []*machine
		if err := json.NewDecoder(resp.Body).Decode(&ms); err != nil {
			return err
		}
		// These fields must be set after the JSON is decoded.
		for _, m := range ms {
			m.addName()
			m.computeTimeLeft()
		}
		if err := c.Send(machines(ms).String()); err != nil {
			return err
		}
		return nil
	})
}
