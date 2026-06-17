package cli

import (
	"fmt"
	"os"

	"github.com/kinleyrabgay/greenlight/internal/config"
	"github.com/kinleyrabgay/greenlight/internal/paths"
	"github.com/kinleyrabgay/greenlight/internal/profile"
	"github.com/spf13/cobra"
)

func newProfilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "profiles",
		Short: "List available framework profiles and show the one detected here",
		Long: "Lists built-in and user-defined framework profiles. A profile bundles\n" +
			"default lint/test/format commands, ignore patterns, and review rules.\n\n" +
			"Select one per repo with `framework: <name>` in .greenlight.yaml, with\n" +
			"`greenlight init --framework <name>`, or let it auto-detect.\n\n" +
			"User profiles live in ~/.greenlight/profiles/<name>.md and override the\n" +
			"built-in of the same name.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return trackCommand("profiles", func() error {
				w := cmd.OutOrStdout()
				userDir := ""
				if p, err := paths.New(); err == nil {
					userDir = p.ProfilesDir()
				}

				fmt.Fprintf(w, "  %s\n", sBold.Render("Available framework profiles"))
				for _, name := range profile.List(userDir) {
					origin := sDim.Render("built-in")
					if _, err := os.Stat(profileUserFile(userDir, name)); err == nil {
						origin = sGreen.Render("user override")
					}
					fmt.Fprintf(w, "    %s  %s\n", sCyan.Render(name), origin)
				}

				cwd, _ := os.Getwd()
				if fw := config.ResolveFramework("", "", cwd); fw != "" {
					fmt.Fprintln(w)
					fmt.Fprintf(w, "  %s %s\n", sDim.Render("detected for this directory:"), sGreen.Render(fw))
				}
				return nil
			})
		},
	}
}

func profileUserFile(userDir, name string) string {
	if userDir == "" {
		return ""
	}
	return userDir + string(os.PathSeparator) + name + ".md"
}
