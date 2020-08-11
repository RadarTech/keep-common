package vault

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/ipfs/go-log"
)

var (
	client *api.Client
	logger = log.Logger("keep-vault")
	re     = regexp.MustCompile(`vault:.*:.*`)
)

func init() {
	var err error
	client, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		logger.Errorf("error connecting to vault: [%v]", err)
	}
}

// IsEnabled returns a bool for whether keypath matches expected
// format "vault:path:field"
func IsEnabled(keyPath string) bool {
	return re.MatchString(keyPath)
}

// ReadKeyFile read the specified path from Vault and returns a []byte
// of the field specified, VAULT_SECRET_PATH and VAULT_SECRET_FIELD
// define the path and field that should be looked up. This relies on
// the VAULT_ADDR and VAULT_TOKEN environment variables being set
func ReadKeyFile(keyPath string) ([]byte, error) {
	config := strings.Split(keyPath, `:`)

	path, field := config[1], config[2]

	fmt.Printf("%s, %s\n", path, field)

	data, err := client.Logical().Read(path)
	if err != nil {
		logger.Errorf("error reading path %s from vault: [%v]", path, err)
		return nil, err
	}

	wallet, _ := data.Data[field].(string)
	return []byte(wallet), nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
