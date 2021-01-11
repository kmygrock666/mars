package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Mysql sql
type Mysql struct {
	ReadHost     string `mapstructure:"READ_HOST" json:"READ_HOST"`
	ReadPort     int    `mapstructure:"READ_PORT" json:"READ_PORT"`
	ReadUser     string `mapstructure:"READ_USER" json:"READ_USER"`
	ReadPassword string `mapstructure:"READ_PASSWORD" json:"READ_PASSWORD"`
	ReadDatabase string `mapstructure:"READ_DATABASE" json:"READ_DATABASE"`

	WriteHost     string `mapstructure:"WRITE_HOST" json:"WRITE_HOST"`
	WritePort     int    `mapstructure:"WRITE_PORT" json:"WRITE_PORT"`
	WriteUser     string `mapstructure:"WRITE_USER" json:"WRITE_USER"`
	WritePassword string `mapstructure:"WRITE_PASSWORD" json:"WRITE_PASSWORD"`
	WriteDatabase string `mapstructure:"WRITE_DATABASE" json:"WRITE_DATABASE"`

	readConnection  *sqlx.DB
	writeConnection *sqlx.DB
}

// SQLMap SQLMap
type SQLMap []map[string]interface{}

// Scan sqlx 可以替 type 寫特定方法 ex. scan/value 讓他可以在查詢寫入資料到資料表時 轉換成對應的資料格式
func (m *SQLMap) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &m)
		return nil
	case string:
		json.Unmarshal([]byte(v), &m)
		return nil
	default:
		return fmt.Errorf("Unsupported type: %T", v)
	}
}

// SetConfig SetConfig
func (mysql *Mysql) SetConfig(env MyEnv) {
	mysql.ReadHost = env.MysqlReadHost
	mysql.ReadPort = env.MysqlReadPort
	mysql.ReadUser = env.MysqlReadUser
	mysql.ReadPassword = env.MysqlReadPassword
	mysql.ReadDatabase = env.MysqlReadDatabase

	mysql.WriteHost = env.MysqlWriteHost
	mysql.WritePort = env.MysqlWritePort
	mysql.WriteUser = env.MysqlWriteUser
	mysql.WritePassword = env.MysqlWritePassword
	mysql.WriteDatabase = env.MysqlWriteDatabase
}

// CreateReadConnection create connection
func (mysql *Mysql) CreateReadConnection() (cost time.Duration, err error) {
	connectionString := fmt.Sprintf("%s:%s@(%s:%d)/%s", mysql.ReadUser, mysql.ReadPassword, mysql.ReadHost, mysql.ReadPort, mysql.ReadDatabase)
	startTime := time.Now()
	mysql.readConnection, err = sqlx.Connect("mysql", connectionString)
	cost = time.Since(startTime)

	return
}

// CreateWriteConnection create connection
func (mysql *Mysql) CreateWriteConnection() (cost time.Duration, err error) {
	connectionString := fmt.Sprintf("%s:%s@(%s:%d)/%s", mysql.WriteUser, mysql.WritePassword, mysql.WriteHost, mysql.WritePort, mysql.WriteDatabase)
	startTime := time.Now()
	mysql.writeConnection, err = sqlx.Connect("mysql", connectionString)
	cost = time.Since(startTime)

	return
}

// CloseReadConnection close read connection
func (mysql *Mysql) CloseReadConnection() (cost time.Duration, err error) {
	startTime := time.Now()
	err = mysql.readConnection.Close()
	cost = time.Since(startTime)

	return
}

// CloseWriteConnection close write connection
func (mysql *Mysql) CloseWriteConnection() (cost time.Duration, err error) {
	startTime := time.Now()
	err = mysql.writeConnection.Close()
	cost = time.Since(startTime)

	return
}

func (mysql *Mysql) printSchema(schema string, values []interface{}) {
	if len(values) > 0 && len(values) <= 50 {
		for _, value := range values {
			str := ""
			switch reflect.TypeOf(value).Kind() {
			case reflect.String:
				str = fmt.Sprintf("'%v'", value)
			default:
				str = fmt.Sprintf("%v", value)
			}
			schema = strings.Replace(schema, "?", str, 1)
		}
	}
	fmt.Println(string(schema))
	// tl := NewTraceLog()
	// tl.Level = TraceLogLevelInfo
	// tl.Type = TraceLogTypeMYSQL
	// tl.TraceID = traceID
	// tl.Message = "LOG_MYSQL_SCHEMA"
	// tl.AddExtraInfo("schema", schema)
	// tl.WriteLog()
}

