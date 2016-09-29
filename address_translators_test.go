package gocql

import (
	"testing"
)

func TestIdentityAddressTranslator_EmptyAddrAndZeroPort(t *testing.T) {
	var tr AddressTranslator = IdentityTranslator()
	addr, port := tr.Translate("", 0)
	assertEqual(t, "translated host", "", addr)
	assertEqual(t, "translated port", 0, port)
}

func TestIdentityAddressTranslator_HostProvided(t *testing.T) {
	var tr AddressTranslator = IdentityTranslator()
	addr, port := tr.Translate("10.1.2.3", 9042)
	assertEqual(t, "translated host", "10.1.2.3", addr)
	assertEqual(t, "translated port", 9042, port)
}

