package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
	wormholetypes "github.com/wormhole-foundation/wormchain/x/wormhole/types"
)

// SaveWormholeConfig allows to store the given config inside the database
func (db *Db) SaveWormholeConfig(config *types.WormholeConfig) error {
	configBz, err := json.Marshal(&config.Config)
	if err != nil {
		return fmt.Errorf("error while marshaling wormhole config: %s", err)
	}

	stmt := `
INSERT INTO wormhole_config (config, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET config = excluded.config,
        height = excluded.height
WHERE wormhole_config.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(configBz), config.Height)
	if err != nil {
		return fmt.Errorf("error while storing wormhole config: %s", err)
	}

	return nil
}

// SaveGuardianValidatorList allows to store the given guardian validators list inside the database
func (db *Db) SaveGuardianValidatorList(guardianvalidatorlist []wormholetypes.GuardianValidator, height int64) error {
	if len(guardianvalidatorlist) == 0 {
		return nil
	}

	stmt := `INSERT INTO guardian_validator (guardian_key, validator_address, height) VALUES `
	var list []interface{}

	for i, entry := range guardianvalidatorlist {
		pi := i * 3
		stmt += fmt.Sprintf("($%d,$%d,$%d),", pi+1, pi+2, pi+3)
		guardianKeyBz, err := json.Marshal(&entry.GuardianKey)
		if err != nil {
			return fmt.Errorf("error while marshaling wormhole guardian key: %s", err)
		}

		validatorAddrBz, err := json.Marshal(&entry.ValidatorAddr)
		if err != nil {
			return fmt.Errorf("error while marshaling wormhole validator address: %s", err)
		}
		list = append(list, string(guardianKeyBz), string(validatorAddrBz), height)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT (guardian_key) DO UPDATE 
	SET validator_address = excluded.validator_address, 
		height = excluded.height
WHERE guardian_validator.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, list...)
	if err != nil {
		return fmt.Errorf("error while storing guardian validator list: %s", err)
	}

	return nil
}
