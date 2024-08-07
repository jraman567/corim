// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/veraison/swid"
)

func Example_encode() {
	comid := NewComid().
		SetLanguage("en-GB").
		SetTagIdentity("my-ns:acme-roadrunner-supplement", 0).
		AddEntity("ACME Ltd.", &TestRegID, RoleCreator, RoleTagCreator).
		AddEntity("EMCA Ltd.", nil, RoleMaintainer).
		AddLinkedTag("my-ns:acme-roadrunner-base", RelSupplements).
		AddLinkedTag("my-ns:acme-roadrunner-old", RelReplaces).
		AddReferenceValue(
			ValueTriple{
				Environment: Environment{
					Class: NewClassOID(TestOID).
						SetVendor("ACME Ltd.").
						SetModel("RoadRunner").
						SetLayer(0).
						SetIndex(1),
					Instance: MustNewUEIDInstance(TestUEID),
					Group:    MustNewUUIDGroup(TestUUID),
				},
				Measurements: *NewMeasurements().
					Add(
						MustNewUUIDMeasurement(TestUUID).
							SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff}).
							SetSVN(2).
							AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00}).
							AddDigest(swid.Sha256_32, []byte{0xff, 0xff, 0xff, 0xff}).
							SetFlagsTrue(FlagIsDebug).
							SetFlagsFalse(FlagIsSecure).
							SetSerialNumber("C02X70VHJHD5").
							SetUEID(TestUEID).
							SetUUID(TestUUID).
							SetMACaddr(MACaddr(TestMACaddr)).
							SetIPaddr(TestIPaddr),
					),
			},
		).
		AddEndorsedValue(
			ValueTriple{
				Environment: Environment{
					Class: NewClassUUID(TestUUID).
						SetVendor("ACME Ltd.").
						SetModel("RoadRunner").
						SetLayer(0).
						SetIndex(1),
					Instance: MustNewUEIDInstance(TestUEID),
					Group:    MustNewUUIDGroup(TestUUID),
				},
				Measurements: *NewMeasurements().
					Add(
						MustNewUUIDMeasurement(TestUUID).
							SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff}).
							SetMinSVN(2).
							AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00}).
							AddDigest(swid.Sha256_32, []byte{0xff, 0xff, 0xff, 0xff}).
							SetFlagsTrue(FlagIsDebug).
							SetFlagsFalse(FlagIsSecure, FlagIsConfigured).
							SetSerialNumber("C02X70VHJHD5").
							SetUEID(TestUEID).
							SetUUID(TestUUID).
							SetMACaddr(MACaddr(TestMACaddr)).
							SetIPaddr(TestIPaddr),
					),
			},
		).
		AddAttestVerifKey(
			KeyTriple{
				Environment: Environment{
					Instance: MustNewUUIDInstance(uuid.UUID(TestUUID)),
				},
				VerifKeys: *NewCryptoKeys().
					Add(
						MustNewPKIXBase64Key(TestECPubKey),
					),
			},
		).AddDevIdentityKey(
		KeyTriple{
			Environment: Environment{
				Instance: MustNewUEIDInstance(TestUEID),
			},
			VerifKeys: *NewCryptoKeys().
				Add(
					MustNewPKIXBase64Key(TestECPubKey),
				),
		},
	)

	cbor, err := comid.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	}

	json, err := comid.ToJSON()
	if err == nil {
		fmt.Printf("%s\n", string(json))
	}

	// Output:
	//a50065656e2d474201a10078206d792d6e733a61636d652d726f616472756e6e65722d737570706c656d656e740282a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c6502820100a20069454d4341204c74642e0281020382a200781a6d792d6e733a61636d652d726f616472756e6e65722d626173650100a20078196d792d6e733a61636d652d726f616472756e6e65722d6f6c64010104a4008182a300a500d86f445502c000016941434d45204c74642e026a526f616452756e6e65720300040101d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa81a200d8255031fb5abf023e4992aa4e95f9c1503bfa01aa01d90228020282820644abcdef00820644ffffffff03a201f403f504d9023044010203040544ffffffff064802005e1000000001075020010db8000000000000000000000068086c43303258373056484a484435094702deadbeefdead0a5031fb5abf023e4992aa4e95f9c1503bfa018182a300a500d8255031fb5abf023e4992aa4e95f9c1503bfa016941434d45204c74642e026a526f616452756e6e65720300040101d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa81a200d8255031fb5abf023e4992aa4e95f9c1503bfa01aa01d90229020282820644abcdef00820644ffffffff03a300f401f403f504d9023044010203040544ffffffff064802005e1000000001075020010db8000000000000000000000068086c43303258373056484a484435094702deadbeefdead0a5031fb5abf023e4992aa4e95f9c1503bfa028182a101d902264702deadbeefdead81d9022a78b12d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d038182a101d8255031fb5abf023e4992aa4e95f9c1503bfa81d9022a78b12d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d
	//{"lang":"en-GB","tag-identity":{"id":"my-ns:acme-roadrunner-supplement"},"entities":[{"name":"ACME Ltd.","regid":"https://acme.example","roles":["creator","tagCreator"]},{"name":"EMCA Ltd.","roles":["maintainer"]}],"linked-tags":[{"target":"my-ns:acme-roadrunner-base","rel":"supplements"},{"target":"my-ns:acme-roadrunner-old","rel":"replaces"}],"triples":{"reference-values":[{"environment":{"class":{"id":{"type":"oid","value":"2.5.2.8192"},"vendor":"ACME Ltd.","model":"RoadRunner","layer":0,"index":1},"instance":{"type":"ueid","value":"At6tvu/erQ=="},"group":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}},"measurements":[{"key":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"value":{"svn":{"type":"exact-value","value":2},"digests":["sha-256-32;q83vAA==","sha-256-32;/////w=="],"flags":{"is-secure":false,"is-debug":true},"raw-value":{"type":"bytes","value":"AQIDBA=="},"raw-value-mask":"/////w==","mac-addr":"02:00:5e:10:00:00:00:01","ip-addr":"2001:db8::68","serial-number":"C02X70VHJHD5","ueid":"At6tvu/erQ==","uuid":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}}]}],"endorsed-values":[{"environment":{"class":{"id":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"vendor":"ACME Ltd.","model":"RoadRunner","layer":0,"index":1},"instance":{"type":"ueid","value":"At6tvu/erQ=="},"group":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}},"measurements":[{"key":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"value":{"svn":{"type":"min-value","value":2},"digests":["sha-256-32;q83vAA==","sha-256-32;/////w=="],"flags":{"is-configured":false,"is-secure":false,"is-debug":true},"raw-value":{"type":"bytes","value":"AQIDBA=="},"raw-value-mask":"/////w==","mac-addr":"02:00:5e:10:00:00:00:01","ip-addr":"2001:db8::68","serial-number":"C02X70VHJHD5","ueid":"At6tvu/erQ==","uuid":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}}]}],"dev-identity-keys":[{"environment":{"instance":{"type":"ueid","value":"At6tvu/erQ=="}},"verification-keys":[{"type":"pkix-base64-key","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"}]}],"attester-verification-keys":[{"environment":{"instance":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}},"verification-keys":[{"type":"pkix-base64-key","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"}]}]}}
}

