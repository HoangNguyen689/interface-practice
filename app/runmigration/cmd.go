package runmigration

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/HoangNguyen689/interface-practice/app/runmigration/sqlmigration"
	"github.com/HoangNguyen689/interface-practice/pkg/postgresql"
)

const schemaDir = "./schema"

type command struct {
	migrationVersion uint
}

func NewCommand() *cobra.Command {
	c := &command{}

	cmd := &cobra.Command{
		Use:   "run-migration",
		Short: "Run migration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().UintVar(&c.migrationVersion, "migration-version", c.migrationVersion, "Migration version number where should be migrated to. Default is up to the latest version.")

	return cmd
}

func (c *command) run() error {
	ctx := context.Background()

	cli, err := postgresql.NewClient(
		ctx,
		"localhost:5432",
		"db",
		"postgres",
		"postgres",
		"run-migration",
	)
	if err != nil {
		return err
	}

	const timeout = 15 * time.Second

	if err := cli.WaitReady(ctx, time.Second, timeout); err != nil {
		return err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	migrator := sqlmigration.NewMigrator(
		cli,
		schemaDir,
		"db",
		logger,
	)

	if err := migrator.Execute(ctx, c.migrationVersion); err != nil {
		return err
	}

	return nil
}
