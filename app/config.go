// Copyright (c) Luke Atherton 2015

package authenticator

import (
	"flag"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config interface {
	GetTopic() string
	GetExchangeAddress() string
	GetAmpqUsername() string
	GetAmpqPassword() string

	GetDbHosts() []string
	GetAuthDb() string
	GetDbUsername() string
	GetDbPassword() string

	GetPrivateKeyPath() string
	GetPublicKeyPath() string
}

type AppConfig struct {
	topic           string
	exchangeAddress string
	ampqUsername    string
	ampqPassword    string
	privateKeyPath  string
	publicKeyPath   string
	dbHosts         []string
	authDb          string
	dbUsername      string
	dbPassword      string
}

func BuildConfig() Config {
	var topic string
	var exchangeAddress string
	var ampqUsername string
	var ampqPassword string

	var dbHostString string
	var dbHosts []string
	var authDb string
	var dbUsername string
	var dbPassword string

	var privateKeyPath string
	var publicKeyPath string
	var configFile string

	flag.StringVar(&configFile, "config", "", "path to yaml config file")

	flag.StringVar(&exchangeAddress, "mq-address", "", "exchange address")
	flag.StringVar(&topic, "mq-topic", "", "exchange topic")
	flag.StringVar(&ampqUsername, "mq-username", "", "ampq username")
	flag.StringVar(&ampqPassword, "mq-password", "", "ampq password")

	flag.StringVar(&dbHostString, "db-hosts", "", "address list of db hosts")
	flag.StringVar(&authDb, "db-auth", "", "db to auth against")
	flag.StringVar(&dbUsername, "db-username", "", "db username")
	flag.StringVar(&dbPassword, "db-password", "", "db password")

	flag.StringVar(&privateKeyPath, "crypto-private-key", "", "path to private key")
	flag.StringVar(&publicKeyPath, "crypto-public-key", "", "path to public key")
	flag.Parse()

	dbHosts, _ = coerceStringSlice(dbHostString)

	var cfgFile map[string]interface{}
	if configFile != "" {
		_, err := toml.DecodeFile(configFile, &cfgFile)
		if err != nil {
			log.Fatalf("ERROR: failed to load config file %s - %s", configFile, err.Error())
		}

		topicSlice, _ := coerceStringSlice(cfgFile["mq-topic"])
		exchangeAddressSlice, _ := coerceStringSlice(cfgFile["mq-address"])
		ampqUsernameSlice, _ := coerceStringSlice(cfgFile["mq-username"])
		ampqPasswordSlice, _ := coerceStringSlice(cfgFile["mq-password"])

		dbHostsSlice, _ := coerceStringSlice(cfgFile["db-hosts"])
		authDbSlice, _ := coerceStringSlice(cfgFile["db-auth"])
		dbUsernameSlice, _ := coerceStringSlice(cfgFile["db-username"])
		dbPasswordSlice, _ := coerceStringSlice(cfgFile["db-password"])

		privateKeyPathSlice, _ := coerceStringSlice(cfgFile["crypto-private-key"])
		publicKeyPathSlice, _ := coerceStringSlice(cfgFile["crypto-public-key"])

		topic = topicSlice[0]
		exchangeAddress = exchangeAddressSlice[0]
		ampqUsername = ampqUsernameSlice[0]
		ampqPassword = ampqPasswordSlice[0]

		dbHosts = dbHostsSlice
		authDb = authDbSlice[0]
		dbUsername = dbUsernameSlice[0]
		dbPassword = dbPasswordSlice[0]

		privateKeyPath = privateKeyPathSlice[0]
		publicKeyPath = publicKeyPathSlice[0]
	}

	if len(topic) == 0 {
		log.Fatalf("--mq-topic required")
	}

	if len(exchangeAddress) == 0 {
		log.Fatalf("--mq-address required")
	}

	if len(ampqUsername) == 0 {
		log.Fatalf("--mq-username required")
	}

	if len(ampqPassword) == 0 {
		log.Fatalf("--mq-password required")
	}

	if len(dbHosts) == 0 {
		log.Fatalf("--db-hosts required")
	}

	// if len(authDb) != 0 {
	// 	// log.Fatalf("--db-auth required")

	// 	if len(dbUsername) == 0 {
	// 		log.Fatalf("--db-username required")
	// 	}

	// 	if len(dbPassword) == 0 {
	// 		log.Fatalf("--db-password required")
	// 	}
	// }

	if len(privateKeyPath) == 0 {
		log.Fatalf("--crypto-private-key required")
	}

	if len(publicKeyPath) == 0 {
		log.Fatalf("--crypto-public-key required")
	}

	log.Println("***** Configuration *****")
	log.Println()
	log.Println("Database")
	log.Println(" ├─ db-hosts -----------> ", dbHosts)
	log.Println(" ├─ db-auth ------------> ", authDb)
	log.Println(" ├─ db-username --------> ", dbUsername)
	log.Println(" └─ db-password --------> ", dbPassword)
	log.Println()
	log.Println("Message Queue")
	log.Println(" ├─ mq-topic -----------> ", topic)
	log.Println(" ├─ mq-address ---------> ", exchangeAddress)
	log.Println(" ├─ mq-username --------> ", ampqUsername)
	log.Println(" └─ mq-password --------> ", ampqPassword)
	log.Println()
	log.Println("Cryptography")
	log.Println(" ├─ crypto-private-key -> ", privateKeyPath)
	log.Println(" └─ crypto-public-key --> ", publicKeyPath)
	log.Println()
	log.Println("*************************")

	config := &AppConfig{topic: topic, exchangeAddress: exchangeAddress, ampqUsername: ampqUsername, ampqPassword: ampqPassword, dbHosts: dbHosts, authDb: authDb, dbUsername: dbUsername, dbPassword: dbPassword, privateKeyPath: privateKeyPath, publicKeyPath: publicKeyPath}

	return config
}

func (config *AppConfig) GetTopic() string {
	return config.topic
}

func (config *AppConfig) GetExchangeAddress() string {
	return config.exchangeAddress
}

func (config *AppConfig) GetAmpqUsername() string {
	return config.ampqUsername
}

func (config *AppConfig) GetAmpqPassword() string {
	return config.ampqPassword
}

func (config *AppConfig) GetDbHosts() []string {
	return config.dbHosts
}

func (config *AppConfig) GetAuthDb() string {
	return config.authDb
}

func (config *AppConfig) GetDbUsername() string {
	return config.dbUsername
}

func (config *AppConfig) GetDbPassword() string {
	return config.dbPassword
}

func (config *AppConfig) GetPrivateKeyPath() string {
	return config.privateKeyPath
}

func (config *AppConfig) GetPublicKeyPath() string {
	return config.publicKeyPath
}

func coerceStringSlice(v interface{}) ([]string, error) {
	var tmp []string
	switch v.(type) {
	case string:
		for _, s := range strings.Split(v.(string), ",") {
			tmp = append(tmp, s)
		}
	case []interface{}:
		for _, si := range v.([]interface{}) {
			tmp = append(tmp, si.(string))
		}
	case []string:
		tmp = v.([]string)
	}
	return tmp, nil
}
