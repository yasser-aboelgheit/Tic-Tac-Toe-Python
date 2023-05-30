package tracer

type Config struct {
	ReceiverEndpoint string `mapstructure:"receiver0endpoint"`
}

func (config *Config) Defaults() {
	config.ReceiverEndpoint = defaultReceiver
}
