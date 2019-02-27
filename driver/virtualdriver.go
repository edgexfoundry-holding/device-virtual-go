// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a implementation of a ProtocolDriver interface.
//
package driver

import (
	"fmt"
	"os"
	"strconv"
	"time"

	sdk "github.com/edgexfoundry/device-sdk-go"
	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logging"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	_ "modernc.org/ql/driver"
)

type VirtualDriver struct {
	lc             logger.LoggingClient
	asyncCh        chan<- *dsModels.AsyncValues
	virtualDevices map[string]*virtualDevice
	db             *db
}

func (d *VirtualDriver) DisconnectDevice(address *models.Addressable) error {
	d.lc.Info(fmt.Sprintf("VirtualDriver.DisconnectDevice: device-virtual driver is disconnecting to %v", address))
	return nil
}

func (d *VirtualDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	d.lc = lc
	d.asyncCh = asyncCh
	d.virtualDevices = make(map[string]*virtualDevice)

	if _, err := os.Stat(qlDatabaseDir); os.IsNotExist(err) {
		os.Mkdir(qlDatabaseDir, os.ModeDir)
	}

	d.db = getDb()
	if err := d.db.openDb(); err != nil {
		d.lc.Info(fmt.Sprintf("Create db connection failed: %v", err))
		return err
	}
	defer func() {
		if err := d.db.closeDb(); err != nil {
			d.lc.Info(fmt.Sprintf("Close db failed: %v", err))
			return
		}
	}()

	err := d.db.startTransaction()
	if err != nil {
		d.lc.Info(fmt.Sprintf("Start a transaction failed: %v", err))
		return err
	}

	if err := d.db.exec(SQL_DROP_TABLE); err != nil {
		d.lc.Info(fmt.Sprintf("Drop table failed: %v", err))
		return err
	}

	if err := d.db.exec(SQL_CREATE_TABLE); err != nil {
		d.lc.Info(fmt.Sprintf("Create table failed: %v", err))
		return err
	}

	service := sdk.RunningService()
	devices := service.Devices()
	for _, device := range devices {
		for _, r := range device.Profile.Resources {
			for _, ro := range r.Get {
				for _, dr := range device.Profile.DeviceResources {
					if ro.Object == dr.Name {
						/*
							d.Name <-> VIRTUAL_RESOURCE.DeviceName
							dr.Name <-> VIRTUAL_RESOURCE.CommandName, VIRTUAL_RESOURCE.ResourceName
							ro.Object <-> VIRTUAL_RESOURCE.DeviceResourceName
							dr.Properties.Value.Type <-> VIRTUAL_RESOURCE.DataType
							dr.Properties.Value.DefaultValue <-> VIRTUAL_RESOURCE.Value
						*/
						if err := d.db.exec(SQL_INSERT, device.Name, dr.Name, dr.Name, true, dr.Properties.Value.Type, dr.Properties.Value.DefaultValue); err != nil {
							d.lc.Info(fmt.Sprintf("Insert one row into db failed: %v", err))
							return err
						}
					}
				}
			}
		}
	}

	if err = d.db.commit(); err != nil {
		d.lc.Info(fmt.Sprintf("Commit transaction failed: %v", err))
		return err
	}

	return nil
}

func (d *VirtualDriver) HandleReadCommands(addr *models.Addressable, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	rd, ok := d.virtualDevices[addr.Name]
	if !ok {
		rd = newVirtualDevice()
		d.virtualDevices[addr.Name] = rd
	}

	res = make([]*dsModels.CommandValue, len(reqs))
	now := time.Now().UnixNano() / int64(time.Millisecond)

	if err := d.db.openDb(); err != nil {
		d.lc.Info(fmt.Sprintf("Create db connection failed: %v", err))
		return nil, err
	}
	defer func() {
		if err := d.db.closeDb(); err != nil {
			d.lc.Info(fmt.Sprintf("Close db failed: %v", err))
			return
		}
	}()

	for i, req := range reqs {
		t := req.DeviceResource.Properties.Value.Type
		v, err := rd.value(addr.Name, req.DeviceResource.Name, req.DeviceResource.Properties.Value.Minimum,
			req.DeviceResource.Properties.Value.Maximum, d.db)
		if err != nil {
			return nil, err
		}
		var cv *dsModels.CommandValue
		switch t {
		case typeBool:
			newValue, _ := strconv.ParseBool(v)
			cv, _ = dsModels.NewBoolValue(&req.RO, now, bool(newValue))
		case typeInt:
			newValue, _ := strconv.ParseInt(v, 10, strconv.IntSize)
			if strconv.IntSize == 32 {
				cv, _ = dsModels.NewInt32Value(&req.RO, now, int32(newValue))
			} else {
				cv, _ = dsModels.NewInt64Value(&req.RO, now, int64(newValue))
			}
		case typeInt8:
			newValue, _ := strconv.ParseInt(v, 10, 8)
			cv, _ = dsModels.NewInt8Value(&req.RO, now, int8(newValue))
		case typeInt16:
			newValue, _ := strconv.ParseInt(v, 10, 16)
			cv, _ = dsModels.NewInt16Value(&req.RO, now, int16(newValue))
		case typeInt32:
			newValue, _ := strconv.ParseInt(v, 10, 32)
			cv, _ = dsModels.NewInt32Value(&req.RO, now, int32(newValue))
		case typeInt64:
			newValue, _ := strconv.ParseInt(v, 10, 64)
			cv, _ = dsModels.NewInt64Value(&req.RO, now, int64(newValue))
		case typeUint:
			newValue, _ := strconv.ParseUint(v, 10, strconv.IntSize)
			if strconv.IntSize == 32 {
				cv, _ = dsModels.NewUint32Value(&req.RO, now, uint32(newValue))
			} else {
				cv, _ = dsModels.NewUint64Value(&req.RO, now, uint64(newValue))
			}
		case typeUint8:
			newValue, _ := strconv.ParseUint(v, 10, 8)
			cv, _ = dsModels.NewUint8Value(&req.RO, now, uint8(newValue))
		case typeUint16:
			newValue, _ := strconv.ParseUint(v, 10, 16)
			cv, _ = dsModels.NewUint16Value(&req.RO, now, uint16(newValue))
		case typeUint32:
			newValue, _ := strconv.ParseUint(v, 10, 32)
			cv, _ = dsModels.NewUint32Value(&req.RO, now, uint32(newValue))
		case typeUint64:
			newValue, _ := strconv.ParseUint(v, 10, 64)
			cv, _ = dsModels.NewUint64Value(&req.RO, now, uint64(newValue))
		case typeFloat32:
			newValue, _ := strconv.ParseFloat(v, 32)
			cv, _ = dsModels.NewFloat32Value(&req.RO, now, float32(newValue))
		case typeFloat64:
			newValue, _ := strconv.ParseFloat(v, 64)
			cv, _ = dsModels.NewFloat64Value(&req.RO, now, float64(newValue))
		}
		res[i] = cv
	}

	return res, nil
}