func Example_encode_PSA() {
	comid := NewComid().
		SetTagIdentity("my-ns:acme-roadrunner-supplement", 0).
		AddEntity("ACME Ltd.", &TestRegID, RoleCreator, RoleTagCreator, RoleMaintainer).
		AddReferenceValue(
			ValueTriple{
				Environment: Environment{
					Class: NewClassImplID(TestImplID).
						SetVendor("ACME Ltd.").
						SetModel("RoadRunner 2.0"),
				},
				Measurements: *NewMeasurements().
					Add(
						MustNewPSAMeasurement(
							MustCreatePSARefValID(
								TestSignerID, "BL", "5.0.5",
							)).AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00}),
					).
					Add(
						MustNewPSAMeasurement(
							MustCreatePSARefValID(
								TestSignerID, "PRoT", "1.3.5",
							)).AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00}),
					),
			},
		).
		AddAttestVerifKey(
			KeyTriple{
				Environment: Environment{
					Instance: MustNewUEIDInstance(TestUEID),
				},
				VerifKeys: *NewCryptoKeys().
					Add(
						MustNewPKIXBase64Key(TestECPubKey),
					),
			},
		)

	cbor, err := comid.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	}

	json, err := comid.ToJSON()
	if err == nil {
		fmt.Printf("%s\n", string(json))
	}

	// Output:
	//a301a10078206d792d6e733a61636d652d726f616472756e6e65722d737570706c656d656e740281a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c65028301000204a2008182a100a300d90258582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031016941434d45204c74642e026e526f616452756e6e657220322e3082a200d90259a30162424c0465352e302e35055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a10281820644abcdef00a200d90259a3016450526f540465312e332e35055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a10281820644abcdef00038182a101d902264702deadbeefdead81d9022a78b12d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d
	//{"tag-identity":{"id":"my-ns:acme-roadrunner-supplement"},"entities":[{"name":"ACME Ltd.","regid":"https://acme.example","roles":["creator","tagCreator","maintainer"]}],"triples":{"reference-values":[{"environment":{"class":{"id":{"type":"psa.impl-id","value":"YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="},"vendor":"ACME Ltd.","model":"RoadRunner 2.0"}},"measurements":[{"key":{"type":"psa.refval-id","value":{"label":"BL","version":"5.0.5","signer-id":"rLsRx+TaIXIFUjzkzhokWuGiOa48a/2eeHH35di66Gs="}},"value":{"digests":["sha-256-32;q83vAA=="]}},{"key":{"type":"psa.refval-id","value":{"label":"PRoT","version":"1.3.5","signer-id":"rLsRx+TaIXIFUjzkzhokWuGiOa48a/2eeHH35di66Gs="}},"value":{"digests":["sha-256-32;q83vAA=="]}}]}],"attester-verification-keys":[{"environment":{"instance":{"type":"ueid","value":"At6tvu/erQ=="}},"verification-keys":[{"type":"pkix-base64-key","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"}]}]}}
}