func (mysql *Mysql) InsertGames(game *gameTable) (rowsAffected int64, err error) {
	timeNow := time.Now().Unix()
	rows := make([]interface{}, 0)
	rows = append(rows, game.UserID)
	rows = append(rows, 0)
	rows = append(rows, timeNow)
	rows = append(rows, time.Now().Format("2006-01-02 15:04:05"))
	schema := `INSERT INTO game(user_id, state, created_at, updated_at) VALUES (?,?,?,?)`

	query, args, err := sqlx.In(schema, rows...)
	if err != nil {
		return 0, err
	}

	mysql.printSchema(query, args)
	res, err := mysql.writeConnection.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (mysql *Mysql) InsertOrders(orders []*orderTable) (rowsAffected int64, err error) {
	var rows []interface{}
	var valuesSymbo []string
	for _, order := range orders {
		rows = append(rows, order)
		valuesSymbo = append(valuesSymbo, "(?)")
	}

	excludeOnDuplicate := []string{
		"order_id",
		"result",
		"game_id",
	}

	fields := make([]string, 0)
	updateOnDuplicate := make([]string, 0)
	for _, field := range structs.Fields(&orderTable{}) {
		fields = append(fields, field.Tag("db"))
		if IsStrInSlice(field.Tag("db"), excludeOnDuplicate) {
			continue
		}
		updateOnDuplicate = append(updateOnDuplicate, field.Tag("db")+`=VALUES(`+field.Tag("db")+`)`)
	}

	// INSERT INTO `orders` ()
	schema := `INSERT INTO orders (` + strings.Join(fields, ", ") + `) VALUES ` +
		strings.Join(valuesSymbo, ", ") +
		` ON DUPLICATE KEY UPDATE ` + strings.Join(updateOnDuplicate, ", ")
	query, args, err := sqlx.In(schema, rows...)
	mysql.printSchema(query, args)
	if err != nil {
		return
	}

	res, err := mysql.writeConnection.Exec(query, args...)
	if err != nil {
		return
	}

	return res.RowsAffected()
}

func (mysql *Mysql) getOrders(gameID int64) ([]*orderTable, error) {
	orders := make([]*orderTable, 0)
	rows := make([]interface{}, 0)
	rows = append(rows, gameID)
	rows = append(rows, 0)
	rows = append(rows, 1)
	schema := `SELECT * FROM order WHERE game_id = ? AND state = ? ORDER BY created_at ASC, order_id ASC  LIMIT ?`
	err := mysql.readConnection.Get(orders, schema, rows)
	mysql.printSchema(schema, []interface{}{gameID})
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return orders, err
}

func (mysql *Mysql) getGamesByUserId(userId int64) (*gameTable, error) {
	game := &gameTable{}
	schema := "SELECT * FROM games WHERE user_id = ?"
	err := mysql.readConnection.Get(game, schema, userId)
	mysql.printSchema(schema, []interface{}{userId})
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return game, err
}

func (mysql *Mysql) updateStatus(game *gameTable) (rowsAffected int64, err error) {
	rows := make([]interface{}, 0)
	rows = append(rows, game.State)
	rows = append(rows, time.Now().Format("2006-01-02 15:04:05"))
	rows = append(rows, game.UserID)

	schema := `UPDATE game SET state = ?, updated_at =? WHERE game_id = ? `

	query, args, err := sqlx.In(schema, rows...)
	if err != nil {
		return 0, err
	}

	mysql.printSchema(query, args)

	res, err := mysql.writeConnection.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// GetMockDataForMockID GetMockDataForMockID
// func (mysql *Mysql) GetMockDataForMockID(mockID string) ([]mockdata, error) {
// 	query := fmt.Sprintf(`SELECT * FROM mock_data WHERE MockID = "%s"`, mockID)
// 	fmt.Println(query)
// 	rows, err := mysql.readConnection.Query(query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	datas := []mockdata{}
// 	for rows.Next() {
// 		var dt mockdata
// 		if err := rows.Scan(&dt.ID, &dt.MockID, &dt.Method, &dt.Route, &dt.Ret); err != nil {
// 			log.Fatal(err)
// 		}
// 		datas = append(datas, dt)
// 	}

// 	return datas, err
// }

// //GetMockDataList GetMockDataList
// func (mysql *Mysql) GetMockDataList() ([]mockdata, error) {
// 	rows, err := mysql.readConnection.Query("SELECT * FROM mock_data")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	datas := []mockdata{}
// 	for rows.Next() {
// 		var dt mockdata
// 		if err := rows.Scan(&dt.ID, &dt.MockID, &dt.Method, &dt.Route, &dt.Ret); err != nil {
// 			log.Fatal(err)
// 		}
// 		datas = append(datas, dt)
// 	}

// 	return datas, err
// }

// //getMockData getMockData
// func (mysql *Mysql) getMockData(mockID string, method string, route string) (string, error) {
// 	query := fmt.Sprintf(`SELECT Ret FROM mock_data WHERE MockID = "%s" AND Route = "%s" AND Method = "%s"`, mockID, route, method)
// 	fmt.Println(query)
// 	var ret string
// 	err := mysql.readConnection.QueryRow(query).Scan(&ret)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return ret, err
// }

func (mysql *Mysql) GetDefaultPricesByGameID(gameID int64, traceID string) (string, error) {
	prices := ""
	atTimestamp := time.Now().Unix()
	schema := "SELECT * FROM x WHERE game_id = ? AND created_at = ( SELECT created_at FROM x WHERE game_id = ? AND created_at <= ? ORDER BY created_at DESC LIMIT 1) ORDER BY level"
	err := mysql.readConnection.Select(&prices, schema, gameID, gameID, atTimestamp)
	mysql.printSchema(schema, []interface{}{gameID, gameID, atTimestamp})
	if err == sql.ErrNoRows {
		return prices, nil
	}

	return prices, err
}
