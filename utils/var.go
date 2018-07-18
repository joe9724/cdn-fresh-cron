package utils

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

func Response200(code int64,msg string) (string) {
	/*fmt.Println(code)
	fmt.Println(msg)
	responseStr := `{
    "code": 201,
    "msg": "ok"
    }`
	return responseStr*/
	var response Response
	response.Code = code
	response.Msg = msg

	out, _ := json.Marshal(response)
	fmt.Println("ous is",out)
	return(string(out))

}

type Response struct{
	Code int64 `json:"code,omitempty"`
	// cover
	Msg string `json:"msg,omitempty"`
}

func GetResourceDomain(filetype string) (string){
	var val string
	if(filetype == "icon") {
		val = "http://tingting-resource.bitekun.xin/resource/image/icon/"
	}else if(filetype == "cover"){
		val = "http://tingting-resource.bitekun.xin/resource/image/cover/"
	}else if(filetype == "m4a"){
		val = "http://tingting-resource.bitekun.xin/resource/mp3/"
	}else if(filetype == "amr"){
		val = "http://tingting-resource.bitekun.xin/resource/amr/"
	}else if(filetype == "other"){
		val = "http://tingting-resource.bitekun.xin/resource/other/"
	}
	return val
}

func OpenConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open("mysql", "root:root811123@tcp(106.14.2.153:3306)/yunfan?charset=utf8&parseTime=True&loc=Local")
	return db,err

	/*switch os.Getenv("GORM_DIALECT") {
	case "mysql":
		// CREATE USER 'gorm'@'localhost' IDENTIFIED BY 'gorm';
		// CREATE DATABASE gorm;
		// GRANT ALL ON gorm.* TO 'gorm'@'localhost';
		fmt.Println("testing mysql...")
		dbhost := os.Getenv("GORM_DBADDRESS")
		if dbhost != "" {
			dbhost = fmt.Sprintf("tcp(%v)", dbhost)
		}
		db, err = gorm.Open("mysql", fmt.Sprintf("gorm:gorm@%v/gorm?charset=utf8&parseTime=True", dbhost))
	case "postgres":
		fmt.Println("testing postgres...")
		dbhost := os.Getenv("GORM_DBHOST")
		if dbhost != "" {
			dbhost = fmt.Sprintf("host=%v ", dbhost)
		}
		db, err = gorm.Open("postgres", fmt.Sprintf("%vuser=gorm password=gorm DB.name=gorm sslmode=disable", dbhost))
	case "foundation":
		fmt.Println("testing foundation...")
		db, err = gorm.Open("foundation", "dbname=gorm port=15432 sslmode=disable")
	case "mssql":
		// CREATE LOGIN gorm WITH PASSWORD = 'LoremIpsum86';
		// CREATE DATABASE gorm;
		// USE gorm;
		// CREATE USER gorm FROM LOGIN gorm;
		// sp_changedbowner 'gorm';
		fmt.Println("testing mssql...")
		db, err = gorm.Open("mssql", "sqlserver://gorm:LoremIpsum86@localhost:1433?database=gorm")
	default:
		fmt.Println("testing sqlite3...")
		db, err = gorm.Open("sqlite3", filepath.Join(os.TempDir(), "gorm.db"))
	}
	// db.SetLogger(Logger{log.New(os.Stdout, "\r\n", 0)})
	// db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}
	db.DB().SetMaxIdleConns(10)
	return*/
}