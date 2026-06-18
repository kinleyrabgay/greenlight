package steps

import (
	"testing"

	"github.com/kinleyrabgay/greenlight/internal/config"
	"github.com/kinleyrabgay/greenlight/internal/db"
	"github.com/kinleyrabgay/greenlight/internal/pipeline"
)

func TestEffectiveBaseBranch(t *testing.T) {
	cases := []struct {
		name string
		base string
		dflt string
		want string
	}{
		{"override wins", "release/2.0", "main", "release/2.0"},
		{"falls back to default", "", "develop", "develop"},
		{"override trims whitespace", "  release/x  ", "main", "release/x"},
		{"defaults to main", "", "", "main"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sctx := &pipeline.StepContext{
				Config: &config.Config{Base: tc.base},
				Repo:   &db.Repo{DefaultBranch: tc.dflt},
			}
			if got := effectiveBaseBranch(sctx); got != tc.want {
				t.Errorf("effectiveBaseBranch() = %q, want %q", got, tc.want)
			}
		})
	}
}
