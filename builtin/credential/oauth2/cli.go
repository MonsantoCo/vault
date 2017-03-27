package oauth2

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	pwd "github.com/hashicorp/vault/helper/password"
)

// CLIHandler struct
type CLIHandler struct{}

// Auth cli method
func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "oauth2"
	}

	username, ok := m["username"]
	if !ok {
		return "", fmt.Errorf("'username' var must be set")
	}
	password, ok := m["password"]
	if !ok {
		fmt.Printf("Password (will be hidden): ")
		var err error
		password, err = pwd.Read(os.Stdin)
		fmt.Println()
		if err != nil {
			return "", err
		}
	}

	data := map[string]interface{}{
		"password": password,
	}

	path := fmt.Sprintf("auth/%s/login/%s", mount, username)
	secret, err := c.Logical().Write(path, data)
	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("empty response from credential provider")
	}

	return secret.Auth.ClientToken, nil
}

func (h *CLIHandler) Help() string {
	help := `
The oauth2 backend allows you to authenticate with an oauth2 service.
To use it, first configure it through the "config" endpoint, and then
login by specifying username and password. If password is not provided
on the command line, it will be read from stdin.

    Example: vault auth -method=oauth2 username=jdoe

    `

	return strings.TrimSpace(help)
}
