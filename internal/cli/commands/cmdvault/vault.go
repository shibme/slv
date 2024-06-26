package cmdvault

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/profiles"
	"oss.amagi.com/slv/internal/core/vaults"
)

func getVault(filePath string) (*vaults.Vault, error) {
	return vaults.Get(filePath)
}

func VaultCommand() *cobra.Command {
	if vaultCmd != nil {
		return vaultCmd
	}
	vaultCmd = &cobra.Command{
		Use:     "vault",
		Aliases: []string{"v", "vaults", "secret", "secrets"},
		Short:   "Manage vaults/secrets with SLV",
		Long:    `Manage vaults/secrets using SLV. SLV Vaults are files that store secrets in a key-value format.`,
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			sealedSecretsMap, err := vault.ListSealedSecrets()
			if err != nil {
				utils.ExitOnError(err)
			}
			accessors, err := vault.ListAccessors()
			if err != nil {
				utils.ExitOnError(err)
			}
			profile, _ := profiles.GetDefaultProfile()
			self := environments.GetSelf()
			envMap := make(map[string]string, len(accessors))
			for _, accessor := range accessors {
				var env *environments.Environment
				accessorStr, err := accessor.String()
				if err != nil {
					utils.ExitOnError(err)
				}
				selfEnv := false
				rootEnv := false
				if self != nil && self.PublicKey == accessorStr {
					env = self
					selfEnv = true
				} else if profile != nil {
					env, err = profile.GetEnv(accessorStr)
					if err != nil {
						utils.ExitOnError(err)
					}
					if env == nil {
						root, err := profile.GetRoot()
						if err != nil {
							utils.ExitOnError(err)
						}
						if root != nil && root.PublicKey == accessorStr {
							rootEnv = true
							env = root
						}
					}
				}
				if env != nil {
					if selfEnv {
						envMap[accessorStr] = accessorStr + "\t(" + color.CyanString("Self   "+": "+env.Name) + ")"
					} else if rootEnv {
						envMap[accessorStr] = accessorStr + "\t(" + color.CyanString("Root   "+": "+env.Name) + ")"
					} else {
						if env.EnvType == environments.USER {
							envMap[accessorStr] = accessorStr + "\t(" + "User   " + ": " + env.Name + ")"
						} else {
							envMap[accessorStr] = accessorStr + "\t(" + "Service" + ": " + env.Name + ")"
						}
						// envMap[accessorStr] = accessorStr + "\t(" + env.Name + ")"
					}
				} else {
					envMap[accessorStr] = accessorStr
				}
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, "Vault ID\t:\t", vault.Config.PublicKey)
			fmt.Fprintln(w, "Secrets:")
			for name, sealedSecret := range sealedSecretsMap {
				hash := sealedSecret.Hash()
				if hash == "" {
					fmt.Fprintln(w, "  -", name, "\t:\t", sealedSecret.EncryptedAt().Format("Jan _2, 2006 03:04:05 PM MST"))
				} else {
					fmt.Fprintln(w, "  -", name, "\t:\t", sealedSecret.EncryptedAt().Format("Jan _2, 2006 03:04:05 PM MST"), "\t(", hash, ")")
				}
			}
			fmt.Fprintln(w, "Accessible by:")
			for _, envDesc := range envMap {
				fmt.Fprintln(w, "  -", envDesc)
			}
			w.Flush()
			utils.SafeExit()
		},
	}
	vaultCmd.PersistentFlags().StringP(vaultFileFlag.Name, vaultFileFlag.Shorthand, "", vaultFileFlag.Usage)
	vaultCmd.MarkPersistentFlagRequired(vaultFileFlag.Name)
	vaultCmd.AddCommand(vaultNewCommand())
	vaultCmd.AddCommand(vaultToK8sCommand())
	vaultCmd.AddCommand(vaultPutCommand())
	vaultCmd.AddCommand(vaultGetCommand())
	vaultCmd.AddCommand(vaultShellCommand())
	vaultCmd.AddCommand(vaultDeleteCommand())
	vaultCmd.AddCommand(vaultRefCommand())
	vaultCmd.AddCommand(vaultDerefCommand())
	vaultCmd.AddCommand(vaultAccessCommand())
	return vaultCmd
}
