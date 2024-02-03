package handlerDb

import (
	"database/sql"
	"errors"
	"fmt"
	"go_test/share"
	"log"
	"regexp"
	"strings"
)

const (
	dataTableName = "go_test"
	flagTableName = "go_flag"
	flagTTL       = 300
	pg_user       = "postgres"
	pg_password   = ""
	pg_port       = 5432
	pg_host       = "db"
)
const createTableSQL = `
create table ` + dataTableName + `
(
    uid       integer
        constraint ` + dataTableName + `_pk
            primary key,
    firstName varchar(200),
    lastName  varchar(200)
);
create index ` + dataTableName + `_firstName_index
    on ` + dataTableName + ` (firstName);
create index ` + dataTableName + `_lastName_index
    on ` + dataTableName + ` (lastName);
create table ` + flagTableName + `
(
    name       varchar(10)
        constraint ` + flagTableName + `_pk
            primary key,
    value integer,
	till timestamp
);
`

var Db *sql.DB

func Init_db() (err error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", pg_host, pg_port, pg_user, pg_password)
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	rows, err := Db.Query("SELECT EXISTS (SELECT FROM information_schema.tables WHERE  table_schema = 'public' AND table_name = $1)", dataTableName)
	if err != nil {
		return err
	}
	if rows.Next() {
		var tableExists bool
		err = rows.Scan(&tableExists)
		if err != nil {
			return err
		}
		if !tableExists {
			_, err := Db.Exec(createTableSQL)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SetFlag(flag string, value int) (err error) {
	interval := fmt.Sprintf("'%d sec'", flagTTL)
	result, err := Db.Exec("insert into "+flagTableName+" (name, value, till) values ($1, $2, now() + interval "+interval+") on conflict (name) do update set value = $2, till = now() + interval "+interval, flag, value)
	if err == nil {
		resultRows, err := result.RowsAffected()
		if err == nil && resultRows != 1 {
			err = errors.New("unable to set flag")
		}
	}
	return err
}
func GetFlag(flag string) (value int, err error) {
	result := 0
	rows, err := Db.Query("select value from "+flagTableName+" where name = $1 and till >= now()", flag)
	if err == nil && rows.Next() {
		err = rows.Scan(&result)
	}
	return result, err
}

func InsertEntry(entry *share.XmlEntry) (err error) {
	_, err = Db.Exec("insert into "+dataTableName+" (uid, firstname, lastname) values ($1, $2, $3) on conflict (uid) do update set firstname=$2, lastname=$3", entry.Uid, entry.FirstName, entry.LastName)
	return err
}

func SearchWeak(names string) (foundEntry []share.JsonEntry, err error) {
	return makeSearchQuery(strings.ReplaceAll(regexp.QuoteMeta(names), " ", "|"))
}

func SearchStrong(names string) (foundEntry []share.JsonEntry, err error) {
	if namesArray := strings.Split(names, " "); len(namesArray) == 2 {
		rows, err := Db.Query("select uid, firstname, lastname from "+dataTableName+" where firstname = $1 and lastname = $2", namesArray[0], namesArray[1])
		if err == nil {
			foundEntry = buildJsonEntryFromRows(rows)
		}
		return foundEntry, err
	} else {
		return makeSearchQuery("^(" + strings.ReplaceAll(regexp.QuoteMeta(names), " ", "|") + ")$")
	}
}

func makeSearchQuery(name string) (foundEntry []share.JsonEntry, err error) {
	rows, err := Db.Query("select uid, firstname, lastname from "+dataTableName+" where firstname ~* $1 or lastname ~* $1", name)
	if err == nil {
		foundEntry = buildJsonEntryFromRows(rows)
	}
	return foundEntry, err
}
func buildJsonEntryFromRows(rows *sql.Rows) (jsonEntryArray []share.JsonEntry) {
	result := []share.JsonEntry{}
	var entry share.JsonEntry
	for rows.Next() {
		err := rows.Scan(&entry.Uid, &entry.FirstName, &entry.LastName)
		if err != nil {
			log.Println(err)
		} else {
			result = append(result, entry)
		}
	}
	return result
}
