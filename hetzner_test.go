package hetzner

import "testing"

func TestCloudInit(t *testing.T) {
	ci := cloudInit("ghcr.io/acme/app:latest")
	if ci == "" || len(ci) < 20 {
		t.Fatal("empty cloud-init")
	}
	if want := "ghcr.io/acme/app:latest"; !contains(ci, want) {
		t.Fatalf("cloud-init missing image %q", want)
	}
}
func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
