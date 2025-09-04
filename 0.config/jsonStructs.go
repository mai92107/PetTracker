package jsonModal

type Config struct{
    Port            string  `json:"port"`
    JsonSecretKey   string  `json:"json_secret_key"`
    CryptoSecretKey string  `json:"crypto_secret_key"`
    DevicePrefix    string  `json:"device_prefix"`
    DeviceSequence  string  `json:"device_sequence"`
}
type Machine struct {
    MariaDB         MariaDbConfig   `json:"mariaDb"`
    Redis           RedisDbConfig   `json:"redisDb"`
    MosquittoBroker MosquittoConfig `json:"mosquittoBroker"`
}

type MariaDbConfig struct {
    InUse       bool   `json:"in_use"`
    Reading     MariaDbSetting   `json:"reading"`
    Writing     MariaDbSetting   `json:"writing"`
}
type MariaDbSetting struct{
    User        string `json:"user"`
    Password    string `json:"password"`
    Host        string `json:"host"`
    Port        string `json:"port"`
    Name        string `json:"name"`
}
type RedisDbConfig struct {
    InUse       bool   `json:"in_use"`
    Reading     RedisDbSetting   `json:"reading"`
    Writing     RedisDbSetting   `json:"writing"`
}

type RedisDbSetting struct{
    Password    string `json:"password"`
    Host        string `json:"host"`
    Port        string `json:"port"`
}

type MosquittoConfig struct{
    InUser      bool    `json:"in_use"`
    BrokerHost  string  `json:"broker_host"`
    BrokerPort  string  `json:"broker_port"`
    VagueTopic  string  `json:"vague_topic"`
    ClientID    string  `json:"client_id"`
}
