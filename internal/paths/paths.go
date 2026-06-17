package paths

import (
	"os"
	"path/filepath"
)

// Paths provides access to all greenlight filesystem locations.
// The root defaults to ~/.greenlight but can be overridden via NM_HOME
// or by using WithRoot (for testing).
type Paths struct {
	root string
}

// New returns Paths rooted at NM_HOME or ~/.greenlight.
func New() (*Paths, error) {
	if env := os.Getenv("NM_HOME"); env != "" {
		return &Paths{root: env}, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Paths{root: filepath.Join(home, ".greenlight")}, nil
}

// WithRoot returns Paths rooted at a custom directory (for testing).
func WithRoot(root string) *Paths {
	return &Paths{root: root}
}

func (p *Paths) Root() string       { return p.root }
func (p *Paths) DB() string         { return filepath.Join(p.root, "state.sqlite") }
func (p *Paths) Socket() string     { return filepath.Join(p.root, "socket") }
func (p *Paths) PIDFile() string    { return filepath.Join(p.root, "daemon.pid") }
func (p *Paths) ConfigFile() string { return filepath.Join(p.root, "config.yaml") }

// ProfilesDir is where user-defined framework profiles (<name>.md) live.
// A profile here overrides the built-in of the same name.
func (p *Paths) ProfilesDir() string { return filepath.Join(p.root, "profiles") }
func (p *Paths) UpdateCheckFile() string {
	return filepath.Join(p.root, "update-check.json")
}

func (p *Paths) ReposDir() string { return filepath.Join(p.root, "repos") }
func (p *Paths) RepoDir(repoID string) string {
	return filepath.Join(p.root, "repos", repoID+".git")
}

func (p *Paths) WorktreesDir() string { return filepath.Join(p.root, "worktrees") }
func (p *Paths) WorktreeDir(repoID, runID string) string {
	return filepath.Join(p.root, "worktrees", repoID, runID)
}

func (p *Paths) LogsDir() string { return filepath.Join(p.root, "logs") }
func (p *Paths) RunLogDir(runID string) string {
	return filepath.Join(p.root, "logs", runID)
}
func (p *Paths) DaemonLog() string { return filepath.Join(p.root, "logs", "daemon.log") }
func (p *Paths) CLILog() string    { return filepath.Join(p.root, "logs", "cli.log") }

// ServerPIDsDir holds PID-tracking files for managed agent servers
// (opencode, rovodev) so a freshly started daemon can reap orphans left
// behind by a crashed predecessor.
func (p *Paths) ServerPIDsDir() string { return filepath.Join(p.root, "servers") }

// EnsureDirs creates all required directories under root.
func (p *Paths) EnsureDirs() error {
	dirs := []string{
		p.root,
		p.ReposDir(),
		p.WorktreesDir(),
		p.LogsDir(),
		p.ServerPIDsDir(),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}
