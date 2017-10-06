package database

import (
	"time"
)

type HasMigrations interface {
	GetMigrations() []Migration
}

type Migration interface {
	Source() string
	Id() string
	Up() []string
	Down() []string
	ModAt() time.Time
	AppliedAt() *time.Time
}

type MigrationItem struct {
	source    string
	id        string
	up        []string
	down      []string
	modAt     time.Time
	appliedAt *time.Time
}

func NewMigration(source, id string, up, down []string, modAt time.Time, appliedAt *time.Time) Migration {
	return &MigrationItem{
		source:    source,
		id:        id,
		up:        up,
		down:      down,
		modAt:     modAt,
		appliedAt: appliedAt,
	}
}

func (m *MigrationItem) Source() string {
	return m.source
}

func (m *MigrationItem) Id() string {
	return m.id
}

func (m *MigrationItem) Up() []string {
	return m.up
}

func (m *MigrationItem) Down() []string {
	return m.down
}

func (m *MigrationItem) ModAt() time.Time {
	return m.modAt
}

func (m *MigrationItem) AppliedAt() *time.Time {
	return m.appliedAt
}
