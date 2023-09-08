package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	OutputDir string `json:"outputDir"`
	ProtoFile string `json:"protoFile"`
	EnvBase   string `json:"envBase"`
}

var globalConfig map[string]*Config

func main() {
	loadConfig()

	var outputDir, protoFile, base string

	var rootCmd = &cobra.Command{Use: "rpc"}

	var cmdGenerate = &cobra.Command{
		Use:   "gen",
		Short: "genCode RPC code",
		Run: func(cmd *cobra.Command, args []string) {
			var conf = globalConfig["main"]
			if len(args) > 0 {

			} else {

			}
			generateGRPCCode(conf)
		},
	}

	var cmdSet = &cobra.Command{
		Use:   "set",
		Short: "Set default values",
		Run: func(cmd *cobra.Command, args []string) {
			var conf = globalConfig["main"]
			if outputDir != "" {
				conf.OutputDir = outputDir
			}
			if protoFile != "" {
				conf.ProtoFile = protoFile
			}
			saveConfig(conf)
		},
	}

	cmdGenerate.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated code (default is read from config)")
	cmdGenerate.Flags().StringVarP(&protoFile, "proto", "p", "", "Path to the .proto file (default is read from config)")
	cmdGenerate.Flags().StringVarP(&base, "base", "b", "", "Base Exe path (default is read from config)")

	cmdSet.Flags().StringVarP(&outputDir, "output", "o", "", "Set default output directory")
	cmdSet.Flags().StringVarP(&protoFile, "proto", "p", "", "Set default .proto file")

	rootCmd.AddCommand(cmdGenerate, cmdSet)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig() {
	file, err := os.OpenFile("config.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error reading config file, ", err)
		os.Exit(1)
	}
	var buf [256]byte
	n, err := file.Read(buf[:])
	if err != nil {
		fmt.Println(err)
	}
	str := strings.Fields(string(buf[:n]))
	fmt.Println("thisBody:", str)
	globalConfig = make(map[string]*Config)
	err = json.Unmarshal(buf[:n], &globalConfig)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println(globalConfig)
	if err != nil {
		fmt.Println("Error parsing config file, ", err)
		os.Exit(1)
	}
}

func saveConfig(config *Config) {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error saving config file, ", err)
		return
	}
	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing config file, ", err)
		return
	}
	fmt.Println("Configuration saved.")
}

func generateGRPCCode(config *Config) {
	cmd := exec.Command(
		fmt.Sprintf("%s/protoc.exe", config.EnvBase),
		fmt.Sprintf("--plugin=protoc-gen-go=%s/protoc-gen-gogofast.exe", config.EnvBase),
		fmt.Sprintf("--go_out=import_path=pb:%s", config.OutputDir),
		fmt.Sprintf("-I=%s", config.ProtoFile),
		fmt.Sprintf("%s/*.proto", config.ProtoFile),
	)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error generating gRPC code: %v\n", err)
		return
	}

	fmt.Println("gRPC code generated successfully")
}
