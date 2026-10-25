package auth

import "testing"

func TestNeedsRehash(t *testing.T) {
	legacy := "$argon2id$v=19$m=65536,t=1,p=20$7Q2+UCMMTnRoU8iL/1rs5A$pmL/C19I/192haTGERtuJtlJfTzPcb0fyMrD96abRtk"
	if !needsRehash(legacy) {
		t.Error("needsRehash(64 MiB hash) = false, want true")
	}
	current, err := newHasher().create(t.Context(), "portal67")
	if err != nil {
		t.Fatal(err)
	}
	if needsRehash(current) {
		t.Error("needsRehash(current hash) = true, want false")
	}
	if needsRehash("not a hash") {
		t.Error("needsRehash(invalid) = true, want false")
	}
}
