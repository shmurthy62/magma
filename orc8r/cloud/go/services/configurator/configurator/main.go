/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
Configurator is a dedicated Magma Cloud service which maintains configurations
and meta data for the network and network entity structures.
*/

package main

import (
	"database/sql"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/servicers"
	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/golang/glog"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, configurator.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}
	db, err := sql.Open(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}

	factory := storage.NewSQLConfiguratorStorageFactory(db, &storage.DefaultIDGenerator{}, sqorc.GetSqlBuilder())
	err = factory.InitializeServiceStorage()

	nbServicer, err := servicers.NewNorthboundConfiguratorServicer(factory)
	if err != nil {
		glog.Fatalf("Failed to instantiate the device servicer: %v", nbServicer)
	}
	protos.RegisterNorthboundConfiguratorServer(srv.GrpcServer, nbServicer)

	sbServicer, err := servicers.NewSouthboundConfiguratorServicer(factory)
	if err != nil {
		glog.Fatalf("Failed to instantiate the device servicer: %v", sbServicer)
	}
	protos.RegisterSouthboundConfiguratorServer(srv.GrpcServer, sbServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Failed to start configurator service: %v", err)
	}
}
