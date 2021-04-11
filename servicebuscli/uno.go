package servicebuscli

import "github.com/rs/xid"

// SetMessageForwardingTopologyProperties sets the custom message properties for the forwarding topology
func SetMessageForwardingTopologyProperties(sender string, name string, domain string, tenantID string, version string) map[string]interface{} {
	if sender == "" {
		sender = "GlobalOutboxSender"
	}
	if tenantID == "" {
		tenantID = "11111111-1111-1111-1111-555555550001"
	}
	if version == "" {
		version = "1.0"
	}

	diagnosticID := xid.New().String()

	properties := map[string]interface{}{
		"X-MsgTypeVersion": version,
		"X-MsgDomain":      domain,
		"X-MsgName":        name,
		"X-Sender":         sender,
		"X-TenantId":       tenantID,
		"Diagnostic-Id":    diagnosticID,
	}
	return properties
}

// SetMessageUnoProperties sets the custom message properties for the uno topology
func SetMessageUnoProperties(serialization string, tenantID string) map[string]interface{} {
	if serialization == "" {
		serialization = "1"
	}
	if tenantID == "" {
		tenantID = "11111111-1111-1111-1111-555555550001"
	}

	properties := map[string]interface{}{
		"Serialization": "1",
		"TenantId":      tenantID,
	}
	return properties
}
