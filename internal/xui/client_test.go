package xui_test

import (
	"testing"

	"github.com/andrejmatveev/3xui-metrics-collector/internal/xui"
)

func TestParseEnvelopeSuccess(t *testing.T) {
	body := []byte(`{"success":true,"msg":"","obj":{"cpu":12.5,"panelVersion":"3.3.1","uptime":3600,"xray":{"state":"running","version":"25.1.1"}}}`)

	var status xui.ServerStatus
	if err := xui.ParseEnvelope(body, &status); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.CPU != 12.5 {
		t.Fatalf("cpu: got %v want 12.5", status.CPU)
	}
	if status.PanelVersion != "3.3.1" {
		t.Fatalf("panelVersion: got %q", status.PanelVersion)
	}
	if status.Xray.State != "running" {
		t.Fatalf("xray state: got %q", status.Xray.State)
	}
}

func TestParseEnvelopeAPIError(t *testing.T) {
	body := []byte(`{"success":false,"msg":"unauthorized","obj":null}`)
	var status xui.ServerStatus
	err := xui.ParseEnvelope(body, &status)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "API error: unauthorized" {
		t.Fatalf("error: got %q", err.Error())
	}
}

func TestParseEnvelopeInboundsArray(t *testing.T) {
	body := []byte(`{"success":true,"msg":"","obj":[{"id":1,"up":100,"down":200,"port":443,"protocol":"vless","nodeId":0,"clientStats":[{"email":"user@example.com","up":50,"down":75}]}]}`)

	var inbounds []xui.Inbound
	if err := xui.ParseEnvelope(body, &inbounds); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inbounds) != 1 {
		t.Fatalf("inbounds len: got %d", len(inbounds))
	}
	if inbounds[0].ClientStats[0].Email != "user@example.com" {
		t.Fatalf("email: got %q", inbounds[0].ClientStats[0].Email)
	}
}

func TestParseEnvelopeOnlines(t *testing.T) {
	body := []byte(`{"success":true,"msg":"","obj":["a@x.com","b@x.com"]}`)
	var onlines []string
	if err := xui.ParseEnvelope(body, &onlines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(onlines) != 2 {
		t.Fatalf("onlines len: got %d", len(onlines))
	}
}

func TestParseEnvelopeLastOnlineMap(t *testing.T) {
	body := []byte(`{"success":true,"msg":"","obj":{"a@x.com":1710000000}}`)
	var lastOnline map[string]int64
	if err := xui.ParseEnvelope(body, &lastOnline); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lastOnline["a@x.com"] != 1710000000 {
		t.Fatalf("timestamp: got %d", lastOnline["a@x.com"])
	}
}
