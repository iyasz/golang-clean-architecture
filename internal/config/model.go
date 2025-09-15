package config

type Config struct {
	App
	Server
	Database 
	// Jwt
	Logrus
}

type App struct {
	Name string
}

type Server struct {
	Host    string
	Port    string
	Prefork string
}
type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Timezone string
	Pool
}

type Pool struct {
	Idle     int
	Max      int
	Lifetime int
}

type Jwt struct {
	RefreshKey string
	AccessKey  string
}

type Logrus struct {
	Level int32
}
