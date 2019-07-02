// Copyright (c) 2018-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package migrations

import (
	"github.com/xzl8028/xenia-server/app"
	tjobs "github.com/xzl8028/xenia-server/jobs/interfaces"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
)

const (
	MIGRATION_STATE_UNSCHEDULED = "unscheduled"
	MIGRATION_STATE_IN_PROGRESS = "in_progress"
	MIGRATION_STATE_COMPLETED   = "completed"

	JOB_DATA_KEY_MIGRATION           = "migration_key"
	JOB_DATA_KEY_MIGRATION_LAST_DONE = "last_done"
)

type MigrationsJobInterfaceImpl struct {
	App *app.App
}

func init() {
	app.RegisterJobsMigrationsJobInterface(func(a *app.App) tjobs.MigrationsJobInterface {
		return &MigrationsJobInterfaceImpl{a}
	})
}

func MakeMigrationsList() []string {
	return []string{
		model.MIGRATION_KEY_ADVANCED_PERMISSIONS_PHASE_2,
	}
}

func GetMigrationState(migration string, store store.Store) (string, *model.Job, *model.AppError) {
	if _, err := store.System().GetByName(migration); err == nil {
		return MIGRATION_STATE_COMPLETED, nil, nil
	}

	jobs, err := store.Job().GetAllByType(model.JOB_TYPE_MIGRATIONS)
	if err != nil {
		return "", nil, err
	}

	for _, job := range jobs {
		if key, ok := job.Data[JOB_DATA_KEY_MIGRATION]; ok {
			if key != migration {
				continue
			}

			switch job.Status {
			case model.JOB_STATUS_IN_PROGRESS, model.JOB_STATUS_PENDING:
				return MIGRATION_STATE_IN_PROGRESS, job, nil
			default:
				return MIGRATION_STATE_UNSCHEDULED, job, nil
			}
		}
	}

	return MIGRATION_STATE_UNSCHEDULED, nil, nil
}
