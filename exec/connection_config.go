package exec

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type execConfig struct {
	WorkingDir  *string  `cty:"working_dir"`
	Interpreter []string `cty:"interpreter"`

	Protocol    *string `cty:"protocol"`
	User        *string `cty:"user"`
	Password    *string `cty:"password"`
	PrivateKey  *string `cty:"private_key"`
	Certificate *string `cty:"certificate"`
	Host        *string `cty:"host"`
	HostKey     *string `cty:"host_key"`
	Port        *int    `cty:"port"`
	Https       *bool   `cty:"https"`
	Insecure    *bool   `cty:"insecure"`
}

/*
	Type string

	User           string
	Password       string
	PrivateKey     string
	Certificate    string
	Host           string
	HostKey        string
	Port           uint16
	Agent          bool
	ScriptPath     string
	TargetPlatform string
	Timeout        string
	TimeoutVal     time.Duration

	ProxyScheme       string
	ProxyHost         string
	ProxyPort         uint16
	ProxyUserName     string
	ProxyUserPassword string

	BastionUser        string
	BastionPassword    string
	BastionPrivateKey  string
	BastionCertificate string
	BastionHost        string
	BastionHostKey     string
	BastionPort        uint16

	AgentIdentity string

	HTTPS    bool
	Insecure bool
	NTLM     bool   `mapstructure:"use_ntlm"`
	CACert   string `mapstructure:"cacert"`
*/

var ConfigSchema = map[string]*schema.Attribute{
	"working_dir": {
		Type: schema.TypeString,
	},
	"interpreter": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},

	"protocol": {
		Type: schema.TypeString,
	},
	"host": {
		Type: schema.TypeString,
	},
	"host_key": {
		Type: schema.TypeString,
	},
	"user": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"private_key": {
		Type: schema.TypeString,
	},
	"certificate": {
		Type: schema.TypeString,
	},
	"port": {
		Type: schema.TypeInt,
	},
	"https": {
		Type: schema.TypeBool,
	},
	"insecure": {
		Type: schema.TypeBool,
	},
}

func ConfigInstance() interface{} {
	return &execConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) execConfig {
	if connection == nil || connection.Config == nil {
		return execConfig{}
	}
	config, _ := connection.Config.(execConfig)
	return config
}
