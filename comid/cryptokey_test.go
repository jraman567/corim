// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/swid"
)

func Test_CryptoKey_NewPKIXBase64Key(t *testing.T) {
	key, err := NewPKIXBase64Key(TestECPubKey)
	require.NoError(t, err)
	assert.Equal(t, TestECPubKey, key.String())
	pub, err := key.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	_, err = NewPKIXBase64Key("")
	assert.EqualError(t, err, "key value not set")

	noPem := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8BlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q=="
	_, err = NewPKIXBase64Key(noPem)
	assert.EqualError(t, err, "could not decode PEM block")

	badBlock := "-----BEGIN CERTIFICATE-----\nDEADBEEF\n-----END CERTIFICATE-----"
	_, err = NewPKIXBase64Key(badBlock)
	assert.Contains(t, err.Error(), "unexpected PEM block type")

	badKey := "-----BEGIN PUBLIC KEY-----\nDEADBEEF\n-----END PUBLIC KEY-----"
	_, err = NewPKIXBase64Key(badKey)
	assert.Contains(t, err.Error(), "unable to parse public key")

	key = MustNewPKIXBase64Key(TestECPubKey)
	assert.Equal(t, TestECPubKey, key.String())

	assert.Panics(t, func() {
		MustNewPKIXBase64Key(badBlock)
	})
}

func Test_CryptoKey_NewPKIXBase64Cert(t *testing.T) {
	cert, err := NewPKIXBase64Cert(TestCert)
	require.NoError(t, err)
	assert.Equal(t, TestCert, cert.String())
	pub, err := cert.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	_, err = NewPKIXBase64Cert("")
	assert.EqualError(t, err, "cert value not set")

	noPem := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8BlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q=="
	_, err = NewPKIXBase64Cert(noPem)
	assert.EqualError(t, err, "could not decode PEM block")

	badBlock := "-----BEGIN PUBLIC KEY-----\nDEADBEEF\n-----END PUBLIC KEY-----"
	_, err = NewPKIXBase64Cert(badBlock)
	assert.Contains(t, err.Error(), "unexpected PEM block type")

	badCert := "-----BEGIN CERTIFICATE-----\nDEADBEEF\n-----END CERTIFICATE-----"
	_, err = NewPKIXBase64Cert(badCert)
	assert.Contains(t, err.Error(), "could not parse x509 cert")

	cert = MustNewPKIXBase64Cert(TestCert)
	assert.Equal(t, TestCert, cert.String())

	assert.Panics(t, func() {
		MustNewPKIXBase64Cert(badBlock)
	})
}

func Test_CryptoKey_NewPKIXBase64CertPath(t *testing.T) {
	certs, err := NewPKIXBase64CertPath(TestCertPath)
	assert.NoError(t, err)
	assert.Equal(t, TestCertPath, certs.String())
	pub, err := certs.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	_, err = NewPKIXBase64CertPath("")
	assert.EqualError(t, err, "cert value not set")

	noPem := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8BlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q=="
	_, err = NewPKIXBase64CertPath(noPem)
	assert.EqualError(t, err, "could not decode PEM block 0")

	badBlock := "-----BEGIN PUBLIC KEY-----\nDEADBEEF\n-----END PUBLIC KEY-----"
	_, err = NewPKIXBase64CertPath(badBlock)
	assert.Contains(t, err.Error(), "unexpected type for PEM block 0")

	badCert := "-----BEGIN CERTIFICATE-----\nDEADBEEF\n-----END CERTIFICATE-----"
	_, err = NewPKIXBase64CertPath(badCert)
	assert.Contains(t, err.Error(), "could not parse x509 cert in PEM block 0")

	certs = MustNewPKIXBase64CertPath(TestCertPath)
	assert.Equal(t, TestCertPath, certs.String())

	assert.Panics(t, func() {
		MustNewPKIXBase64CertPath(badBlock)
	})
}

