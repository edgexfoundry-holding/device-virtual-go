// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	//defMinInt8, defMaxInt8   = math.MinInt8, math.MaxInt8
	defMinInt8, defMaxInt8 = -100, 100
	//defMinInt16, defMaxInt16 = math.MinInt16, math.MaxInt16
	defMinInt16, defMaxInt16 = -100, 100
	//defMinInt32, defMaxInt32 = math.MinInt32, math.MaxInt32
	defMinInt32, defMaxInt32 = -100, 100
	//defMinInt64, defMaxInt64 = math.MinInt64, math.MaxInt64
	defMinInt64, defMaxInt64 = -100, 100
	//defMinUint8, defMaxUint8 = 0, math.MaxUint8
	defMinUint8, defMaxUint8 = 0, 100
	//defMinUint16, defMaxUint16 = 0, math.MaxUint16
	defMinUint16, defMaxUint16 = 0, 100
	//defMinUint32, defMaxUint32 = 0, math.MaxUint32
	defMinUint32, defMaxUint32 = 0, 100
	//defMinUint64, defMaxUint64 = 0, math.MaxUint64
	defMinUint64, defMaxUint64 = 0, 100
	//defMinFloat32, defMaxFloat32 = -math.MaxFloat32, math.MaxFloat32
	defMinFloat32, defMaxFloat32 = -100, 100
	//defMinFloat64, defMaxFloat64 = -math.MaxFloat64, math.MaxFloat64
	defMinFloat64, defMaxFloat64 = -100, 100
	typeBool                     = "Bool"
	typeInt8                     = "Int8"
	typeInt16                    = "Int16"
	typeInt32                    = "Int32"
	typeInt64                    = "Int64"
	typeUint8                    = "Uint8"
	typeUint16                   = "Uint16"
	typeUint32                   = "Uint32"
	typeUint64                   = "Uint64"
	typeFloat32                  = "Float32"
	typeFloat64                  = "Float64"
)

type virtualDevice struct {
	minInt8    int64
	maxInt8    int64
	minUint8   uint64
	maxUint8   uint64
	minInt16   int64
	maxInt16   int64
	minUint16  uint64
	maxUint16  uint64
	minInt32   int64
	maxInt32   int64
	minUint32  uint64
	maxUint32  uint64
	minInt64   int64
	maxInt64   int64
	minUint64  uint64
	maxUint64  uint64
	minFloat32 float64
	maxFloat32 float64
	minFloat64 float64
	maxFloat64 float64
}

