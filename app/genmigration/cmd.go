package genmigration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type command struct {
	gracePeriod  time.Duration
	migrationDir string
	name         string
}

func NewCommand() *cobra.Command {
	a := &command{
		migrationDir: "schema",
		gracePeriod:  30 * time.Second,
	}

	cmd := &cobra.Command{
		Use:   "gen-migration",
		Short: "Generate migration file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&a.name, "name", "n", a.name, "The name of migration file.")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

	return cmd
}

func (c *command) run() error {
	files := make([]string, 0)
	err := filepath.Walk(c.migrationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, info.Name())

		return nil
	})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return createFiles(1, c.name, c.migrationDir)
	}

	sort.Strings(files)

	lastFileName := files[len(files)-1]
	nameParts := strings.Split(lastFileName, "_")
	lastVersion, err := strconv.Atoi(nameParts[0])
	if err != nil {
		return fmt.Errorf("could not parse version %s from file name %s", nameParts[0], lastFileName)
	}

	return createFiles(lastVersion+1, c.name, c.migrationDir)
}

const fileTemplate = `-- Best practices:
-- 1. Do not change old migration files.
-- 2. Making your migrations idempotent.
-- 3. Use a clear and readable name.
-- Further notes: https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md
`

func createFiles(version int, name, dir string) error {
	var (
		upFile   = fmt.Sprintf("%s/%06d_%s.up.sql", dir, version, name)
		downFile = fmt.Sprintf("%s/%06d_%s.down.sql", dir, version, name)
	)

	if err := os.WriteFile(upFile, []byte(fileTemplate), 0o644); err != nil { //nolint:gosec
		return fmt.Errorf("could not write to %s: %w", upFile, err)
	}

	if err := os.WriteFile(downFile, []byte(fileTemplate), 0o644); err != nil { //nolint:gosec
		return fmt.Errorf("could not write to %s: %w", downFile, err)
	}

	return nil
}
