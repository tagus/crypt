package cobracli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/store"
)

var diffCmd = &cobra.Command{
	Use:   "diff [cryptfile]",
	Short: "compares cryptfiles",
	Long: `compares the current cryptfile with the one given using
crypt and credential fingerprints.`,
	Args: cobra.ExactArgs(1),
	RunE: diff,
}

func diff(cmd *cobra.Command, args []string) error {
	asker := asker.DefaultAsker()
	pwd, err := asker.AskSecret(color.YellowString("second cryptfile pwd"), false)
	if err != nil {
		return err
	}

	alt, err := store.Decrypt(args[0], pwd)
	if err != nil {
		return err
	}

	st, err := getStore()
	if err != nil {
		return err
	}

	if alt.Fingerprint == st.Fingerprint {
		color.Green("given cryptfiles are the exact same")
		return nil
	}

	for svc := range alt.Credentials {
		other := alt.GetCredential(svc)
		cur := st.GetCredential(svc)
		if cur == nil {
			color.Red("%s not in cryptfile", other.Service)
			continue
		}

		if cur.Fingerprint != other.Fingerprint {
			color.Yellow("%s discrepancy", other.Service)
		}
	}

	return nil
}
