/*
	Copyright (C) 2024  Pagefault Games

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var handle *sql.DB

func Init(username, password, protocol, address, database string) error {
	var err error

	handle, err = sql.Open("mysql", username+":"+password+"@"+protocol+"("+address+")/"+database)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %s", err)
	}

	handle.SetMaxIdleConns(256)
	handle.SetMaxOpenConns(256)
	handle.SetConnMaxIdleTime(time.Second * 30)
	handle.SetConnMaxLifetime(time.Minute)

	tx, err := handle.Begin()
	if err != nil {
		panic(err)
	}
	tx.Exec("CREATE TABLE IF NOT EXISTS systemSaveData (uuid BINARY(16) PRIMARY KEY, data LONGBLOB, timestamp TIMESTAMP)")
	tx.Exec("CREATE TABLE IF NOT EXISTS sessionSaveData (uuid BINARY(16), slot TINYINT, data LONGBLOB, timestamp TIMESTAMP, PRIMARY KEY (uuid, slot))")
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return nil
}

func PrepareTables() error {
	err := PrepareAccountTables()
	if err != nil {
		return fmt.Errorf("failed to prepare account tables: %s", err)
	}

	err = PrepareDailyTables()
	if err != nil {
		return fmt.Errorf("failed to prepare daily tables: %s", err)
	}

	err = PrepareGameTables()
	if err != nil {
		return fmt.Errorf("failed to prepare game tables: %s", err)
	}

	err = PrepareSaveDataTables()
	if err != nil {
		return fmt.Errorf("failed to prepare savedata tables: %s", err)
	}

	return nil
}
