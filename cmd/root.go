// Package cmd is the root of all thing Munge
/*
Copyright © 2023 Nino Stephen <ninostephen.me>
*/
package cmd

import (
	"os"

	"github.com/ninostephen/munge/models"
	"github.com/ninostephen/munge/worker"
	"github.com/spf13/cobra"
)

var (
	flagvals models.Flags
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "munge",
	Short: "Dirty little word munger",
	Long: `
 ______                              
|  ___ \                             
| | _ | | _   _  ____    ____   ____ 
| || || || | | ||  _ \  / _  | / _  )
| || || || |_| || | | |( ( | |( (/ / 
|_||_||_| \____||_| |_| \_|| | \____)
                       (_____|       
  ______           ______                                      _  _  _  _  _ 
 / _____)         (____  \                                    | || || || || |
| /  ___   ___     ____)  )  ____   ____   ____   ____   ____ | || || || || |
| | (___) / _ \   |  __  (  / ___) / ___) / ___) / ___) / ___)|_||_||_||_||_|
| \____/|| |_| |  | |__)  )| |    | |    | |    | |    | |     _  _  _  _  _ 
 \_____/  \___/   |______/ |_|    |_|    |_|    |_|    |_|    |_||_||_||_||_|
                              Copyright © 2023 Nino Stephen <ninostephen.me>   
	
	A golang port of Mudge by Th3S3cr3tAg3nt
	
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		worker.Start(cmd, flagvals)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&flagvals.Word, "word", "w", "", "word to munge")
	rootCmd.Flags().IntVarP(&flagvals.Level, "level", "l", 2, "munge level [1-3] (default 2)")
	rootCmd.Flags().StringVarP(&flagvals.Input, "input", "i", "", "input file")
	rootCmd.Flags().StringVarP(&flagvals.Output, "output", "o", "", "output file")
}
