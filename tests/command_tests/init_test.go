package command_tests

import (
	"os"
	"testing"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/core/utils"
	"github.com/KambojRajan/ship/tests/helpers"
)

func TestInit_forCurrentDir_itShouldCreateBasicDirStructure(t *testing.T) {
	dir := t.TempDir()
	err := os.Chdir(dir)
	if err != nil {
		return
	}

	err = commands.Init(".")
	if err != nil {
		t.Error(err)
	}

	helpers.AssertExists(t, utils.RootShipDir)
	helpers.AssertExists(t, utils.RootIndexPath)
	helpers.AssertExists(t, utils.RootObjectDir)
	helpers.AssertExists(t, utils.RootHEADPath)
}

func TestInit_forExistingDir_initShouldBeIdempotent(t *testing.T) {
	dir := t.TempDir()

	err := commands.Init(dir)
	helpers.AssertNil(err)

	err = commands.Init(dir)
	helpers.AssertNil(err)
}

func TestInit_forNonExistingDir_itShouldFail(t *testing.T) {
	err := commands.Init("pata")
	helpers.AssertNotNil(err)
}

func TestInit_passingNonDirPath_itShouldFail(t *testing.T) {
	err := commands.Init("mx.md")
	helpers.AssertNotNil(err)
}

func TestInit_forEmptyDir_itShouldPass(t *testing.T) {
	dir := t.TempDir()
	err := commands.Init(dir)
	helpers.AssertNil(err)
}
