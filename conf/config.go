package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var Config *Conf

type Conf struct {
	System        *System                 `yaml:"system"`
	Oss           *Oss                    `yaml:"oss"`
	MySql         map[string]*MySql       `yaml:"mysql"`
	Email         *Email                  `yaml:"email"`
	Redis         *Redis                  `yaml:"redis"`
	EncryptSecret *EncryptSecret          `yaml:"encryptSecret"`
	Cache         *Cache                  `yaml:"cache"`
	KafKa         map[string]*KafkaConfig `yaml:"kafKa"`
	RabbitMq      *RabbitMq               `yaml:"rabbitMq"`
	Es            *Es                     `yaml:"es"`
	PhotoPath     *LocalPhotoPath         `yaml:"photoPath"`
	LogConfig     *LogConfig              `yaml:"logConfig"`
}

type RabbitMq struct {
	RabbitMQ         string `yaml:"rabbitMq"`
	RabbitMQUser     string `yaml:"rabbitMqUser"`
	RabbitMQPassWord string `yaml:"rabbitMqPassWord"`
	RabbitMQHost     string `yaml:"rabbitMqHost"`
	RabbitMQPort     string `yaml:"rabbitMqPort"`
}

type KafkaConfig struct {
	DisableConsumer bool   `yaml:"disableConsumer"`
	Debug           bool   `yaml:"debug"`
	Address         string `yaml:"address"`
	RequiredAck     int    `yaml:"requiredAck"`
	ReadTimeout     int64  `yaml:"readTimeout"`
	WriteTimeout    int64  `yaml:"writeTimeout"`
	MaxOpenRequests int    `yaml:"maxOpenRequests"`
	Partition       int    `yaml:"partition"`
}

type System struct {
	AppEnv      string `yaml:"appEnv"`
	Domain      string `yaml:"domain"`
	Version     string `yaml:"version"`
	HttpPort    string `yaml:"httpPort"`
	Host        string `yaml:"host"`
	UploadModel string `yaml:"uploadModel"`
	Mode        string `yaml:"mode"`
}

type Oss struct {
	BucketName      string `yaml:"bucketName"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Endpoint        string `yaml:"endPoint"`
	EndpointOut     string `yaml:"endpointOut"`
	QiNiuServer     string `yaml:"qiNiuServer"`
}

type MySql struct {
	Dialect  string `yaml:"dialect"`
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Email struct {
	ValidEmail string `yaml:"validEmail"`
	SmtpHost   string `yaml:"smtpHost"`
	SmtpEmail  string `yaml:"smtpEmail"`
	SmtpPass   string `yaml:"smtpPass"`
}

type Redis struct {
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisUsername string `yaml:"redisUsername"`
	RedisPassword string `yaml:"redisPwd"`
	RedisDbName   int    `yaml:"redisDbName"`
	RedisNetwork  string `yaml:"redisNetwork"`
}

// EncryptSecret 加密的东西
type EncryptSecret struct {
	JwtSecret   string `yaml:"jwtSecret"`
	EmailSecret string `yaml:"emailSecret"`
	PhoneSecret string `yaml:"phoneSecret"`
	MoneySecret string `yaml:"moneySecret"`
}

type LocalPhotoPath struct {
	PhotoHost   string `yaml:"photoHost"`
	ProductPath string `yaml:"productPath"`
	AvatarPath  string `yaml:"avatarPath"`
}

type Cache struct {
	CacheType    string `yaml:"cacheType"`
	CacheExpires int64  `yaml:"cacheExpires"`
	CacheWarmUp  bool   `yaml:"cacheWarmUp"`
	CacheServer  string `yaml:"cacheServer"`
}

type Es struct {
	EsHost  string `yaml:"esHost"`
	EsPort  string `yaml:"esPort"`
	EsIndex string `yaml:"esIndex"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}

func InitConfig() {
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("加载配置文件出错了：", err)
		return
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		fmt.Println("反序列化初始文件出错了：", err)
		return
	}
}

func GetExpiresTime() int64 {
	if Config.Cache.CacheExpires == 0 {
		return int64(30 * time.Minute) // 默认 30min
	}

	if Config.Cache.CacheExpires == -1 {
		return -1 // Redis.KeepTTL = -1
	}

	return int64(time.Duration(Config.Cache.CacheExpires) * time.Minute)
}