func Example_encode_PSA_attestation_verification() {
	comid := NewComid().
		SetTagIdentity("my-ns:acme-roadrunner-supplement", 0).
		AddEntity("ACME Ltd.", &TestRegID, RoleCreator, RoleTagCreator, RoleMaintainer).
		AddAttestVerifKey(
			KeyTriple{
				Environment: Environment{
					Instance: MustNewUEIDInstance(TestUEID),
				},
				VerifKeys: *NewCryptoKeys().
					Add(
						MustNewPKIXBase64Key(TestECPubKey),
					),
			},
		)

	cbor, err := comid.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	}

	json, err := comid.ToJSON()
	if err == nil {
		fmt.Printf("%s", string(json))
	}

	// Output:
	//a301a10078206d792d6e733a61636d652d726f616472756e6e65722d737570706c656d656e740281a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c65028301000204a1038182a101d902264702deadbeefdead81d9022a78b12d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d
	//{"tag-identity":{"id":"my-ns:acme-roadrunner-supplement"},"entities":[{"name":"ACME Ltd.","regid":"https://acme.example","roles":["creator","tagCreator","maintainer"]}],"triples":{"attester-verification-keys":[{"environment":{"instance":{"type":"ueid","value":"At6tvu/erQ=="}},"verification-keys":[{"type":"pkix-base64-key","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"}]}]}}
}

func Example_decode_JSON() {
	j := `
{
	"lang": "en-GB",
	"tag-identity": {
		"id": "43BBE37F-2E61-4B33-AED3-53CFF1428B16",
		"version": 1
	},
	"entities": [
		{
			"name": "ACME Ltd.",
			"regid": "https://acme.example",
			"roles": [ "tagCreator" ]
		},
		{
			"name": "EMCA Ltd.",
			"regid": "https://emca.example",
			"roles": [ "maintainer", "creator" ]
		}
	],
	"linked-tags": [
		{
			"target": "6F7D8D2F-EAEC-4A15-BB46-1E4DCB85DDFF",
			"rel": "replaces"
		}
	],
	"triples": {
		"reference-values": [
			{
				"environment": {
					"class": {
						"id": {
							"type": "uuid",
							"value": "83294297-97EB-42EF-8A72-AE9FEA002750"
						},
						"vendor": "ACME",
						"model": "RoadRunner Boot ROM",
						"layer": 0,
						"index": 0
					},
					"instance": {
						"type": "ueid",
						"value": "Ad6tvu/erb7v3q2+796tvu8="
					}
				},
				"measurements": [
					{
						"value": {
							"digests": [
								"sha-256:3q2+7w=="
							]
						}
					}
				]
			},
			{
				"environment": {
					"class": {
						"id": {
							"type": "psa.impl-id",
							"value": "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
						},
						"vendor": "PSA-X",
						"model": "Turbo PRoT"
					}
				},
				"measurements": [
					{
						"key": {
							"type": "psa.refval-id",
							"value": {
								"label": "PRoT",
								"version": "1.3.5",
								"signer-id": "rLsRx+TaIXIFUjzkzhokWuGiOa48a/2eeHH35di66Gs="
							}
						},
						"value": {
							"digests": [
								"sha-256:3q2+7w=="
							],
							"svn": {
								"type": "exact-value",
								"value": 1
							},
							"mac-addr": "00:00:5e:00:53:01"
						}
					}
				]
			}
		],
		"endorsed-values": [
			{
				"environment": {
					"class": {
						"id": {
							"type": "oid",
							"value": "2.16.840.1.101.3.4.2.1"
						}
					},
					"instance": {
						"type": "uuid",
						"value": "9090B8D3-3B17-474C-A0B9-6F54731CAB72"
					}
				},
				"measurements": [
					{
						"value": {
							"mac-addr": "00:00:5e:00:53:01",
							"ip-addr": "2001:4860:0:2001::68",
							"serial-number": "C02X70VHJHD5",
							"ueid": "Ad6tvu/erb7v3q2+796tvu8=",
							"uuid": "9090B8D3-3B17-474C-A0B9-6F54731CAB72",
							"raw-value": {
								"type": "bytes",
								"value": "cmF3dmFsdWUKcmF3dmFsdWUK"
							},
							"raw-value-mask": "qg==",
							"op-flags": [ "notSecure" ],
							"digests": [
								"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=",
								"sha-384;S1bPoH+usqtX3pIeSpfWVRRLVGRw66qrb3HA21GN31tKX7KPsq0bSTQmRCTrHlqG"
							],
							"version": {
								"scheme": "semaver",
								"value": "1.2.3beta4"
							},
							"svn": {
								"type": "min-value",
								"value": 10
							}
						}
					}
				]
			}
		],
		"attester-verification-keys": [
			{
				"environment": {
					"group": {
						"type": "uuid",
						"value": "83294297-97EB-42EF-8A72-AE9FEA002750"
					}
				},
				"verification-keys": [
					{
						"type": "pkix-base64-key",
						"value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"
					}
				]
			}
		],
		"dev-identity-keys": [
			{
				"environment": {
					"instance": {
						"type": "uuid",
						"value": "4ECCE47C-85F2-4FD9-9EC6-00DEB72DA707"
					}
				},
				"verification-keys": [
					{
						"type": "pkix-base64-key",
						"value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"
					},
					{
						"type": "pkix-base64-key",
						"value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"
					}
				]
			}
		]
	}
}
`
	comid := Comid{}
	err := comid.FromJSON([]byte(j))

	if err != nil {
		fmt.Printf("FAIL: %v", err)
	} else {
		fmt.Println("OK")
	}

	// Output: OK
}

