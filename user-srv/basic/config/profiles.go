package config

// Profiles Interface
type IProfiles interface {
	GetInclude() string
}

// profilesConfig struct
type profilesConfig struct {
	Include string `json:"include"`
}

func (config profilesConfig) GetInclude() string {
	return config.Include
}