func Test_CryptoKey_NewCOSEKey(t *testing.T) {
	key, err := NewCOSEKey(TestCOSEKey)
	require.NoError(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString(TestCOSEKey), key.String())
	pub, err := key.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	_, err = NewCOSEKey([]byte{})
	assert.EqualError(t, err, "empty COSE_Key bytes")

	_, err = NewCOSEKey([]byte("DEADBEEF"))
	assert.Contains(t, err.Error(), "cbor: 3 bytes of extraneous data starting at index 5")

	badKey := []byte{ // taken from go-cose unit tests
		0xa2,       // map(2)
		0x01, 0x01, // kty: OKP
		0x03, 0x41, 0x01, // alg: bstr(1)
	}
	_, err = NewCOSEKey(badKey)
	assert.Contains(t, err.Error(), "alg: invalid type")

	keySet, err := NewCOSEKey(TestCOSEKeySetOne)
	require.NoError(t, err)
	pub, err = keySet.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	keySet, err = NewCOSEKey(TestCOSEKeySetMulti)
	require.NoError(t, err)
	_, err = keySet.PublicKey()
	assert.Contains(t, err.Error(), "COSE_KeySet contains more than one key")

	key = MustNewCOSEKey(TestCOSEKey)
	assert.Equal(t, base64.StdEncoding.EncodeToString(TestCOSEKey), key.String())

	assert.Panics(t, func() {
		MustNewCOSEKey(badKey)
	})
}

func Test_CryptoKey_NewThumbprint(t *testing.T) {
	type newKeyFunc func(swid.HashEntry) (*CryptoKey, error)

	for _, newFunc := range []newKeyFunc{
		NewThumbprint,
		NewCertThumbprint,
		NewCertPathThumbprint,
	} {
		key, err := newFunc(TestThumbprint)
		require.NoError(t, err)
		assert.Equal(t, TestThumbprint.String(), key.String())
		_, err = key.PublicKey()
		assert.EqualError(t, err, "cannot get PublicKey from a digest")

		badAlg := swid.HashEntry{
			HashAlgID: 99,
			HashValue: MustHexDecode(nil, `deadbeef`),
		}
		_, err = newFunc(badAlg)
		assert.Contains(t, err.Error(), "unknown hash algorithm 99")

		badHash := swid.HashEntry{
			HashAlgID: 1,
			HashValue: MustHexDecode(nil, `deadbeef`),
		}
		_, err = newFunc(badHash)
		assert.Contains(t, err.Error(), "length mismatch for hash algorithm")
	}

	type mustNewKeyFunc func(swid.HashEntry) *CryptoKey

	for _, mustNewFunc := range []mustNewKeyFunc{
		MustNewThumbprint,
		MustNewCertThumbprint,
		MustNewCertPathThumbprint,
	} {
		key := mustNewFunc(TestThumbprint)
		assert.Equal(t, TestThumbprint.String(), key.String())

		assert.Panics(t, func() {
			mustNewFunc(swid.HashEntry{})
		})
	}
}

func Test_CryptoKey_JSON_roundtrip(t *testing.T) {
	for _, tv := range []struct {
		Type string
		In   any
		Out  string
	}{
		{
			Type: PKIXBase64KeyType,
			In:   TestECPubKey,
			Out:  TestECPubKey,
		},
		{
			Type: PKIXBase64CertType,
			In:   TestCert,
			Out:  TestCert,
		},
		{
			Type: PKIXBase64CertPathType,
			In:   TestCertPath,
			Out:  TestCertPath,
		},
		{
			Type: COSEKeyType,
			In:   TestCOSEKey,
			Out:  base64.StdEncoding.EncodeToString(TestCOSEKey),
		},
		{
			Type: ThumbprintType,
			In:   TestThumbprint,
			Out: fmt.Sprintf("%s;%s",
				TestThumbprint.AlgIDToString(),
				base64.StdEncoding.EncodeToString(TestThumbprint.HashValue),
			),
		},
		{
			Type: CertThumbprintType,
			In:   TestThumbprint,
			Out: fmt.Sprintf("%s;%s",
				TestThumbprint.AlgIDToString(),
				base64.StdEncoding.EncodeToString(TestThumbprint.HashValue),
			),
		},
		{
			Type: CertPathThumbprintType,
			In:   TestThumbprint,
			Out: fmt.Sprintf("%s;%s",
				TestThumbprint.AlgIDToString(),
				base64.StdEncoding.EncodeToString(TestThumbprint.HashValue),
			),
		},
	} {
		key := MustNewCryptoKey(tv.In, tv.Type)
		data, err := json.Marshal(key)
		require.NoError(t, err)

		expected := fmt.Sprintf(`{"type": %q, "value": %q}`, tv.Type, tv.Out)
		assert.JSONEq(t, expected, string(data))

		var key2 CryptoKey
		err = json.Unmarshal(data, &key2)
		require.NoError(t, err)
		assert.Equal(t, *key, key2)
	}
}

