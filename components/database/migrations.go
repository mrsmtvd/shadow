package database

import (
	"time"
)

type HasMigrations interface {
	DatabaseMigrations() []Migration
}

type Migration interface {
	ID() string
	Up() []string
	Down() []string
	ModAt() time.Time
}

type MigrationSimple struct {
	id    string
	up    []string
	down  []string
	modAt time.Time
}

func NewMigration(id string, up, down []string, modAt time.Time) *MigrationSimple {
	return &MigrationSimple{
		id:    id,
		up:    up,
		down:  down,
		modAt: modAt,
	}
}

func (m *MigrationSimple) ID() string {
	return m.id
}

func (m *MigrationSimple) Up() []string {
	return m.up
}

func (m *MigrationSimple) Down() []string {
	return m.down
}

func (m *MigrationSimple) ModAt() time.Time {
	return m.modAt
}
