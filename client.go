package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"os"

	"github.com/spf13/cobra"
)

func main1() {
	var hashType string

	rootCmd := &cobra.Command{
		Use:   "hash [string]",
		Short: "Calculate hash value of a string",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			str := args[0]
			var h hash.Hash

			switch hashType {
			case "md5":
				h = md5.New()
			case "sha1":
				h = sha1.New()
			case "sha256":
				h = sha256.New()
			default:
				fmt.Printf("Invalid hash type: %s\n", hashType)
				os.Exit(1)
			}

			h.Write([]byte(str))
			hashBytes := h.Sum(nil)
			hashString := hex.EncodeToString(hashBytes)

			fmt.Printf("%s: %s\n", hashType, hashString)
		},
	}

	rootCmd.Flags().StringVarP(&hashType, "type", "t", "md5", "Hash type (md5, sha1, sha256)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