func Example_decode_CBOR_1() {
	// https://github.com/ietf-rats/ietf-corim-cddl/blob/main/examples/comid-1.diag
	in := []byte{
		0xa3, 0x01, 0xa1, 0x00, 0x50, 0x3f, 0x06, 0xaf, 0x63, 0xa9, 0x3c, 0x11,
		0xe4, 0x97, 0x97, 0x00, 0x50, 0x56, 0x90, 0x77, 0x3f, 0x02, 0x81, 0xa3,
		0x00, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x01,
		0xd8, 0x20, 0x74, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61,
		0x63, 0x6d, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x02,
		0x81, 0x00, 0x04, 0xa1, 0x00, 0x81, 0x82, 0xa1, 0x00, 0xa4, 0x00, 0xd8,
		0x25, 0x50, 0x67, 0xb2, 0x8b, 0x6c, 0x34, 0xcc, 0x40, 0xa1, 0x91, 0x17,
		0xab, 0x5b, 0x05, 0x91, 0x1e, 0x37, 0x01, 0x69, 0x41, 0x43, 0x4d, 0x45,
		0x20, 0x49, 0x6e, 0x63, 0x2e, 0x02, 0x6f, 0x41, 0x43, 0x4d, 0x45, 0x20,
		0x52, 0x6f, 0x61, 0x64, 0x52, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x03, 0x01,
		0x81, 0xa1, 0x01, 0xa2, 0x00, 0xa2, 0x00, 0x65, 0x31, 0x2e, 0x30, 0x2e,
		0x30, 0x01, 0x19, 0x40, 0x00, 0x02, 0x81, 0x82, 0x01, 0x58, 0x20, 0x44,
		0xaa, 0x33, 0x6a, 0xf4, 0xcb, 0x14, 0xa8, 0x79, 0x43, 0x2e, 0x53, 0xdd,
		0x65, 0x71, 0xc7, 0xfa, 0x9b, 0xcc, 0xaf, 0xb7, 0x5f, 0x48, 0x82, 0x59,
		0x26, 0x2d, 0x6e, 0xa3, 0xa4, 0xd9, 0x1b,
	}

	comid := Comid{}
	err := comid.FromCBOR(in)
	if err != nil {
		fmt.Printf("FAIL: %v", err)
	} else {
		fmt.Println("OK")
	}

	// Output: OK
}

