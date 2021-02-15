package variables

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

var (
	// Mode this service
	// ex: development, staging, production
	Mode string

	// Version this service
	Version string
)

// MongoConfig data type
var MongoConfig struct {
	Host     string
	Database string
	User     string
	Password string
}

// JWTConfig datatype
var JWTConfig struct {
	Key string

	TokenName string
	Path      string
	Domain    string
	HTTPS     bool
}

// InitializationVariable environment
func InitializationVariable() {
	Mode = os.Getenv("MODE")
	if Mode == "" {
		log.Fatal("Mode variable is null")
	}

	// reading version file
	ver, err := ioutil.ReadFile("./VERSION")
	if err != nil {
		log.Fatalf("read version file %v \n", err)
	}
	Version = string(ver)

	MongoConfig.Host = os.Getenv("MONGO_HOST")
	MongoConfig.User = os.Getenv("MONGO_USER")
	MongoConfig.Password = os.Getenv("MONGO_PASS")
	MongoConfig.Database = os.Getenv("MONGO_DATABASE")

	JWTConfig.Key = os.Getenv("JWT_SECRET_KEY")
	JWTConfig.TokenName = "auth_token"
	switch Mode {
	case "development":
		JWTConfig.Domain = "localhost:8080"
		JWTConfig.Path = "/api"
		JWTConfig.HTTPS = false
	default:
		log.Fatal(errors.New("Invalid MODE, must be: local, staging, production"))
	}
}
