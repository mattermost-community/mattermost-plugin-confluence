package main

import (
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/bot/logger"
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/telemetry"
)

const (
	keysPerPage           = 1000
	ConfluenceUserInfoKey = "_userinfo"
)

func (p *Plugin) TrackEvent(event string, properties map[string]interface{}) {
	err := p.tracker.TrackEvent(event, properties)
	if err != nil {
		p.API.LogDebug("Error sending telemetry event", "event", event, "error", err.Error())
	}
}

func (p *Plugin) TrackUserEvent(event, userID string, properties map[string]interface{}) {
	err := p.tracker.TrackUserEvent(event, userID, properties)
	if err != nil {
		p.API.LogDebug("Error sending user telemetry event", "event", event, "error", err.Error())
	}
}

// Initialize telemetry setups the tracker/clients needed to send telemetry data.
// The telemetry.NewTrackerConfig(...) param will take care of extract/parse the config to set rge right settings.
// If you don't want the default behavior you still can pass a different telemetry.TrackerConfig data.
func (p *Plugin) initializeTelemetry() {
	var err error

	// Telemetry client
	p.telemetryClient, err = telemetry.NewRudderClient()
	if err != nil {
		p.API.LogWarn("Telemetry client not started", "error", err.Error())
		return
	}

	// Get config values
	p.tracker = telemetry.NewTracker(
		p.telemetryClient,
		p.API.GetTelemetryId(),
		p.API.GetServerVersion(),
		manifest.Id,
		manifest.Version,
		"confluence",
		telemetry.NewTrackerConfig(p.API.GetConfig()),
		logger.New(p.API),
	)
}
