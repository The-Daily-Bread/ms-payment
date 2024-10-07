package utils

type Configuration struct {
	Database DatabaseSetting
	Server   ServerSettings
	App      Application
}

type DatabaseSetting struct {
	Url    string
	DbName string
}

type ServerSettings struct {
	Port string
}

type Application struct {
	Name string
}
