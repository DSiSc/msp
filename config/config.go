package config

import (
	"github.com/DSiSc/msp/bccsp/factory"
)

type LocalMspConfig struct {
	MspConfigPath string
	LocalMspId string
	LocalMspType string
	BccspConfig *factory.FactoryOpts
	KeystorePath string

}
