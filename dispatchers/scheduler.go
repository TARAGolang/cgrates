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

package dispatchers

import (
	"time"

	"github.com/cgrates/cgrates/utils"
)

func (dS *DispatcherService) SchedulerSv1Ping(args *utils.CGREventWithArgDispatcher, reply *string) (err error) {
	if args.ArgDispatcher == nil {
		return utils.NewErrMandatoryIeMissing("ArgDispatcher")
	}
	if dS.attrS != nil {
		if err = dS.authorize(utils.SchedulerSv1Ping,
			args.CGREvent.Tenant,
			args.APIKey, args.CGREvent.Time); err != nil {
			return
		}
	}
	return dS.Dispatch(args.CGREvent, utils.MetaScheduler, args.RouteID,
		utils.SchedulerSv1Ping, args.CGREvent, reply)
}

func (dS *DispatcherService) SchedulerSv1Reload(args *StringkWithApiKey, reply *string) (err error) {
	if args.ArgDispatcher == nil {
		return utils.NewErrMandatoryIeMissing("ArgDispatcher")
	}
	if dS.attrS != nil {
		if err = dS.authorize(utils.SchedulerSv1Ping,
			args.TenantArg.Tenant,
			args.APIKey, utils.TimePointer(time.Now())); err != nil {
			return
		}
	}
	return dS.Dispatch(&utils.CGREvent{Tenant: args.TenantArg.Tenant}, utils.MetaScheduler,
		args.RouteID, utils.SchedulerSv1Reload, args.Arg, reply)
}
