package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/util/log"
)

var (
	err error
)

var (
	defaultRootPath         = "app"
	defaultConfigFilePrefix = "application-"
	etcdConf              	etcdConfig
	mysqlConf             	mysqlConfig
	profilesConf            profilesConfig
	m                       sync.RWMutex
	inited                  bool
	sp                      = string(filepath.Separator)
)

// Init Configs
func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Logf("[Init] config inited")
		return
	}

	// loading .yml files
	// get project path
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("."+sp, sp)))

	pt := filepath.Join(appPath, "conf")
	os.Chdir(appPath)

	// find application.yml file
	if err = config.Load(file.NewSource(file.WithPath(pt + sp + "application.yml"))); err != nil {
		panic(err)
	}

	// get profiles 
	if err = config.Get(defaultRootPath, "profiles").Scan(&profilesConf); err != nil {
		panic(err)
	}

	log.Logf("[Init] Loading config file, path: %s, %+v\n", pt+sp+"application.yml", profilesConf)

	// loading related config files
	if len(profilesConf.GetInclude()) > 0 {
		//split by "," 
		includes := strings.Split(profilesConf.GetInclude(), ",")

		sources := make([]source.Source, len(includes))
		for i := 0; i < len(includes); i++ {
			filePath := pt + string(filepath.Separator) + defaultConfigFilePrefix + strings.TrimSpace(includes[i]) + ".yml"

			log.Logf("[Init] Loading config file, path: %s\n", filePath)

			sources[i] = file.NewSource(file.WithPath(filePath))
		}

		
		if err = config.Load(sources...); err != nil {
			panic(err)
		}
	}

	// get value 
	config.Get(defaultRootPath, "etcd").Scan(&etcdConf)
	config.Get(defaultRootPath, "mysql").Scan(&mysqlConf)

	// intted
	inited = true
}

// GetMysqlConfig
func GetMysqlConfig() (ret mysqlConfig) {
	return mysqlConf
}

// GetEtcdConfig
func GetEtcdConfig() (ret etcdConfig) {
	return etcdConf
}
