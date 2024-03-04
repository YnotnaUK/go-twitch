package twitch

import "testing"

func TestLoad(t *testing.T) {
	repo, err := NewAuthFilesystemRepository("auth.json")
	if err != nil {
		t.Error(err)
	}
	_, err = repo.Load()
	if err != nil {
		t.Error(err)
	}
}
