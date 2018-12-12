/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package viperutil

import (
	"encoding/json"
	//"github.com/hyperledger/fabric/common/flogging"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

//var logger = flogging.MustGetLogger("viperutil")

type viperGetter func(key string) interface{}

func getKeysRecursively(base string, getKey viperGetter, nodeKeys map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key := range nodeKeys {
		fqKey := base + key
		val := getKey(fqKey)
		if m, ok := val.(map[interface{}]interface{}); ok {
			//logger.Debugf("Found map[interface{}]interface{} value for %s", fqKey)
			tmp := make(map[string]interface{})
			for ik, iv := range m {
				cik, ok := ik.(string)
				if !ok {
					panic("Non string key-entry")
				}
				tmp[cik] = iv
			}
			result[key] = getKeysRecursively(fqKey+".", getKey, tmp)
		} else if m, ok := val.(map[string]interface{}); ok {
			//logger.Debugf("Found map[string]interface{} value for %s", fqKey)
			result[key] = getKeysRecursively(fqKey+".", getKey, m)
		} else if m, ok := unmarshalJSON(val); ok {
			//logger.Debugf("Found real value for %s setting to map[string]string %v", fqKey, m)
			result[key] = m
		} else {
			if val == nil {
				fileSubKey := fqKey + ".File"
				fileVal := getKey(fileSubKey)
				if fileVal != nil {
					result[key] = map[string]interface{}{"File": fileVal}
					continue
				}
			}
			//logger.Debugf("Found real value for %s setting to %T %v", fqKey, val, val)
			result[key] = val

		}
	}
	return result
}

func unmarshalJSON(val interface{}) (map[string]string, bool) {
	mp := map[string]string{}
	s, ok := val.(string)
	if !ok {
		//logger.Debugf("Unmarshal JSON: value is not a string: %v", val)
		return nil, false
	}
	err := json.Unmarshal([]byte(s), &mp)
	if err != nil {
		//logger.Debugf("Unmarshal JSON: value cannot be unmarshalled: %s", err)
		return nil, false
	}
	return mp, true
}

// EnhancedExactUnmarshalKey is intended to unmarshal a config file subtreee into a structure
func EnhancedExactUnmarshalKey(baseKey string, output interface{}, conf *viper.Viper) error {
	m := make(map[string]interface{})
	m[baseKey] = nil
	leafKeys := getKeysRecursively("", conf.Get, m)

	//logger.Debugf("%+v", leafKeys)

	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(leafKeys[baseKey])
}
