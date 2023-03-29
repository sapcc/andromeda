/*
 *   Copyright 2023 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package housekeeping

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/xo/dburl"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
)

func HouseKeeping() error {
	if !config.Global.HouseKeeping.Enabled {
		logger.Fatal("Housekeeping disabled")
	}
	logger.Info("Running housekeeping")

	// Database
	u, err := dburl.Parse(config.Global.Database.Connection)
	if err != nil {
		return err
	}
	db := sqlx.MustConnect(u.Driver, u.DSN)

	// Mapper function for SQL name mapping, snake_case table names
	db.MapperFunc(strcase.ToSnake)

	sql := fmt.Sprintf(`
	DELETE FROM domain 
	WHERE provisioning_status = 'DELETED' AND updated_at < (NOW() - INTERVAL %d SECOND)`,
		config.Global.HouseKeeping.DeleteAfter)
	res := db.MustExec(sql)

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	logger.Infof("Cleaned up %d domains", rowsAffected)
	return db.Close()
}
