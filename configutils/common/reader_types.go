package common

// ConfigReaderName represents the type of config reader.
type ConfigReaderName string

const (
	FileConfigReader      ConfigReaderName = "file"
	ConsulConfigReader    ConfigReaderName = "consul"
	ZookeeperConfigReader ConfigReaderName = "zookeeper"
	DatabaseConfigReader  ConfigReaderName = "database"
)
