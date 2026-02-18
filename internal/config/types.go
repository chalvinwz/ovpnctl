package config

type Profile struct {
	Name           string `yaml:"name" mapstructure:"name"`
	ConfigFile     string `yaml:"config_file" mapstructure:"config_file"`
	Username       string `yaml:"username" mapstructure:"username"`
	Password       string `yaml:"password" mapstructure:"password"`
	PrivateKeyPass string `yaml:"private_key_pass" mapstructure:"private_key_pass"`
}

type Config struct {
	Profiles []Profile `yaml:"profiles" mapstructure:"profiles"`
}
