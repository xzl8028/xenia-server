// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"fmt"
	"reflect"

	"github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/utils"
)

const ADVANCED_PERMISSIONS_MIGRATION_KEY = "AdvancedPermissionsMigrationComplete"
const EMOJIS_PERMISSIONS_MIGRATION_KEY = "EmojisPermissionsMigrationComplete"
const GUEST_ROLES_CREATION_MIGRATION_KEY = "GuestRolesCreationMigrationComplete"

// This function migrates the default built in roles from code/config to the database.
func (a *App) DoAdvancedPermissionsMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := a.Srv.Store.System().GetByName(ADVANCED_PERMISSIONS_MIGRATION_KEY); err == nil {
		return
	}

	mlog.Info("Migrating roles to database.")
	roles := model.MakeDefaultRoles()
	roles = utils.SetRolePermissionsFromConfig(roles, a.Config(), a.License() != nil)

	allSucceeded := true

	for _, role := range roles {
		_, err := a.Srv.Store.Role().Save(role)
		if err == nil {
			continue
		}

		// If this failed for reasons other than the role already existing, don't mark the migration as done.
		fetchedRole, err := a.Srv.Store.Role().GetByName(role.Name)
		if err != nil {
			mlog.Critical("Failed to migrate role to database.", mlog.Err(err))
			allSucceeded = false
			continue
		}

		// If the role already existed, check it is the same and update if not.
		if !reflect.DeepEqual(fetchedRole.Permissions, role.Permissions) ||
			fetchedRole.DisplayName != role.DisplayName ||
			fetchedRole.Description != role.Description ||
			fetchedRole.SchemeManaged != role.SchemeManaged {
			role.Id = fetchedRole.Id
			if _, err = a.Srv.Store.Role().Save(role); err != nil {
				// Role is not the same, but failed to update.
				mlog.Critical("Failed to migrate role to database.", mlog.Err(err))
				allSucceeded = false
			}
		}
	}

	if !allSucceeded {
		return
	}

	config := a.Config()
	if *config.ServiceSettings.DEPRECATED_DO_NOT_USE_AllowEditPost == model.ALLOW_EDIT_POST_ALWAYS {
		*config.ServiceSettings.PostEditTimeLimit = -1
		if err := a.SaveConfig(config, true); err != nil {
			mlog.Error("Failed to update config in Advanced Permissions Phase 1 Migration.", mlog.Err(err))
		}
	}

	system := model.System{
		Name:  ADVANCED_PERMISSIONS_MIGRATION_KEY,
		Value: "true",
	}

	if err := a.Srv.Store.System().Save(&system); err != nil {
		mlog.Critical("Failed to mark advanced permissions migration as completed.", mlog.Err(err))
	}
}

func (a *App) SetPhase2PermissionsMigrationStatus(isComplete bool) error {
	if !isComplete {
		if _, err := a.Srv.Store.System().PermanentDeleteByName(model.MIGRATION_KEY_ADVANCED_PERMISSIONS_PHASE_2); err != nil {
			return err
		}
	}
	a.Srv.phase2PermissionsMigrationComplete = isComplete
	return nil
}

func (a *App) DoEmojisPermissionsMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := a.Srv.Store.System().GetByName(EMOJIS_PERMISSIONS_MIGRATION_KEY); err == nil {
		return
	}

	var role *model.Role = nil
	var systemAdminRole *model.Role = nil
	var err *model.AppError = nil

	mlog.Info("Migrating emojis config to database.")
	switch *a.Config().ServiceSettings.DEPRECATED_DO_NOT_USE_RestrictCustomEmojiCreation {
	case model.RESTRICT_EMOJI_CREATION_ALL:
		role, err = a.GetRoleByName(model.SYSTEM_USER_ROLE_ID)
		if err != nil {
			mlog.Critical("Failed to migrate emojis creation permissions from xenia config.", mlog.Err(err))
			return
		}
	case model.RESTRICT_EMOJI_CREATION_ADMIN:
		role, err = a.GetRoleByName(model.TEAM_ADMIN_ROLE_ID)
		if err != nil {
			mlog.Critical("Failed to migrate emojis creation permissions from xenia config.", mlog.Err(err))
			return
		}
	case model.RESTRICT_EMOJI_CREATION_SYSTEM_ADMIN:
		role = nil
	default:
		mlog.Critical("Failed to migrate emojis creation permissions from xenia config. Invalid restrict emoji creation setting")
		return
	}

	if role != nil {
		role.Permissions = append(role.Permissions, model.PERMISSION_CREATE_EMOJIS.Id, model.PERMISSION_DELETE_EMOJIS.Id)
		if _, err = a.Srv.Store.Role().Save(role); err != nil {
			mlog.Critical("Failed to migrate emojis creation permissions from xenia config.", mlog.Err(err))
			return
		}
	}

	systemAdminRole, err = a.GetRoleByName(model.SYSTEM_ADMIN_ROLE_ID)
	if err != nil {
		mlog.Critical("Failed to migrate emojis creation permissions from xenia config.", mlog.Err(err))
		return
	}

	systemAdminRole.Permissions = append(systemAdminRole.Permissions, model.PERMISSION_CREATE_EMOJIS.Id, model.PERMISSION_DELETE_EMOJIS.Id)
	systemAdminRole.Permissions = append(systemAdminRole.Permissions, model.PERMISSION_DELETE_OTHERS_EMOJIS.Id)
	if _, err := a.Srv.Store.Role().Save(systemAdminRole); err != nil {
		mlog.Critical("Failed to migrate emojis creation permissions from xenia config.", mlog.Err(err))
		return
	}

	system := model.System{
		Name:  EMOJIS_PERMISSIONS_MIGRATION_KEY,
		Value: "true",
	}

	if err := a.Srv.Store.System().Save(&system); err != nil {
		mlog.Critical("Failed to mark emojis permissions migration as completed.", mlog.Err(err))
	}
}

