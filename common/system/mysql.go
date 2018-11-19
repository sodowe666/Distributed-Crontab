package system

import (
	"errors"
	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
	"reflect"
	"sync"
	"time"
)

type DBStore struct {
	databases  interface{}
	dbMapStore map[string]*xorm.EngineGroup
	lock       sync.RWMutex
}

//可能会有多个数据源
type database struct {
	DB *xorm.EngineGroup `db:"default"`
}

var dbStore *DBStore

func getDB() *DBStore {
	if dbStore == nil {
		panic("you must init DB first. call the InitDB function at main package.")
	}
	return dbStore
}

func GetDbByName(confName string) *xorm.EngineGroup {
	if dbStore == nil {
		panic("you must init DB first. call the InitDB function at main package.")
	}
	if dbStore.dbMapStore[confName] == nil {
		panic("DB " + confName + "did not open")
	}
	return getDB().dbMapStore[confName]
}

func (dbStore *DBStore) setDatabases(databases interface{}) {
	dbStore.databases = databases
}

//初始化db
//参数暂定interface
func (dbStore *DBStore) InitDBWithConf(conf interface{}) error {
	if allDBs := viper.Get("db"); allDBs == nil {
		return errors.New("must config db")
	}
	if dbStore.databases == nil {
		return errors.New("databases element must be set")
	}
	//反射机制获取多个db属性
	refDb := reflect.Indirect(reflect.ValueOf(dbStore.databases))
	for i := 0; i < refDb.NumField(); i++ {
		//for里面lock会怎么？
		field := refDb.Type().Field(i)
		value, ok := field.Tag.Lookup("db")
		if !(ok && viper.Get("db."+value) != nil) {
			continue
		}
		//拼接tag获取对应配置
		dbConf := viper.Get("db." + value)

		if dbConf == nil {
			return errors.New("Database struct element:" + field.Name + " must have " + value + " config")
		}
		//duck typing 然后打开对应的数据库服务
		engineGroup, err := dbStore.openDB(dbConf.(map[string]interface{}))
		if err != nil {
			return errors.New("db init fail." + err.Error())
		}
		//连接成功并给有db tag的元素设置值
		refDb.Field(i).Set(reflect.ValueOf(engineGroup))
		//把engineGroup放入map容器里，key为对应的conf里的db_name
		dbStore.dbMapStore[value] = engineGroup
	}
	return nil
}

func (dbstore *DBStore) openDB(conf map[string]interface{}) (*xorm.EngineGroup, error) {
	//初始化主库
	dbType := conf["type"].(string)
	masterDB, err := dbStore.openMaster(dbType, conf)
	if err != nil {
		return nil, err
	}
	//初始化从库,从库可能有多个，若用HAProxy做了负载均衡，只有一个slave入口。也要以数组的方式在配置里加一个slave
	slaveDBArray := make([]*xorm.Engine, 0)

	if slave := conf["slave"]; slave != nil && len(slave.([]interface{})) > 0 {
		slaveConfArray, ok := slave.([]map[string]interface{})
		if !ok {
			return nil, errors.New("slave`s config must be an Array")
		}
		for _, slaveConf := range slaveConfArray {
			slaveDB, err := dbStore.openSlave(dbType, slaveConf)
			if err != nil {
				return nil, err
			}
			slaveDBArray = append(slaveDBArray, slaveDB)
		}
	}
	engineGroup, err := xorm.NewEngineGroup(masterDB, slaveDBArray, xorm.LeastConnPolicy())
	if err != nil {
		return nil, err
	}
	if err := engineGroup.Ping(); err != nil {
		return nil, err
	}
	return engineGroup, nil
}

//打开主库
func (dbstore *DBStore) openMaster(dbType string, conf map[string]interface{}) (*xorm.Engine, error) {
	user := conf["user"].(string)
	password := conf["password"].(string)
	host := conf["host"].(string)
	port := conf["port"].(string)
	database := conf["database"].(string)
	charset := conf["charset"].(string)
	masterSource := user + ":" + password + "@" + "tcp(" + host + ":" + port + ")" + "/" + database + "?charset" + charset
	masterDB, err := xorm.NewEngine(dbType, masterSource)
	if err != nil {
		return nil, err
	}
	maxIdleConnsNum := conf["maxIdleConns"].(int)
	if maxIdleConnsNum == 0 {
		maxIdleConnsNum = 100
	}
	maxOpenConnsNum := conf["maxOpenConns"].(int)
	if maxOpenConnsNum == 0 {
		maxOpenConnsNum = 100
	}
	idleTime := conf["idleTime"].(int)
	if idleTime == 0 {
		idleTime = 300
	}
	masterDB.SetMaxIdleConns(maxIdleConnsNum)
	masterDB.SetMaxOpenConns(maxOpenConnsNum)
	masterDB.SetConnMaxLifetime(time.Duration(idleTime))
	return masterDB, nil
}

//打开从库
func (dbstore *DBStore) openSlave(dbType string, conf map[string]interface{}) (*xorm.Engine, error) {
	user := conf["user"].(string)
	password := conf["password"].(string)
	host := conf["host"].(string)
	port := conf["port"].(string)
	database := conf["database"].(string)
	charset := conf["charset"].(string)
	slaveSource := user + ":" + password + "@" + "tcp(" + host + ":" + port + ")" + "/" + database + "?charset" + charset
	slaveDB, err := xorm.NewEngine(dbType, slaveSource)
	if err != nil {
		return nil, err
	}
	maxIdleConnsNum := conf["maxIdleConns"].(int)
	if maxIdleConnsNum == 0 {
		maxIdleConnsNum = 100
	}
	maxOpenConnsNum := conf["maxOpenConns"].(int)
	if maxOpenConnsNum == 0 {
		maxOpenConnsNum = 100
	}
	idleTime := conf["idleTime"].(int)
	if idleTime == 0 {
		idleTime = 300
	}
	slaveDB.SetMaxIdleConns(maxIdleConnsNum)
	slaveDB.SetMaxOpenConns(maxOpenConnsNum)
	slaveDB.SetConnMaxLifetime(time.Duration(idleTime))
	return slaveDB, nil
}

//关闭所有的数据库连接
func (db *DBStore) Close() {
	if db == nil {
		return
	}
	for _, val := range db.dbMapStore {
		val.Close()
	}
}
