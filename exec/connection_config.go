package exec

import (
	"errors"

	communicator "github.com/turbot/go-exec-communicator"
	"github.com/turbot/go-exec-communicator/shared"
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

	BastionUser       *string `cty:"bastion_user"`
	BastionPassword   *string `cty:"bastion_password"`
	BastionPrivateKey *string `cty:"bastion_private_key"`
	BastionHost       *string `cty:"bastion_host"`
	BastionHostKey    *string `cty:"bastion_host_key"`
	BastionPort       *int    `cty:"bastion_port"`

	ProxyHost         *string `cty:"proxy_host"`
	ProxyPort         *int    `cty:"proxy_port"`
	ProxyUserName     *string `cty:"proxy_user_name"`
	ProxyUserPassword *string `cty:"proxy_user_password"`
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
	"bastion_user": {
		Type: schema.TypeString,
	},
	"bastion_password": {
		Type: schema.TypeString,
	},
	"bastion_private_key": {
		Type: schema.TypeString,
	},
	"bastion_host": {
		Type: schema.TypeString,
	},
	"bastion_host_key": {
		Type: schema.TypeString,
	},
	"bastion_port": {
		Type: schema.TypeInt,
	},
	"proxy_host": {
		Type: schema.TypeString,
	},
	"proxy_port": {
		Type: schema.TypeInt,
	},
	"proxy_user_name": {
		Type: schema.TypeString,
	},
	"proxy_user_password": {
		Type: schema.TypeString,
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

// GetCommunicator :: creates a communicator from config
func GetCommunicator(connection *plugin.Connection) (communicator.Communicator, error) {
	conf := GetConfig(connection)

	config := shared.ConnectionInfo{
		Timeout: "10s",
		Port:    22,
	}

	// Bastion settings
	if conf.BastionUser != nil {
		config.BastionUser = *conf.BastionUser
	}
	if conf.BastionPassword != nil {
		config.BastionPassword = *conf.BastionPassword
	}
	if conf.BastionPrivateKey != nil {
		config.BastionPrivateKey = *conf.BastionPrivateKey
	}
	if conf.BastionHost != nil {
		config.BastionHost = *conf.BastionHost
	}
	if conf.BastionHostKey != nil {
		config.BastionHostKey = *conf.BastionHostKey
	}
	if conf.BastionPort != nil {
		config.BastionPort = uint16(*conf.BastionPort)
	}

	// PROXY settings
	if conf.ProxyHost != nil {
		config.ProxyHost = *conf.ProxyHost
	}
	if conf.ProxyPort != nil {
		config.ProxyPort = uint16(*conf.ProxyPort)
	}
	if conf.ProxyUserName != nil {
		config.ProxyUserName = *conf.ProxyUserName
	}
	if conf.ProxyUserPassword != nil {
		config.ProxyUserPassword = *conf.ProxyUserPassword
	} else {
		if conf.ProxyUserName != nil {
			return nil, errors.New("password is required when proxy username is set")
		}
	}

	if conf.Protocol != nil {
		config.Type = *conf.Protocol
	} else {
		return nil, errors.New("protocol is required in config")
	}
	if conf.Host != nil {
		config.Host = *conf.Host
	} else {
		return nil, errors.New("host is required in config")
	}
	if conf.Port != nil {
		config.Port = uint16(*conf.Port)
	}
	if conf.Https != nil {
		config.HTTPS = *conf.Https
	}
	if conf.Insecure != nil {
		config.Insecure = *conf.Insecure
	}
	if config.Type == "ssh" {
		if conf.User != nil {
			config.User = *conf.User
		} else {
			return nil, errors.New("user is required for SSH connections")
		}
		if conf.Password != nil {
			config.Password = *conf.Password
		} else if conf.PrivateKey != nil {
			config.PrivateKey = *conf.PrivateKey
		} else {
			return nil, errors.New("password or private_key is required for SSH connections")
		}
	}
	if config.Type == "winrm" {
		if conf.User != nil {
			config.User = *conf.User
		} else {
			return nil, errors.New("user is required for WinRM connections")
		}
		if conf.Password != nil {
			config.Password = *conf.Password
		} else {
			return nil, errors.New("password is required for WinRM connections")
		}
	}

	return communicator.New(config)
}
