package mysqlPersistentStorage

import (
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/configuration"
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/credentials"
	. "bitbucket.org/HeilaSystems/persistentstorage"
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
	config configuration.Config) PersistentStorage {
	conStr := credentials.GetCredentials().SqlConnectionString
	dbName, err := config.Get("SQL_DB_NAME").String()
	if err != nil {
		panic("env variable SQL_DB_NAME must be defined")
	}
	conStr = strings.Replace(conStr, "<dbname>", dbName, 1)
	mySQLDb := &MySQLDb{updateStampGetter: updateStampGetter}
	db, err := gorm.Open(mysql.Open(conStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mySQLDb.db = db
	return mySQLDb
}

type MySQLDb struct {
	db                *gorm.DB
	updateStampGetter UpdateStampGetter
}

func (repo MySQLDb) QueryOne(c context.Context, target QueryGetter, params map[string]interface{}) error {
	query := target.GetQuery()
	return repo.db.WithContext(c).Raw(query, params).First(target).Error
}

func (repo MySQLDb) QueryMany(c context.Context, target QueryGetter, params map[string]interface{}) error {
	query := target.GetQuery()
	return repo.db.WithContext(c).Raw(query, params).Find(target).Error
}

func (repo MySQLDb) QueryInt(c context.Context, query QueryGetter, params map[string]interface{}) (int64, error) {
	var result map[string]interface{}
	err := repo.db.WithContext(c).Raw(query.GetQuery(), params).First(&result).Error
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
	err := repo.db.WithContext(c).Raw(query.GetQuery(), params).Take(&result).Error
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

func (repo MySQLDb) GetOne(c context.Context, target interface{}, params map[string]interface{}) error {
	return repo.db.WithContext(c).Where(cancelDateField + " IS NOT NULL").Where(params).First(target).Error
}

func (repo MySQLDb) GetMany(c context.Context, target interface{}, params map[string]interface{}) error {
	return repo.db.WithContext(c).Where(cancelDateField + " IS NOT NULL").Where(params).Find(target).Error
}

func (repo MySQLDb) Insert(c context.Context, target BaseModelSetter, now time.Time) error {
	updateStamp, err := repo.updateStampGetter.GetUpdateStamp(c, "update from persistent storage lib")
	if err != nil{
		return err
	}
	target.SetUpdateStamp(updateStamp)
	target.SetCreateDate(now)
	return repo.db.WithContext(c).Create(target).Error
}

func (repo MySQLDb) Update(c context.Context, model interface{}, update map[string]interface{}, params map[string]interface{}) error {
	updateStamp, err := repo.updateStampGetter.GetUpdateStamp(c, "update from persistent storage lib")
	if err != nil{
		return err
	}
	update[updateStampField] = updateStamp
	return repo.db.WithContext(c).Model(model).Where(params).Updates(update).Error
}

func (repo MySQLDb) Delete(c context.Context, model interface{}, params map[string]interface{}, now time.Time) error {
	update := make(map[string]interface{})
	updateStamp, err := repo.updateStampGetter.GetUpdateStamp(c, "update from persistent storage lib")
	if err != nil{
		return err
	}
	update[updateStampField] = updateStamp
	update[cancelDateField] = now
	return repo.db.WithContext(c).Model(model).Where(params).Updates(update).Error
}

func (repo MySQLDb) Exec(c context.Context, queryGetter QueryGetter, params map[string]interface{}) error {
	query := queryGetter.GetQuery()
	return repo.db.WithContext(c).Exec(query, params).Error
}
