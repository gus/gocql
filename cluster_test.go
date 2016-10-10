package gocql

import (
	"testing"
	"time"
	"net"
)

func TestNewCluster_Defaults(t *testing.T) {
	cfg := NewCluster()
	assertEqual(t, "cluster config cql version", "3.0.0", cfg.CQLVersion)
	assertEqual(t, "cluster config proto version", 2, cfg.ProtoVersion)
	assertEqual(t, "cluster config timeout", 600 * time.Millisecond, cfg.Timeout)
	assertEqual(t, "cluster config port", 9042, cfg.Port)
	assertEqual(t, "cluster config num-conns", 2, cfg.NumConns)
	assertEqual(t, "cluster config consistency", Quorum, cfg.Consistency)
	assertEqual(t, "cluster config max prepared statements", defaultMaxPreparedStmts, cfg.MaxPreparedStmts)
	assertEqual(t, "cluster config max routing key info", 1000, cfg.MaxRoutingKeyInfo)
	assertEqual(t, "cluster config page-size", 5000, cfg.PageSize)
	assertEqual(t, "cluster config default timestamp", true, cfg.DefaultTimestamp)
	assertEqual(t, "cluster config max wait schema agreement", 60 * time.Second, cfg.MaxWaitSchemaAgreement)
	assertEqual(t, "cluster config reconnect interval", 60 * time.Second, cfg.ReconnectInterval)
}

func TestNewCluster_WithHosts(t *testing.T) {
	cfg := NewCluster("addr1", "addr2")
	assertEqual(t, "cluster config hosts length", 2, len(cfg.Hosts))
	assertEqual(t, "cluster config host 0", "addr1", cfg.Hosts[0])
	assertEqual(t, "cluster config host 1", "addr2", cfg.Hosts[1])
}

func TestClusterConfig_translate_NilTranslator(t *testing.T) {
	cfg := NewCluster()
	assertNil(t, "cluster config address translator", cfg.AddressTranslator)
	host, port := cfg.translateAddress(net.ParseIP("10.0.0.1"), 1234)
	assertEqual(t, "translated address", "10.0.0.1", host.String())
	assertEqual(t, "translated port", 1234, port)
}

func TestClusterConfig_translate_WithTranslator(t *testing.T) {
	cfg := NewCluster()
	cfg.AddressTranslator = staticAddressTranslator(net.ParseIP("10.10.10.10"), 5432)
	addr, port := cfg.translateAddress(net.ParseIP("10.0.0.1"), 1234)
	assertEqual(t, "translated address", "10.10.10.10", addr.String())
	assertEqual(t, "translated port", 5432, port)
}
