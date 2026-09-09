package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "devbox",
	Short: "Devbox é a sua CLI de produtividade",
	Long:  ` Uma CLI para automatizar o setup de projetos, git e tarefas do dia a dia.                                              `,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		HandleError(err, "Falha Crítica")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "arquivo de configuração (padrão: $HOME/.devbox.yaml)")

	rootCmd.SetHelpFunc(helpFunc)
}

func initConfig() {
	viper.SetDefault("author", "Devbox User")
	viper.SetDefault("default-port", "8080")
	viper.SetDefault("update-channel", "stable")
	viper.SetDefault("template-style", "clean")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".devbox")
		configPath := filepath.Join(home, ".devbox.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// Cria uma vez; erro aqui não é fatal (ex: $HOME readonly).
			_ = os.WriteFile(configPath, []byte(""), 0644)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("DEVBOX")

	_ = viper.ReadInConfig()
}
