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
package engine

import (
	"reflect"
	"testing"
	"time"

	"github.com/cgrates/cgrates/utils"
)

var sq *StatQueue

func TestStatQueuesSort(t *testing.T) {
	sInsts := StatQueues{
		&StatQueue{sqPrfl: &StatQueueProfile{ID: "FIRST", Weight: 30.0}},
		&StatQueue{sqPrfl: &StatQueueProfile{ID: "SECOND", Weight: 40.0}},
		&StatQueue{sqPrfl: &StatQueueProfile{ID: "THIRD", Weight: 30.0}},
		&StatQueue{sqPrfl: &StatQueueProfile{ID: "FOURTH", Weight: 35.0}},
	}
	sInsts.Sort()
	eSInst := StatQueues{
		&StatQueue{sqPrfl: &StatQueueProfile{ID: "SECOND", Weight: 40.0}},
		&StatQueue{sqPrfl: &StatQueueProfile{ID: "FOURTH", Weight: 35.0}},
		&StatQueue{sqPrfl: &StatQueueProfile{ID: "FIRST", Weight: 30.0}},
		&StatQueue{sqPrfl: &StatQueueProfile{ID: "THIRD", Weight: 30.0}},
	}
	if !reflect.DeepEqual(eSInst, sInsts) {
		t.Errorf("expecting: %+v, received: %+v", eSInst, sInsts)
	}
}

func TestStatRemEventWithID(t *testing.T) {
	sq = &StatQueue{
		SQMetrics: map[string]StatMetric{
			utils.MetaASR: &StatASR{
				Answered: 1,
				Count:    2,
				Events: map[string]bool{
					"cgrates.org:TestRemEventWithID_1": true,
					"cgrates.org:TestRemEventWithID_2": false,
				},
			},
		},
	}
	asrMetric := sq.SQMetrics[utils.MetaASR].(*StatASR)
	if asrMetricIf := asrMetric.GetValue(); asrMetricIf.(float64) != 50 {
		t.Errorf("received asrMetric: %v", asrMetricIf)
	}
	sq.remEventWithID("cgrates.org:TestRemEventWithID_1")
	if asrMetricIf := asrMetric.GetValue(); asrMetricIf.(float64) != 0 {
		t.Errorf("received asrMetric: %v", asrMetricIf)
	} else if len(asrMetric.Events) != 1 {
		t.Errorf("unexpected Events in asrMetric: %+v", asrMetric.Events)
	}
	sq.remEventWithID("cgrates.org:TestRemEventWithID_5") // non existent
	if asrMetricIf := asrMetric.GetValue(); asrMetricIf.(float64) != 0 {
		t.Errorf("received asrMetric: %v", asrMetricIf)
	} else if len(asrMetric.Events) != 1 {
		t.Errorf("unexpected Events in asrMetric: %+v", asrMetric.Events)
	}
	sq.remEventWithID("cgrates.org:TestRemEventWithID_2")
	if asrMetricIf := asrMetric.GetValue(); asrMetricIf.(float64) != -1 {
		t.Errorf("received asrMetric: %v", asrMetricIf)
	} else if len(asrMetric.Events) != 0 {
		t.Errorf("unexpected Events in asrMetric: %+v", asrMetric.Events)
	}
	sq.remEventWithID("cgrates.org:TestRemEventWithID_2")
	if asrMetricIf := asrMetric.GetValue(); asrMetricIf.(float64) != -1 {
		t.Errorf("received asrMetric: %v", asrMetricIf)
	} else if len(asrMetric.Events) != 0 {
		t.Errorf("unexpected Events in asrMetric: %+v", asrMetric.Events)
	}
}

func TestStatRemExpired(t *testing.T) {
	sq = &StatQueue{
		SQMetrics: map[string]StatMetric{
			utils.MetaASR: &StatASR{
				Answered: 2,
				Count:    3,
				Events: map[string]bool{
					"cgrates.org:TestStatRemExpired_1": true,
					"cgrates.org:TestStatRemExpired_2": false,
					"cgrates.org:TestStatRemExpired_3": true,
				},
			},
		},
		SQItems: []struct {
			EventID    string
			ExpiryTime *time.Time
		}{
			struct {
				EventID    string     // Bounded to the original StatEvent
				ExpiryTime *time.Time // Used to auto-expire events
			}{"cgrates.org:TestStatRemExpired_1", utils.TimePointer(time.Now())},
			{"cgrates.org:TestStatRemExpired_2", utils.TimePointer(time.Now())},
			{"cgrates.org:TestStatRemExpired_3", utils.TimePointer(time.Now().Add(time.Duration(time.Minute)))},
		},
	}
	asrMetric := sq.SQMetrics[utils.MetaASR].(*StatASR)
	if asrMetricIf := asrMetric.GetValue(); asrMetricIf.(float64) != 66.66667 {
		t.Errorf("received asrMetric: %v", asrMetricIf)
	}
	sq.remExpired()
	if asrMetricIf := asrMetric.GetValue(); asrMetricIf.(float64) != 100 {
		t.Errorf("received asrMetric: %v", asrMetricIf)
	} else if len(asrMetric.Events) != 1 {
		t.Errorf("unexpected Events in asrMetric: %+v", asrMetric.Events)
	}
	if len(sq.SQItems) != 1 {
		t.Errorf("Unexpected items: %+v", sq.SQItems)
	}
}