package config

type QueueType string

const (
	RabbitMQ   = QueueType("RABBIT_MQ")
	GoChannels = QueueType("GO_CHANNELS")
)
