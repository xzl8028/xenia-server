// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package einterfaces

import "github.com/xzl8028/xenia-server/model"

type AccountMigrationInterface interface {
	MigrateToLdap(fromAuthService string, forignUserFieldNameToMatch string, force bool, dryRun bool) *model.AppError
	MigrateToSaml(fromAuthService string, usersMap map[string]string, auto bool, dryRun bool) *model.AppError
}