func Example_decode_CBOR_2() {
	// https://github.com/ietf-rats/ietf-corim-cddl/blob/main/examples/comid-2.diag
	in := []byte{
		0xa3, 0x01, 0xa1, 0x00, 0x50, 0x3f, 0x06, 0xaf, 0x63, 0xa9, 0x3c, 0x11,
		0xe4, 0x97, 0x97, 0x00, 0x50, 0x56, 0x90, 0x77, 0x3f, 0x02, 0x81, 0xa3,
		0x00, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x01,
		0xd8, 0x20, 0x74, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61,
		0x63, 0x6d, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x02,
		0x81, 0x00, 0x04, 0xa2, 0x00, 0x83, 0x82, 0xa1, 0x00, 0xa4, 0x00, 0xd8,
		0x25, 0x50, 0x67, 0xb2, 0x8b, 0x6c, 0x34, 0xcc, 0x40, 0xa1, 0x91, 0x17,
		0xab, 0x5b, 0x05, 0x91, 0x1e, 0x37, 0x01, 0x69, 0x41, 0x43, 0x4d, 0x45,
		0x20, 0x49, 0x6e, 0x63, 0x2e, 0x02, 0x78, 0x18, 0x41, 0x43, 0x4d, 0x45,
		0x20, 0x52, 0x6f, 0x61, 0x64, 0x52, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x20,
		0x46, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x03, 0x01, 0x81, 0xa1,
		0x01, 0xa1, 0x02, 0x81, 0x82, 0x01, 0x58, 0x20, 0x44, 0xaa, 0x33, 0x6a,
		0xf4, 0xcb, 0x14, 0xa8, 0x79, 0x43, 0x2e, 0x53, 0xdd, 0x65, 0x71, 0xc7,
		0xfa, 0x9b, 0xcc, 0xaf, 0xb7, 0x5f, 0x48, 0x82, 0x59, 0x26, 0x2d, 0x6e,
		0xa3, 0xa4, 0xd9, 0x1b, 0x82, 0xa1, 0x00, 0xa5, 0x00, 0xd8, 0x25, 0x50,
		0xa7, 0x1b, 0x3e, 0x38, 0x8d, 0x45, 0x4a, 0x05, 0x81, 0xf3, 0x52, 0xe5,
		0x8c, 0x83, 0x2c, 0x5c, 0x01, 0x6a, 0x57, 0x59, 0x4c, 0x49, 0x45, 0x20,
		0x49, 0x6e, 0x63, 0x2e, 0x02, 0x77, 0x57, 0x59, 0x4c, 0x49, 0x45, 0x20,
		0x43, 0x6f, 0x79, 0x6f, 0x74, 0x65, 0x20, 0x54, 0x72, 0x75, 0x73, 0x74,
		0x65, 0x64, 0x20, 0x4f, 0x53, 0x03, 0x02, 0x04, 0x00, 0x81, 0xa1, 0x01,
		0xa1, 0x02, 0x81, 0x82, 0x01, 0x58, 0x20, 0xbb, 0x71, 0x19, 0x8e, 0xd6,
		0x0a, 0x95, 0xdc, 0x3c, 0x61, 0x9e, 0x55, 0x5c, 0x2c, 0x0b, 0x8d, 0x75,
		0x64, 0xa3, 0x80, 0x31, 0xb0, 0x34, 0xa1, 0x95, 0x89, 0x25, 0x91, 0xc6,
		0x53, 0x65, 0xb0, 0x82, 0xa1, 0x00, 0xa5, 0x00, 0xd8, 0x25, 0x50, 0xa7,
		0x1b, 0x3e, 0x38, 0x8d, 0x45, 0x4a, 0x05, 0x81, 0xf3, 0x52, 0xe5, 0x8c,
		0x83, 0x2c, 0x5c, 0x01, 0x6a, 0x57, 0x59, 0x4c, 0x49, 0x45, 0x20, 0x49,
		0x6e, 0x63, 0x2e, 0x02, 0x77, 0x57, 0x59, 0x4c, 0x49, 0x45, 0x20, 0x43,
		0x6f, 0x79, 0x6f, 0x74, 0x65, 0x20, 0x54, 0x72, 0x75, 0x73, 0x74, 0x65,
		0x64, 0x20, 0x4f, 0x53, 0x03, 0x02, 0x04, 0x01, 0x81, 0xa1, 0x01, 0xa1,
		0x02, 0x81, 0x82, 0x01, 0x58, 0x20, 0xbb, 0x71, 0x19, 0x8e, 0xd6, 0x0a,
		0x95, 0xdc, 0x3c, 0x61, 0x9e, 0x55, 0x5c, 0x2c, 0x0b, 0x8d, 0x75, 0x64,
		0xa3, 0x80, 0x31, 0xb0, 0x34, 0xa1, 0x95, 0x89, 0x25, 0x91, 0xc6, 0x53,
		0x65, 0xb0, 0x01, 0x81, 0x82, 0xa1, 0x00, 0xa4, 0x00, 0xd8, 0x25, 0x50,
		0x67, 0xb2, 0x8b, 0x6c, 0x34, 0xcc, 0x40, 0xa1, 0x91, 0x17, 0xab, 0x5b,
		0x05, 0x91, 0x1e, 0x37, 0x01, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x49,
		0x6e, 0x63, 0x2e, 0x02, 0x72, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x52, 0x6f,
		0x6f, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x54, 0x72, 0x75, 0x73, 0x74, 0x03,
		0x00, 0x81, 0xa1, 0x01, 0xa1, 0x01, 0xd9, 0x02, 0x28, 0x01,
	}

	comid := Comid{}
	err := comid.FromCBOR(in)
	if err != nil {
		fmt.Printf("FAIL: %v", err)
	} else {
		fmt.Println("OK")
	}

	// Output: OK
}

