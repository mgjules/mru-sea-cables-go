package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	jsoniter "github.com/json-iterator/go"
	"github.com/mgjules/mru-sea-cables-go/cable"
	"github.com/mgjules/mru-sea-cables-go/config"
	"github.com/mgjules/mru-sea-cables-go/logger"
	"github.com/mgjules/speedtest"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	configfile = "mru-cables.toml"
	datafile   = "realtime.json"
	writePerm  = 0o600
)

// Result represents a speedtest result
type Result struct {
	Latency   float64   `json:"latency"`
	DLSpeed   float64   `json:"download"`
	ULSpeed   float64   `json:"upload"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	// if config file does not exist, create from .dist version
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %v", err)
	}

	cfg, err := config.LoadConfig(configfile)
	if err != nil {
		log.Fatal(err)
	}

	sugaredLogger := logger.New(cfg.DEV, cfg.Debug)

	client, err := speedtest.NewDefaultClient()
	if err != nil {
		sugaredLogger.Fatalf("error creating client: %v", err)
	}

	if len(cfg.DLSizes) > 0 {
		client.DLSizes = cfg.DLSizes
	}

	if len(cfg.ULSizes) > 0 {
		client.ULSizes = cfg.ULSizes
	}

	results := make(map[string]Result)

	for i := range cfg.Cables {
		c := &cfg.Cables[i]

		cbl, cblErr := cable.New(c.Name, client, sugaredLogger)
		if cblErr != nil {
			sugaredLogger.Fatalf("Failed creating new cable: %v", cblErr)
		}

		for _, s := range c.Servers {
			if sErr := cbl.AddServer(s); err != nil {
				sugaredLogger.Errorf("Failed adding new server: %v", sErr)
			}
		}

		latency := cbl.Latency()
		dlspeed := cbl.DLSpeed()
		ulspeed := cbl.ULSpeed()

		results[cbl.Name()] = Result{
			Latency:   latency,
			DLSpeed:   dlspeed,
			ULSpeed:   ulspeed,
			Timestamp: time.Now(),
		}

		sugaredLogger.Debugf(
			"[%s] Latency: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps",
			cbl.Name(),
			latency,
			dlspeed,
			ulspeed,
		)
	}

	data, err := json.Marshal(results)
	if err != nil {
		sugaredLogger.Fatalf("error marshalling results: %v", err)
	}

	if err := ioutil.WriteFile("data/"+datafile, data, writePerm); err != nil {
		sugaredLogger.Fatalf("error writing results: %v", err)
	}

	if err := saveToGist(context.TODO(), cfg.GistID, cfg.GithubToken, data, sugaredLogger); err != nil {
		sugaredLogger.Fatal(err)
	}
}

// saveToGist saves data using given id and token
func saveToGist(ctx context.Context, id, token string, data []byte, sugaredLogger *zap.SugaredLogger) error {
	if id != "" && token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		client := github.NewClient(oauth2.NewClient(ctx, ts))

		gistLogger := sugaredLogger.With("gist id", id)

		// Retrieve gist by ID
		gistLogger.Debug("retrieving gist...")

		gist, _, err := client.Gists.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("error retrieving gist: %w", err)
		}

		gistLogger.Debug("retrieved gist: %s", gist.ID)

		// Change relevant file in gist
		for i := range gist.Files {
			gistFile := gist.Files[i]
			if *gistFile.Filename != datafile {
				continue
			}

			*gistFile.Content = string(data)
		}

		// Saving changes in gist
		gistLogger.Debug("saving gist...")

		if _, _, err := client.Gists.Edit(ctx, *gist.ID, gist); err != nil {
			return fmt.Errorf("error saving gist: %w", err)
		}

		gistLogger.Debug("saved gist")
	}

	return nil
}
