package velonetics

import (
	"encoding/json"
	"errors"

	botdetector "github.com/velonetics/velonetics-botdetector/v2"
	"github.com/velonetics/lura/v2/config"
)

// Namespace is the key used to store the bot detector config at the ExtraConfig struct
const Namespace = "github_com/velonetics/velonetics-botdetector"

// ErrNoConfig is returned when there is no config defined for the module
var ErrNoConfig = errors.New("no config defined for the module")

// ParseConfig extracts the module config from the ExtraConfig and returns a struct
// suitable for using the botdetector package
func ParseConfig(cfg config.ExtraConfig) (botdetector.Config, error) {
	res := botdetector.Config{}
	e, ok := cfg[Namespace]
	if !ok {
		return res, ErrNoConfig
	}
	b, err := json.Marshal(e)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(b, &res)
	return res, err
}
