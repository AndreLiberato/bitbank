package cmd

import (
	"fmt"

	"github.com/AndreLiberato/bitbank/internal/api"
	"github.com/AndreLiberato/bitbank/internal/repository"
	"github.com/AndreLiberato/bitbank/internal/service"
	"github.com/spf13/cobra"
)

var servePort string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia a API REST do BitBank",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := repository.NewAccountRepository()
		if err != nil {
			return fmt.Errorf("erro ao inicializar repositório: %w", err)
		}
		svc := service.NewAccountService(repo)
		return api.Serve(":"+servePort, svc)
	},
}

func init() {
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8080", "porta HTTP da API REST")
	rootCmd.AddCommand(serveCmd)
}
