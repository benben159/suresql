package suresql

import (
	"fmt"
	"os"
	"strconv"
	"time"

	orm "github.com/medatechnology/simpleorm"

	utils "github.com/medatechnology/goutil"
	"github.com/medatechnology/goutil/object"
)

const (
	// DEFAULT LEADER NODE
	LEADER_NODE_NUMBER = 1

	// for readibility
	INTERNAL_MODE = false
	NODE_MODE     = true

	// ConfigTable Categories and keys
	CONFIG_TOKEN_CATEGORY  = "token"
	CONFIG_TOKEN_EXP_KEY   = "token_exp"   // value int: in minutes
	CONFIG_REFRESH_EXP_KEY = "refresh_exp" // value int: in minutes
	CONFIG_TOKEN_TTL_KEY   = "token_ttl"   // value int: in minutes, beat for checking expiration

	CONFIG_CONNECTION_CATEGORY = "connection"
	CONFIG_MAX_POOL_KEY        = "max_pool" // value int: 0 overwrite pool_on, meaning no pooling, automatically pool_on=false
	CONFIG_ENABLE_POOL_KEY     = "pool_on"  // value string: true or false

	CONFIG_NODES_CATEGORY = "nodes"
	CONFIG_NODE_NAME_KEY  = "node_name" // value string: node_number;hostname;ip;mode
	CONFIG_NODE_DELIMITER = "|"

	CONFIG_EMPTY_CATEGORY = "nocategory"
)

// map ConfigTable by the key (string) which is same as ConfigTable.ConfigKey
// instead of using array, this is faster to search for specific config key
type ConfigMap map[string]ConfigTable

// Map ConfigMap by the category, which is the same as inside ConfigTable.Category
// Finding key based on category: Configs[category][Key] ie: Configs[token][token_exp].IntValue =
type ConfigCategory map[string]ConfigMap

// This is config needed by SureSQL to connect to Internal DB (DBMS), at this point only RQLite
type SureSQLConfig struct {
	Host         string        `json:"host,omitempty"            db:"host"`
	Port         string        `json:"port,omitempty"            db:"port"`
	Username     string        `json:"username,omitempty"        db:"username"`
	Password     string        `json:"password,omitempty"        db:"password"`
	Database     string        `json:"database,omitempty"        db:"database"`
	SSL          bool          `json:"ssl,omitempty"             db:"ssl"`
	Options      string        `json:"options,omitempty"         db:"options"`
	Consistency  string        `json:"consistency,omitempty"     db:"consistency"`
	URL          string        `json:"url,omitempty"             db:"url"`
	Token        string        `json:"token,omitempty"           db:"token"`
	RefreshToken string        `json:"refresh_token,omitempty"   db:"refresh_token"`
	JWEKey       string        `json:"jwe_key,omitempty"         db:"jwe_key"`
	APIKey       string        `json:"api_key,omitempty"         db:"api_key"`
	ClientID     string        `json:"client_id,omitempty"       db:"client_id"`
	HttpTimeout  time.Duration `json:"http_timeout,omitempty"    db:"http_timeout"`
	RetryTimeout time.Duration `json:"retry_timeout,omitempty"   db:"retry_timeout"`
	MaxRetries   int           `json:"max_retries,omitempty"     db:"max_retries"`
}

func (sc *SureSQLConfig) PrintDebug(secure bool) {
	fmt.Println("Loading from environment")
	fmt.Println("Host          : ", sc.Host)
	fmt.Println("Port          : ", sc.Port)
	fmt.Println("UserName      : ", sc.Username)
	fmt.Println("Database      : ", sc.Database)
	fmt.Println("SSL           : ", sc.SSL)
	fmt.Println("Options       : ", sc.Options)
	fmt.Println("URL           : ", sc.URL)
	fmt.Println("HTTP Timeout  : ", sc.URL)
	fmt.Println("Retry Timeout : ", sc.URL)
	fmt.Println("Max Retries   : ", sc.URL)
	if secure {
		fmt.Println("Password      : ", sc.Password)
		fmt.Println("Token         : ", sc.Token)
		fmt.Println("Refresh       : ", sc.RefreshToken)
		fmt.Println("JWEKey        : ", sc.JWEKey)
		fmt.Println("APIKey        : ", sc.APIKey)
		fmt.Println("ClientID      : ", sc.ClientID)
	}
}

// If using direct-rqlite (our own) implementation, then no need, because when direct-rqlite connects to
// RQLite server it will use basic-auth format for the username and password.
func (sc *SureSQLConfig) GenerateRQLiteURL() {
	tmpURL := "http://"
	if sc.SSL {
		tmpURL = "https://"
	}
	if len(sc.Host) > 0 {
		tmpURL += sc.Host
	} else {
		tmpURL += "localhost"
		fmt.Println("ERROR! No Host defined in environment")
	}
	if len(sc.Port) > 0 {
		tmpURL += ":" + sc.Port
	}
	sc.URL = tmpURL
}

