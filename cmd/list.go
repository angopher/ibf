package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list IBF",
	Short: "List available keys from the set.",
	Run: func(cmd *cobra.Command, args []string) {
		var path = args[0]

		file, err := os.Open(path)
		cannot(err)

		decoder := json.NewDecoder(file)

		ibf := ibf.NewEmptyIBF()
		err = decoder.Decode(ibf)
		cannot(err)
		file.Close()

		leftEmpty := true
		for val, err := ibf.Pop(); err == nil; val, err = ibf.Pop() {
			if !cfg.suppressLeft {
				fmt.Printf("%s\n", string(val.Bytes()))
			}
		}
		if !cfg.suppressLeft {
			leftEmpty = ibf.IsEmpty()
		}

		rightEmpty := true
		ibf.Invert()
		for val, err := ibf.Pop(); err == nil; val, err = ibf.Pop() {
			if !cfg.suppressRight {
				fmt.Printf("%s\n", string(val.Bytes()))
			}
		}
		if !cfg.suppressRight {
			rightEmpty = ibf.IsEmpty()
		}

		// Incomplete listing?
		if !leftEmpty || !rightEmpty {
			// Which side was empty?
			side := ""
			switch {
			case leftEmpty && rightEmpty:
				side = "left and right"
			case leftEmpty && !rightEmpty:
				side = "left"
			case !leftEmpty && rightEmpty:
				side = "right"
			}

			fmt.Fprintf(os.Stderr, "Unable to list all elements (%s).\n", side)
			os.Exit(1)
		}
	},
}

func init() {
	listCmd.Flags().BoolVarP(&cfg.suppressLeft, "left", "1", false, "Suppress values unique to left-side (positive count).")
	listCmd.Flags().BoolVarP(&cfg.suppressRight, "right", "2", false, "Suppress values unique to right-side (negative count).")

	RootCmd.AddCommand(listCmd)
}
