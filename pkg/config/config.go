package config

type MozaikConfig struct {
	BasePath string
}

var (
	Config MozaikConfig
)

func init() {
	Config = MozaikConfig{
		BasePath: "./mozaikz",
	}
}