func (d *virtualDevice) value(deviceName, deviceResourceName, minimum, maximum string, db *db) (string, error) {
	rows, err := db.query(SQL_SELECT, deviceName, deviceResourceName)
	if err != nil {
		return "0", err
	}

	var data struct {
		DeviceName          string
		CommandName         string
		DeviceResourceName  string
		EnableRandomization string
		DataType            string
		Value               string
	}

	if rows.Next() {
		if err = rows.Scan(&data.DeviceName, &data.CommandName, &data.DeviceResourceName, &data.EnableRandomization, &data.DataType, &data.Value); err != nil {
			rows.Close()
			return "0", err
		}
		rows.Close()
	}

	if er, _ := strconv.ParseBool(data.EnableRandomization); er {
		var newValueInt int64
		var newValueUint uint64
		var newValueFloat float64
		var newValueBool bool
		var para1 string

		switch data.DataType {
		case typeBool:
			newValueBool = randomBool()
			para1 = strconv.FormatBool(newValueBool)
		case typeInt8:
			min, err := parseStrToInt(minimum, typeInt8)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToInt(maximum, typeInt8)
			if err != nil {
				return "0", err
			}
			if min < d.minInt8 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %d ", typeInt8, d.minInt8)
			}
			if max > d.maxInt8 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %d ", typeInt8, d.maxInt8)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %d <= minimum : %d", max, min)
			} else {
				newValueInt = randomInt(min, max)
				para1 = strconv.FormatInt(newValueInt, 10)
			}
		case typeUint8:
			min, err := parseStrToUint(minimum, typeUint8)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToUint(maximum, typeUint8)
			if err != nil {
				return "0", err
			}
			if min < d.minUint16 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %d ", typeUint8, d.minUint8)
			}
			if max > d.maxUint16 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %d ", typeUint8, d.maxUint8)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %d <= minimum : %d", max, min)
			} else {
				newValueUint = randomUint(min, max)
				para1 = strconv.FormatUint(newValueUint, 10)
			}
		case typeInt16:
			min, err := parseStrToInt(minimum, typeInt16)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToInt(maximum, typeInt16)
			if err != nil {
				return "0", err
			}
			if min < d.minInt16 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %d ", typeInt16, d.minInt16)
			}
			if max > d.maxInt16 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %d ", typeInt16, d.maxInt16)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %d <= minimum : %d", max, min)
			} else {
				newValueInt = randomInt(min, max)
				para1 = strconv.FormatInt(newValueInt, 10)
			}
		case typeUint16:
			min, err := parseStrToUint(minimum, typeUint16)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToUint(maximum, typeUint16)
			if err != nil {
				return "0", err
			}
			if min < d.minUint16 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %d ", typeUint16, d.minUint16)
			}
			if max > d.maxUint16 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %d ", typeUint16, d.maxUint16)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %d <= minimum : %d", max, min)
			} else {
				newValueUint = randomUint(min, max)
				para1 = strconv.FormatUint(newValueUint, 10)
			}
		case typeInt32:
			min, err := parseStrToInt(minimum, typeInt32)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToInt(maximum, typeInt32)
			if err != nil {
				return "0", err
			}
			if min < d.minInt32 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %d ", typeInt32, d.minInt32)
			}
			if max > d.maxInt32 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %d ", typeInt32, d.maxInt32)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %d <= minimum : %d", max, min)
			} else {
				newValueInt = randomInt(min, max)
				para1 = strconv.FormatInt(newValueInt, 10)
			}
		case typeUint32:
			min, err := parseStrToUint(minimum, typeUint32)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToUint(maximum, typeUint32)
			if err != nil {
				return "0", err
			}
			if min < d.minUint32 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %d ", typeUint32, d.minUint32)
			}
			if max > d.maxUint32 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %d ", typeUint32, d.maxUint32)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %d <= minimum : %d", max, min)
			} else {
				newValueUint = randomUint(min, max)
				para1 = strconv.FormatUint(newValueUint, 10)
			}
		case typeInt64:
			min, err := parseStrToInt(minimum, typeInt64)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToInt(maximum, typeInt64)
			if err != nil {
				return "0", err
			}
			if min < d.minInt64 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %d ", typeInt64, d.minInt64)
			}
			if max > d.maxInt64 {

				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %d ", typeInt64, d.maxInt64)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %d <= minimum : %d", max, min)
			} else {
				newValueInt = randomInt(min, max)
				para1 = strconv.FormatInt(newValueInt, 10)
			}
		case typeUint64:
			min, err := parseStrToUint(minimum, typeUint64)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToUint(maximum, typeUint64)
			if err != nil {
				return "0", err
			}
			if min < d.minUint64 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %d ", typeUint64, d.minUint64)
			}
			if max > d.maxUint64 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %d ", typeUint64, d.maxUint64)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %d <= minimum : %d", max, min)
			} else {
				newValueUint = randomUint(min, max)
				para1 = strconv.FormatUint(newValueUint, 10)
			}
		case typeFloat32:
			min, err := parseStrToFloat(minimum, typeFloat32)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToFloat(maximum, typeFloat32)
			if err != nil {
				return "0", err
			}
			if min < d.minFloat32 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %f ", typeFloat32, d.minFloat32)
			}
			if max > d.maxFloat32 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %f ", typeFloat32, d.maxFloat32)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %f <= minimum : %f", max, min)
			} else {
				newValueFloat = randomFloat(min, max)
				para1 = strconv.FormatFloat(newValueFloat, 'f', -1, 32)
			}
		case typeFloat64:
			min, err := parseStrToFloat(minimum, typeFloat64)
			if err != nil {
				return "0", err
			}
			max, err := parseStrToFloat(maximum, typeFloat64)
			if err != nil {
				return "0", err
			}
			if min < d.minFloat64 {
				return "0", fmt.Errorf("virtualDevice.value: the minimum value of %s is: %f ", typeFloat64, d.minFloat64)
			}
			if max > d.maxFloat64 {
				return "0", fmt.Errorf("virtualDevice.value: the maximum value of %s is: %f ", typeFloat64, d.maxFloat64)
			}

			if max <= min {
				return "0", fmt.Errorf("virtualDevice.value: maximum: %f <= minimum : %f", max, min)
			} else {
				newValueFloat = randomFloat(min, max)
				para1 = strconv.FormatFloat(newValueFloat, 'f', -1, 64)
			}
		default:
			return "0", fmt.Errorf("virtualDevice.value: wrong value type: %s", deviceResourceName)
		}

		if err := db.exec(SQL_UPDATE_VALUE, para1, data.DeviceName, data.CommandName); err != nil {
			return "0", err
		}
		return para1, nil
	} else {
		return data.Value, nil
	}

}

