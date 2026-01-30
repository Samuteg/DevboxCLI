package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const AppVersion = "0.0.1"

var cfgFile string

// rootCmd representa o comando base quando chamado sem subcomandos
var rootCmd = &cobra.Command{
	Use:   "DevboxCLI",
	Short: "Sua caixa de ferramentas pessoal para desenvolvimento",
	Long:  `Uma CLI para automatizar o setup de projetos, git e tarefas do dia a dia.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Erro no Cobra:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Permite passar --config manualmente
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "arquivo de config (padrão é $HOME/.devbox.yaml)")

	// Altera o template de versão para algo mais limpo (opcional)
	rootCmd.SetVersionTemplate("DevBox CLI version {{.Version}}\n")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Procura na home do usuário
		home, err := os.UserHomeDir()
		if err != nil {
			cobra.CheckErr(err)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".") // Procura também na pasta atual
		viper.SetConfigType("yaml")
		viper.SetConfigName(".devbox")
	}

	viper.AutomaticEnv() // Lê variáveis de ambiente (ex: DEV_TOKEN)

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Usando config:", viper.ConfigFileUsed()) // Debug opcional
	}
}
