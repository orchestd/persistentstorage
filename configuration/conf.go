package configuration

type PersistentStorageConfiguration struct {
	SqlDBName *string `json:"SQL_DB_NAME,omitempty"`
	SqlHost   *string `json:"SQL_HOST,omitempty"`
}
