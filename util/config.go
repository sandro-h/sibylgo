package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

// Config holds configuration values and has convenience methods to access these.
type Config struct {
	cfg map[interface{}]interface{}
}

// ConfigPath defines a path to a specific config value, where each part is a key in the next sub-config.
type ConfigPath struct {
	Parts []string
}

// GetString returns the string value for key, or defaultVal if the key is not found.
func (cfg Config) GetString(key string, defaultVal string) string {
	v, found := cfg.cfg[key]
	if found {
		return v.(string)
	}
	return defaultVal
}

// GetStringOrFail returns the string value for key, or exits the program with error code 1 if it doesn't exist.
func (cfg Config) GetStringOrFail(key string) string {
	return cfg.getOrFail(key).(string)
}

func (cfg Config) getOrFail(key string) interface{} {
	v, found := cfg.cfg[key]
	if !found {
		fmt.Fprintf(os.Stderr, "%s must be set\n", key)
		os.Exit(1)
	}
	return v
}

// SetString sets the config key to string val.
func (cfg Config) SetString(key string, val string) {
	cfg.cfg[key] = val
}

// GetStringList returns the value for key, or defaultVal if the key is not found.
func (cfg Config) GetStringList(key string, defaultVal []string) []string {
	v, found := cfg.cfg[key]
	if found {
		il := v.([]interface{})
		var sl []string
		for _, ie := range il {
			sl = append(sl, ie.(string))
		}
		return sl
	}
	return defaultVal
}

// GetInt returns the int value for key, or defaultVal if the key is not found.
func (cfg Config) GetInt(key string, defaultVal int) int {
	v, found := cfg.cfg[key]
	if found {
		return v.(int)
	}
	return defaultVal
}

// GetIntOrFail returns the int value for key, or exits the program with error code 1 if it doesn't exist.
func (cfg Config) GetIntOrFail(key string) int {
	return cfg.getOrFail(key).(int)
}

// HasKey returns true if the config contains the key.
func (cfg Config) HasKey(key string) bool {
	_, found := cfg.cfg[key]
	return found
}

// GetSubConfig returns the part of the config under key as a new Config object.
// If the key is not found an empty config is returned.
func (cfg Config) GetSubConfig(key string) *Config {
	v, found := cfg.cfg[key]
	if found {
		return &Config{v.(map[interface{}]interface{})}
	}
	return &Config{}
}

// GetByPath returns the value for the config path, or defaultVal if the key is not found.
func (cfg Config) GetByPath(path ConfigPath, defaultVal interface{}) interface{} {
	v, found := getByPath(cfg.cfg, 0, path.Parts)
	if !found {
		return defaultVal
	}
	return v
}

// SetByPath sets the value for the config path, if the path points to a correct sub-config.
// If the sub-config is not found, nothing is written - it does not create intermediate sub-configs.
// If the sub-config is found but doesn't have the final part of the path as a key yet, that key is created.
func (cfg Config) SetByPath(path ConfigPath, val interface{}) bool {
	pathWithoutLast := path.Parts[:len(path.Parts)-1]
	pathLast := path.Parts[len(path.Parts)-1]
	v, found := getByPath(cfg.cfg, 0, pathWithoutLast)
	if found {
		subCfg, ok := v.(map[interface{}]interface{})
		if ok {
			subCfg[pathLast] = val
			return true
		}
	}
	return false
}

func getByPath(cur map[interface{}]interface{}, curPathIndex int, path []string) (interface{}, bool) {
	v, found := cur[path[curPathIndex]]
	if !found {
		return nil, false
	}
	if curPathIndex == len(path)-1 {
		return v, true
	}
	subCfg, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, false
	}
	return getByPath(subCfg, curPathIndex+1, path)
}

// LoadConfig loads a Config object from a YAML file.
func LoadConfig(cfgFile string) (*Config, error) {
	if _, err := os.Stat(cfgFile); !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return nil, err
		}
		return loadConfigBytes(data)
	}
	return &Config{}, nil
}

// LoadConfigString loads a Config object from a YAML string.
func LoadConfigString(str string) (*Config, error) {
	return loadConfigBytes([]byte(str))
}

func loadConfigBytes(data []byte) (*Config, error) {
	var cfg map[interface{}]interface{}
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &Config{cfg}, nil
}

// SaveConfig saves the config object to a YAML file.
func SaveConfig(cfg *Config, cfgFile string) error {
	data, err := saveConfigBytes(cfg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cfgFile, data, 0644)
	return err
}

// SaveConfigString saves the config object to a YAML string.
func SaveConfigString(cfg *Config) (string, error) {
	data, err := saveConfigBytes(cfg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func saveConfigBytes(cfg *Config) ([]byte, error) {
	return yaml.Marshal(cfg.cfg)
}

// ToString converts the config path to a dot-separated string
func (p *ConfigPath) ToString() string {
	return strings.Join(p.Parts, ".")
}

// FromString sets the config path to the path denoted by the dot-separated string.
func (p *ConfigPath) FromString(str string) {
	p.Parts = strings.Split(str, ".")
}
