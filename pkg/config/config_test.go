package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/IBM/argocd-vault-plugin/pkg/config"
	"github.com/spf13/viper"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		environment  map[string]interface{}
		expectedType string
	}{
		{
			map[string]interface{}{
				"AVP_TYPE":         "vault",
				"AVP_AUTH_TYPE":    "github",
				"AVP_GITHUB_TOKEN": "token",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"AVP_TYPE":      "vault",
				"AVP_AUTH_TYPE": "approle",
				"AVP_ROLE_ID":   "role_id",
				"AVP_SECRET_ID": "secret_id",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"AVP_TYPE":        "secretmanager",
				"AVP_AUTH_TYPE":   "iam",
				"AVP_IBM_API_KEY": "token",
			},
			"*backends.IBMSecretManager",
		},
	}
	for _, tc := range testCases {
		for k, v := range tc.environment {
			os.Setenv(k, v.(string))
		}
		viper := viper.New()
		config, err := config.New(viper)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		xType := fmt.Sprintf("%T", config.Backend)
		if xType != tc.expectedType {
			t.Errorf("expected: %s, got: %s.", tc.expectedType, xType)
		}
		for k := range tc.environment {
			os.Unsetenv(k)
		}
	}
}

func TestNewConfigNoType(t *testing.T) {
	viper := viper.New()
	_, err := config.New(viper)
	expectedError := "Must provide a supported Vault Type"

	if err.Error() != expectedError {
		t.Errorf("expected error %s to be thrown, got %s", expectedError, err)
	}
}

func TestNewConfigNoAuthType(t *testing.T) {
	os.Setenv("AVP_TYPE", "vault")
	viper := viper.New()
	_, err := config.New(viper)
	expectedError := "Must provide a supported Authentication Type"

	if err.Error() != expectedError {
		t.Errorf("expected error %s to be thrown, got %s", expectedError, err)
	}
	os.Unsetenv("AVP_TYPE")
}

func TestNewConfigMissingParameter(t *testing.T) {
	testCases := []struct {
		environment  map[string]interface{}
		expectedType string
	}{
		{
			map[string]interface{}{
				"AVP_TYPE":      "vault",
				"AVP_AUTH_TYPE": "github",
				"AVP_GH_TOKEN":  "token",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"AVP_TYPE":      "vault",
				"AVP_AUTH_TYPE": "approle",
				"AVP_ROLEID":    "role_id",
				"AVP_SECRET_ID": "secret_id",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"AVP_TYPE":        "secretmanager",
				"AVP_AUTH_TYPE":   "iam",
				"AVP_IAM_API_KEY": "token",
			},
			"*backends.IBMSecretManager",
		},
	}
	for _, tc := range testCases {
		for k, v := range tc.environment {
			os.Setenv(k, v.(string))
		}
		viper := viper.New()
		_, err := config.New(viper)
		if err == nil {
			t.Fatalf("%s should not instantiate", tc.expectedType)
		}
		for k := range tc.environment {
			os.Unsetenv(k)
		}
	}

}
