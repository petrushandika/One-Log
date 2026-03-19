package storage

import (
	"fmt"
)

// ArchiveStorage defines the interface for log archiving
type ArchiveStorage interface {
	IsEnabled() bool
	ArchiveLogs(date string, logs []byte) error
	RetrieveLogs(date string) ([]byte, error)
	ListArchivedDates(prefix string) ([]string, error)
	DeleteArchivedLogs(date string) error
}

// LogArchiver handles archiving old logs to object storage
type LogArchiver struct {
	storage ArchiveStorage
}

// NewLogArchiver creates a new log archiver
func NewLogArchiver(storage ArchiveStorage) *LogArchiver {
	return &LogArchiver{storage: storage}
}

// IsEnabled returns whether archiving is enabled
func (a *LogArchiver) IsEnabled() bool {
	return a.storage != nil && a.storage.IsEnabled()
}

// ArchiveLogs archives logs for a specific date
func (a *LogArchiver) ArchiveLogs(date string, logs []byte) error {
	if !a.IsEnabled() {
		return fmt.Errorf("log archiving is not enabled")
	}
	return a.storage.ArchiveLogs(date, logs)
}

// RetrieveLogs retrieves archived logs for a specific date
func (a *LogArchiver) RetrieveLogs(date string) ([]byte, error) {
	if !a.IsEnabled() {
		return nil, fmt.Errorf("log archiving is not enabled")
	}
	return a.storage.RetrieveLogs(date)
}

// ListArchivedDates returns list of archived dates
func (a *LogArchiver) ListArchivedDates() ([]string, error) {
	if !a.IsEnabled() {
		return nil, fmt.Errorf("log archiving is not enabled")
	}
	return a.storage.ListArchivedDates("logs/")
}

// DeleteArchivedLogs deletes archived logs for a specific date
func (a *LogArchiver) DeleteArchivedLogs(date string) error {
	if !a.IsEnabled() {
		return fmt.Errorf("log archiving is not enabled")
	}
	return a.storage.DeleteArchivedLogs(date)
}
