package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kinleyrabgay/greenlight/internal/config"
	"github.com/kinleyrabgay/greenlight/internal/daemon"
	"github.com/kinleyrabgay/greenlight/internal/gate"
	"github.com/kinleyrabgay/greenlight/internal/profile"
	"github.com/kinleyrabgay/greenlight/internal/skill"
	"github.com/spf13/cobra"
)

const banner = `____ ____ ____ ____ _  _ _    _ ____ _  _ ___
| __ |__/ |___ |___ |\ | |    | | __ |__|  |
|__] |  \ |___ |___ | \| |___ | |__] |  |  |`

func newInitCmd() *cobra.Command {
	var framework string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize greenlight gate for the current repository",
		Long: "Sets up or refreshes a local bare repo as a gate, installs a post-receive hook,\n" +
			"best-effort isolates the gate hook path from shared local git config writes when Git supports `config --worktree`,\n" +
			"adds or repairs the \"greenlight\" git remote, and records the repo in the database.\n\n" +
			"Run this from inside a git repository that has an \"origin\" remote.\n\n" +
			"Pass --framework <name> to scaffold a .greenlight.yaml that selects a\n" +
			"framework profile (default commands + review rules). Omit it to rely on\n" +
			"auto-detection at run time.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return trackCommand("init", func() error {
				p, d, err := openResources()
				if err != nil {
					return err
				}
				defer d.Close()

				if framework != "" {
					if _, perr := profile.Load(framework, p.ProfilesDir()); perr != nil {
						return perr
					}
				}

				repo, created, err := gate.Init(cmd.Context(), d, p, ".")
				if err != nil {
					return fmt.Errorf("init: %w", err)
				}

				// Seed a GREENLIGHT.md audit guide in the project. The framework
				// (explicit flag, else auto-detected) picks the seed content;
				// with none, a generic template is written. The file is the
				// repo's own, so an existing one is never overwritten.
				guideStatus := ""
				fw := config.ResolveFramework(framework, "", repo.WorkingPath)
				if framework != "" && fw != "" {
					if werr := writeFrameworkConfig(repo.WorkingPath, fw); werr != nil {
						fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s\n", sDim.Render("config"), sYellow.Render("framework not written: "+werr.Error()))
					}
				}
				var prof *profile.Profile
				if fw != "" {
					prof, _ = profile.Load(fw, p.ProfilesDir())
				}
				if wrote, werr := writeGreenlightGuide(repo.WorkingPath, fw, prof); werr != nil {
					guideStatus = sYellow.Render("GREENLIGHT.md not written: " + werr.Error())
				} else if wrote {
					seed := "generic template"
					if fw != "" {
						seed = fw + " template"
					}
					guideStatus = sGreen.Render("GREENLIGHT.md") + sDim.Render(" created ("+seed+") — edit it to control what greenlight audits")
				} else {
					guideStatus = sDim.Render("GREENLIGHT.md already present (left as-is)")
				}
				if err := daemon.EnsureDaemon(p); err != nil {
					// Only roll back a gate we created in this run; a re-init
					// must never eject a user's pre-existing gate.
					if created {
						if _, ejectErr := gate.Eject(cmd.Context(), d, p, "."); ejectErr != nil {
							return fmt.Errorf("start daemon: %w, rollback init: %v", err, ejectErr)
						}
					}
					return fmt.Errorf("start daemon: %w", err)
				}

				// Install the agent skill at user level so agents can drive
				// greenlight via `/greenlight` in any repo. Best-effort: a
				// skill write failure must not undo a successful gate setup.
				_, skillErr := skill.InstallUser()

				w := cmd.OutOrStdout()
				fmt.Fprintln(w, sCyan.Render(banner))
				fmt.Fprintln(w)
				headline := "Gate initialized"
				if !created {
					headline = "Gate already initialized (refreshed)"
				}
				fmt.Fprintf(w, "  %s %s\n", sGreen.Render("✓"), headline)
				fmt.Fprintln(w)
				fmt.Fprintf(w, "  %s  %s\n", sDim.Render("  repo"), repo.WorkingPath)
				fmt.Fprintf(w, "  %s  greenlight → %s\n", sDim.Render("  gate"), p.RepoDir(repo.ID))
				fmt.Fprintf(w, "  %s  %s\n", sDim.Render("remote"), repo.UpstreamURL)
				if skillErr != nil {
					fmt.Fprintf(w, "  %s  %s\n", sDim.Render(" skill"), sYellow.Render("skipped: "+skillErr.Error()))
				} else {
					fmt.Fprintf(w, "  %s  %s %s\n", sDim.Render(" skill"), sGreen.Render("/greenlight"), sDim.Render("installed for agents at user level"))
				}
				if guideStatus != "" {
					fmt.Fprintf(w, "  %s  %s\n", sDim.Render(" guide"), guideStatus)
				}
				if legacy := skill.Vendored(repo.WorkingPath); len(legacy) > 0 {
					fmt.Fprintf(w, "  %s  %s\n", sDim.Render("  note"), sDim.Render("vendored skill copy ("+strings.Join(legacy, ", ")+") is no longer needed and can be removed"))
				}
				fmt.Fprintln(w)
				fmt.Fprintf(w, "  %s\n", sDim.Render("Push through the gate with:"))
				fmt.Fprintf(w, "  %s\n", sBold.Render("git push greenlight <branch>"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&framework, "framework", "", "framework profile to scaffold into .greenlight.yaml (e.g. angular, react, node, go, rails)")
	return cmd
}

// writeFrameworkConfig creates a minimal .greenlight.yaml selecting the given
// framework profile. It never clobbers an existing config: if the file is
// present, the caller is expected to add `framework:` manually.
func writeFrameworkConfig(repoDir, framework string) error {
	path := filepath.Join(repoDir, ".greenlight.yaml")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf(".greenlight.yaml already exists; add `framework: %s` manually", framework)
	} else if !os.IsNotExist(err) {
		return err
	}
	content := fmt.Sprintf("# greenlight config — see `greenlight profiles`\nframework: %s\n", framework)
	return os.WriteFile(path, []byte(content), 0o644)
}
