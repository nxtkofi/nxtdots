package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const expectedCavaTheme = `[color]
gradient = 1
gradient_count = 3
gradient_color_1 = '#111111'
gradient_color_2 = '#222222'
gradient_color_3 = '#333333'
`

func TestCavaGradientGenerators_write_valid_theme_and_reload_cava(t *testing.T) {
	for _, generator := range []struct {
		name string
		run  func(t *testing.T, homeDir string) error
	}{
		{
			name: "Go",
			run: func(_ *testing.T, homeDir string) error {
				return UpdateCavaGradient(homeDir, map[string]string{
					"color1": "#111111",
					"color2": "#222222",
					"color3": "#333333",
				})
			},
		},
		{name: "shell", run: runCavaGradientShell},
	} {
		t.Run(generator.name, func(t *testing.T) {
			// Given
			homeDir := t.TempDir()
			pkillLog := installFakePkill(t)

			// When
			err := generator.run(t, homeDir)

			// Then
			if err != nil {
				t.Fatal(err)
			}
			theme, err := os.ReadFile(filepath.Join(homeDir, ".config", "cava", "themes", "pywal.generated.local"))
			if err != nil {
				t.Fatal(err)
			}
			if string(theme) != expectedCavaTheme {
				t.Errorf("generated theme = %q, want %q", theme, expectedCavaTheme)
			}
			pkillArgs, err := os.ReadFile(pkillLog)
			if err != nil {
				t.Fatal(err)
			}
			if string(pkillArgs) != "-USR2 -x cava\n" {
				t.Errorf("pkill arguments = %q, want %q", pkillArgs, "-USR2 -x cava\\n")
			}
		})
	}
}

func TestUpdateCavaGradient_writes_single_quote_layer_when_colors_are_shell_quoted(t *testing.T) {
	// Given
	homeDir := t.TempDir()
	installFakePkill(t)

	// When
	err := UpdateCavaGradient(homeDir, map[string]string{
		"color1": "'#111111'",
		"color2": "'#222222'",
		"color3": "'#333333'",
	})

	// Then
	if err != nil {
		t.Fatal(err)
	}
	theme, err := os.ReadFile(filepath.Join(homeDir, ".config", "cava", "themes", "pywal.generated.local"))
	if err != nil {
		t.Fatal(err)
	}
	if string(theme) != expectedCavaTheme {
		t.Errorf("generated theme = %q, want %q", theme, expectedCavaTheme)
	}
}

func TestCavaGradientGenerators_return_pkill_errors_except_no_matching_process(t *testing.T) {
	for _, generator := range []struct {
		name string
		run  func(t *testing.T, homeDir string) error
	}{
		{
			name: "Go",
			run: func(_ *testing.T, homeDir string) error {
				return UpdateCavaGradient(homeDir, map[string]string{
					"color1": "#111111",
					"color2": "#222222",
					"color3": "#333333",
				})
			},
		},
		{name: "shell", run: runCavaGradientShell},
	} {
		t.Run(generator.name, func(t *testing.T) {
			for _, exitCode := range []struct {
				name     string
				exitCode string
				wantErr  bool
			}{
				{name: "no matching process", exitCode: "1"},
				{name: "failure", exitCode: "2", wantErr: true},
			} {
				t.Run(exitCode.name, func(t *testing.T) {
					// Given
					homeDir := t.TempDir()
					installFakePkill(t)
					t.Setenv("PKILL_EXIT", exitCode.exitCode)

					// When
					err := generator.run(t, homeDir)

					// Then
					if (err != nil) != exitCode.wantErr {
						t.Errorf("generator error = %v, want error = %t", err, exitCode.wantErr)
					}
				})
			}
		})
	}
}

func TestCavaGradientShell_does_not_reload_when_theme_generation_fails(t *testing.T) {
	// Given
	homeDir := filepath.Join(t.TempDir(), "home-file")
	if err := os.WriteFile(homeDir, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	pkillLog := installFakePkill(t)

	// When
	command := exec.Command("bash", cavaGradientShellPath(t))
	command.Env = append(os.Environ(), "HOME="+homeDir)
	err := command.Run()

	// Then
	if err == nil {
		t.Fatal("runCavaGradientShell() error = nil, want error")
	}
	if _, err := os.Stat(pkillLog); !os.IsNotExist(err) {
		t.Errorf("pkill was called after failed generation, stat error = %v", err)
	}
}

func runCavaGradientShell(t *testing.T, homeDir string) error {
	t.Helper()

	cacheDir := filepath.Join(homeDir, ".cache", "wal")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "colors"), []byte("#000000\n#111111\n#222222\n#333333\n"), 0o644); err != nil {
		return err
	}

	command := exec.Command("bash", cavaGradientShellPath(t))
	command.Env = append(os.Environ(), "HOME="+homeDir)
	return command.Run()
}

func cavaGradientShellPath(t *testing.T) string {
	t.Helper()

	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("determine test source location")
	}
	resolvedTestFile, err := filepath.EvalSymlinks(testFile)
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Join(filepath.Dir(resolvedTestFile), "..", "..", "settings", "executable_update-cava-gradient.sh")
}

func installFakePkill(t *testing.T) string {
	t.Helper()

	binDir := t.TempDir()
	pkillLog := filepath.Join(t.TempDir(), "pkill.log")
	pkill := filepath.Join(binDir, "pkill")
	pkillScript := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$PKILL_LOG\"\nexit \"${PKILL_EXIT:-0}\"\n"
	if err := os.WriteFile(pkill, []byte(pkillScript), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", strings.Join([]string{binDir, os.Getenv("PATH")}, string(os.PathListSeparator)))
	t.Setenv("PKILL_LOG", pkillLog)
	return pkillLog
}