func Test_CryptoKey_UnmarshalJSON_negative(t *testing.T) {
	var key CryptoKey

	for _, tv := range []struct {
		Val    string
		ErrMsg string
	}{
		{
			Val:    `@@`,
			ErrMsg: "invalid character",
		},
		{
			Val:    `{"value":"deadbeef"}`,
			ErrMsg: "key type not set",
		},
		{
			Val:    `{"type": "cose-key", "value":";;;"}`,
			ErrMsg: "base64 decode error",
		},
		{
			Val:    `{"type": "thumbprint", "value":"deadbeef"}`,
			ErrMsg: "swid.HashEntry decode error",
		},
		{
			Val:    `{"type": "random-key", "value":"deadbeef"}`,
			ErrMsg: "unexpected ICryptoKeyValue type",
		},
	} {
		err := key.UnmarshalJSON([]byte(tv.Val))
		assert.ErrorContains(t, err, tv.ErrMsg)
	}
}

func Test_CryptoKey_CBOR_roundtrip(t *testing.T) {
	for _, tv := range []struct {
		Type string
		In   any
		Out  string
	}{
		{
			Type: PKIXBase64KeyType,
			In:   TestECPubKey,
			Out:  "d9022a78b12d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d",
		},
		{
			Type: PKIXBase64CertType,
			In:   TestCert,
			Out:  "d9022b7902c82d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d4949423454434341596567417749424167495547687241394d337951494671636b413276366e5165776d4633304977436759494b6f5a497a6a3045417749770a5254454c4d416b474131554542684d4351565578457a415242674e564241674d436c4e766257557455335268644755784954416642674e5642416f4d47456c750a64475679626d5630494664705a47647064484d6755485235494578305a444167467730794d7a41354d4451784d5441784e446861474138794d4455784d4445780a4f5445784d4445304f466f775254454c4d416b474131554542684d4351565578457a415242674e564241674d436c4e76625755745533526864475578495441660a42674e5642416f4d47456c7564475679626d5630494664705a47647064484d6755485235494578305a44425a4d424d4742797147534d343941674547434371470a534d3439417745484130494142467451623668667636387641566d75325244464e63574742784550415a53302b4442527a6d6a7669416b37534c367848684b610a2f37483077442b7568396f516833572b527846797255537148656b433344697a774e576a557a42524d423047413155644467515742425157704e5062366557440a534d2f2b6a7770627a6f4f33694867344c54416642674e5648534d454744415767425157704e506236655744534d2f2b6a7770627a6f4f33694867344c5441500a42674e5648524d4241663845425441444151482f4d416f4743437147534d343942414d43413067414d455543494161794e49463065434a445a6d637271526a480a663968384778654944556e4c716c646549764e66612b39534169454139554c4254506a6e545568596c653232364f416a67327364686b587462334d75304530460a6e75556d7349513d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d",
		},
		{
			Type: PKIXBase64CertPathType,
			In:   TestCertPath,
			Out:  "d9022c7919272d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943656a4343416979674177494241674955497065567756684e2f71594c67744e4a6c775a484a6a2b49542f7777425159444b3256774d444d784d5441760a42674e56424155544b4464684d445a6c5a5755304d5749334f446c6d4e4467324d3251344e6d49344e7a6334596a46684d6a417859545a6d5a57526b4e5459770a4942634e4d5467774d7a49794d6a4d314f545535576867504f546b354f5445794d7a45794d7a55354e546c614d444d784d54417642674e56424155544b4449790a4f5463354e574d784e5467305a475a6c59545977596a67795a444d304f546b334d4459304e7a49324d3259344f44526d5a6d4d774b6a414642674d725a5841440a495141565569377856796e4d3835554a366c77566f6d767053654f494236584362766b6f4649667653755a3752714f43415534776767464b4d423847413155640a497751594d42614146486f4737755162654a394959396872683369786f67476d2f7431574d4230474131556444675157424251696c3558425745332b706775430a30306d58426b636d503468502f44414f42674e56485138424166384542414d434167517744775944565230544151482f42415577417745422f7a43423567594b0a4b7759424241485765514942474145422f77534231444342306142434245414d7476656f414a783476335175394b386749304a326b70376967427634566439420a48454d746c665a4f5164793679444f6c6350725147596a537a7776416a4668384441343662712b78476a33314e46557936706b686f304945514e436b6c2f4b660a6245534c5a364f684563644f6e417a69533567783554714a6d463232794b436a49764c524956784e49685a4e32456e4d41746d3464703145477550424c5572410a747a58686c7a755a754b31785636536b516752415365536f486d6e4e674c686e6e454b54574b7a634c326a7a506a4f41465551544e52792b694f67686c7549330a66696347364e4237634d624c415a6b665631326c696853562b2f37694b33544a3062554e6a51675770615944436745424d4155474179746c63414e42414875360a447475504e4f7572634158632b3431515932336859384b526b42434b434537706873694977526662784b4d4c6c6446474e354f79745166524f5161576f4163760a4957547156394a527a47516147596e6c4c77453d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a2d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943656a434341697967417749424167495545626165377531635037732b472f434462324e75366e505264596b77425159444b3256774d444d784d5441760a42674e56424155544b4449794f5463354e574d784e5467305a475a6c59545977596a67795a444d304f546b334d4459304e7a49324d3259344f44526d5a6d4d770a4942634e4d5467774d7a49794d6a4d314f545535576867504f546b354f5445794d7a45794d7a55354e546c614d444d784d54417642674e56424155544b4445780a596a59355a57566c5a575131597a4e6d596d497a5a5446695a6a41344d7a5a6d4e6a4d325a5756684e7a4e6b4d5463314f446b774b6a414642674d725a5841440a4951447a676b54523775766f50394e427a5345423967752f6c70642b4e4c33384f5659516c30666569574b582b614f43415534776767464b4d423847413155640a497751594d42614146434b586c6346595466366d43344c54535a63475279592f69452f384d423047413155644467515742425152747037753756772f757a34620a38494e7659323771633946316954414f42674e56485138424166384542414d434167517744775944565230544151482f42415577417745422f7a43423567594b0a4b7759424241485765514942474145422f775342314443423061424342454262734b67636176452b75793141786b49646c376c4e39696648793348452b4c69750a386c4d453237434d59396b557479772f6c657331483876706d537978684f3461545767777577516137596e39486f4777654548736f304945514862383730484e0a3162556e396e46696831315342416a396c6f6270754a354772492f6d2b673648776d6f517a35556c79306f584d4e6e78454d4137664c327a613031796e4770490a2f757a383272554932764c57536c476b5167524133584667496f56496d6f736441677675504856616f6276334a476a476c332b41444f543163366454366451450a646e4f62524e7564593871687a547666455752346553364f4a746679724f65527958656b324f564a68365944436745424d4155474179746c63414e42414a496a0a7946717764725a435375596d43342b5a555563414e4b514b41314b63524669496c4b63672f7070774b56796b5058624168736e365343567157474137763743650a4c6935684f72482f566c6a41514163645967633d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a2d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943656a43434169796741774942416749555955584d636656334d756370756e6c313933777558784242722f6777425159444b3256774d444d784d5441760a42674e56424155544b444578596a59355a57566c5a575131597a4e6d596d497a5a5446695a6a41344d7a5a6d4e6a4d325a5756684e7a4e6b4d5463314f446b770a4942634e4d5467774d7a49794d6a4d314f545535576867504f546b354f5445794d7a45794d7a55354e546c614d444d784d54417642674e56424155544b4459780a4e44566a597a63785a6a55334e7a4d795a5463794f574a684e7a6b334e5759334e324d795a54566d4d5441304d57466d5a6a67774b6a414642674d725a5841440a4951424a6f3950677665486a30616876384d6b574851554753785a2f77535464614e4e5a6264425a4e61314c30614f43415534776767464b4d423847413155640a497751594d426141464247326e75377458442b37506876776732396a6275707a3058574a4d423047413155644467515742425268526378783958637935796d360a6558583366433566454547762b44414f42674e56485138424166384542414d434167517744775944565230544151482f42415577417745422f7a43423567594b0a4b7759424241485765514942474145422f775342314443423061424342454132616e736566305362524e386a37367735687a57352f54435846497351634552730a62534b51594e6e71756731726a4543506e686533412f385a3657477861444b316568452b6e72637643394252677257705536374a6f30494551495a795243484b0a394855692f387936563950305a754e45766d64704564496d51303952552f6c4e507358587879763056456d693657447334654679706d4252394c5658425875640a7243647575767953367442577353366b516752416257525443625872642f716c4c504949383549504238705a3975582b586749484934735348662b33463673650a68412f38307a55427a5369364f7a6330442b496259594259786472585a456b6e386955575364516f6b4b5944436745424d4155474179746c63414e42414b6c4a0a2f335659616c5a6d3958624547544b725652616f43566f55785156483375644d726b39796f716a466f77433465336b6453426c474766386d59454937787673410a6172316b6632624758542f634565464749774d3d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a2d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943656a4343416979674177494241674955582b69765048544f6d76566b744d6e514759516a754e6c6b2f445577425159444b3256774d444d784d5441760a42674e56424155544b4459784e44566a597a63785a6a55334e7a4d795a5463794f574a684e7a6b334e5759334e324d795a54566d4d5441304d57466d5a6a67770a4942634e4d5467774d7a49794d6a4d314f545535576867504f546b354f5445794d7a45794d7a55354e546c614d444d784d54417642674e56424155544b44566d0a5a5468685a6a4e6a4e7a526a5a546c685a6a55324e474930597a6c6b4d4445354f4451794d3249345a446b324e475a6a4d7a55774b6a414642674d725a5841440a495143367533626c77453442317864504d65554a50363537502f6d376953742b486572677647626b6b53784d72714f43415534776767464b4d423847413155640a497751594d426141464746467a484831647a4c6e4b627035646664384c6c385151612f344d423047413155644467515742425266364b3838644d3661395753300a7964415a68434f34325754384e54414f42674e56485138424166384542414d434167517744775944565230544151482f42415577417745422f7a43423567594b0a4b7759424241485765514942474145422f7753423144434230614243424543347a326a754a4978356a4437783649754d4e55693754556f6d577843516639516e0a434a39316f7a586b30764a396e4a4f33526476654a7662765a686f5066445149593854695a7038554b447834652b7a573063486b6f304945514f68704d4a36470a45584c5a67487452416d38316f5858414345462b6e6576324d437636434f6875527446797047394233666f526d32726e46556261565a7330704c66424d4738730a7353524a526361775843696d57344f6b5167524132374667783741343231327170714c61786150643974492b7a70664b57724c59634c7832302b444c6663716e0a4249497055434e3330537541753731736534782f696c634b7561574f4f307144673334534a45774679715944436745424d4155474179746c63414e42414374540a3558727836353971476e79776d6c4b48646c484f364264376650626f797a794951686f4574464e75694433576a44672f56777a38634e43556b552b74684737660a432b575a68637041636b446c6461692b5041633d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a2d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943656a434341697967417749424167495558583044564f796c6745697037727a7679614e4665676344545a4977425159444b3256774d444d784d5441760a42674e56424155544b44566d5a5468685a6a4e6a4e7a526a5a546c685a6a55324e474930597a6c6b4d4445354f4451794d3249345a446b324e475a6a4d7a55770a4942634e4d5467774d7a49794d6a4d314f545535576867504f546b354f5445794d7a45794d7a55354e546c614d444d784d54417642674e56424155544b44566b0a4e3251774d7a55305a574e684e5467774e4468684f57566c596d4e6c5a6d4d3559544d304e5464684d4463774d7a526b4f5449774b6a414642674d725a5841440a495143696143326748684d4f31706265516255674c48685367464250442f7a584e414777414873573237322b63364f43415534776767464b4d423847413155640a497751594d42614146462f6f727a78307a7072315a4c544a30426d4549376a5a5a5077314d42304741315564446751574242526466514e55374b5741534b6e750a764f2f4a6f30563642774e4e6b6a414f42674e56485138424166384542414d434167517744775944565230544151482f42415577417745422f7a43423567594b0a4b7759424241485765514942474145422f7753423144434230614243424542595877532f6d7272582b44344d717a4d384a546d49484339584871734a664f47630a623266714259505830555172694c44526c316170484e32327131452b4665614c4857424532755864613151366c596b51416148696f30494551474b483845414e0a4d7631504d4d62577364645a657732472f44522b4139746253693748363830794253653943652b6774616242617251444870673942384c65626d6f50706458740a4154762b6f537a7a6b2b5a7565564b6b516752417a74625532517a614a62634735747745596a59416746757443626e67706732742f32657a3751544e6e344e6d0a723934704f4178384c4970753643662f577a6376642f346b4c4f76577853622f62754d716247767273715944436745424d4155474179746c63414e42414d64780a66466e6b35327275346656354a3167734642746c6879356d4662516e49526947484c78424661544a6b39692b69784f35714662526a71763748512f6a475573490a7372554a513465324a456e534e616b4e6367453d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a2d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943656a434341697967417749424167495541316f576f57507756644335474f3362516f4433726f446531536f77425159444b3256774d444d784d5441760a42674e56424155544b44566b4e3251774d7a55305a574e684e5467774e4468684f57566c596d4e6c5a6d4d3559544d304e5464684d4463774d7a526b4f5449770a4942634e4d5467774d7a49794d6a4d314f545535576867504f546b354f5445794d7a45794d7a55354e546c614d444d784d54417642674e56424155544b44417a0a4e5745784e6d45784e6a4e6d4d4455315a4442694f5445345a57526b596a51794f44426d4e32466c4f44426b5a5751314d6d45774b6a414642674d725a5841440a495144534e45593167624c4d4e414f432b33656f6b2b5279513666684e384632336f32647836315162734d3054714f43415534776767464b4d423847413155640a497751594d4261414646313941315473705942497165363837386d6a52586f48413032534d42304741315564446751574242514457686168592f4256304c6b590a3764744367506575674e37564b6a414f42674e56485138424166384542414d434167517744775944565230544151482f42415577417745422f7a43423567594b0a4b7759424241485765514942474145422f7753423144434230614243424541396d5070416d572b49454f58584f4267537933727935334935363244394f5a485a0a2b4447312f4d396d577869556b524131556369714d704767366e6779717033384a354f7055497546736f53564471465650796a786f3049455149472f376831370a416d3737474c6d51316e534d425a6a74724a2b46726d5754635a6a784a3963583043504a7175377775674c3554636a31493863396e424e71736f6b46783870450a74526f71697a377274365a353244326b516752415a7276714664796a347256636a74566b4a624d6c702f386a6d664765614b682f5247363457724b32754e6b390a79684b4f706b695152307035557354616d2b586445767172786a4c61343373723064692f704b45625a715944436745424d4155474179746c63414e42414f75510a71585a553532314c7a4454585878324559715675574379555a494a5a67526c2f4a477332526d5950594a435a756e304b6a3159497658356d425a3370433835770a30664a466d4d3142322b414373702b703651673d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a2d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943656a4343416979674177494241674955533356553735327371445566593630452f684571536e313432415577425159444b3256774d444d784d5441760a42674e56424155544b44417a4e5745784e6d45784e6a4e6d4d4455315a4442694f5445345a57526b596a51794f44426d4e32466c4f44426b5a5751314d6d45770a4942634e4d5467774d7a49794d6a4d314f545535576867504f546b354f5445794d7a45794d7a55354e546c614d444d784d54417642674e56424155544b4452690a4e7a55314e47566d4f575268593245344d7a55785a6a597a595751774e475a6c4d544579595452684e3251334f4751344d4455774b6a414642674d725a5841440a4951422f6f475854363775635978396c7078465a4652597674676d437942483232692f4c6e554e304b46364c73614f43415534776767464b4d423847413155640a497751594d42614146414e614671466a3846585175526a7432304b4139363641337455714d42304741315564446751574242524c645654766e61796f4e52396a0a7251542b4553704b66586a594254414f42674e56485138424166384542414d434167517744775944565230544151482f42415577417745422f7a43423567594b0a4b7759424241485765514942474145422f77534231444342306142434245444a3657313561697079433255764d6971364943322f77466b7673466339504f72540a314e6e675a47666b65384a6c6e4f37385652555a6373463775687471797265796a697135695a485339684d304a3576494f78756a6f304945514a635a3738616c0a6e43744a69577143486a5467475a6a6f572b6c514a6a4a3955783530545478526545703365454f44394f3374347967645348347254467569754c3674646c5a380a72432f304b544334473576456f77476b51675241675150494251656d5a316973516f4635724b70506f747048584e38595978475935574651497a6b39647a37500a7a78496e5131716e4741736a51505353392b4a4d79774444416937584b754677526630576b32543954615944436745424d4155474179746c63414e424146556c0a5572545135717043634266504765546163584e776c35793357544667706a464b722b4d77367175736a2b62645a366c2b4e334378764f784a396d2b6939364d780a727054366b69536e417a6b2b327a67536941343d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d",
		},
		{
			Type: COSEKeyType,
			In:   TestCOSEKey,
			Out:  "d9022ea501020258246d65726961646f632e6272616e64796275636b406275636b6c616e642e6578616d706c65200121582065eda5a12577c2bae829437fe338701a10aaa375e1bb5b5de108de439c08551d2258201e52ed75701163f7f9e40ddf9f341b3dc9ba860af7e0ca7ca7e9eecd0084d19c",
		},
		{
			Type: ThumbprintType,
			In:   TestThumbprint,
			Out:  "d9022d8201582068e656b251e67e8358bef8483ab0d51c6619f3e7a1a9f0e75838d41ff368f728",
		},
		{
			Type: CertThumbprintType,
			In:   TestThumbprint,
			Out:  "d9022f8201582068e656b251e67e8358bef8483ab0d51c6619f3e7a1a9f0e75838d41ff368f728",
		},
		{
			Type: CertPathThumbprintType,
			In:   TestThumbprint,
			Out:  "d902318201582068e656b251e67e8358bef8483ab0d51c6619f3e7a1a9f0e75838d41ff368f728",
		},
	} {
		key := MustNewCryptoKey(tv.In, tv.Type)
		data, err := em.Marshal(key)
		require.NoError(t, err)

		expected := MustHexDecode(t, tv.Out)
		assert.Equal(t, expected, data)

		var key2 CryptoKey
		err = dm.Unmarshal(data, &key2)
		require.NoError(t, err)
		assert.Equal(t, key.String(), key2.String())
	}
}

