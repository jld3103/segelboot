package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type Config struct {
	Entries []*Entry
}

type Entry struct {
	ID            string
	Name          string
	PartitionUUID string
	Loader        string
	CmdLine       string
}

func ReadConfigFile(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}
	c := &Config{}

	for _, section := range cfg.Sections() {
		if section.Name() != "DEFAULT" {
			for _, key := range []string{"name", "partition", "loader", "cmdline"} {
				if section.Key(key).String() == "" {
					//nolint:goerr113
					return nil, fmt.Errorf("section '%s' is missing key '%s'", section.Name(), key)
				}
			}
			c.Entries = append(c.Entries, &Entry{
				ID:            section.Name(),
				Name:          section.Key("name").String(),
				PartitionUUID: section.Key("partition").String(),
				Loader:        section.Key("loader").String(),
				CmdLine:       section.Key("cmdline").String(),
			})
		}
	}

	return c, nil
}
