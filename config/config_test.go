package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestGlobal(t *testing.T) {
	c := Config{
		env: Environment{
			"MEOW": "DEPRECATED",
		},
		Global: Binary{
			Launcher: "meow",
			Env: Environment{
				"MEOW": "GLOBAL",
			},
		},
		Player: Binary{
			Env: Environment{
				"MEOW": "PLAYER",
			},
		},
	}

	if err := c.globalize(); err != nil {
		t.Fatal(err)
	}

	if c.Global.Env["MEOW"] != "DEPRECATED" {
		t.Error("expected env overrides global env")
	}

	if c.Player.Launcher != "meow" && c.Studio.Launcher != "meow" {
		t.Error("expected player or/and studio applies global launcher")
	}

	if c.Player.Env["MEOW"] != "PLAYER" {
		t.Error("expected player overrides global env")
	}

	if c.Studio.Env["MEOW"] != "DEPRECATED" {
		t.Error("expected studio applies global env")
	}
}

func TestBinarySetup(t *testing.T) {
	b := Binary{
		Env: Environment{
			"MEOW": "MEOW",
		},
	}

	if err := b.setup(); err != nil {
		t.Fatal(err)
	}

	b.Renderer = "Meow"
	if err := b.setup(); !errors.Is(err, ErrInvalidRenderer) {
		t.Error("expected renderer check")
	}

	b.Dxvk = true
	b.Renderer = "Vulkan"
	if err := b.setup(); !errors.Is(err, ErrNeedDXVKRenderer) {
		t.Error("expected dxvk appropiate renderer check")
	}

	b.Renderer = "D3D11FL10"
	if err := b.setup(); errors.Is(err, ErrNeedDXVKRenderer) {
		t.Error("expected not dxvk appropiate renderer check")
	}

	if os.Getenv("MEOW") == "MEOW" {
		t.Error("expected no change in environment")
	}
}

func TestSetup(t *testing.T) {
	wr := t.TempDir()
	c := Config{
		WineRoot: wr,
	}

	if err := c.setup(); !errors.Is(err, ErrWineRootInvalid) {
		t.Error("expected wine root wine check")
	}

	c.WineRoot = filepath.Join(".", wr)
	if err := c.setup(); !errors.Is(err, ErrWineRootAbs) {
		t.Error("expected wine root absolute path")
	}
}