func newVirtualDevice() *virtualDevice {
	return &virtualDevice{
		minInt8:    defMinInt8,
		maxInt8:    defMaxInt8,
		minUint8:   defMinUint8,
		maxUint8:   defMaxUint8,
		minInt16:   defMinInt16,
		maxInt16:   defMaxInt16,
		minUint16:  defMinUint16,
		maxUint16:  defMaxUint16,
		minInt32:   defMinInt32,
		maxInt32:   defMaxInt32,
		minUint32:  defMinUint32,
		maxUint32:  defMaxUint32,
		minInt64:   defMinInt64,
		maxInt64:   defMaxInt64,
		minUint64:  defMinUint64,
		maxUint64:  defMaxUint64,
		minFloat32: defMinFloat32,
		maxFloat32: defMaxFloat32,
		minFloat64: defMinFloat64,
		maxFloat64: defMaxFloat64,
	}
}

func randomInt(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	if max-min == -1 {
		return (rand.Int63n(max/2-min/2) + min/2) * 2
	} else {
		return rand.Int63n(max-min+1) + min
	}
}

func randomUint(min, max uint64) uint64 {
	rand.Seed(time.Now().UnixNano())
	if max-min < math.MaxInt64 {
		return uint64(rand.Int63n(int64(max-min+1))) + min
	}
	x := rand.Uint64()
	for x > max-min {
		x = rand.Uint64()
	}
	return x
}

func randomFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

func randomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()%2 == 0
}

func parseStrToInt(s, dt string) (int64, error) {
	v := int64(0)
	var err error

	switch dt {
	case typeInt8:
		v, err = strconv.ParseInt(s, 10, 8)
	case typeInt16:
		v, err = strconv.ParseInt(s, 10, 16)
	case typeInt32:
		v, err = strconv.ParseInt(s, 10, 32)
	case typeInt64:
		v, err = strconv.ParseInt(s, 10, 64)
	default:
		return v, fmt.Errorf("virtualDevice.value: unknown data type: %s ", dt)
	}

	if err != nil {
		return v, fmt.Errorf("virtualDevice.value: error when parsing string to %s : %s ", dt, s)
	} else {
		return v, nil
	}
}

func parseStrToUint(s, dt string) (uint64, error) {
	v := uint64(0)
	var err error

	switch dt {
	case typeUint8:
		v, err = strconv.ParseUint(s, 10, 8)
	case typeUint16:
		v, err = strconv.ParseUint(s, 10, 16)
	case typeUint32:
		v, err = strconv.ParseUint(s, 10, 32)
	case typeUint64:
		v, err = strconv.ParseUint(s, 10, 64)
	default:
		return v, fmt.Errorf("virtualDevice.value: unknown data type: %s ", dt)
	}

	if err != nil {
		return v, fmt.Errorf("virtualDevice.value: error when parsing string to %s : %s ", dt, s)
	} else {
		return v, nil
	}
}

func parseStrToFloat(s, dt string) (float64, error) {
	v := float64(0)
	var err error

	switch dt {
	case typeFloat32:
		v, err = strconv.ParseFloat(s, 32)
	case typeFloat64:
		v, err = strconv.ParseFloat(s, 64)
	default:
		return v, fmt.Errorf("virtualDevice.value: unknown data type: %s ", dt)
	}

	if err != nil {
		return v, fmt.Errorf("virtualDevice.value: error when parsing string to %s : %s ", dt, s)
	} else {
		return v, nil
	}
}
