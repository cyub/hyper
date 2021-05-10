package mysql

import (
	"fmt"
	"time"

	"github.com/cyub/hyper/app"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// Provider use for mount to app bootstrap
func Provider() app.ComponentMount {
	return func(app *app.App) (err error) {

		dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?&loc=%s&parseTime=%s&collation=%s",
			app.Config.GetString("mysql.user", "homestead"),
			app.Config.GetString("mysql.passwd", "secret"),
			app.Config.GetString("mysql.host", "localhost"),
			app.Config.GetInt("mysql.port", 3306),
			app.Config.GetString("mysql.db", "hyper"),
			app.Config.GetString("mysql.loc", "UTC"),
			app.Config.GetString("mysql.parse_time", "true"),
			app.Config.GetString("mysql.collation", "utf8mb4_general_ci"),
		)
		// all parameters see https://github.com/go-sql-driver/mysql#parameters
		// https://github.com/go-sql-driver/mysql#charset
		charset := app.Config.GetString("mysql.charset", "")
		if len(charset) > 0 {
			dsn += "&charset=" + charset
		}

		db, err = gorm.Open("mysql", dsn)
		if err != nil {
			panic(fmt.Errorf("mysql open failure %s, %s", dsn, err.Error()))
		}

		switch app.RunMode {
		case "prod", "production", "release":
			db.LogMode(false)
		default:
			db.LogMode(true)
		}

		db.DB().SetMaxIdleConns(app.Config.GetInt("mysql.max_idle_conns", 10))
		db.DB().SetMaxOpenConns(app.Config.GetInt("mysql.max_open_conns", 50))
		db.DB().SetConnMaxLifetime(time.Duration(app.Config.GetInt("mysql.max_lift_time", 5)) * 30)

		if err = db.DB().Ping(); err != nil {
			app.Logger.Errorf("mysql ping failure %s, %s", dsn, err.Error())
			return
		}
		app.Logger.Info("mysql ping is ok")
		return
	}
}

// Instance return the instance of gorm.DB
func Instance() *gorm.DB {
	return db
}
