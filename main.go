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

	"github.com/adamiklukasz/go-ping"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/runner"
)

const (
	pluginName    = "ping"
	pluginVersion = "0.0.2"

	targetAddr = "www.solarwinds.com"
)

type pingCollector struct {
}

func (c *pingCollector) Collect(ctx plugin.CollectContext) error {
	dur, err := ping.Ping(targetAddr)
	if err != nil {
		return err
	}

	durMs := int(float64(dur) / 1000000)

	ctx.AddMetric("/solarwinds/ping/time", durMs, plugin.MetricTag("target", targetAddr), plugin.MetricUnit("ms"))

	return nil
}

func main() {
	runner.StartCollectorWithContext(context.Background(), &pingCollector{}, pluginName, pluginVersion)
}