func Example_decode_CBOR_3() {
	// https://github.com/ietf-rats/ietf-corim-cddl/blob/main/examples/comid-design-cd.diag
	in := []byte{
		0xa4, 0x01, 0xa1, 0x00, 0x50, 0x1e, 0xac, 0xd5, 0x96, 0xf4, 0xa3, 0x4f,
		0xb6, 0x99, 0xbf, 0xae, 0xb5, 0x8e, 0x0a, 0x4e, 0x47, 0x02, 0x81, 0xa3,
		0x00, 0x71, 0x46, 0x50, 0x47, 0x41, 0x20, 0x44, 0x65, 0x73, 0x69, 0x67,
		0x6e, 0x73, 0x2d, 0x52, 0x2d, 0x55, 0x73, 0x01, 0xd8, 0x20, 0x78, 0x1e,
		0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x66, 0x70, 0x67, 0x61,
		0x64, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x73, 0x72, 0x75, 0x73, 0x2e, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x02, 0x81, 0x00, 0x03, 0x81, 0xa2,
		0x00, 0x50, 0x97, 0xf5, 0xa7, 0x07, 0x1c, 0x6f, 0x43, 0x8f, 0x87, 0x7a,
		0x4a, 0x02, 0x07, 0x80, 0xeb, 0xe9, 0x01, 0x00, 0x04, 0xa2, 0x00, 0x84,
		0x82, 0xa1, 0x00, 0xa3, 0x00, 0xd8, 0x6f, 0x4b, 0x60, 0x86, 0x48, 0x01,
		0x86, 0xf8, 0x4d, 0x01, 0x0f, 0x04, 0x01, 0x01, 0x76, 0x66, 0x70, 0x67,
		0x61, 0x64, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x73, 0x72, 0x75, 0x73, 0x2e,
		0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03, 0x02, 0x81, 0xa1, 0x01,
		0xa2, 0x04, 0xd9, 0x02, 0x30, 0x48, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x05, 0x48, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x82, 0xa1, 0x00, 0xa3, 0x00, 0xd8, 0x6f, 0x4b, 0x60, 0x86, 0x48, 0x01,
		0x86, 0xf8, 0x4d, 0x01, 0x0f, 0x04, 0x02, 0x01, 0x76, 0x66, 0x70, 0x67,
		0x61, 0x64, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x73, 0x72, 0x75, 0x73, 0x2e,
		0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03, 0x02, 0x81, 0xa1, 0x01,
		0xa1, 0x02, 0x81, 0x82, 0x07, 0x58, 0x30, 0x3f, 0xe1, 0x8e, 0xca, 0x40,
		0x53, 0x87, 0x9e, 0x01, 0x7e, 0xf5, 0xeb, 0x7a, 0x3e, 0x51, 0x57, 0x65,
		0x9c, 0x5f, 0x9b, 0xb1, 0x5b, 0x7d, 0x09, 0x95, 0x9b, 0x8b, 0x86, 0x47,
		0x82, 0x2a, 0x4c, 0xc2, 0x1c, 0x3a, 0xa6, 0x72, 0x1c, 0xef, 0x87, 0xf5,
		0xbf, 0xa5, 0x34, 0x95, 0xdb, 0x08, 0x33, 0x82, 0xa1, 0x00, 0xa3, 0x00,
		0xd8, 0x6f, 0x4b, 0x60, 0x86, 0x48, 0x01, 0x86, 0xf8, 0x4d, 0x01, 0x0f,
		0x04, 0x03, 0x01, 0x76, 0x66, 0x70, 0x67, 0x61, 0x64, 0x65, 0x73, 0x69,
		0x67, 0x6e, 0x73, 0x72, 0x75, 0x73, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
		0x6c, 0x65, 0x03, 0x02, 0x81, 0xa1, 0x01, 0xa1, 0x02, 0x81, 0x82, 0x07,
		0x58, 0x30, 0x20, 0xff, 0x68, 0x1a, 0x08, 0x82, 0xe2, 0x9b, 0x48, 0x19,
		0x53, 0x88, 0x89, 0x36, 0x20, 0x9c, 0xb5, 0x3d, 0xf9, 0xc5, 0xaa, 0xec,
		0x60, 0x6a, 0x2c, 0x24, 0xa0, 0xfb, 0x13, 0x85, 0x95, 0x12, 0x4b, 0x8e,
		0x3f, 0x24, 0xa1, 0x27, 0x71, 0xbc, 0x38, 0x54, 0xcc, 0x68, 0xb4, 0x03,
		0x61, 0xad, 0x82, 0xa1, 0x00, 0xa2, 0x00, 0xd8, 0x6f, 0x4c, 0x60, 0x86,
		0x48, 0x01, 0x86, 0xf8, 0x4d, 0x01, 0x0f, 0x04, 0x63, 0x01, 0x01, 0x76,
		0x66, 0x70, 0x67, 0x61, 0x64, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x73, 0x72,
		0x75, 0x73, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x81, 0xa1,
		0x01, 0xa2, 0x04, 0xd9, 0x02, 0x30, 0x58, 0x30, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x58, 0x30, 0x46,
		0x62, 0x24, 0x34, 0x3d, 0x68, 0x18, 0x02, 0xc1, 0x50, 0x6b, 0xbe, 0xd7,
		0xd7, 0xf0, 0x0b, 0x96, 0x9b, 0xad, 0xdd, 0x63, 0x46, 0xe4, 0xf2, 0xe7,
		0xce, 0x14, 0x66, 0x92, 0x99, 0x6f, 0x22, 0xa4, 0x58, 0x14, 0xde, 0x81,
		0xd2, 0x48, 0xf5, 0x83, 0xb6, 0x5f, 0x81, 0x7b, 0x5f, 0xce, 0xab, 0x01,
		0x81, 0x82, 0xa1, 0x00, 0xa2, 0x00, 0xd8, 0x6f, 0x4c, 0x60, 0x86, 0x48,
		0x01, 0x86, 0xf8, 0x4d, 0x01, 0x0f, 0x04, 0x63, 0x02, 0x01, 0x76, 0x66,
		0x70, 0x67, 0x61, 0x64, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x73, 0x72, 0x75,
		0x73, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x81, 0xa1, 0x01,
		0xa2, 0x04, 0xd9, 0x02, 0x30, 0x48, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x05, 0x48, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
	}

	comid := Comid{}
	err := comid.FromCBOR(in)
	if err != nil {
		fmt.Printf("FAIL: %v", err)
	} else {
		fmt.Println("OK")
	}

	// Output: OK
}

