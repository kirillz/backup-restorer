package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

const (
	pgdataDir    = "/srv/pg1/pgdata"
	walBackupDir = pgdataDir + "/pg_wal_backup"
)

// RestoreBackupMsg is a message sent when the backup restoration is complete.
type RestoreBackupMsg struct{}

func RestoreBackupCmd(backupDir string) tea.Cmd {
	return func() tea.Msg {
		err := restoreBackup(backupDir)
		if err != nil {
			logrus.Error("Error:", err)
		} else {
			logrus.Info("Backup restored successfully.")
		}
		return RestoreBackupMsg{}
	}
}

func restoreBackup(backupDir string) error {
	// Проверка наличия утилиты pg_waldump
	if _, err := exec.LookPath("pg_waldump"); err != nil {
		return fmt.Errorf("pg_waldump utility not found. Is Postgres installed correctly?")
	}

	// Проверка свободного места
	if err := checkFreeSpace(); err != nil {
		return err
	}

	// Остановка PostgreSQL
	if err := stopPostgres(); err != nil {
		return err
	}

	// Восстановление базы данных
	if err := copyBackup(backupDir); err != nil {
		return err
	}

	// Запуск PostgreSQL
	if err := startPostgres(); err != nil {
		return err
	}

	return nil
}

func checkFreeSpace() error {
	// Проверка свободного места на разделе
	dfOutput, err := exec.Command("df", "/usr").Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(dfOutput), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("unexpected output from df command")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return fmt.Errorf("unexpected output from df command")
	}
	freeSpace := fields[3]
	if freeSpace == "0" {
		return fmt.Errorf("not enough free space on the partition")
	}
	return nil
}

func stopPostgres() error {
	return exec.Command("service", "postgresql", "stop").Run()
}

func startPostgres() error {
	return exec.Command("service", "postgresql", "start").Run()
}

func copyBackup(backupDir string) error {
	// Удаление старых файлов
	if err := os.RemoveAll(pgdataDir); err != nil {
		return err
	}

	// Копирование файлов из резервной копии
	if err := filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destPath := filepath.Join(pgdataDir, strings.TrimPrefix(path, backupDir))
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		return os.Rename(path, destPath)
	}); err != nil {
		return err
	}

	// Копирование WAL файлов
	if err := filepath.Walk(walBackupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destPath := filepath.Join(pgdataDir, "pg_wal", strings.TrimPrefix(path, walBackupDir))
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		return os.Rename(path, destPath)
	}); err != nil {
		return err
	}

	return nil
}
