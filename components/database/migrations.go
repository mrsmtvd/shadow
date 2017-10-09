package database

import (
	"strings"
	"time"
)

type HasMigrations interface {
	GetMigrations() Migrations
}

type Migrations []Migration

type Migration interface {
	Source() string
	Id() string
	Up() []string
	Down() []string
	ModAt() time.Time
	AppliedAt() *time.Time
}

func (m Migrations) Len() int {
	return len(m)
}

func (m Migrations) Less(i, j int) bool {
	return strings.Compare(m[i].Id(), m[j].Id()) < 0
}

func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
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
