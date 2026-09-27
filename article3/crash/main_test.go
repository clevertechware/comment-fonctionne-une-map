package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCrashHelper n'est pas un vrai test : il relance main() dans le
// sous-processus lancé par TestConcurrentMapWritesCrashesPastRecover.
func TestCrashHelper(t *testing.T) {
	if os.Getenv("CRASH_HELPER") != "1" {
		t.Skip("exécuté seulement en sous-processus de TestConcurrentMapWritesCrashesPastRecover")
	}
	main()
}

func TestConcurrentMapWritesCrashesPastRecover(t *testing.T) {
	if raceDetectorEnabled {
		t.Skip("sous -race, le timing change : le race detector signale la race (code de sortie 66) et le fatal error du runtime ne se déclenche généralement plus ; ce test vise le fatal error, à lancer sans -race")
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCrashHelper")
	cmd.Env = append(os.Environ(), "CRASH_HELPER=1")

	output, err := cmd.CombinedOutput()

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("le programme aurait dû se terminer en erreur, err=%v, sortie=%s", err, output)
	}
	if exitErr.Success() {
		t.Fatalf("le programme aurait dû sortir avec un code non nul, sortie=%s", output)
	}

	if !strings.Contains(string(output), "fatal error: concurrent map writes") {
		t.Errorf("sortie attendue avec 'fatal error: concurrent map writes', obtenu:\n%s", output)
	}
	if strings.Contains(string(output), "recover a intercepté") {
		t.Errorf("recover n'aurait jamais dû intercepter le fatal error, sortie:\n%s", output)
	}
}
