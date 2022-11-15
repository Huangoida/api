package Init

import (
	"api/config"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	constant "github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func InitNacos() {
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NamespaceId:         "ab38261c-8904-4e9e-bbbf-9947e970f57b",
		CacheDir:            "cache",
		NotLoadCacheAtStart: false,
		LogDir:              "log",
		LogLevel:            "debug",
	}

	serverConfig := []constant.ServerConfig{
		{
			IpAddr:      config.GetConf().Nacos.Ip,
			ContextPath: config.GetConf().Nacos.ContextPath,
			Port:        config.GetConf().Nacos.Port,
			Scheme:      config.GetConf().Nacos.Scheme,
		},
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfig,
		},
	)
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "path.test",
		Group:  "path",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(content)

}
