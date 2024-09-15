package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"spacet/internal/domain"
	log "spacet/pkg/logger"
	pubsub "spacet/pkg/pubsub"
)

var attributes = map[string]string{
	"origin": "spacet-service",
	"source": "pubsub-notifier",
}

type pubsubEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type topics struct {
	launchesTopic pubsub.Topic
}

type gcpPubSubNotifier struct {
	l      log.Interface
	client pubsub.Interface
	topics *topics
}

// NewPubSubNotifierService implements NotificationService using gcp pubsub
func NewPubSubNotifierService(pubsubClient pubsub.Interface, logger log.Interface, launchesTopic string) domain.NotificationService {
	return &gcpPubSubNotifier{logger, pubsubClient, &topics{
		launchesTopic: pubsubClient.Topic(launchesTopic),
	}}
}

func (n *gcpPubSubNotifier) getTopic(event_type string) (pubsub.Topic, error) {
	switch event_type {
	case "BookingsCancelled":
		return n.topics.launchesTopic, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", event_type)
	}
}

// Publish uses the configured pubsub client to publish the notification.
func (n *gcpPubSubNotifier) Publish(ctx context.Context, event *domain.Event) error {
	topic, err := n.getTopic(event.Type)
	if err != nil {
		n.l.Error("PubSubNotifier Failed to getTopic: %s", err.Error())
		return domain.ErrNotificationNotSent
	}

	jsonEvent, err := json.Marshal(pubsubEvent{
		Type:    event.Type,
		Payload: json.RawMessage(event.Payload),
	})
	if err != nil {
		n.l.Error("PubSubNotifier Failed to marshall event: %s", err.Error())
		return domain.ErrFailedToProcessData
	}
	msg := &pubsub.Message{
		Data:       jsonEvent,
		Attributes: attributes,
	}
	result := topic.Publish(ctx, msg)

	// block until the result is returned and a server-generated ID is returned for the published message.
	if _, err = result.Get(ctx); err != nil {
		n.l.Debug("PubSubNotifier Failed to publish message: %s", err.Error())
		return domain.ErrNotificationNotSent
	}
	n.l.Debug("PubSubNotifier Published: %v", string(jsonEvent))
	return nil
}
