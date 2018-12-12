/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package viperutil

import (
	"bytes"
	"github.com/spf13/viper"
	"os"
	"strings"
	"testing"
)

const Prefix = "VIPERUTIL"



func TestEnhancedExactUnmarshalKey(t *testing.T) {
	type Nested struct {
		Key     string
		BoolVar bool
	}

	type nestedKey struct {
		Nested Nested
	}

	yaml := "---\n" +
		"Top:\n" +
		"  Nested:\n" +
		"    Nested:\n" +
		"      Key: BAD\n" +
		"      BoolVar: true\n"

	envVar := "VIPERUTIL_TOP_NESTED_NESTED_KEY"
	envVal := "GOOD"
	os.Setenv(envVar, envVal)
	defer os.Unsetenv(envVar)

	viper.SetEnvPrefix(Prefix)
	defer viper.Reset()
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewReader([]byte(yaml))); err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	var uconf nestedKey
	if err := EnhancedExactUnmarshalKey("top.Nested", &uconf); err != nil {
		t.Fatalf("Failed to unmarshall: %s", err)
	}

	if uconf.Nested.Key != envVal {
		t.Fatalf(`Expected: "%s", Actual: "%s"`, envVal, uconf.Nested.Key)
	}

	if uconf.Nested.BoolVar != true {
		t.Fatalf(`Expected: "%t", Actual: "%t"`, true, uconf.Nested.BoolVar)
	}

}
