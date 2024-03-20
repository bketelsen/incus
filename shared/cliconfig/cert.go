package cliconfig

import (
	"fmt"
	"io"
	"os"

	localtls "github.com/lxc/incus/shared/tls"
	"github.com/lxc/incus/shared/util"
)

// HasClientCertificate will return true if a client certificate has already been generated.
func (c *Config) HasClientCertificate() bool {
	certf := c.ConfigPath("client.crt")
	keyf := c.ConfigPath("client.key")
	if !util.PathExists(certf) || !util.PathExists(keyf) {
		return false
	}

	return true
}

// GenerateClientCertificate will generate the needed client.crt and client.key if needed.
func (c *Config) GenerateClientCertificate() error {
	if c.HasClientCertificate() {
		return nil
	}

	certf := c.ConfigPath("client.crt")
	keyf := c.ConfigPath("client.key")

	return localtls.FindOrGenCert(certf, keyf, true, false)
}

// CopyGlobalCert will copy global (system-wide) certificates to the user config path.
func (c *Config) CopyGlobalCert(src string, dst string) error {
	copyFile := func(oldPath string, newPath string) error {
		sourceFile, err := os.Open(oldPath)
		if err != nil {
			return err
		}

		defer sourceFile.Close()

		// Create new file.
		newFile, err := os.Create(newPath)
		if err != nil {
			return err
		}

		defer newFile.Close()

		_, err = io.Copy(newFile, sourceFile)
		if err != nil {
			return err
		}

		return nil
	}

	// Server certificate.
	oldPath := c.GlobalConfigPath("servercerts", fmt.Sprintf("%s.crt", src))
	if util.PathExists(oldPath) {
		newPath := c.ConfigPath("servercerts", fmt.Sprintf("%s.crt", dst))

		err := copyFile(oldPath, newPath)
		if err != nil {
			return err
		}
	}

	// Client certificate.
	oldPath = c.GlobalConfigPath("clientcerts", fmt.Sprintf("%s.crt", src))
	if util.PathExists(oldPath) {
		newPath := c.ConfigPath("clientcerts", fmt.Sprintf("%s.crt", dst))

		err := copyFile(oldPath, newPath)
		if err != nil {
			return err
		}
	}

	// Client key.
	oldPath = c.GlobalConfigPath("clientcerts", fmt.Sprintf("%s.key", src))
	if util.PathExists(oldPath) {
		newPath := c.ConfigPath("clientcerts", fmt.Sprintf("%s.key", dst))

		err := copyFile(oldPath, newPath)
		if err != nil {
			return err
		}
	}

	// Client CA.
	oldPath = c.GlobalConfigPath("clientcerts", fmt.Sprintf("%s.ca", src))
	if util.PathExists(oldPath) {
		newPath := c.ConfigPath("clientcerts", fmt.Sprintf("%s.ca", dst))

		err := copyFile(oldPath, newPath)
		if err != nil {
			return err
		}
	}

	return nil
}
