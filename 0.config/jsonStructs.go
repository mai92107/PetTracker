package jsonModal

type Config struct {
	Port                      string `json:"port"`
	JwtSecretKey              string `json:"jwt_secret_key"`
	CryptoSecretKey           string `json:"crypto_secret_key"`
	DevicePrefix              string `json:"device_prefix"`
	DeviceSequence            string `json:"device_sequence"`
	DefaultSecretKey          string `json:"default_secret_key"`
	TrackingDataSurvivingDays int    `json:"tracking_data_surviving_days"`
}
type Machine struct {
	MariaDB         MariaDbConfig   `json:"mariaDb"`
	MongoDB         MongoDbConfig   `json:"mongoDb"`
	Redis           RedisDbConfig   `json:"redisDb"`
	MosquittoBroker MosquittoConfig `json:"mosquittoBroker"`
}

type MariaDbConfig struct {
	InUse   bool           `json:"in_use"`
	Reading MariaDbSetting `json:"reading"`
	Writing MariaDbSetting `json:"writing"`
}
type MongoDbConfig struct {
	InUse   bool           `json:"in_use"`
	Reading MongoDbSetting `json:"reading"`
	Writing MongoDbSetting `json:"writing"`
}
type MariaDbSetting struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

type MongoDbSetting struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Name         string `json:"name"`
	TimeoutRange int    `json:"timeout_range"`
}
type RedisDbConfig struct {
	InUse   bool           `json:"in_use"`
	Reading RedisDbSetting `json:"reading"`
	Writing RedisDbSetting `json:"writing"`
}

type RedisDbSetting struct {
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type MosquittoConfig struct {
	InUser          bool     `json:"in_use"`
	BrokerHostCloud string   `json:"broker_host_cloud"`
	BrokerHostLocal string   `json:"broker_host_local"`
	BrokerPort      string   `json:"broker_port"`
	Username        string   `json:"username"`
	Password        string   `json:"password"`
	VagueTopic      []string `json:"vague_topic"`
	ClientID        string   `json:"client_id"`
}
