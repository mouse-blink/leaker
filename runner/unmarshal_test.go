package runner

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/vflame6/leaker/runner/sources"
)

// keyRecorder is a Source that remembers the API keys handed to it.
type keyRecorder struct {
	name string
	keys []string
	adds int
}

func (k *keyRecorder) Run(context.Context, string, sources.ScanType, *sources.Session) <-chan sources.Result {
	return nil
}
func (k *keyRecorder) Name() string          { return k.name }
func (k *keyRecorder) UsesKey() bool         { return true }
func (k *keyRecorder) NeedsKey() bool        { return true }
func (k *keyRecorder) AddApiKeys(s []string) { k.keys = s; k.adds++ }
func (k *keyRecorder) RateLimit() int        { return 1 }

func TestUnmarshalFrom_HandsConfiguredKeysToTheirSource(t *testing.T) {
	withKeys := &keyRecorder{name: "WithKeys"}
	withoutKeys := &keyRecorder{name: "WithoutKeys"}

	saved := AllSources
	t.Cleanup(func() { AllSources = saved })
	AllSources[0], AllSources[1] = withKeys, withoutKeys

	path := filepath.Join(t.TempDir(), "provider-config.yaml")
	cfg := "withkeys:\n  - key-one\n"
	if err := os.WriteFile(path, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := UnmarshalFrom(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if want := []string{"key-one"}; !reflect.DeepEqual(withKeys.keys, want) {
		t.Errorf("keys for the configured source = %v, want %v", withKeys.keys, want)
	}
	if withoutKeys.adds != 0 {
		t.Errorf("a source with no configured keys was given keys: %v", withoutKeys.keys)
	}
}