// If using the gorqlite implementation, then we need to put username+password in the URL
// then gorqlite use this to connect to the rqlite server
func (sc *SureSQLConfig) GenerateGoRQLiteURL() {
	tmpURL := "http://"
	if sc.SSL {
		tmpURL = "https://"
	}
	if len(sc.Username) > 0 {
		tmpURL += sc.Username
	}
	if len(sc.Password) > 0 {
		tmpURL += ":" + sc.Password
	}
	if len(sc.Username) > 0 || len(sc.Password) > 0 {
		tmpURL += "@"
	}
	if len(sc.Host) > 0 {
		tmpURL += sc.Host
	} else {
		fmt.Println("ERROR! No Host defined in environment")
	}
	if len(sc.Port) > 0 {
		tmpURL += ":" + sc.Port
	}
	tmpURL += "/"
	if len(sc.Options) > 0 {
		tmpURL += "?" + sc.Options
	}
	sc.URL = tmpURL
}

// Reading internal DB configuration for this SureSQL Node, from environment
// TODO: maybe add second return parameter: error, so caller can check if error then quit the app
func LoadConfigFromEnvironment() SureSQLConfig {
	tmpConfig := SureSQLConfig{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASSWORD"),
		Database:     os.Getenv("DB_DATABASE"),
		SSL:          utils.GetEnvBool("DB_SSL", false),
		Options:      os.Getenv("DB_OPTIONS"),
		Consistency:  os.Getenv("DB_CONSISTENCY"),
		Token:        os.Getenv("DB_TOKEN"),
		RefreshToken: os.Getenv("DB_TOKEN_REFRESH"),
		JWEKey:       os.Getenv("DB_JWE_KEY"),
		APIKey:       os.Getenv("DB_API_KEY"),
		ClientID:     os.Getenv("DB_CLIENT_ID"),
		HttpTimeout:  utils.GetEnvDuration("DB_HTTP_TIMEOUT", DEFAULT_TIMEOUT),
		RetryTimeout: utils.GetEnvDuration("DB_RETRY_TIMEOUT", DEFAULT_RETRY_TIMEOUT),
		MaxRetries:   utils.GetEnvInt("DB_MAX_RETRIES", DEFAULT_RETRY),
	}
	return tmpConfig
}

// if DB settings is not there, get from environment. DB's settings table always wins
func RecheckSettingsFromEnvironment() {
	if CurrentNode.IP == "" {
		CurrentNode.IP = os.Getenv("SURESQL_IP")
	}
	if CurrentNode.Settings.Host == "" {
		CurrentNode.Settings.Host = os.Getenv("SURESQL_HOST")
	}
	if CurrentNode.Settings.Port == "" {
		CurrentNode.Settings.Port = os.Getenv("SURESQL_PORT")
	}
	// CurrentNode.Settings.SSL = os.Getenv("SURESQL_SSL")
	if CurrentNode.Settings.DBMS == "" {
		CurrentNode.Settings.DBMS = os.Getenv("SURESQL_DBMS")
	}
	if CurrentNode.InternalAPI == "" {
		CurrentNode.InternalAPI = os.Getenv("SURESQL_INTERNAL_API")
	}
	tmpBool, _ := strconv.ParseBool(os.Getenv("SURESQL_SSL"))
	CurrentNode.Settings.SSL = tmpBool

}

// LoadSettings loads settings from _settings table
func LoadSettings(db *SureSQLDB) error {
	record, err := (*db).SelectOne(CurrentNode.Settings.TableName())
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	// Get from database
	CurrentNode.Settings = object.MapToStructSlow[SettingsTable](record.Data)
	CurrentNode.IsEncrypted = CurrentNode.Settings.EncryptionMethod != "none"
	RecheckSettingsFromEnvironment()
	// TODO: Get the peers and leader
	return nil
}

func LoadConfigFromDB(db *SureSQLDB) error {
	records, err := (*db).SelectMany(ConfigTable{}.TableName())
	if err != nil {
		if err != orm.ErrSQLNoRows {
			return nil
		}
		return fmt.Errorf("failed to load configs from DB: %s", err)
	}
	for _, r := range records {
		tmp := object.MapToStruct[ConfigTable](r.Data)
		if tmp.Category == "" {
			tmp.Category = CONFIG_EMPTY_CATEGORY
		}
		tmpConfigMap, ok := CurrentNode.DBConfigs[tmp.Category]
		if !ok {
			tmpConfigMap = make(ConfigMap)
		}
		tmpConfigMap[tmp.ConfigKey] = tmp
		CurrentNode.DBConfigs[tmp.Category] = tmpConfigMap
	}
	// fmt.Println("DEBUG: reading configs table:", len(records), " rows")
	// fmt.Println("DEBUG: current node configs :", len(CurrentNode.DBConfigs), " category")
	return err
}

// By category and key
func (c ConfigCategory) ConfigExist(category, key string) (ConfigTable, bool) {
	if category == "" {
		category = CONFIG_EMPTY_CATEGORY
	}
	if tmp, ok := c[category]; ok {
		if conf, ok := tmp.ConfigExist(key); ok {
			return conf, true
		}
	}
	return ConfigTable{}, false
}

func (c ConfigMap) ConfigExist(key string) (ConfigTable, bool) {
	if conf, ok := c[key]; ok {
		return conf, true
	}
	return ConfigTable{}, false
}
