package pubsub

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/open-source-cloud/realtime/channels"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type (
	KafkaConfig struct {
		Username   string
		Password   string
		ServerAddr string
		Topic      string
		GroupId    string
	}
	KafkaAdapter struct {
		channels.ProducerAdapter
		channels.ConsumerAdapter
		ctx    context.Context
		config *KafkaConfig
		writer *kafka.Writer
		reader *kafka.Reader
	}
)

func NewKafkaAdapter(ctx context.Context, config *KafkaConfig) (*KafkaAdapter, error) {
	k := &KafkaAdapter{
		ctx:    ctx,
		config: config,
		writer: nil,
		reader: nil,
	}
	return k, nil
}

func (k *KafkaAdapter) Send(message *channels.Message) error {
	if k.writer == nil {
		mechanism, err := scram.Mechanism(scram.SHA256, k.config.Username, k.config.Password)
		if err != nil {
			return err
		}
		k.writer = &kafka.Writer{
			Addr:  kafka.TCP(k.config.ServerAddr),
			Topic: k.config.Topic,
			Transport: &kafka.Transport{
				SASL: mechanism,
				TLS:  &tls.Config{},
			},
		}
	}
	key := []byte(message.ID)
	val, err := message.MessageToJSON()
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Value: []byte(val),
		Key:   key,
	}
	err = k.writer.WriteMessages(context.Background(), msg)
	log.Printf("message sent to kafka, topic %s, message: %s, err: %v", k.config.Topic, val, err)
	return err
}

func (k *KafkaAdapter) Subscribe(client *channels.Client) error {
	if k.reader == nil {
		mechanism, err := scram.Mechanism(scram.SHA512, k.config.Username, k.config.Password)
		if err != nil {
			return err
		}
		k.reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{k.config.ServerAddr},
			GroupID: k.config.GroupId,
			Topic:   k.config.Topic,
			Dialer: &kafka.Dialer{
				SASLMechanism: mechanism,
				TLS:           &tls.Config{},
			},
		})
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		for {
			message, err := k.reader.ReadMessage(ctx)
			if err != nil {
				log.Panicf("error reading message from kafka, topic %s, err: %v", k.config.Topic, err)
			}
			rawMsg := string(message.Value)
			msg, err := channels.MessageFromJSON(rawMsg)
			if err != nil {
				log.Panicf("error deserializing message from json, topic %s, err: %v", k.config.Topic, err)
			}
			client.MessageChan() <- msg
		}
	}()
	return nil
}
