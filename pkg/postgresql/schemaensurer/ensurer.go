package schemaensurer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/HoangNguyen689/interface-practice/pkg/postgresql"
)

type Ensurer struct {
	client    *postgresql.Client
	schemaDir string
	logger    *log.Logger
}

func NewEnsurer(client *postgresql.Client, schemaDir string, logger log.Logger) (*Ensurer, error) {
	return &Ensurer{
		client:    client,
		schemaDir: schemaDir,
		logger:    log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

func (e *Ensurer) Ensure(ctx context.Context) error {
	files := make([]string, 0)
	err := filepath.Walk(e.schemaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, filepath.Join(e.schemaDir, info.Name()))

		return nil
	})
	if err != nil {
		return err
	}

	sort.Strings(files)
	e.logger.Println(fmt.Sprintf("there are %d schema files to apply", len(files)))

	for _, file := range files {
		if err := e.ensureFile(ctx, file); err != nil {
			return err
		}
	}

	return nil
}

func (e *Ensurer) ensureFile(ctx context.Context, path string) error {
	statements, err := loadStatements(path)
	if err != nil {
		return err
	}

	for _, stmt := range statements {
		e.logger.Println("start applying a statement")
		e.logger.Println(stmt)

		if _, err := e.client.Exec(ctx, stmt); err != nil {
			return err
		}
		e.logger.Println("successfully applied statement")
	}

	return nil
}

func loadStatements(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	items := strings.Split(strings.TrimSpace(string(data)), ";")
	statements := make([]string, 0, len(items))
	for _, item := range items {
		// Ignore dummy statement.
		if item == "" {
			continue
		}
		statements = append(statements, strings.TrimSpace(item))
	}

	return statements, nil
}
