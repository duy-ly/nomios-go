package kafkapublisher

import "github.com/spf13/viper"

type KafkaPublisherConfig struct {
	Addrs []string
	Topic string
}

func loadConfig() KafkaPublisherConfig {
	return KafkaPublisherConfig{
		Addrs: viper.GetStringSlice("publisher.kafka.addrs"),
		Topic: viper.GetString("publisher.kafka.topic"),
	}
}
