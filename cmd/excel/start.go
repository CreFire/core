package excel

import (
	"fmt"
	"github.com/core/tools"
	"github.com/core/tools/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// 创建生成excel配置
type excelConfig struct {
	BasePath            string
	ExcelPath           string
	GenCsPath           string
	GenGoPath           string
	GenJsonPath         string
	GenProtobufPath     string
	GenProtobufDataPath string
	GenProtobufCSPath   string
	CustomTypePath      string
	CsvPath             string
	Md5Path             string
}

var Conf = &excelConfig{}
var cfgFile string
var ExcelCmd = &cobra.Command{
	Use:   "gengo",
	Short: "generator excel to go json protoBuff",
	Long:  `excel配置生成`,
	Run:   GenStart,
}

// GenStart 开始生成Excel
func init() {
	cobra.OnInitialize(initConfig)
	ExcelCmd.PersistentFlags().StringVar(&cfgFile, "config", "./Conf/excel.json",
		"config file (default is $HOME/.cobra.yaml)")
}

// excelConfig 从JSON配置文件中加载配置
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&Conf)
	if err != nil {
		log.Error("Error unmarshaling config: %s", log.Err(err))
	}

	Conf.BasePath = filepath.FromSlash(Conf.BasePath)
	Conf.ExcelPath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.ExcelPath))
	Conf.GenCsPath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.GenCsPath))
	Conf.GenGoPath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.GenGoPath))
	Conf.GenJsonPath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.GenJsonPath))
	Conf.GenProtobufPath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.GenProtobufPath))
	Conf.GenProtobufDataPath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.GenProtobufDataPath))
	Conf.GenProtobufCSPath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.GenProtobufCSPath))
	Conf.CustomTypePath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.CustomTypePath))
	Conf.CsvPath = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.CsvPath))
	Conf.Md5Path = filepath.Join(Conf.BasePath, filepath.FromSlash(Conf.Md5Path))
}

func GenStart(cmd *cobra.Command, args []string) {
	tools.TryCreateDir(Conf.ExcelPath)
	tools.TryCreateDir(Conf.GenCsPath)
	tools.TryCreateDir(Conf.GenGoPath)
	tools.TryCreateDir(Conf.GenJsonPath)
	tools.TryCreateDir(Conf.GenProtobufPath)
	tools.TryCreateDir(Conf.GenProtobufDataPath)
	tools.TryCreateDir(Conf.GenProtobufCSPath)
	tools.TryCreateDir(Conf.CustomTypePath)
	tools.TryCreateDir(Conf.CsvPath)
	tools.TryCreateDir(Conf.Md5Path)
	getAllExcel(Conf.ExcelPath)
}
func getAllExcel(excelPath string) {
	//matchPath := filepath.Join(excelPath, "*.xls?")
	//fileNames, err := filepath.Glob(matchPath)
	//if err != nil {
	//	log.Error("getAllExcel err", log.Err(err))
	//}
	//
}
