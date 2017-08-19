package setup

import (
	"database/sql"
	"fmt"
	"strings"
)

/*
CreateUser : Creates a new user in the database
*/
func CreateUser(connDetail ConnectionDetails, userName string, options []string) error {
	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}

	var exists bool

	if exists, err = checkIfUserExists(connDetail, userName); err != nil {
		return err
	} else if exists {
		// returns if the user already exists
		return nil
	}

	_, err = db.Exec("CREATE USER dbview " + strings.Join(options, " ") + ";")
	return err
}

func checkIfUserExists(connDetail ConnectionDetails, userName string) (bool, error) {

	var db *sql.DB
	var err error

	found := false

	if db, err = connect(connDetail); err != nil {
		return found, err
	}

	totalRows := 0
	err = db.QueryRow("SELECT count(1) FROM pg_roles WHERE rolname = $1", userName).Scan(&totalRows)

	if totalRows > 0 {
		found = true
	}

	return found, err
}

/*
GrantRolesToUser : Grant some roles privilieges for a user
*/
func GrantRolesToUser(connDetail ConnectionDetails, userName string, roles []string) error {
	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}

	var exists bool

	if exists, err = checkIfUserExists(connDetail, userName); err != nil {
		return err
	} else if !exists {
		// returns if the user not exists
		return nil
	}

	for _, role := range roles { // GRANT dbview TO ${DBVIEW_USER};
		_, err = db.Exec(fmt.Sprintf("GRANT %s to %s;", role, userName))
		if err != nil {
			return err
		}
	}

	return nil
}
