package database

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Storage() Storage
	Migration(id, source string) Migration
	Migrations() Migrations
}
