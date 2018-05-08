package cache

import "github.com/go-gl/mathgl/mgl32"

type Config struct {
	Ambient mgl32.Vec3
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config {
			mgl32.Vec3{0,0,0},
		}
	}

	return config
}