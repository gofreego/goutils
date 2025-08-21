package common

// ConfigReaderName represents the type of config reader.
type ConfigReaderName string

const (
	ConsulConfigReader    ConfigReaderName = "consul"
	ZookeeperConfigReader ConfigReaderName = "zookeeper"
)
