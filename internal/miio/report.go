package miio

import (
	"encoding/json"
	"github.com/AlexxIT/openmiio_agent/internal/app"
	"github.com/AlexxIT/openmiio_agent/pkg/rpc"
)

var report struct {
	CloudStarts int             `json:"cloud_starts,omitempty"`
	CloudState  json.RawMessage `json:"cloud_state,omitempty"`
	CloudUptime *app.Uptime     `json:"cloud_uptime,omitempty"`
}

var cloudState string
type Params map[string]json.RawMessage

func miioReport(to int, req rpc.Message, res *rpc.Message) bool {
	if string(req["method"]) == `"local.query_status"` || string(req["method"]) == `"basis.network"` {
		if state := string((*res)["params"]); state != cloudState {
			cloudState = state

			// params is bytes slice with quotes
			report.CloudState = (*res)["params"]

			if state == `"cloud_connected"` {
				report.CloudStarts++
				report.CloudUptime = app.NewUptime()
			} else {
				report.CloudState = nil
			}
		}
		if app.IsAiot() {
			var params Params
			if err := json.Unmarshal(req["params"], &params); err == nil {
				if state := string(params["name"]); state != cloudState {
					cloudState = state
					// params is bytes slice with quotes
					report.CloudState = params["name"]

					if state == `"network_signal"` {
						report.CloudStarts++
						report.CloudUptime = app.NewUptime()
					} else {
						report.CloudState = nil
					}
				}
			}
		}
	}

	return false // because we don't change response
}
