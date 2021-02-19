package variables

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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
	Connector string

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
	HTTPOnly  bool
	MaxAge    int
	SameSite  http.SameSite
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
	JWTConfig.HTTPOnly = true
	JWTConfig.MaxAge = 24 * 3600

	switch Mode {
	case "development":
		MongoConfig.Connector = "mongodb"
		JWTConfig.Domain = "localhost"
		JWTConfig.Path = "/"
		JWTConfig.HTTPS = false
		JWTConfig.SameSite = http.SameSiteNoneMode

	case "staging-local":
		MongoConfig.Connector = "mongodb+srv"
		JWTConfig.Domain = "weeber.id" // for https://web-localhost.weeber.id:3000
		JWTConfig.Path = "/"
		JWTConfig.HTTPS = true
		JWTConfig.SameSite = http.SameSiteNoneMode

	case "staging":
		MongoConfig.Connector = "mongodb+srv"
		JWTConfig.Domain = "staging-tanjungbunga.weeber.id" // for https://staging-tanjungbunga.weeber.id
		JWTConfig.Path = "/"
		JWTConfig.HTTPS = true
		JWTConfig.SameSite = http.SameSiteNoneMode

	case "production":
		MongoConfig.Connector = "mongodb+srv"
		JWTConfig.Domain = "tanjungbunga.id"
		JWTConfig.Path = "/"
		JWTConfig.HTTPS = true
		JWTConfig.SameSite = http.SameSiteNoneMode

	default:
		log.Fatal(errors.New("Invalid MODE, must be: development, staging-local, staging, production"))
	}
}