func (d *VirtualDriver) HandleWriteCommands(addr *models.Addressable, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {
	rd, ok := d.virtualDevices[addr.Name]
	if !ok {
		rd = newVirtualDevice()
		d.virtualDevices[addr.Name] = rd
	}

	if err := d.db.openDb(); err != nil {
		d.lc.Info(fmt.Sprintf("Create db connection failed: %v", err))
		return err
	}
	defer func() {
		if err := d.db.closeDb(); err != nil {
			d.lc.Info(fmt.Sprintf("Close db failed: %v", err))
			return
		}
	}()

	var err error
	if err = d.db.startTransaction(); err != nil {
		d.lc.Info(fmt.Sprintf("Start a transaction failed: %v", err))
		return err
	}

	for _, param := range params {
		if param.RO.Object == "Enable_Randomization" {
			v, err := param.BoolValue()
			if err != nil {
				return fmt.Errorf("VirtualDriver.HandleWriteCommands: %v", err)
			}
			if err := d.db.exec(SQL_UPDATE_ENABLERANDOMIZATION, v, addr.Name, param.RO.Resource); err != nil {
				d.lc.Info(fmt.Sprintf("VirtualDriver.HandleWriteCommands: %v", err))
				return err
			}
			continue
		}

		switch param.Type {
		case dsModels.Bool:
			_, err = param.BoolValue()
		case dsModels.Int8:
			_, err = param.Int8Value()
		case dsModels.Int16:
			_, err = param.Int16Value()
		case dsModels.Int32:
			_, err = param.Int32Value()
		case dsModels.Int64:
			_, err = param.Int64Value()
		case dsModels.Uint8:
			_, err = param.Uint8Value()
		case dsModels.Uint16:
			_, err = param.Uint16Value()
		case dsModels.Uint32:
			_, err = param.Uint32Value()
		case dsModels.Uint64:
			_, err = param.Uint64Value()
		case dsModels.Float32:
			_, err = param.Float32Value()
		case dsModels.Float64:
			_, err = param.Float64Value()
		default:
			return fmt.Errorf("VirtualDriver.HandleWriteCommands: there is no matched device resource for %s", param.String())
		}
		if err != nil {
			return fmt.Errorf("VirtualDriver.HandleWriteCommands: %v", err)
		}

		switch param.Type {
		case dsModels.Float32:
			v, _ := param.Float32Value()
			if err = d.db.exec(SQL_UPDATE_VALUE, strconv.FormatFloat(float64(v), 'f', -1, 32), addr.Name, param.RO.Resource); err != nil {
				return err
			}
		case dsModels.Float64:
			v, _ := param.Float64Value()
			if err = d.db.exec(SQL_UPDATE_VALUE, strconv.FormatFloat(float64(v), 'f', -1, 64), addr.Name, param.RO.Resource); err != nil {
				return err
			}
		default:
			if err = d.db.exec(SQL_UPDATE_VALUE, param.ValueToString(), addr.Name, param.RO.Resource); err != nil {
				return err
			}
		}
	}
	if err := d.db.commit(); err != nil {
		return err
	}
	return nil
}

func (d *VirtualDriver) Stop(force bool) error {
	d.lc.Info("VirtualDriver.Stop: device-virtual driver is stopping...")
	return nil
}
