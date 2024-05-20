package configutils

import "github.com/spf13/viper"

func ReadConfig(filename string, config any) error {
	// Read the YAML file
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(config)
	if err != nil {
		return err
	}
	return nil
}
