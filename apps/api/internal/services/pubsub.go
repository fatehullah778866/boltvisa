package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type PubSubClient struct {
	client    *pubsub.Client
	topic     *pubsub.Topic
	projectID string
	topicName string
}

type NotificationMessage struct {
	UserID   uint                   `json:"user_id"`
	Type     string                 `json:"type"`
	Title    string                 `json:"title"`
	Message  string                 `json:"message"`
	Email    string                 `json:"email,omitempty"`
	Phone    string                 `json:"phone,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func NewPubSubClient(ctx context.Context, projectID, topicName, credentialsPath string) (*PubSubClient, error) {
	var client *pubsub.Client
	var err error

	if credentialsPath != "" {
		client, err = pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	} else {
		// Use default credentials (for Cloud Run/GKE)
		client, err = pubsub.NewClient(ctx, projectID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Pub/Sub client: %w", err)
	}

	topic := client.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check topic existence: %w", err)
	}

	if !exists {
		// Create topic if it doesn't exist
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to create topic: %w", err)
		}
		log.Printf("Created Pub/Sub topic: %s", topicName)
	}

	return &PubSubClient{
		client:    client,
		topic:     topic,
		projectID: projectID,
		topicName: topicName,
	}, nil
}

func (p *PubSubClient) PublishNotification(ctx context.Context, msg NotificationMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	result := p.topic.Publish(ctx, &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"type":    msg.Type,
			"user_id": fmt.Sprintf("%d", msg.UserID),
		},
	})

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published notification to Pub/Sub: user_id=%d, type=%s", msg.UserID, msg.Type)
	return nil
}

func (p *PubSubClient) Close() error {
	return p.client.Close()
}
