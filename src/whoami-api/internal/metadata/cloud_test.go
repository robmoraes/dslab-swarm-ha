package metadata

import "testing"

func TestRegionFromAvailabilityZone(t *testing.T) {
	got := regionFromAvailabilityZone("us-east-1a")
	if got != "us-east-1" {
		t.Fatalf("expected us-east-1, got %q", got)
	}
}

func TestRegionFromAvailabilityZoneRejectsInvalidValue(t *testing.T) {
	got := regionFromAvailabilityZone("us-east-1")
	if got != "" {
		t.Fatalf("expected empty region for invalid az, got %q", got)
	}
}
