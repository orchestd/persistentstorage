package mysqlPersistentStorage

import (
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/configuration"
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/contextData"
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/credentials"
	. "bitbucket.org/HeilaSystems/persistentstorage"
	"bitbucket.org/HeilaSystems/persistentstorage/baseHeila"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

const cancelDateField = "cancelDate"
const updateStampField = "updateStamp"

func NewMySQLDb(updateStampGetter UpdateStampGetter, credentials credentials.CredentialsGetter,
	config configuration.Config, ctxResolver contextData.ContextDataResolver) PersistentStorage {
	//TODO: add trace
	conStr := credentials.GetCredentials().SqlConnectionString
	dbName, err := config.Get("SQL_DB_NAME").String()
	if err != nil {
		panic("env variable SQL_DB_NAME must be defined")
	}
	conStr = strings.Replace(conStr, "<dbname>", dbName, 1)
	if !strings.Contains(conStr, "parseTime") {
		if strings.Contains(conStr, "?") {
			conStr += "&parseTime=true"
		} else {
			conStr += "?parseTime=true"
		}
	}

	mySQLDb := &MySQLDb{}
	db, err := gorm.Open(mysql.Open(conStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mySQLDb.db = db
	mySQLDb.ctxResolver = ctxResolver
	baseHeila.UpdStampGetter = updateStampGetter
	return mySQLDb
}

type MySQLDb struct {
	db          *gorm.DB
	ctxResolver contextData.ContextDataResolver
}

func (repo MySQLDb) setTimeNow(c context.Context, db *gorm.DB) {
	db.NowFunc = func() time.Time {
		now, ok, err := repo.ctxResolver.GetDateNow(c)
		if err != nil {
			fmt.Println(err) // TODO: add log
		}
		if !ok {
			fmt.Println(err) // TODO: add log
		}
		return now
	}
}

func (repo MySQLDb) getDbWithContext(c context.Context, db *gorm.DB) *gorm.DB {
	db = repo.db.WithContext(c)
	repo.setTimeNow(c, db)
	return db
}

func (repo MySQLDb) QueryOne(c context.Context, target QueryGetter, params map[string]interface{}) error {
	query := target.GetQuery()
	return repo.getDbWithContext(c, repo.db).Raw(query, params).First(target).Error
}

func (repo MySQLDb) QueryMany(c context.Context, target QueryGetter, params map[string]interface{}) error {
	query := target.GetQuery()
	return repo.getDbWithContext(c, repo.db).Raw(query, params).Find(target).Error
}

func (repo MySQLDb) QueryInt(c context.Context, query QueryGetter, params map[string]interface{}) (int64, error) {
	var result map[string]interface{}
	err := repo.getDbWithContext(c, repo.db).Raw(query.GetQuery(), params).First(&result).Error
	if err != nil {
		return 0, err
	}

	var keys []string
	for k := range result {
		keys = append(keys, k)
	}

	if len(keys) != 1 {
		return 0, fmt.Errorf("Query must return single int value")
	}

	switch i := result[keys[0]].(type) {
	case int64:
		return i, nil
	default:
		return 0, fmt.Errorf("Query must return single int value")
	}
}

func (repo MySQLDb) QueryString(c context.Context, query QueryGetter, params map[string]interface{}) (string, error) {
	var result map[string]interface{}
	err := repo.getDbWithContext(c, repo.db).Raw(query.GetQuery(), params).Take(&result).Error
	if err != nil {
		return "", err
	}

	var keys []string
	for k := range result {
		keys = append(keys, k)
	}

	if len(keys) != 1 {
		return "", fmt.Errorf("Query must return single string value")
	}
	switch i := result[keys[0]].(type) {
	case string:
		return i, nil
	default:
		return "", fmt.Errorf("Query must return single string value")
	}
}

func (repo MySQLDb) GetOne(c context.Context, target interface{}, params interface{}) error {
	return repo.getDbWithContext(c, repo.db).Where(params).First(target).Error
}

func (repo MySQLDb) GetMany(c context.Context, target interface{}, params interface{}) error {
	return repo.getDbWithContext(c, repo.db).Where(params).Find(target).Error
}

func (repo MySQLDb) Insert(c context.Context, target interface{}) error {
	return repo.getDbWithContext(c, repo.db).Model(target).Create(target).Error
}

func (repo MySQLDb) Update(c context.Context, update interface{}, query interface{}) error {
	return repo.getDbWithContext(c, repo.db).Model(update).Where(query).Updates(update).Error
}

func (repo MySQLDb) Delete(c context.Context, model interface{}, params interface{}) error {
	return repo.getDbWithContext(c, repo.db).Where(params).Delete(model).Error
}

func (repo MySQLDb) Exec(c context.Context, queryGetter QueryGetter, params map[string]interface{}) error {
	query := queryGetter.GetQuery()
	return repo.getDbWithContext(c, repo.db).Exec(query, params).Error
}