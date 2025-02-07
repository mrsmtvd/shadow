package database

import (
	"github.com/mrsmtvd/shadow"
)

type Component interface {
	shadow.Component

	Storage() Storage
	Migration(id, source string) Migration
	Migrations() []Migration
}
