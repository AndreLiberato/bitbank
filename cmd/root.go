package cmd

import (
	"fmt"
	"os"

	"github.com/AndreLiberato/bitbank/internal/repository"
	"github.com/AndreLiberato/bitbank/internal/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bitbank",
	Short: "BitBank — sistema bancário CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := repository.NewAccountRepository()
		if err != nil {
			return fmt.Errorf("erro ao inicializar repositório: %w", err)
		}
		svc := service.NewAccountService(repo)
		RunInteractive(svc)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
