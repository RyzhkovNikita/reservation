package conf

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"time"
)

type Mode string

const (
	Dev  Mode = "dev"
	Prod      = "prod"
	Test      = "test"
)

type Configuration struct {
	Secret               string
	DriverName           string
	DataSourceUrl        string
	Mode                 Mode
	FileStoragePath      string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	ServerPrefix         string
}

var AppConfig Configuration

func init() {
	errors := make([]error, 0, 10)
	secret, err := beego.AppConfig.String("secret")
	errors = append(errors, err)
	driverName, err := beego.AppConfig.String("driverName")
	errors = append(errors, err)
	dataSourceUrl, err := beego.AppConfig.String("dataSourceUrl")
	errors = append(errors, err)
	mode, err := beego.AppConfig.String("runmode")
	errors = append(errors, err)
	fileStoragePath, err := beego.AppConfig.String("filesStoragePath")
	errors = append(errors, err)
	accessTokenLifetimeSeconds, err := beego.AppConfig.Int("accessTokenLifetimeSeconds")
	errors = append(errors, err)
	refreshTokenLifetimeSeconds, err := beego.AppConfig.Int("refreshTokenLifetimeSeconds")
	errors = append(errors, err)
	serverPrefix, err := beego.AppConfig.String("serverPrefix")
	errors = append(errors, err)
	for _, er := range errors {
		if er != nil {
			panic(fmt.Errorf("config is not valid: %v", er))
		}
	}
	AppConfig = Configuration{
		Secret:               secret,
		DriverName:           driverName,
		DataSourceUrl:        dataSourceUrl,
		Mode:                 Mode(mode),
		FileStoragePath:      fileStoragePath,
		AccessTokenLifetime:  time.Duration(accessTokenLifetimeSeconds) * time.Second,
		RefreshTokenLifetime: time.Duration(refreshTokenLifetimeSeconds) * time.Second,
		ServerPrefix:         serverPrefix,
	}
}
