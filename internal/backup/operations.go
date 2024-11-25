package backup

import (
	"log"
	"os"
	"path/filepath"
)

func DeleteBackupFiles(backupFileNames, consistencyCheckReports []string) error {
	if value, present := os.LookupEnv("KEEP_BACKUP_FILES"); present && value == "false" {
		backupDir := "/backups" // default
		if dir, exists := os.LookupEnv("BACKUP_DIR"); exists {
			backupDir = dir
		}

		for _, backupFileName := range backupFileNames {
			filePath := filepath.Join(backupDir, backupFileName)
			log.Printf("Deleting file %s", filePath)
			err := os.Remove(filePath)
			if err != nil {
				return err
			}
		}
		for _, consistencyCheckReportName := range consistencyCheckReports {
			filePath := filepath.Join(backupDir, consistencyCheckReportName)
			log.Printf("Deleting file %s", filePath)
			err := os.Remove(filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
