package migrations

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/fileutils"
	"github.com/stackrox/rox/pkg/migrations/internal"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	// Current is the current database in use.
	Current = "current"
	// PreviousReplica is the symbolic link pointing to the previous databases.
	PreviousReplica = ".previous"
)

// DBMountPath is the directory path (within a container) where database storage device is mounted.
func DBMountPath() string {
	return internal.DBMountPath
}

// CurrentPath is the link (within a container) to current migration directory. This directory contains
// databases and other migration related contents.
func CurrentPath() string {
	return filepath.Join(internal.DBMountPath, Current)
}

// SafeRemoveDBWithSymbolicLink removes databases in path if it exists, it protects current database and remove only the databases that is not in use.
// If path is a symbolic link, remove it and the database it points to.
func SafeRemoveDBWithSymbolicLink(path string) error {
	currentDB, err := fileutils.ResolveIfSymlink(CurrentPath())
	utils.Must(errors.Wrap(err, "no current database"))

	switch path {
	case CurrentPath(), currentDB:
		utils.Must(errors.Errorf("Database in use. Cannot remove %s", path))
	default:
		if exists, err := fileutils.Exists(path); err != nil {
			return err
		} else if exists {
			log.Infof("Removing database %s", path)
			linkTo, err := fileutils.ResolveIfSymlink(path)
			if err != nil {
				return err
			}
			if err = os.RemoveAll(path); err != nil {
				return err
			}
			// Remove linked database if it is not the current database
			if linkTo != CurrentPath() && linkTo != currentDB {
				return os.RemoveAll(linkTo)
			}
		}
	}
	return nil
}
