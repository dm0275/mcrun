package minecraft

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type ServerConfig struct {
	MaxTickTime                 int    `prop:"max-tick-time"`
	GeneratorSettings           string `prop:"generator-settings"`
	AllowNether                 bool   `prop:"allow-nether"`
	ForceGamemode               bool   `prop:"force-gamemode"`
	Gamemode                    int    `prop:"gamemode"`
	EnableQuery                 bool   `prop:"enable-query"`
	PlayerIdleTimeout           int    `prop:"player-idle-timeout"`
	Difficulty                  int    `prop:"difficulty"`
	SpawnMonsters               bool   `prop:"spawn-monsters"`
	OpPermissionLevel           int    `prop:"op-permission-level"`
	Pvp                         bool   `prop:"pvp"`
	SnooperEnabled              bool   `prop:"snooper-enabled"`
	LevelType                   string `prop:"level-type"`
	Hardcore                    bool   `prop:"hardcore"`
	EnableCommandBlock          bool   `prop:"enable-command-block"`
	MaxPlayers                  int    `prop:"max-players"`
	NetworkCompressionThreshold int    `prop:"network-compression-threshold"`
	ResourcePackSha1            string `prop:"resource-pack-sha1"`
	MaxWorldSize                int    `prop:"max-world-size"`
	ServerPort                  int    `prop:"server-port"`
	ServerIp                    string `prop:"server-ip"`
	SpawnNpcs                   bool   `prop:"spawn-npcs"`
	AllowFlight                 bool   `prop:"allow-flight"`
	LevelName                   string `prop:"level-name"`
	ViewDistance                int    `prop:"view-distance"`
	ResourcePack                string `prop:"resource-pack"`
	SpawnAnimals                bool   `prop:"spawn-animals"`
	WhiteList                   bool   `prop:"white-list"`
	GenerateStructures          bool   `prop:"generate-structures"`
	OnlineMode                  bool   `prop:"online-mode"`
	MaxBuildHeight              int    `prop:"max-build-height"`
	LevelSeed                   string `prop:"level-seed"`
	PreventProxyConnections     bool   `prop:"prevent-proxy-connections"`
	UseNativeTransport          bool   `prop:"use-native-transport"`
	Motd                        string `prop:"motd"`
	EnableRcon                  bool   `prop:"enable-rcon"`
}

func (s ServerConfig) MarshalAndWriteToFile(filePath string) error {
	// Create the parameters file
	paramsFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer paramsFile.Close()

	// Get the value of the Parameters struct
	val := reflect.ValueOf(s)

	// Get the type of the struct
	typ := val.Type()

	// Iterate over the fields of the struct
	for i := 0; i < typ.NumField(); i++ {
		// Get the field tag and default value field
		tag := val.Type().Field(i).Tag.Get("prop")
		defaultVal := val.Type().Field(i).Tag.Get("default")

		// Use the fieldTag as the keyName
		key := tag

		// Get the field value
		fieldVal := val.Field(i)

		// Get the field value as a string
		var value string
		switch fieldVal.Kind() {
		case reflect.Bool:
			value = strconv.FormatBool(fieldVal.Bool())
		case reflect.String:
			value = fieldVal.String()
		case reflect.Int:
			value = strconv.FormatInt(fieldVal.Int(), 10)
		}

		// Attempt to set a default value
		if value == "" && defaultVal != "" {
			value = defaultVal
		}

		if value == "" {
			continue
		}

		// Write the key and value to the file, if the field value is not empty
		if _, err := fmt.Fprintf(paramsFile, "%s=%s\n", key, value); err != nil {
			return err
		}
	}

	return nil
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		MaxTickTime:                 60000,
		AllowNether:                 true,
		ForceGamemode:               false,
		Gamemode:                    0,
		EnableQuery:                 false,
		PlayerIdleTimeout:           0,
		Difficulty:                  1,
		SpawnMonsters:               true,
		OpPermissionLevel:           4,
		Pvp:                         true,
		SnooperEnabled:              true,
		LevelType:                   "DEFAULT",
		Hardcore:                    false,
		EnableCommandBlock:          false,
		MaxPlayers:                  20,
		NetworkCompressionThreshold: 256,
		MaxWorldSize:                29999984,
		ServerPort:                  25565,
		ServerIp:                    "",
		SpawnNpcs:                   true,
		AllowFlight:                 false,
		LevelName:                   "world",
		ViewDistance:                10,
		SpawnAnimals:                true,
		WhiteList:                   false,
		GenerateStructures:          true,
		OnlineMode:                  true,
		MaxBuildHeight:              256,
		LevelSeed:                   "",
		PreventProxyConnections:     false,
		UseNativeTransport:          true,
		Motd:                        "Docker Minecraft Server",
		EnableRcon:                  false,
	}
}
