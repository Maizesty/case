// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conf

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ascii = `
	_____                    _____                    _____                    _____          
	/\    \                  /\    \                  /\    \                  /\    \         
 /::\    \                /::\    \                /::\    \                /::\    \        
/::::\    \              /::::\    \              /::::\    \              /::::\    \       
/::::::\    \            /::::::\    \            /::::::\    \            /::::::\    \      
/:::/\:::\    \          /:::/\:::\    \          /:::/\:::\    \          /:::/\:::\    \     
/:::/  \:::\    \        /:::/__\:::\    \        /:::/__\:::\    \        /:::/__\:::\    \    
/:::/    \:::\    \      /::::\   \:::\    \       \:::\   \:::\    \      /::::\   \:::\    \   
/:::/    / \:::\    \    /::::::\   \:::\    \    ___\:::\   \:::\    \    /::::::\   \:::\    \  
/:::/    /   \:::\    \  /:::/\:::\   \:::\    \  /\   \:::\   \:::\    \  /:::/\:::\   \:::\    \ 
/:::/____/     \:::\____\/:::/  \:::\   \:::\____\/::\   \:::\   \:::\____\/:::/__\:::\   \:::\____\
\:::\    \      \::/    /\::/    \:::\  /:::/    /\:::\   \:::\   \::/    /\:::\   \:::\   \::/    /
\:::\    \      \/____/  \/____/ \:::\/:::/    /  \:::\   \:::\   \/____/  \:::\   \:::\   \/____/ 
\:::\    \                       \::::::/    /    \:::\   \:::\    \       \:::\   \:::\    \     
\:::\    \                       \::::/    /      \:::\   \:::\____\       \:::\   \:::\____\    
\:::\    \                      /:::/    /        \:::\  /:::/    /        \:::\   \::/    /    
\:::\    \                    /:::/    /          \:::\/:::/    /          \:::\   \/____/     
\:::\    \                  /:::/    /            \::::::/    /            \:::\    \         
\:::\____\                /:::/    /              \::::/    /              \:::\____\        
 \::/    /                \::/    /                \::/    /                \::/    /        
	\/____/                  \/____/                  \/____/                  \/____/         
																
`
	Conf Config
)

// Config configuration details of balancer
type Config struct {
	Location            []*Location `yaml:"location"`
	Port                int         `yaml:"port"`
	PartitionFile				string			`yaml:"partition_files"`
	CacheSize						int64					`yaml:"cache_size"`
}

// Location routing details of balancer
type Location struct {
	ProxyPass   []string `yaml:"proxy_pass"`
	BalanceMode string   `yaml:"balance_mode"`
}

// ReadConfig read configuration from `fileName` file
func ReadConfig(fileName string) (*Config, error) {
	in, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(in, &Conf)
	if err != nil {
		return nil, err
	}
	return &Conf, nil
}

// Print print config details
func (c *Config) Print() {
	fmt.Printf("%s\nPort: %d\n\nLocation:\n",
		ascii,  c.Port, )
	for _, l := range c.Location {
		fmt.Printf("\tRoute: %s\n\t\tMode: %s\n\n",
			 l.ProxyPass, l.BalanceMode)
	}
}

// Validation verify the configuration details of the balancer
func (c *Config) Validation() error {
	if len(c.Location) == 0 {
		return errors.New("the details of location cannot be null")
	}
	return nil
}
