package main

import (
	"encoding/json"
)


type JSONPreset map[string]PresetEntry

type PresetEntry struct {
	Enabled bool `json:"en"`
	
	// May be nil
	Params map[string]string `json:"p,omitempty"`
}

func UnmarshalJSONPreset(raw string) (JSONPreset, error) {
	var entries JSONPreset
    err := json.Unmarshal([]byte(raw), &entries)
    return entries, err
}

func GetDefaultPreset() JSONPreset {
	return JSONPreset{
		"os-packages": {Enabled: true},
		"package-zsh": {Enabled: true},
		"package-tmux": {Enabled: true},
		"package-vim": {Enabled: true},
		"package-curl": {Enabled: true},
		"package-wget": {Enabled: true},
	}
}


type jsonUninstallPresetRaw []string
type JSONUninstallPreset map[string]bool

func UnmarshalJSONUninstallPreset(raw string) (JSONUninstallPreset, error) {
	var entries jsonUninstallPresetRaw
    err := json.Unmarshal([]byte(raw), &entries)
	if err != nil {
		return nil, err
	}

	resultPreset := make(JSONUninstallPreset, 0)
	for _, v := range entries {
		resultPreset[v] = true
	}

    return resultPreset, err
}
