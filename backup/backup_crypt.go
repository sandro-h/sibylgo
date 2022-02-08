package backup

import (
	"io"
	"io/ioutil"

	vault "github.com/sosedoff/ansible-vault-go"
)

// Cryptor provides methods to encrypt and decrypt backup content.
type Cryptor interface {
	EncryptContent(in io.Reader, out io.Writer) error
	DecryptContent(in io.Reader, out io.Writer) error
}

// AnsibleCryptor uses the same mechanism as Ansible Vaults to
// encrypt backup content using a user-provided password.
type AnsibleCryptor struct {
	Password string
}

// EncryptContent encrypts the backup content using Ansible Vault style encryption
// with the user-provided password.
func (c *AnsibleCryptor) EncryptContent(in io.Reader, out io.Writer) error {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	secret := string(data)

	str, err := vault.Encrypt(secret, c.Password)
	if err != nil {
		return err
	}

	out.Write([]byte(str))
	return nil
}

// DecryptContent decrypts the backup content using Ansible Vault style encryption
// with the user-provided password.
func (c *AnsibleCryptor) DecryptContent(in io.Reader, out io.Writer) error {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	str, err := vault.Decrypt(string(data), c.Password)
	if err == nil {
		out.Write([]byte(str))
	} else {
		out.Write(data)
	}

	return nil
}
