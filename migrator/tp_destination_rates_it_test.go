// +build integration

/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package migrator

/*
import (
	"log"
	"path"
	"reflect"
	"testing"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

var (
	tpDstRtPathIn     string
	tpDstRtPathOut    string
	tpDstRtCfgIn      *config.CGRConfig
	tpDstRtCfgOut     *config.CGRConfig
	tpDstRtMigrator   *Migrator
	tpDestinationRate []*utils.TPDestinationRate
)

var sTestsTpDstRtIT = []func(t *testing.T){
	testTpDstRtITConnect,
	testTpDstRtITFlush,
	testTpDstRtITPopulate,
	testTpDstRtITMove,
	testTpDstRtITCheckData,
}

func TestTpDstRtMove(t *testing.T) {
	for _, stest := range sTestsTpDstRtIT {
		t.Run("TestTpDstRtMove", stest)
	}
}

func testTpDstRtITConnect(t *testing.T) {
	var err error
	tpDstRtPathIn = path.Join(*dataDir, "conf", "samples", "tutmongo")
	tpDstRtCfgIn, err = config.NewCGRConfigFromFolder(tpDstRtPathIn)
	if err != nil {
		t.Fatal(err)
	}
	tpDstRtPathOut = path.Join(*dataDir, "conf", "samples", "tutmysql")
	tpDstRtCfgOut, err = config.NewCGRConfigFromFolder(tpDstRtPathOut)
	if err != nil {
		t.Fatal(err)
	}
	storDBIn, err := engine.ConfigureStorDB(tpDstRtCfgIn.StorDBType, tpDstRtCfgIn.StorDBHost,
		tpDstRtCfgIn.StorDBPort, tpDstRtCfgIn.StorDBName,
		tpDstRtCfgIn.StorDBUser, tpDstRtCfgIn.StorDBPass,
		config.CgrConfig().StorDBMaxOpenConns,
		config.CgrConfig().StorDBMaxIdleConns,
		config.CgrConfig().StorDBConnMaxLifetime,
		config.CgrConfig().StorDBCDRSIndexes)
	if err != nil {
		log.Fatal(err)
	}
	storDBOut, err := engine.ConfigureStorDB(tpDstRtCfgOut.StorDBType,
		tpDstRtCfgOut.StorDBHost, tpDstRtCfgOut.StorDBPort, tpDstRtCfgOut.StorDBName,
		tpDstRtCfgOut.StorDBUser, tpDstRtCfgOut.StorDBPass,
		config.CgrConfig().StorDBMaxOpenConns,
		config.CgrConfig().StorDBMaxIdleConns,
		config.CgrConfig().StorDBConnMaxLifetime,
		config.CgrConfig().StorDBCDRSIndexes)
	if err != nil {
		log.Fatal(err)
	}
	tpDstRtMigrator, err = NewMigrator(nil, nil, tpDstRtCfgIn.DataDbType,
		tpDstRtCfgIn.DBDataEncoding, storDBIn, storDBOut, tpDstRtCfgIn.StorDBType, nil,
		tpDstRtCfgIn.DataDbType, tpDstRtCfgIn.DBDataEncoding, nil,
		tpDstRtCfgIn.StorDBType, false, false, false, false, false)
	if err != nil {
		log.Fatal(err)
	}
}

func testTpDstRtITFlush(t *testing.T) {
	if err := tpDstRtMigrator.storDBIn.Flush(
		path.Join(tpDstRtCfgIn.DataFolderPath, "storage", tpDstRtCfgIn.StorDBType)); err != nil {
		t.Error(err)
	}

	if err := tpDstRtMigrator.storDBOut.Flush(
		path.Join(tpDstRtCfgOut.DataFolderPath, "storage", tpDstRtCfgOut.StorDBType)); err != nil {
		t.Error(err)
	}
}

func testTpDstRtITPopulate(t *testing.T) {
	tpDestinationRate = []*utils.TPDestinationRate{
		&utils.TPDestinationRate{
			TPid: utils.TEST_SQL,
			ID:   "DR_FREESWITCH_USERS",
			DestinationRates: []*utils.DestinationRate{
				&utils.DestinationRate{
					DestinationId:    "FS_USERS",
					RateId:           "RT_FS_USERS",
					RoundingMethod:   "*up",
					RoundingDecimals: 2},
			},
		},
	}
	if err := tpDstRtMigrator.storDBIn.SetTPDestinationRates(tpDestinationRate); err != nil {
		t.Error("Error when setting TpDestinationRate ", err.Error())
	}
	currentVersion := engine.CurrentStorDBVersions()
	err := tpDstRtMigrator.storDBOut.SetVersions(currentVersion, false)
	if err != nil {
		t.Error("Error when setting version for TpDestinationRate ", err.Error())
	}
}

func testTpDstRtITMove(t *testing.T) {
	err, _ := tpDstRtMigrator.Migrate([]string{utils.MetaTpDestinationRates})
	if err != nil {
		t.Error("Error when migrating TpDestinationRate ", err.Error())
	}
}

func testTpDstRtITCheckData(t *testing.T) {
	result, err := tpDstRtMigrator.storDBOut.GetTPDestinationRates(
		tpDestinationRate[0].TPid, tpDestinationRate[0].ID, nil)
	if err != nil {
		t.Error("Error when getting TpDestinationRate ", err.Error())
	}
	if !reflect.DeepEqual(tpDestinationRate[0], result[0]) {
		t.Errorf("Expecting: %+v, received: %+v",
			utils.ToJSON(tpDestinationRate[0]), utils.ToJSON(result[0]))
	}
	result, err = tpDstRtMigrator.storDBIn.GetTPDestinationRates(
		tpDestinationRate[0].TPid, tpDestinationRate[0].ID, nil)
	if err != utils.ErrNotFound {
		t.Error(err)
	}
}
*/