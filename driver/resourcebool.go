package driver

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type resourceBool struct{}

func (rb *resourceBool) value(db *db, ro *models.ResourceOperation, deviceName, deviceResourceName string) (*dsModels.CommandValue, error) {
	result := &dsModels.CommandValue{}

	enableRandomization, currentValue, _, err := db.getVirtualResourceData(deviceName, deviceResourceName)
	if err != nil {
		return result, err
	}

	var newValueBool bool
	if enableRandomization {
		rand.Seed(time.Now().UnixNano())
		newValueBool = rand.Int()%2 == 0
	} else {
		if newValueBool, err = strconv.ParseBool(currentValue); err != nil {
			return result, err
		}
	}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	if result, err = dsModels.NewBoolValue(ro, now, newValueBool); err != nil {
		return result, err
	}
	if err := db.updateResourceValue(result.ValueToString(), deviceName, deviceResourceName, false); err != nil {
		return result, err
	}

	return result, nil
}

func (rb *resourceBool) write(param *dsModels.CommandValue, deviceName string, db *db) error {
	switch param.RO.Object {
	case deviceResourceEnableRandomization:
		v, err := param.BoolValue()
		if err != nil {
			return fmt.Errorf("resourceBool.write: %v", err)
		}
		if err := db.updateResourceRandomization(v, deviceName, param.RO.Resource); err != nil {
			return fmt.Errorf("resourceBool.write: %v", err)
		} else {
			return nil
		}
	case deviceResourceBool:
		if _, err := param.BoolValue(); err == nil {
			if err := db.updateResourceValue(param.ValueToString(), deviceName, param.RO.Resource, true); err != nil {
				return fmt.Errorf("resourceBool.write: %v", err)
			} else {
				return nil
			}
		} else {
			return fmt.Errorf("resourceBool.write: %v", err)
		}
	default:
		return fmt.Errorf("resourceBool.write: unknown device resource: %s", param.RO.Object)
	}
}
