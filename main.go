// Copyright © 2016 Christian Kniep <christian@qnib.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

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
	"github.com/qnib/qframe-handler-elasticsearch/lib"
	"github.com/qnib/qframe-collector-file/lib"
	"github.com/qnib/qframe-collector-gelf/lib"
)


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
	hEsLog := qframe_handler_elasticsearch.NewElasticsearch(qChan, *cfg, "es_logs")
	go hEsLog.Run()
	hEsEvents := qframe_handler_elasticsearch.NewElasticsearch(qChan, *cfg, "es_events")
	go hEsEvents.Run()
	//// Filters
	fi := qframe_filter_id.New(qChan, *cfg, "id")
	go fi.Run()
	//// Inputs
	cf := qframe_collector_file.NewPlugin(qChan, *cfg, "file")
	go cf.Run()
	cg := qframe_collector_gelf.NewPlugin(qChan, *cfg, "gelf")
	go cg.Run()
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


func main() {
	app := cli.NewApp()
	app.Name = "qwatch-static"
	app.Usage = "Statically compiled ETL framework for logs/events"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "qwatch.yml",
			Usage: "Config file, will overwrite flag default if present.",
		},
	}
	app.Action = Run
	app.Run(os.Args)
}