func Example_decode_CBOR_4() {
	// https://github.com/ietf-rats/ietf-corim-cddl/blob/main/examples/comid-firmware-cd.diag
	in := []byte{
		0xa3, 0x01, 0xa1, 0x00, 0x50, 0xaf, 0x1c, 0xd8, 0x95, 0xbe, 0x78, 0x4a,
		0xdb, 0xb7, 0xe9, 0xad, 0xd4, 0x4a, 0x65, 0xab, 0xf3, 0x02, 0x81, 0xa3,
		0x00, 0x71, 0x46, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x20, 0x4d,
		0x46, 0x47, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x01, 0xd8, 0x20, 0x78, 0x18,
		0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x66, 0x77, 0x6d, 0x66,
		0x67, 0x69, 0x6e, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
		0x02, 0x81, 0x00, 0x04, 0xa2, 0x00, 0x82, 0x82, 0xa1, 0x00, 0xa4, 0x01,
		0x70, 0x66, 0x77, 0x6d, 0x66, 0x67, 0x69, 0x6e, 0x63, 0x2e, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x02, 0x67, 0x66, 0x77, 0x59, 0x5f, 0x6e,
		0x35, 0x78, 0x03, 0x00, 0x04, 0x00, 0x81, 0xa1, 0x01, 0xa2, 0x01, 0xd9,
		0x02, 0x28, 0x01, 0x02, 0x81, 0x82, 0x07, 0x58, 0x30, 0x15, 0xe7, 0x7d,
		0x6f, 0x13, 0x32, 0x52, 0xf1, 0xdb, 0x70, 0x44, 0x90, 0x13, 0x13, 0x88,
		0x4f, 0x29, 0x77, 0xd2, 0x10, 0x9b, 0x33, 0xc7, 0x9f, 0x33, 0xe0, 0x79,
		0xbf, 0xc7, 0x88, 0x65, 0x25, 0x5c, 0x0f, 0xb7, 0x33, 0xc2, 0x40, 0xfd,
		0xda, 0x54, 0x4b, 0x82, 0x15, 0xd7, 0xb8, 0xf8, 0x15, 0x82, 0xa1, 0x00,
		0xa4, 0x01, 0x70, 0x66, 0x77, 0x6d, 0x66, 0x67, 0x69, 0x6e, 0x63, 0x2e,
		0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x02, 0x67, 0x66, 0x77, 0x58,
		0x5f, 0x6e, 0x35, 0x78, 0x03, 0x01, 0x04, 0x00, 0x81, 0xa1, 0x01, 0xa2,
		0x01, 0xd9, 0x02, 0x28, 0x01, 0x02, 0x81, 0x82, 0x07, 0x58, 0x30, 0x3d,
		0x90, 0xb6, 0xbf, 0x00, 0x3d, 0xa2, 0xd9, 0x4e, 0xa5, 0x46, 0x3f, 0x97,
		0xfb, 0x3c, 0x53, 0xdd, 0xc5, 0x1c, 0xfb, 0xa1, 0xe3, 0xe3, 0x8e, 0xef,
		0x7a, 0xf0, 0x71, 0xa6, 0x79, 0x86, 0x59, 0x5d, 0x22, 0x72, 0x91, 0x31,
		0xdf, 0x9f, 0xe8, 0x0f, 0x54, 0x51, 0xee, 0xf1, 0x54, 0xf8, 0x5e, 0x01,
		0x81, 0x82, 0xa1, 0x00, 0xa2, 0x00, 0xd8, 0x6f, 0x4c, 0x60, 0x86, 0x48,
		0x01, 0x86, 0xf8, 0x4d, 0x01, 0x0f, 0x04, 0x63, 0x01, 0x01, 0x70, 0x66,
		0x77, 0x6d, 0x66, 0x67, 0x69, 0x6e, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d,
		0x70, 0x6c, 0x65, 0x81, 0xa1, 0x01, 0xa2, 0x04, 0xd9, 0x02, 0x30, 0x48,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x48, 0xff, 0xff,
		0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
	}

	comid := Comid{}
	err := comid.FromCBOR(in)
	if err != nil {
		fmt.Printf("FAIL: %v", err)
	} else {
		fmt.Println("OK")
	}

	// Output: OK
}

