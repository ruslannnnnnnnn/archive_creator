package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	ArchiveStorageDirPath string   `mapstructure:"archive_storage_dir_path"`
	ObjectsInArchiveLimit int      `mapstructure:"objects_in_archive_limit"`
	ArchivesLimit         int      `mapstructure:"archives_limit"`
	Hostname              string   `mapstructure:"hostname"`
	Scheme                string   `mapstructure:"scheme"`
	Port                  int      `mapstructure:"port"`
	AvailableMimeTypes    []string `mapstructure:"available_mime_types"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения конфига: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ошибка маппинга конфига: %w", err)
	}

	return &config, nil
}
