/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/adamiklukasz/go-ping"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/runner"
)

const (
	pluginName    = "ping"
	pluginVersion = "0.1.0"

	configKey = "config"
)

type config struct {
	TargetAddresses []string `json:"target_addresses"`
	Requests        int      `json:"requests"`
	IntervalSec     int      `json:"interval_sec"`
}

type pingCollector struct {
}

func (*pingCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	def.DefineMetric("/solarwinds/ping/avg_time", "ms", true, "Average duration of ping requests")
	def.DefineMetric("/solarwinds/ping/min_time", "ms", true, "Minimum duration of ping requests")
	def.DefineMetric("/solarwinds/ping/max_time", "ms", true, "Maximum duration of ping requests")
	def.DefineMetric("/solarwinds/ping/availability", "", true, "Target address is available via ping [0-no, 1-yes]")

	return nil
}

func (*pingCollector) Load(ctx plugin.Context) error {
	var cfg config
	b := ctx.RawConfig()

	err := json.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}

	ctx.Store(configKey, cfg)

	return nil
}

func (*pingCollector) Collect(ctx plugin.CollectContext) error {
	var cfg config
	wg := sync.WaitGroup{}

	err := ctx.LoadTo(configKey, &cfg)
	if err != nil {
		return err
	}

	for _, targetAddr := range cfg.TargetAddresses {
		wg.Add(1)

		go func(targetAddr string) {
			defer wg.Done()

			pingOk := 1
			tags := plugin.MetricTag("target", targetAddr)

			stat, err := ping.PingN(targetAddr, cfg.Requests, time.Second*time.Duration(cfg.IntervalSec))
			if err != nil {
				pingOk = 0
			} else {
				_ = ctx.AddMetric("/solarwinds/ping/avg_time", toMs(stat.Avg), tags)
				_ = ctx.AddMetric("/solarwinds/ping/min_time", toMs(stat.Min), tags)
				_ = ctx.AddMetric("/solarwinds/ping/max_time", toMs(stat.Max), tags)
			}

			_ = ctx.AddMetric("/solarwinds/ping/availability", pingOk, tags)
		}(targetAddr)
	}

	wg.Wait()
	return nil
}

func toMs(d time.Duration) int {
	return int(float64(d) / 1000000)
}

func main() {
	runner.StartCollectorWithContext(context.Background(), &pingCollector{}, pluginName, pluginVersion)
}
