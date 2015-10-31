package config

type Config struct {
	CfgName              string
	StorageInterfaceName string
	Storage              struct {
		Region string
	}
	Version      string
	OutputFormat string
}
