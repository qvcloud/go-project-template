package provider

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

func NewViper(pathStr string) *viper.Viper {
	var v = viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	if pathStr != "" {
		v.SetConfigFile(pathStr)
	} else {
		// for unit test
		if IsDevelopment() {
			_, current, _, _ := runtime.Caller(0)
			root := path.Dir(path.Dir(path.Dir(path.Dir(current))))
			fmt.Printf("root: %s\n", root)
			v.AddConfigPath(root)
		}

		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("warn: not found config: %v\n", err)
	}

	fmt.Printf("config file: %s\n", v.ConfigFileUsed())
	return v
}

func Root() string {
	_, current, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(path.Dir(current))))
	return root
}

func IsDevelopment() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	if strings.Contains(exe, "__debug") {
		return true
	}
	return strings.Contains(exe, os.TempDir())
}