func (a *App) DoGuestRolesCreationMigration() {
	// If the migration is already marked as completed, don't do it again.
	if _, err := a.Srv.Store.System().GetByName(GUEST_ROLES_CREATION_MIGRATION_KEY); err == nil {
		return
	}

	roles := model.MakeDefaultRoles()

	allSucceeded := true
	if _, err := a.Srv.Store.Role().GetByName(model.CHANNEL_GUEST_ROLE_ID); err != nil {
		if _, err := a.Srv.Store.Role().Save(roles[model.CHANNEL_GUEST_ROLE_ID]); err != nil {
			mlog.Critical("Failed to create new guest role to database.", mlog.Err(err))
			allSucceeded = false
		}
	}
	if _, err := a.Srv.Store.Role().GetByName(model.TEAM_GUEST_ROLE_ID); err != nil {
		if _, err := a.Srv.Store.Role().Save(roles[model.TEAM_GUEST_ROLE_ID]); err != nil {
			mlog.Critical("Failed to create new guest role to database.", mlog.Err(err))
			allSucceeded = false
		}
	}
	if _, err := a.Srv.Store.Role().GetByName(model.SYSTEM_GUEST_ROLE_ID); err != nil {
		if _, err := a.Srv.Store.Role().Save(roles[model.SYSTEM_GUEST_ROLE_ID]); err != nil {
			mlog.Critical("Failed to create new guest role to database.", mlog.Err(err))
			allSucceeded = false
		}
	}

	resultSchemes := <-a.Srv.Store.Scheme().GetAllPage("", 0, 1000000)
	if resultSchemes.Err != nil {
		mlog.Critical("Failed to get all schemes.", mlog.Err(resultSchemes.Err))
		allSucceeded = false
	}
	schemes := resultSchemes.Data.([]*model.Scheme)
	for _, scheme := range schemes {
		if scheme.DefaultTeamGuestRole == "" || scheme.DefaultChannelGuestRole == "" {
			// Team Guest Role
			teamGuestRole := &model.Role{
				Name:          model.NewId(),
				DisplayName:   fmt.Sprintf("Team Guest Role for Scheme %s", scheme.Name),
				Permissions:   roles[model.TEAM_GUEST_ROLE_ID].Permissions,
				SchemeManaged: true,
			}

			if savedRole, err := a.Srv.Store.Role().Save(teamGuestRole); err != nil {
				mlog.Critical("Failed to create new guest role for custom scheme.", mlog.Err(err))
				allSucceeded = false
			} else {
				scheme.DefaultTeamGuestRole = savedRole.Name
			}

			// Channel Guest Role
			channelGuestRole := &model.Role{
				Name:          model.NewId(),
				DisplayName:   fmt.Sprintf("Channel Guest Role for Scheme %s", scheme.Name),
				Permissions:   roles[model.CHANNEL_GUEST_ROLE_ID].Permissions,
				SchemeManaged: true,
			}

			if savedRole, err := a.Srv.Store.Role().Save(channelGuestRole); err != nil {
				mlog.Critical("Failed to create new guest role for custom scheme.", mlog.Err(err))
				allSucceeded = false
			} else {
				scheme.DefaultChannelGuestRole = savedRole.Name
			}

			result := <-a.Srv.Store.Scheme().Save(scheme)
			if result.Err != nil {
				mlog.Critical("Failed to update custom scheme.", mlog.Err(result.Err))
				allSucceeded = false
			}
		}
	}

	if !allSucceeded {
		return
	}

	system := model.System{
		Name:  GUEST_ROLES_CREATION_MIGRATION_KEY,
		Value: "true",
	}

	if err := a.Srv.Store.System().Save(&system); err != nil {
		mlog.Critical("Failed to mark guest roles creation migration as completed.", mlog.Err(err))
	}
}

func (a *App) DoAppMigrations() {
	a.DoAdvancedPermissionsMigration()
	a.DoEmojisPermissionsMigration()
	a.DoGuestRolesCreationMigration()
	// This migration always must be the last, because can be based on previous
	// migrations. For example, it needs the guest roles migration.
	a.DoPermissionsMigrations()
}
