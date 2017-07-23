package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/devopracy/devopracy-cli/command"
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/kardianos/osext"
	// "github.com/devopracy/devopracy-cli/command"
	"github.com/devopracy/devopracy-cli/devo"
	// "github.com/devopracy/devopracy-cli/devo/plugin"
)

// DEVOSPACE is used to represent the spaces that separate args for a command
// without being confused with spaces in the path to the command itself.
const DEVOSPACE = "-DEVOSPACE-"

type config struct {
	DisableCheckpoint          bool `json:"disable_checkpoint"`
	DisableCheckpointSignature bool `json:"disable_checkpoint_signature"`
	PluginMinPort              uint
	PluginMaxPort              uint

	Builders       map[string]string
	PostProcessors map[string]string `json:"post-processors"`
	Provisioners   map[string]string
}

// Decodes configuration in JSON format from the given io.Reader into
// the config object pointed to.
func decodeConfig(r io.Reader, c *config) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(c)
}

// Discover discovers plugins.
//
// Search the directory of the executable, then the plugins directory, and
// finally the CWD, in that order. Any conflicts will overwrite previously
// found plugins, in that order.
// Hence, the priority order is the reverse of the search order - i.e., the
// CWD has the highest priority.
func (c *config) Discover() error {
	// If we are already inside a plugin process we should not need to
	// discover anything.
	if os.Getenv(plugin.MagicCookieKey) == plugin.MagicCookieValue {
		return nil
	}

	// First, look in the same directory as the executable.
	exePath, err := osext.Executable()
	if err != nil {
		log.Printf("[ERR] Error loading exe directory: %s", err)
	} else {
		if err := c.discover(filepath.Dir(exePath)); err != nil {
			return err
		}
	}

	// Next, look in the plugins directory.
	dir, err := devo.ConfigDir()
	if err != nil {
		log.Printf("[ERR] Error loading config directory: %s", err)
	} else {
		if err := c.discover(filepath.Join(dir, "plugins")); err != nil {
			return err
		}
	}

	// Next, look in the CWD.
	if err := c.discover("."); err != nil {
		return err
	}

	// Finally, try to use an internal plugin. Note that this will not override
	// any previously-loaded plugins.
	if err := c.discoverInternal(); err != nil {
		return err
	}

	return nil
}

func (c *config) discover(path string) error {
	var err error

	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
	}

	err = c.discoverSingle(
		filepath.Join(path, "packer-builder-*"), &c.Builders)
	if err != nil {
		return err
	}

	err = c.discoverSingle(
		filepath.Join(path, "packer-post-processor-*"), &c.PostProcessors)
	if err != nil {
		return err
	}

	return c.discoverSingle(
		filepath.Join(path, "packer-provisioner-*"), &c.Provisioners)
}

func (c *config) discoverSingle(glob string, m *map[string]string) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	if *m == nil {
		*m = make(map[string]string)
	}

	prefix := filepath.Base(glob)
	prefix = prefix[:strings.Index(prefix, "*")]
	for _, match := range matches {
		file := filepath.Base(match)

		// One Windows, ignore any plugins that don't end in .exe.
		// We could do a full PATHEXT parse, but this is probably good enough.
		if runtime.GOOS == "windows" && strings.ToLower(filepath.Ext(file)) != ".exe" {
			log.Printf(
				"[DEBUG] Ignoring plugin match %s, no exe extension",
				match)
			continue
		}

		// If the filename has a ".", trim up to there
		if idx := strings.Index(file, "."); idx >= 0 {
			file = file[:idx]
		}

		// Look for foo-bar-baz. The plugin name is "baz"
		plugin := file[len(prefix):]
		log.Printf("[DEBUG] Discovered plugin: %s = %s", plugin, match)
		(*m)[plugin] = match
	}

	return nil
}

func (c *config) discoverInternal() error {
	// Get the packer binary path
	packerPath, err := osext.Executable()
	if err != nil {
		log.Printf("[ERR] Error loading exe directory: %s", err)
		return err
	}

	for builder := range command.Builders {
		_, found := (c.Builders)[builder]
		if !found {
			log.Printf("Using internal plugin for %s", builder)
			(c.Builders)[builder] = fmt.Sprintf("%s%splugin%spacker-builder-%s",
				packerPath, DEVOSPACE, DEVOSPACE, builder)
		}
	}

	for provisioner := range command.Provisioners {
		_, found := (c.Provisioners)[provisioner]
		if !found {
			log.Printf("Using internal plugin for %s", provisioner)
			(c.Provisioners)[provisioner] = fmt.Sprintf(
				"%s%splugin%spacker-provisioner-%s",
				packerPath, DEVOSPACE, DEVOSPACE, provisioner)
		}
	}

	for postProcessor := range command.PostProcessors {
		_, found := (c.PostProcessors)[postProcessor]
		if !found {
			log.Printf("Using internal plugin for %s", postProcessor)
			(c.PostProcessors)[postProcessor] = fmt.Sprintf(
				"%s%splugin%spacker-post-processor-%s",
				packerPath, DEVOSPACE, DEVOSPACE, postProcessor)
		}
	}

	return nil
}
