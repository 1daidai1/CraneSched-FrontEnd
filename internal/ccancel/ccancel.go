/**
 * Copyright (c) 2024 Peking University and Peking University
 * Changsha Institute for Computing and Digital Economy
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package ccancel

import (
	"context"
	"fmt"
	"os"

	"github.com/1daidai1/CraneSched-FrontEnd/generated/protos"
	"github.com/1daidai1/CraneSched-FrontEnd/pkg/util"

	log "github.com/sirupsen/logrus"
)

var (
	stub protos.CraneCtldClient
)

func CancelTask(args []string) util.CraneCmdError {
	req := &protos.CancelTaskRequest{
		OperatorUid: uint32(os.Getuid()),

		FilterPartition: FlagPartition,
		FilterAccount:   FlagAccount,
		FilterTaskName:  FlagJobName,
		FilterState:     protos.TaskStatus_Invalid,
		FilterUsername:  FlagUserName,
	}

	if len(args) > 0 {
		taskIds, err := util.ParseJobIdList(args[0], ",")
		if err != nil {
			log.Errorf("Invalid job list specified: %v.\n", err)
			return util.ErrorCmdArg
		}
		req.FilterTaskIds = taskIds
	}

	if FlagState != "" {
		stateList, err := util.ParseInRamTaskStatusList(FlagState)
		if err != nil {
			log.Errorln(err)
			return util.ErrorCmdArg
		}
		if len(stateList) == 1 {
			req.FilterState = stateList[0]
		}
	}

	req.FilterNodes = FlagNodes

	reply, err := stub.CancelTask(context.Background(), req)
	if err != nil {
		util.GrpcErrorPrintf(err, "Failed to cancel tasks")
		os.Exit(util.ErrorNetwork)
	}

	if FlagJson {
		fmt.Println(util.FmtJson.FormatReply(reply))
		if len(reply.NotCancelledTasks) > 0 {
			return util.ErrorBackend
		} else {
			return util.ErrorSuccess
		}
	}

	if len(reply.CancelledTasks) > 0 {
		fmt.Printf("Jobs %v cancelled successfully.\n", reply.CancelledTasks)
	}

	if len(reply.NotCancelledTasks) > 0 {
		for i := 0; i < len(reply.NotCancelledTasks); i++ {
			log.Errorf("Failed to cancel job: %d. Reason: %s.\n", reply.NotCancelledTasks[i], reply.NotCancelledReasons[i])
		}
		os.Exit(util.ErrorBackend)
	}
	return util.ErrorSuccess
}
