package configuration

type PersistentStorageConfiguration struct {
	SqlDBName         *string `json:"SQL_DB_NAME,omitempty"`
	SqlHost           *string `json:"SQL_HOST,omitempty"`
	SqlReadOnlyDBName *string `json:"SQL_READ_ONLY_DB_NAME,omitempty"`
	SqlReadOnlyHost   *string `json:"SQL_READ_ONLY_HOST,omitempty"`
}
