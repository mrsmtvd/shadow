package internal

import (
	"strings"
)

type MigrationsCollection []*MigrationItem

func (m MigrationsCollection) Len() int {
	return len(m)
}

func (m MigrationsCollection) Less(i, j int) bool {
	return strings.Compare(m[i].ID(), m[j].ID()) < 0
}

func (m MigrationsCollection) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
