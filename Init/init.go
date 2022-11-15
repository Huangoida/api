package Init

import "api/config"

func InitConfig() {
	config.ParseConf()
	InitNacos()
}
