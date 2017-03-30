/*
 * Copyright (c) 2016, 2017
 *     PC-Game of Qihu.360. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. Neither the name of the University nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:  "<SLATER-GAME>",
	Long: "Slater game",
}

var configFile string

// Get : Get configuration varible
/* {{{ [config.Get] Get varible */
func Get(key string) interface{} {
	val := viper.Get(key)
	if nil == val {
		return ""
	}

	return val
}

/* }}} */

// SetDefault : Set default configuration variable
/* {{{ [config.SetDefault] Set variable */
func SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)

	return
}

/* }}} */

// Load : Load all configurations from file and enviroment
func Load() {
	// Fetch flag
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "", "", "Path of configuration file")

	if 0 == len(configFile) {
		viper.SetConfigFile("slater")
		viper.SetConfigType("json")
		viper.AddConfigPath("/etc/slater/")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(configFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		// Read config failed
	}

	// Enviroments
	viper.SetEnvPrefix("slater")
	viper.AutomaticEnv()

	return
}

/* }}} */

// init : Initialize, set global default values
/* {{{ [init] */
func init() {
	viper.SetDefault("listen_addr", "0.0.0.0")
	viper.SetDefault("listen_port", 9797)
	viper.SetDefault("room_service_addr", "127.0.0.1")
	viper.SetDefault("room_service_port", 9798)
	viper.SetDefault("room_service_ssl", true)
	viper.SetDefault("storage_service_addr", "127.0.0.1")
	viper.SetDefault("storage_service_port", 9799)
	viper.SetDefault("storage_service_ssl", true)

	return
}

/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
