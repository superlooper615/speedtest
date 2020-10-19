package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/superlooper615/speedtest/config"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// speedtest configuration.
type speedtest struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	lastIndexTime time.Time

}

// New creates an instance of speedtest.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &speedtest{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}




func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func setTimeout() {
	if *timeoutOpt != 0 {
		timeout = *timeoutOpt
	}
}

var (
	showList   = kingpin.Flag("list", "Show available speedtest.net servers").Short('l').Bool()
	serverIds  = kingpin.Flag("server", "Select server id to speedtest").Short('s').Ints()
	timeoutOpt = kingpin.Flag("timeout", "Define timeout seconds. Default: 10 sec").Short('t').Int()
	timeout    = 10
)


// Run starts speedtest.
func (bt *speedtest) Run(b *beat.Beat) error {
	logp.Info("speedtest is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		now := time.Now()
		kingpin.Version("1.0.3")
		kingpin.Parse()
	
		setTimeout()
	
		user := fetchUserInfo()
		user.Show()
	
		list := fetchServerList(user)
		if *showList {
			list.Show()
			// return
		}
	
		targets := list.FindServer(*serverIds)
		targets.StartTest()
		download, upload := targets.ShowResult()
		bt.lastIndexTime = now
		logp.Info("Event sent")
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"download": download,
				"upload": upload,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
		counter++
	}
}

// Stop stops speedtest.
func (bt *speedtest) Stop() {
	bt.client.Close()
	close(bt.done)
}
