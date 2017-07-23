package devo

import (
	"os"
	"path/filepath"
)

// ConfigFile returns the default path to the configuration file. On
// Unix-like systems this is the ".Devoconfig" file in the home directory.
// On Windows, this is the "Devo.config" file in the application data
// directory.
func ConfigFile() (string, error) {
	return configFile()
}

// ConfigDir returns the configuration directory for Devo.
func ConfigDir() (string, error) {
	return configDir()
}

// ConfigTmpDir returns the configuration tmp directory for Devo
func ConfigTmpDir() (string, error) {
	if tmpdir := os.Getenv("DEVO_TMP_DIR"); tmpdir != "" {
		return filepath.Abs(tmpdir)
	}
	configdir, err := configDir()
	if err != nil {
		return "", err
	}
	td := filepath.Join(configdir, "tmp")
	_, err = os.Stat(td)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(td, 0755); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return td, nil
}
