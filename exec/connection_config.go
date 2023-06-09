package exec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
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
func GetCommunicator(connection *plugin.Connection) (communicator.Communicator, *exec.Cmd, bool, error) {
	conf := GetConfig(connection)

	config := shared.ConnectionInfo{
		Timeout: "10s",
	}

	// If no other connection info is provided, assume local connection
	localConnection := true

	// Bastion settings
	if conf.BastionHost != nil {
		config.BastionHost = *conf.BastionHost
		localConnection = false

		if conf.BastionUser != nil {
			config.BastionUser = *conf.BastionUser
		} else {
			return nil, nil, localConnection, errors.New("bastion_user is required when using bastion host")
		}
		if conf.BastionPassword != nil {
			config.BastionPassword = *conf.BastionPassword
		} else if conf.BastionPrivateKey != nil {
			content, err := PathOrContents(*conf.BastionPrivateKey)
			if err != nil {
				return nil, nil, localConnection, err
			}
			config.BastionPrivateKey = content
		} else {
			return nil, nil, localConnection, errors.New("either bastion_password or bastion_private_key is required when using bastion host")
		}
	}
	if conf.BastionHostKey != nil {
		config.BastionHostKey = *conf.BastionHostKey
		localConnection = false
	}
	if conf.BastionPort != nil {
		config.BastionPort = uint16(*conf.BastionPort)
		localConnection = false
	}

	// PROXY settings
	if conf.ProxyHost != nil {
		config.ProxyHost = *conf.ProxyHost
		localConnection = false
	}
	if conf.ProxyPort != nil {
		config.ProxyPort = uint16(*conf.ProxyPort)
		localConnection = false
	}
	if conf.ProxyUserName != nil {
		config.ProxyUserName = *conf.ProxyUserName
		localConnection = false
	}
	if conf.ProxyUserPassword != nil {
		config.ProxyUserPassword = *conf.ProxyUserPassword
		localConnection = false
	} else {
		if conf.ProxyUserName != nil {
			return nil, nil, localConnection, errors.New("password is required when proxy username is set")
		}
	}

	if conf.Protocol != nil {
		config.Type = *conf.Protocol
		localConnection = false
	} else {
		if !localConnection {
			return nil, nil, localConnection, errors.New("protocol is required in config")
		}
	}
	if conf.Host != nil {
		config.Host = *conf.Host
		localConnection = false
	} else {
		if !localConnection {
			return nil, nil, localConnection, errors.New("host is required in config")
		}
	}
	if conf.Port != nil {
		config.Port = uint16(*conf.Port)
		localConnection = false
	}
	if conf.Https != nil {
		config.HTTPS = *conf.Https
		localConnection = false
	}
	if conf.Insecure != nil {
		config.Insecure = *conf.Insecure
		localConnection = false
	}
	if config.Type == "ssh" {
		localConnection = false
		if conf.User != nil {
			config.User = *conf.User
		} else {
			return nil, nil, localConnection, errors.New("user is required for SSH connections")
		}
		if conf.Password != nil {
			config.Password = *conf.Password
		} else if conf.PrivateKey != nil {
			content, err := PathOrContents(*conf.PrivateKey)
			if err != nil {
				return nil, nil, localConnection, err
			}
			config.PrivateKey = content
		} else {
			return nil, nil, localConnection, errors.New("either password or private_key is required for SSH connections")
		}
	}
	if config.Type == "winrm" {
		localConnection = false
		if conf.User != nil {
			config.User = *conf.User
		} else {
			return nil, nil, localConnection, errors.New("user is required for WinRM connections")
		}
		if conf.Password != nil {
			config.Password = *conf.Password
		} else {
			return nil, nil, localConnection, errors.New("password is required for WinRM connections")
		}
	}

	if localConnection {
		return nil, nil, localConnection, nil
	}

	comm, err := communicator.New(config)
	return comm, nil, localConnection, err
}

// PathOrContents :: returns the contents of a file if the parameter is a file path, otherwise returns the parameter itself
func PathOrContents(poc string) (string, error) {
	if len(poc) == 0 {
		return poc, nil
	}

	path := poc
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, err
		}
	}

	// Check for valid file path
	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return string(contents), err
		}
		return string(contents), nil
	}

	// Return error if content is a file path and the file doesn't exist
	if len(path) > 1 && (path[0] == '/' || path[0] == '\\') {
		return "", fmt.Errorf("%s: no such file or dir", path)
	}

	// Return the inline content
	return poc, nil
}
