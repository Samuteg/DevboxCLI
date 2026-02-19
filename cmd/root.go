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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "ficheiro de config (por omissão é $HOME/.devbox.yaml)")

	rootCmd.SetHelpFunc(helpFunc)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Define o nome e o caminho do ficheiro
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".devbox")
		viper.AutomaticEnv()
		viper.SetEnvPrefix("devbox")

		// --- DEFININDO DEFAULTS ---
		viper.SetDefault("author", "Devbox User")
		viper.SetDefault("default-port", "8080")
		viper.SetDefault("update-channel", "stable")
		viper.SetDefault("template-style", "clean")

		configPath := filepath.Join(home, ".devbox.yaml")

		// Se o ficheiro não existir, cria um vazio
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			os.WriteFile(configPath, []byte(""), 0644)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("devbox")

	if viper.ReadInConfig() == nil {
		// Pode descomentar a linha abaixo para debug inicial
		// fmt.Println("A usar ficheiro de configuração:", viper.ConfigFileUsed())
	}
}