func Test_NewCryptoKey_negative(t *testing.T) {
	for _, tv := range []struct {
		Type   string
		In     any
		ErrMsg string
	}{
		{
			Type:   PKIXBase64KeyType,
			In:     7,
			ErrMsg: "value must be a string; found int",
		},
		{
			Type:   PKIXBase64CertType,
			In:     7,
			ErrMsg: "value must be a string; found int",
		},
		{
			Type:   PKIXBase64CertPathType,
			In:     7,
			ErrMsg: "value must be a string; found int",
		},
		{
			Type:   COSEKeyType,
			In:     7,
			ErrMsg: "value must be a []byte; found int",
		},
		{
			Type:   ThumbprintType,
			In:     7,
			ErrMsg: "value must be a swid.HashEntry; found int",
		},
		{
			Type:   CertThumbprintType,
			In:     7,
			ErrMsg: "value must be a swid.HashEntry; found int",
		},
		{
			Type:   CertPathThumbprintType,
			In:     7,
			ErrMsg: "value must be a swid.HashEntry; found int",
		},
		{
			Type:   "random-key",
			In:     7,
			ErrMsg: "unexpected CryptoKey type: random-key",
		},
	} {

		_, err := NewCryptoKey(tv.In, tv.Type)
		assert.ErrorContains(t, err, tv.ErrMsg)
	}
}
