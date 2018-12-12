package msp

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	mspmgmt "github.com/DSiSc/msp/msp/mgmt"
	"github.com/DSiSc/msp/config"
)

func InitCrypto(mspConfig config.LocalMspConfig) error {
	var err error
	// Check whether msp folder exists
	fi, err := os.Stat(mspConfig.MspConfigPath)
	if os.IsNotExist(err) || !fi.IsDir() {
		// No need to try to load MSP from folder which is not available
		return errors.Errorf("cannot init crypto, folder \"%s\" does not exist", mspConfig.MspConfigPath)
	}
	// Check whether localMSPID exists
	if mspConfig.LocalMspId == "" {
		return errors.New("the local MSP must have an ID")
	}

	// Init the BCCSP
	SetBCCSPKeystorePath(mspConfig.KeystorePath)
	if err != nil {
		return errors.WithMessage(err, "could not parse YAML config")
	}

	err = mspmgmt.LoadLocalMspWithType(mspConfig.MspConfigPath, mspConfig.BccspConfig, mspConfig.LocalMspId, mspConfig.LocalMspType)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("error when setting up MSP of type %s from directory %s", mspConfig.LocalMspType, mspConfig.MspConfigPath))
	}

	return nil
}

// SetBCCSPKeystorePath sets the file keystore path for the SW BCCSP provider
// to an absolute path relative to the config file
func SetBCCSPKeystorePath(path string) {
	viper.Set("general.msp.BCCSP.SW.FileKeyStore.KeyStore",
		path)
}
/*// GetDefaultSigner return a default Signer(Default/PERR) for cli
func GetDefaultSigner() (msp.SigningIdentity, error) {
	signer, err := mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		return nil, errors.WithMessage(err, "error obtaining the default signing identity")
	}

	return signer, err
}
*/