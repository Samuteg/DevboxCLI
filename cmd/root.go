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
	Long:  `Uma interface de linha de comandos para automatizar e padronizar fluxos de trabalho.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		HandleError(err, "Falha Crítica")
		os.Exit(1)
	}
}

func init() {
	// Diz ao Cobra para correr esta função quando arrancar
	cobra.OnInitialize(initConfig)

	// Permite passar um ficheiro de configuração customizado via flag
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "ficheiro de config (por omissão é $HOME/.devbox.yaml)")

	rootCmd.SetHelpFunc(helpFunc)
}

// initConfig lê o ficheiro de configuração guardado
func initConfig() {
	if cfgFile != "" {
		// Usa o ficheiro passado pela flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Encontra a diretoria Home do utilizador
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
		viper.SetDefault("update-channel", "stable") // pode ser 'stable' ou 'beta'
		viper.SetDefault("template-style", "clean")

		// Caminho completo para podermos criar o ficheiro se não existir
		configPath := filepath.Join(home, ".devbox.yaml")

		// Se o ficheiro não existir, cria um vazio
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			os.WriteFile(configPath, []byte(""), 0644)
		}
	}

	// Lê as variáveis de ambiente com o prefixo DEVBOX_ (ex: DEVBOX_AUTHOR)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("devbox")

	// Se encontrar o ficheiro, lê os dados
	if err := viper.ReadInConfig(); err == nil {
		// Pode descomentar a linha abaixo para debug inicial
		// fmt.Println("A usar ficheiro de configuração:", viper.ConfigFileUsed())
	}
}
