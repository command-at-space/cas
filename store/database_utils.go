/* */

package store

import (
	"encoding/json"
	"log"
)

func loadConfigJSON(c *configDB) {
	err := json.Unmarshal(getConfigJSON(), &c)
	if err != nil {
		log.Fatal("Error parsing JSON config => \n", err)
	}
}
