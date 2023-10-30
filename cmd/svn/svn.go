package svn

import (
	"fmt"
	"github.com/core/tools/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
)

var SvnUpCmd = &cobra.Command{
	Use:   "svn",
	Short: "svn update",
	Run:   svnUpdateFn,
}
var svnPath, fileExe string

func init() {
	SvnUpCmd.PersistentFlags().StringVarP(&svnPath, "svnPath", "sp", "E:\\remoteWork\\server\\Config\\Excels\\Main",
		"svn update path")
	err := viper.BindPFlag("svnPath", SvnUpCmd.PersistentFlags().Lookup("svnPath"))
	if err != nil {
		log.Error("svn err")
		return
	}
}
func svnUpdateFn(cmd *cobra.Command, args []string) {
	// 执行 svn svn 命令
	cmde := exec.Command("svn", "update")
	cmde.Dir = filepath.Dir(svnPath)
	cmde.Stdout = os.Stdout
	cmde.Stderr = os.Stderr
	err := cmde.Run()
	if err != nil {
		fmt.Printf("Error executing proto svn svn: %s\n", err)
		os.Exit(1)
	}
}

var ExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "exec command",
	Run:   runExec,
}

func runExec(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gen exec <exePath> [arguments...]")
		return
	}
	exePath := args[0]
	exeArgs := args[1:]

	cmde := exec.Command(exePath, exeArgs...)
	cmdRun(cmde)
}

func cmdRun(cmde *exec.Cmd) {
	cmde.Stdout = os.Stdout
	cmde.Stderr = os.Stderr
	err := cmde.Run()
	if err != nil {
		fmt.Printf("Error executing proto svn svn: %s\n", err)
		os.Exit(1)
	}
}