func Example_decode_CBOR_5() {
	// Taken from https://github.com/ietf-corim-cddl/blob/main/examples/comid-3.diag
	in := []byte{
		0xa3, 0x01, 0xa1, 0x00, 0x78, 0x20, 0x6d, 0x79, 0x2d, 0x6e, 0x73, 0x3a,
		0x61, 0x63, 0x6d, 0x65, 0x2d, 0x72, 0x6f, 0x61, 0x64, 0x72, 0x75, 0x6e,
		0x6e, 0x65, 0x72, 0x2d, 0x73, 0x75, 0x70, 0x70, 0x6c, 0x65, 0x6d, 0x65,
		0x6e, 0x74, 0x02, 0x81, 0xa3, 0x00, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20,
		0x49, 0x6e, 0x63, 0x2e, 0x01, 0xd8, 0x20, 0x74, 0x68, 0x74, 0x74, 0x70,
		0x73, 0x3a, 0x2f, 0x2f, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x65, 0x78, 0x61,
		0x6d, 0x70, 0x6c, 0x65, 0x02, 0x83, 0x01, 0x00, 0x02, 0x04, 0xa1, 0x00,
		0x81, 0x82, 0xa1, 0x00, 0xa3, 0x00, 0xd8, 0x6f, 0x44, 0x55, 0x02, 0xc0,
		0x00, 0x01, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x49, 0x6e, 0x63, 0x2e,
		0x02, 0x78, 0x18, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x52, 0x6f, 0x61, 0x64,
		0x52, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x20, 0x46, 0x69, 0x72, 0x6d, 0x77,
		0x61, 0x72, 0x65, 0x81, 0xa2, 0x00, 0x19, 0x02, 0xbc, 0x01, 0xa1, 0x02,
		0x81, 0x82, 0x06, 0x44, 0xab, 0xcd, 0xef, 0x00,
	}
	comid := Comid{}
	err := comid.FromCBOR(in)
	if err != nil {
		fmt.Printf("FAIL: %v", err)
	} else {
		fmt.Println("OK")
	}

	// Output: OK
}
