package internal

import (
	"sync"
	"time"

	"github.com/kihamo/shadow/components/database"
)

type MigrationItem struct {
	database.Migration

	mutex     sync.RWMutex
	source    string
	migration database.Migration
	appliedAt *time.Time
}

func NewMigrationItem(migration database.Migration, source string) *MigrationItem {
	if source == "" {
		source = "unknown"
	}

	return &MigrationItem{
		migration: migration,
		source:    source,
	}
}

func (m *MigrationItem) Source() string {
	return m.source
}

func (m *MigrationItem) Id() string {
	return m.migration.Id()
}

func (m *MigrationItem) Up() []string {
	return m.migration.Up()
}

func (m *MigrationItem) Down() []string {
	return m.migration.Down()
}

func (m *MigrationItem) ModAt() time.Time {
	return m.migration.ModAt()
}

func (m *MigrationItem) AppliedAt() *time.Time {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.appliedAt
}

func (m *MigrationItem) SetAppliedAt(appliedAt *time.Time) {
	m.mutex.Lock()
	m.appliedAt = appliedAt
	m.mutex.Unlock()
}
