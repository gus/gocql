package gocql

// AddressTranslator provides a way to translate node addresses (and ports) that are
// discovered or received as a node event. This is especially useful in ec2 (when
// using the EC2MultiRegionAddressTranslator) to translate public IPs to private IPs
// when possible.
type AddressTranslator interface {
	// Translate will translate the provided address and/or port to another
	// address and/or port
	Translate(addr string, port int) (string, int)
}

type AddressTranslatorFunc func(addr string, port int) (string, int)

func (fn AddressTranslatorFunc) Translate(addr string, port int) (string, int) {
	return fn(addr, port)
}

// IdentityTranslator will do nothing but return what it was provided. It is essentially a no-op.
func IdentityTranslator() AddressTranslator {
	return AddressTranslatorFunc(func(addr string, port int) (string, int) {
		return addr, port
	})
}
