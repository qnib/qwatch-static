package qserver

import (
	"log"
	"os"
	"os/signal"
	"time"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-filter-id/lib"
	//"github.com/qnib/qframe-handler-log/lib"
	"github.com/qnib/qframe-handler-elasticsearch/lib"
	"github.com/qnib/qframe-collector-file/lib"
)

// Run start daemon
func Run(ctx *cli.Context) {
	// Create conf
	log.Printf("[II] Start Version: %s", ctx.App.Version)
	cfg := config.NewConfig([]config.Provider{})
	if _, err := os.Stat(ctx.String("config")); err == nil {
		log.Printf("[II] Use config file: %s", ctx.String("config"))
		cfg.Providers = append(cfg.Providers, config.NewYAMLFile(ctx.String("config")))
	} else {
		log.Printf("[II] No config file found")
	}
	cfg.Providers = append(cfg.Providers, config.NewCLI(ctx, false))
	// Create chan
	qChan := qtypes.NewQChan()
	// Create ticker
	i, _ := cfg.Int("ticker.interval")
	interval := time.Duration(i) * time.Millisecond
	ticker := time.NewTicker(interval).C
	// Create Broadcaster goroutine
	qChan.Broadcast()
	// fetches interrupt and closes
	signal.Notify(qChan.Done, os.Interrupt)
	// instanciate handlers,filters,collectors
	//// Handlers
	hEsLog := qframe_handler_elasticsearch.NewElasticsearch(qChan, *cfg, "es_log")
	go hEsLog.Run()
	hEsEvents := qframe_handler_elasticsearch.NewElasticsearch(qChan, *cfg, "es_events")
	go hEsEvents.Run()
	//// Filters
	fi := qframe_filter_id.New(qChan, *cfg, "id")
	go fi.Run()
	//// Inputs
	cf := qframe_collector_file.NewPlugin(qChan, *cfg, "file")
	go cf.Run()
	// Inserts tick to get Inventory started
	var tickCnt int64
	var endTick int64
	eTick, _ := cfg.Int("ticks")
	endTick = int64(eTick)
	qChan.Tick.Send(tickCnt)
	time.Sleep(100 * time.Millisecond)
	for {
		select {
		case <-qChan.Done:
			fmt.Printf("\nDone\n")
			os.Exit(0)
		case <-ticker:
			tickCnt++
			if endTick != 0 && tickCnt == endTick {
				log.Printf("[II] End loop as tick-cnt '%d' reaches ticks '%d'", tickCnt, endTick)
			}
			qChan.Tick.Send(tickCnt)

		}
	}
}
