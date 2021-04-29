package pkg

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestGoodMD5Check(t *testing.T) {
	pf := NewPreflight(&NoLookup{})
	res := pf.Check("echo 'hello'", "md5=a849639cc38d82e3c0ac4e4dfd8186dd")
	assert.Equal(t, true, res.Ok)
}
func TestGoodSHA1Check(t *testing.T) {
	pf := NewPreflight(&NoLookup{})
	res := pf.Check("echo 'hello'", "sha1=098f8f78f1e13e2a2eee10d6974daebf892e4a71")
	assert.Equal(t, true, res.Ok)
}
func TestGoodSHA256Check(t *testing.T) {
	pf := NewPreflight(&NoLookup{})
	res := pf.Check("echo 'hello'", "sha256=3b084aa6ad2246428c9270825d8631e077b7e7c9bb16f6cafb482bc7fd63e348")
	assert.Equal(t, true, res.Ok)
}
func TestGoodSHA256DefaultCheck(t *testing.T) {
	pf := NewPreflight(&NoLookup{})
	res := pf.Check("echo 'hello'", "3b084aa6ad2246428c9270825d8631e077b7e7c9bb16f6cafb482bc7fd63e348")
	assert.Equal(t, true, res.Ok)
}
func TestBadCheck(t *testing.T) {
	pf := NewPreflight(&NoLookup{})
	res := pf.Check("abcd", "123")
	sig := "88d4266fd4e6338d13b845fcf289579d209c897823b9217da3e161936f031589"
	assert.Equal(t, res.ActualDigest, sig)
	assert.Equal(t, res.ExpectedDigest, "123")
	assert.Equal(t, false, res.Ok)
}

func TestGoodCheck(t *testing.T) {
	pf := NewPreflight(&NoLookup{})
	sig := "88d4266fd4e6338d13b845fcf289579d209c897823b9217da3e161936f031589"
	res := pf.Check("abcd", sig)
	assert.Equal(t, res.ActualDigest, sig)
	assert.Equal(t, res.ExpectedDigest, sig)
	assert.Equal(t, true, res.Ok)
}

type FakeLookup struct {
}

func (f *FakeLookup) Name() string {
	return "Fake"
}

func (f *FakeLookup) Hash(digest Digest) LookupResult {
	return LookupResult{Vulnerable: true, Message: "vuln", Link: "https://example.com/1"}
}
func TestVulnerableCheck(t *testing.T) {
	pf := NewPreflight(&FakeLookup{})
	sig := "88d4266fd4e6338d13b845fcf289579d209c897823b9217da3e161936f031589"
	res := pf.Check("abcd", sig)
	assert.Equal(t, res.ActualDigest, sig)
	assert.Equal(t, res.ExpectedDigest, sig)
	assert.Equal(t, res.Lookup.Vulnerable, true)
	assert.Equal(t, false, res.Ok)
}
func ExampleExecBadDigest() {
	pf := NewPreflight(&FakeLookup{})
	pf.Exec([]string{"../test.sh"}, "123")

	// Output:
	// ⌛️ Preflight starting with Fake
	// ❌ Preflight failed: Digest does not match.
	//
	//    Expected: 123
	//    Actual: fe6d02cf15642ff8d5f61cad6d636a62fd46a5e5a49c06733fece838f5fa9d85
}

func ExampleExecVuln() {
	pf := NewPreflight(&FakeLookup{})
	pf.Exec([]string{"../test.sh"}, "fe6d02cf15642ff8d5f61cad6d636a62fd46a5e5a49c06733fece838f5fa9d85")
	// Output:
	// ⌛️ Preflight starting with Fake
	// ❌ Preflight failed: Digest matches but marked as vulnerable.
	//
	// Information:
	//   Vulnerability: vuln
	//   More: https://example.com/1
}

func ExampleExecOk() {
	pf := NewPreflight(&NoLookup{})
	pf.Exec([]string{"../test.sh"}, "fe6d02cf15642ff8d5f61cad6d636a62fd46a5e5a49c06733fece838f5fa9d85")

	// Output:
	// ⌛️ Preflight starting
	// ✅ Preflight verified
	// hello
}

func ExampleExecPipedBadDigest() {
	pf := NewPreflight(&FakeLookup{})
	pf.ExecPiped("echo 'hello'", "123")

	// Output:
	// ⌛️ Preflight starting with Fake
	// ❌ Preflight failed: Digest does not match.
	//
	//    Expected: 123
	//    Actual: 3b084aa6ad2246428c9270825d8631e077b7e7c9bb16f6cafb482bc7fd63e348
}

func ExampleExecPipedVuln() {
	pf := NewPreflight(&FakeLookup{})
	pf.ExecPiped("echo 'hello'", "3b084aa6ad2246428c9270825d8631e077b7e7c9bb16f6cafb482bc7fd63e348")
	// Output:
	// ⌛️ Preflight starting with Fake
	// ❌ Preflight failed: Digest matches but marked as vulnerable.
	//
	// Information:
	//   Vulnerability: vuln
	//   More: https://example.com/1
}

func ExampleExecPipedOk() {
	pf := NewPreflight(&NoLookup{})
	pf.ExecPiped("echo 'hello'", "3b084aa6ad2246428c9270825d8631e077b7e7c9bb16f6cafb482bc7fd63e348")

	// Output:
	// ⌛️ Preflight starting
	// ✅ Preflight verified
	// hello
}

func ExampleFileLookup() {
	lookup, _ := NewFileLookup("../file_lookup_list.txt")
	pf := NewPreflight(lookup)
	pf.ExecPiped("echo 'hello'", "3b084aa6ad2246428c9270825d8631e077b7e7c9bb16f6cafb482bc7fd63e348")

	// Output:
	// ⌛️ Preflight starting with file lookup: ../file_lookup_list.txt
	// ✅ Preflight verified
	// hello
}