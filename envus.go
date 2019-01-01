package envus

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"strings"
)

var (
	defaults   map[string]string
	properties map[string]string
	registry   map[string]*string
	envusFile  string
	loaded     bool
	useTier2   bool
)

func init() {
	refresh()
}

func refresh() {
	defaults = map[string]string{}
	properties = map[string]string{}
	registry = map[string]*string{}
	envusFile = ".envus"
	loaded = false
	useTier2 = false
}

// OverrideEnvusFilename overrides the default envus filename [.envus] with provided value
func OverrideEnvusFilename(filename string) {
	envusFile = filename
}

// UseLocalEnvironment directs envus to enable the local environment variables as a 2nd tier
func UseLocalEnvironment() {
	useTier2 = true
}

// Register enables the given key and default value to target string pointer. If the key is defined in the configured environment file
// (default is .envus) then a call to Load() will update string pointer.
func Register(key string, target *string, defaultValue string) {
	defaults[key] = defaultValue
	registry[key] = target
	*registry[key] = defaultValue
}

// Load parses through the defined environment file (defaults to .envus) and the local environment (if UseLocalEnvironment was called) and matches property
// values to registered and unregistered keys using defaults if registered
func Load() {
	//TODO this feature is not well tested. Currently I am supporting toml first, and json second but I am unsure whether
	// this might lead to problems or not in the future, and right now I can't do the math in my head

	var e envus
	// If the file fails to decode / unmarshall, that is a feature, so the envus object will merely contain an empty list of env items
	if _, err := toml.DecodeFile(envusFile, &e); err != nil {
		// try json
		if body, err := ioutil.ReadFile(envusFile); err != nil {
			// this is a feature, don't force the file to be right or existing
		} else {
			json.Unmarshal(body, &e)
		}
	}

	tier1 := map[string]string{}
	for _, item := range e.Prop {
		tier1[strings.ToUpper(item.Name)] = item.Value
	}

	for key, value := range defaults {
		if tier1Value, keyFoundInTier1 := tier1[key]; keyFoundInTier1 {
			properties[key] = tier1Value
			*registry[key] = tier1Value
		} else if !keyFoundInTier1 && useTier2 {
			if tier2 := os.Getenv(strings.ToUpper(key)); tier2 != "" {
				properties[key] = tier2
				*registry[key] = tier2
			} else {
				properties[key] = value
				*registry[key] = value
			}
		} else {
			properties[key] = value
			*registry[key] = value
		}
	}

	for _, item := range e.Prop {
		if _, propertyExists := properties[item.Name]; !propertyExists {
			properties[item.Name] = item.Value
		}
	}
	loaded = true
}

// Properties returns a map of strings representing all registered properties and unregistered properties found within
// environment file. This map does not contain local environment variables that exist outside of the environment file
// or within a call to Register
func Properties() map[string]string {
	if !loaded {
		Load()
	}
	return properties
}

// Fetch retrieves the value of a given key from calculated environment properties [See Properties]
func Fetch(key string) string {
	return properties[key]
}

type envus struct {
	Prop        []envusProperty `toml:"prop" json:"prop"`
}

type envusProperty struct {
	Name  string `toml:"key" json:"key"`
	Value string `toml:"value" json:"value"`
}
