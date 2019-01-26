package config

import (
	"encoding/json"
	"fmt"
	"os"
	fp "path/filepath"
)

// Profile is struct containing path to qmake, moc, rcc, C compiler and C++ compiler
type Profile struct {
	OS     string
	Arch   string
	Static bool

	Qmake   string
	Moc     string
	Rcc     string
	Gcc     string
	Gxx     string
	Windres string
	Objdump string
}

// LoadProfiles load all profiles inside config file in ${XDG_CONFIG_HOME}/qamel/config
func LoadProfiles(configPath string) (map[string]Profile, error) {
	// If config file doesn't exist, return empty map
	if !fileExists(configPath) {
		return map[string]Profile{}, nil
	}

	// Open file
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	// Decode JSON
	profiles := map[string]Profile{}
	err = json.NewDecoder(configFile).Decode(&profiles)
	return profiles, err
}

// LoadProfile loads profile with specified name from config file
func LoadProfile(configPath string, name string) (Profile, error) {
	profiles, err := LoadProfiles(configPath)
	if err != nil {
		return Profile{}, err
	}

	if prof, ok := profiles[name]; ok {
		return prof, nil
	}

	return Profile{}, fmt.Errorf("profile %s doesn't exist", name)
}

// SaveProfiles saves the profile as JSON in ${XDG_CONFIG_HOME}/qamel/config
func SaveProfiles(configPath string, profiles map[string]Profile) error {
	// Make sure config dir is exists
	os.MkdirAll(fp.Dir(configPath), os.ModePerm)

	// Create file
	configFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// Encode profiles to JSON
	return json.NewEncoder(configFile).Encode(&profiles)
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return !os.IsNotExist(err) && !info.IsDir()
}
