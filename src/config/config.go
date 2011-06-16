// config implements functionality working with configuration file
// and make it easy to access its properties.
package config

// TODO: Clean all pathes and check if they are exist and accessable.

import (
	"os"
	"fmt"
)

// Configuration parameters parsed from config file and command line parameters.
var Configurations Config

// Possible types of the config item value.
const (
	typeString = iota
	typeInt    = iota
)

// item is the one configuration item representation struct.
type item struct {
	t     int
	value interface{}
}

// Map with default values.
var defaults = map[string]item{
	"fs.root": item{typeString, "/"},
}

// Config represents configuration file.
type Config struct {
	items map[string]item
}

// Parse parses configuration file and returns its representation.
func Parse(filename string) (config *Config, err os.Error) {
	config = new(Config)
	config.items = defaults

	// TODO: Parse file and fill the map.

	root := config.items["fs.root"]
	root.value = "foo"

	return config, nil
}

// GetString returns string value for the given key.
func (config *Config) GetString(key string) (value string, err os.Error) {
	i, present := config.items[key]
	if !present {
		return "", os.NewError(fmt.Sprintf("Key '%s' not found", key))
	}
	if i.t != typeString {
		return "", os.NewError(fmt.Sprintf("Key '%s' associated with not a string value", key))
	}

	return i.value.(string), nil
}

// SetString sets new value for the configuration item.
func (config *Config) SetString(key string, value string) os.Error {
	i, present := config.items[key]
	if !present {
		return os.NewError(fmt.Sprintf("Key '%s' not found", key))
	}
	if i.t != typeString {
		return os.NewError(fmt.Sprintf("Key '%s' associated with not a string value", key))
	}

	config.items[key] = item{i.t, value}

	return nil
}

// File initialization function.
func init() {
	filename := "/etc/chubd.conf"

	cfg, err := Parse(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse configuration file '%s'. %s", filename, err))
	}

	Configurations = *cfg
	Configurations.SetString("fs.root", "/home/viacheslav/projects/chubd_music")
}
