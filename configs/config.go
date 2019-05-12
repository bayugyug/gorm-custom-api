package configs

import (
	"encoding/json"
	"flag"
	"log"
)

const (
	//status
	usageConfig = "use to set the config file parameter with HTTP-port"
)

var (
	//Settings of the app
	Settings *APISettings
)

// ParameterConfig optional parameter structure
type ParameterConfig struct {
	Port    string `json:"port"`
	DSN     string `json:"dsn"`
	Verbose bool   `json:"showlog"`
}

// APISettings is a config mapping
type APISettings struct {
	Config    *ParameterConfig
	CmdParams string
}

// Setup options settings
type Setup func(*APISettings)

// WithSetupConfig for cfg
func WithSetupConfig(r *ParameterConfig) Setup {
	return func(args *APISettings) {
		args.Config = r
	}
}

// WithSetupCmdParams for the json params
func WithSetupCmdParams(r string) Setup {
	return func(args *APISettings) {
		args.CmdParams = r
	}
}

// NewAppSettings main entry for config
func NewAppSettings(setters ...Setup) *APISettings {
	//set default
	cfg := &APISettings{}

	//chk the passed params
	for _, setter := range setters {
		setter(cfg)
	}
	//start
	cfg.Initializer()
	return cfg
}

//InitRecov is for dumpIng segv in
func (g *APISettings) InitRecov() {
	//might help u
	defer func() {
		recvr := recover()
		if recvr != nil {
			log.Println("MAIN-RECOV-INIT: ", recvr)
		}
	}()
}

//InitEnvParams enable all OS envt vars to reload internally
func (g *APISettings) InitEnvParams() {
	//get options
	flag.StringVar(&g.CmdParams, "config", g.CmdParams, usageConfig)
	flag.Parse()
}

//Initializer set defaults for initial reqmts
func (g *APISettings) Initializer() {
	//prepare
	g.InitRecov()
	g.InitEnvParams()

	//set default maybe
	if g.CmdParams == "" {
		g.CmdParams = `{"port":"8989"}`
	}

	//try to reconfigure if there is passed params, otherwise use show err
	if g.CmdParams != "" {
		g.Config = g.FormatParameterConfig(g.CmdParams)
	}

	//check defaults
	if g.Config == nil {
		return
	}
}

//FormatParameterConfig new ParameterConfig
func (g *APISettings) FormatParameterConfig(s string) *ParameterConfig {
	var cfg ParameterConfig
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		log.Println("FormatParameterConfig", err)
		return nil
	}
	return &cfg
}
