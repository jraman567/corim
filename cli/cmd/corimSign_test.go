// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

var (
	// note: embedded CoSWIDs are not validated {0: h'5C57E8F446CD421B91C908CF93E13CFC', 1: [505(h'deadbeef')]}
	testCorimValid = comid.MustHexDecode(nil, "a200505c57e8f446cd421b91c908cf93e13cfc0181d901f944deadbeef")
	// {0: h'5C57E8F446CD421B91C908CF93E13CFC'}
	testCorimInvalid = comid.MustHexDecode(nil, "a100505c57e8f446cd421b91c908cf93e13cfc")
	testMetaInvalid  = []byte("{}")
	testMetaValid    = []byte(`{
		"signer": {
			"name": "ACME Ltd signing key",
			"uri": "https://acme.example"
		},
		"validity": {
			"not-before": "2021-12-31T00:00:00Z",
			"not-after": "2025-12-31T00:00:00Z"
		}
	}`)
	testECKey = []byte(`{
		"kty": "EC",
		"crv": "P-256",
		"x": "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y": "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
		"d": "870MB6gfuTJ4HtUnUvYMyJpr5eUZNP4Bk43bVdj3eAE",
		"use": "enc",
		"kid": "1"
	  }`)
)

func Test_CorimSignCmd_unknown_argument(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CorimSignCmd_mandatory_args_missing_corim_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--key-file=ignored.jwk",
		"--meta-file=ignored.json",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM supplied")
}

func Test_CorimSignCmd_mandatory_args_missing_meta_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ignored.cbor",
		"--key-file=ignored.jwk",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM Meta supplied")
}

func Test_CorimSignCmd_mandatory_args_missing_key_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ignored.cbor",
		"--meta-file=ignored.json",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no key supplied")
}

func Test_CorimSignCmd_non_existent_unsigned_corim_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=nonexistent.cbor",
		"--key-file=ignored.jwk",
		"--meta-file=ignored.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()

	err := cmd.Execute()
	assert.EqualError(t, err, "error loading unsigned CoRIM from nonexistent.cbor: open nonexistent.cbor: file does not exist")
}

func Test_CorimSignCmd_bad_unsigned_corim(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=bad.txt",
		"--key-file=ignored.jwk",
		"--meta-file=ignored.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "bad.txt", []byte("hello!"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding unsigned CoRIM from bad.txt: unexpected EOF")
}

func Test_CorimSignCmd_invalid_unsigned_corim(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=invalid.cbor",
		"--key-file=ignored.jwk",
		"--meta-file=ignored.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "invalid.cbor", testCorimInvalid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error validating CoRIM: tags validation failed: no tags")
}

func Test_CorimSignCmd_non_existent_meta_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ok.cbor",
		"--key-file=ignored.jwk",
		"--meta-file=nonexistent.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading CoRIM Meta from nonexistent.json: open nonexistent.json: file does not exist")
}

func Test_CorimSignCmd_bad_meta_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ok.cbor",
		"--key-file=ignored.jwk",
		"--meta-file=bad.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "bad.json", []byte("{"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding CoRIM Meta from bad.json: unexpected end of JSON input")
}

func Test_CorimSignCmd_invalid_meta_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ok.cbor",
		"--key-file=ignored.jwk",
		"--meta-file=invalid.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid.json", testMetaInvalid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error validating CoRIM Meta: invalid meta: empty name")
}

func Test_CorimSignCmd_non_existent_key_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ok.cbor",
		"--key-file=nonexistent.jwk",
		"--meta-file=ok.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.json", testMetaValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading signing key from nonexistent.jwk: open nonexistent.jwk: file does not exist")
}

func Test_CorimSignCmd_invalid_key_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ok.cbor",
		"--key-file=invalid.jwk",
		"--meta-file=ok.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.json", testMetaValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid.jwk", []byte("{}"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading signing key from invalid.jwk: failed to unmarshal JWK set: failed to unmarshal key from JSON headers: invalid key type from JSON ()")
}

func Test_CorimSignCmd_ok_with_default_output_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ok.cbor",
		"--key-file=ok.jwk",
		"--meta-file=ok.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.json", testMetaValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.jwk", testECKey, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("signed-ok.cbor")
	assert.NoError(t, err)
}

func Test_CorimSignCmd_ok_with_custom_output_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--corim-file=ok.cbor",
		"--key-file=ok.jwk",
		"--meta-file=ok.json",
		"--output-file=my-signed-corim.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.json", testMetaValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.jwk", testECKey, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("my-signed-corim.cbor")
	assert.NoError(t, err)
}
