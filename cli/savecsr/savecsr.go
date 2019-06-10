// Package savecsr implements the savecsr command.
package savecsr

import (
//	"encoding/hex"
	"encoding/pem"
	"errors"
//	"time"

	"github.com/cloudflare/cfssl/certdb"
	"github.com/cloudflare/cfssl/certdb/dbconf"
	"github.com/cloudflare/cfssl/certdb/sql"
	"github.com/cloudflare/cfssl/cli"
	cferr "github.com/cloudflare/cfssl/errors"
//	"github.com/cloudflare/cfssl/helpers"
//	"github.com/cloudflare/cfssl/log"
//	"github.com/cloudflare/cfssl/ocsp"
)

// Usage text of 'cfssl savecsr'
var savecsrUsageText = `cfssl savecsr -- saves certificate signing request
for later retrieval by signer instance

Usage of savecsr:
        cfssl savecsr -db-config db-config CSR

Arguments:
        CSR:        PEM file for certificate request, use '-' for reading PEM from stdin.

Note: CSR can also be supplied via flag values; flag value will take precedence over the argument.

Flags:
`

// Flags of 'cfssl savecsr'
var savecsrFlags = []string{"csr", "db-config"}

func savecsrMain(args []string, c cli.Config) (err error) {
	if c.CSRFile == "" {
		c.CSRFile, args, err = cli.PopFirstArgument(args)
		if err != nil {
			return errors.New("need CSR file(provide with -csr flag or CSR argument)")
		}
	}

	if len(args) > 0 {
		return errors.New("too many arguments are provided, please check with usage")
	}

	csr, err := cli.ReadStdin(c.CSRFile)
	if err != nil {
		return errors.New("failed to read CSR file")
	}

	block, _ := pem.Decode(csr)
	if block == nil {
		return cferr.New(cferr.CSRError, cferr.DecodeFailed)
	}

	if block.Type != "NEW CERTIFICATE REQUEST" && block.Type != "CERTIFICATE REQUEST" {
		return cferr.Wrap(cferr.CSRError,
			cferr.BadRequest, errors.New("not a csr"))
	}

	if c.DBConfigFile == "" {
		return errors.New("need DB config file (provide with -db-config)")
	}

	db, err := dbconf.DBFromConfig(c.DBConfigFile)
	if err != nil {
		return err
	}

	dbAccessor := sql.NewAccessor(db)

	var csrRecord = certdb.CSRRecord{
		CSR: string(csr),
	}

	err = dbAccessor.InsertCSR(csrRecord)
	if err != nil {
		return err
	}

	return nil
}

// Command assembles the definition of Command 'savecsr'
var Command = &cli.Command{UsageText: savecsrUsageText, Flags: savecsrFlags, Main: savecsrMain}
