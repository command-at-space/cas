/* */

package main

func getGlobalConfigJSON() (configjson []byte) {
	configjson = []byte(`
	{
		"config": {
			"mode": "production",
			"host": "localhost",
			"port": 6900,
			"errorsLogFile": "logs/errors.log",
			"infoLogFile":"logs/info.log",
			"chatLogFile":"logs/chat.log"
		}
	}
	`)
	return
}
