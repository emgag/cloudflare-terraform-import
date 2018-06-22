package cmd

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/emgag/cloudflare-terraform-import/internal/lib/dns"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import <zone>",
	Short: "Generate import files for zone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))

		if err != nil {
			log.Fatal(err)
		}

		e := dns.NewExporter(api)

		if err != nil {
			log.Fatal(err)
		}

		tf, err := os.Create("import.tf")

		if err != nil {
			log.Fatal(err)
		}

		sh, err := os.Create("import.sh")

		if err != nil {
			log.Fatal(err)
		}

		e.DumpZone(args[0], tf, sh)
	},
}
