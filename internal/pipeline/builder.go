package pipeline

import (
	"io"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/throttle"
)

// BuilderConfig holds high-level knobs used by Builder to construct a
// ready-to-use Pipeline without callers needing to wire every dependency.
type BuilderConfig struct {
	Output        io.Writer
	Format        string        // "text" | "csv"
	DedupeWindow  time.Duration // 0 disables dedupe
	ThrottleCool  time.Duration // 0 disables throttle
	AlertThresh   int           // 0 disables alerting
	WebhookURL    string        // "" disables webhook
	ExtraStages   []Stage
}

// Builder constructs a Pipeline from a BuilderConfig.
func Builder(bc BuilderConfig) *Pipeline {
	cfg := Config{Extra: bc.ExtraStages}

	if bc.Output != nil {
		cfg.Reporter = reporter.New(reporter.Config{
			Writer: bc.Output,
			Format: bc.Format,
		})
	}

	if bc.DedupeWindow > 0 {
		cfg.Dedupe = dedupe.New(dedupe.Config{Window: bc.DedupeWindow})
	}

	if bc.ThrottleCool > 0 {
		cfg.Throttle = throttle.New(throttle.Config{Cooldown: bc.ThrottleCool})
	}

	if bc.AlertThresh > 0 {
		cfg.Alert = alert.New(alert.Config{Threshold: bc.AlertThresh})
	}

	notifyCfg := notify.Config{}
	if bc.WebhookURL != "" {
		notifyCfg.Channels = []notify.Channel{
			notify.NewWebhookChannel(bc.WebhookURL),
		}
	}
	cfg.Notifier = notify.New(notifyCfg)

	return New(cfg)
}
