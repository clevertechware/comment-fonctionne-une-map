package article3

import "testing"

func TestSyncMapStoreAndLoadExampleReturnsStoredUser(t *testing.T) {
	want := User{ID: "42", Name: "Marie"}

	got, ok := syncMapStoreAndLoadExample(want)
	if !ok {
		t.Fatal("utilisateur non trouvé dans la sync.Map après Store")
	}
	if got != want {
		t.Errorf("utilisateur = %+v, voulu %+v", got, want)
	}
}
