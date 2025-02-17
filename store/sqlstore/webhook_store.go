// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package sqlstore

import (
	"database/sql"
	"net/http"

	"github.com/xzl8028/xenia-server/einterfaces"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
	"github.com/xzl8028/xenia-server/utils"
)

type SqlWebhookStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

const (
	WEBHOOK_CACHE_SIZE = 25000
	WEBHOOK_CACHE_SEC  = 900 // 15 minutes
)

var webhookCache = utils.NewLru(WEBHOOK_CACHE_SIZE)

func (s SqlWebhookStore) ClearCaches() {
	webhookCache.Purge()

	if s.metrics != nil {
		s.metrics.IncrementMemCacheInvalidationCounter("Webhook - Purge")
	}
}

func NewSqlWebhookStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.WebhookStore {
	s := &SqlWebhookStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.IncomingWebhook{}, "IncomingWebhooks").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(26)
		table.ColMap("UserId").SetMaxSize(26)
		table.ColMap("ChannelId").SetMaxSize(26)
		table.ColMap("TeamId").SetMaxSize(26)
		table.ColMap("DisplayName").SetMaxSize(64)
		table.ColMap("Description").SetMaxSize(500)

		tableo := db.AddTableWithName(model.OutgoingWebhook{}, "OutgoingWebhooks").SetKeys(false, "Id")
		tableo.ColMap("Id").SetMaxSize(26)
		tableo.ColMap("Token").SetMaxSize(26)
		tableo.ColMap("CreatorId").SetMaxSize(26)
		tableo.ColMap("ChannelId").SetMaxSize(26)
		tableo.ColMap("TeamId").SetMaxSize(26)
		tableo.ColMap("TriggerWords").SetMaxSize(1024)
		tableo.ColMap("CallbackURLs").SetMaxSize(1024)
		tableo.ColMap("DisplayName").SetMaxSize(64)
		tableo.ColMap("Description").SetMaxSize(500)
		tableo.ColMap("ContentType").SetMaxSize(128)
		tableo.ColMap("TriggerWhen").SetMaxSize(1)
		tableo.ColMap("Username").SetMaxSize(64)
		tableo.ColMap("IconURL").SetMaxSize(1024)
	}

	return s
}

func (s SqlWebhookStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_incoming_webhook_user_id", "IncomingWebhooks", "UserId")
	s.CreateIndexIfNotExists("idx_incoming_webhook_team_id", "IncomingWebhooks", "TeamId")
	s.CreateIndexIfNotExists("idx_outgoing_webhook_team_id", "OutgoingWebhooks", "TeamId")

	s.CreateIndexIfNotExists("idx_incoming_webhook_update_at", "IncomingWebhooks", "UpdateAt")
	s.CreateIndexIfNotExists("idx_incoming_webhook_create_at", "IncomingWebhooks", "CreateAt")
	s.CreateIndexIfNotExists("idx_incoming_webhook_delete_at", "IncomingWebhooks", "DeleteAt")

	s.CreateIndexIfNotExists("idx_outgoing_webhook_update_at", "OutgoingWebhooks", "UpdateAt")
	s.CreateIndexIfNotExists("idx_outgoing_webhook_create_at", "OutgoingWebhooks", "CreateAt")
	s.CreateIndexIfNotExists("idx_outgoing_webhook_delete_at", "OutgoingWebhooks", "DeleteAt")
}

func (s SqlWebhookStore) InvalidateWebhookCache(webhookId string) {
	webhookCache.Remove(webhookId)
	if s.metrics != nil {
		s.metrics.IncrementMemCacheInvalidationCounter("Webhook - Remove by WebhookId")
	}
}

func (s SqlWebhookStore) SaveIncoming(webhook *model.IncomingWebhook) (*model.IncomingWebhook, *model.AppError) {

	if len(webhook.Id) > 0 {
		return nil, model.NewAppError("SqlWebhookStore.SaveIncoming", "store.sql_webhooks.save_incoming.existing.app_error", nil, "id="+webhook.Id, http.StatusBadRequest)
	}

	webhook.PreSave()
	if err := webhook.IsValid(); err != nil {
		return nil, err
	}

	if err := s.GetMaster().Insert(webhook); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.SaveIncoming", "store.sql_webhooks.save_incoming.app_error", nil, "id="+webhook.Id+", "+err.Error(), http.StatusInternalServerError)
	}

	return webhook, nil

}

func (s SqlWebhookStore) UpdateIncoming(hook *model.IncomingWebhook) (*model.IncomingWebhook, *model.AppError) {
	hook.UpdateAt = model.GetMillis()

	if _, err := s.GetMaster().Update(hook); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.UpdateIncoming", "store.sql_webhooks.update_incoming.app_error", nil, "id="+hook.Id+", "+err.Error(), http.StatusInternalServerError)
	}
	return hook, nil
}

func (s SqlWebhookStore) GetIncoming(id string, allowFromCache bool) (*model.IncomingWebhook, *model.AppError) {
	if allowFromCache {
		if cacheItem, ok := webhookCache.Get(id); ok {
			if s.metrics != nil {
				s.metrics.IncrementMemCacheHitCounter("Webhook")
			}
			return cacheItem.(*model.IncomingWebhook), nil
		}
		if s.metrics != nil {
			s.metrics.IncrementMemCacheMissCounter("Webhook")
		}
	}

	var webhook model.IncomingWebhook
	if err := s.GetReplica().SelectOne(&webhook, "SELECT * FROM IncomingWebhooks WHERE Id = :Id AND DeleteAt = 0", map[string]interface{}{"Id": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, model.NewAppError("SqlWebhookStore.GetIncoming", "store.sql_webhooks.get_incoming.app_error", nil, "id="+id+", err="+err.Error(), http.StatusNotFound)
		}
		return nil, model.NewAppError("SqlWebhookStore.GetIncoming", "store.sql_webhooks.get_incoming.app_error", nil, "id="+id+", err="+err.Error(), http.StatusInternalServerError)
	}

	webhookCache.AddWithExpiresInSecs(id, &webhook, WEBHOOK_CACHE_SEC)

	return &webhook, nil
}

func (s SqlWebhookStore) DeleteIncoming(webhookId string, time int64) *model.AppError {
	_, err := s.GetMaster().Exec("Update IncomingWebhooks SET DeleteAt = :DeleteAt, UpdateAt = :UpdateAt WHERE Id = :Id", map[string]interface{}{"DeleteAt": time, "UpdateAt": time, "Id": webhookId})
	if err != nil {
		return model.NewAppError("SqlWebhookStore.DeleteIncoming", "store.sql_webhooks.delete_incoming.app_error", nil, "id="+webhookId+", err="+err.Error(), http.StatusInternalServerError)
	}

	s.InvalidateWebhookCache(webhookId)
	return nil
}

func (s SqlWebhookStore) PermanentDeleteIncomingByUser(userId string) *model.AppError {
	_, err := s.GetMaster().Exec("DELETE FROM IncomingWebhooks WHERE UserId = :UserId", map[string]interface{}{"UserId": userId})
	if err != nil {
		return model.NewAppError("SqlWebhookStore.DeleteIncomingByUser", "store.sql_webhooks.permanent_delete_incoming_by_user.app_error", nil, "id="+userId+", err="+err.Error(), http.StatusInternalServerError)
	}

	s.ClearCaches()
	return nil
}

func (s SqlWebhookStore) PermanentDeleteIncomingByChannel(channelId string) *model.AppError {
	_, err := s.GetMaster().Exec("DELETE FROM IncomingWebhooks WHERE ChannelId = :ChannelId", map[string]interface{}{"ChannelId": channelId})
	if err != nil {
		return model.NewAppError("SqlWebhookStore.DeleteIncomingByChannel", "store.sql_webhooks.permanent_delete_incoming_by_channel.app_error", nil, "id="+channelId+", err="+err.Error(), http.StatusInternalServerError)
	}

	s.ClearCaches()

	return nil
}

func (s SqlWebhookStore) GetIncomingList(offset, limit int) ([]*model.IncomingWebhook, *model.AppError) {
	var webhooks []*model.IncomingWebhook

	if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM IncomingWebhooks WHERE DeleteAt = 0 LIMIT :Limit OFFSET :Offset", map[string]interface{}{"Limit": limit, "Offset": offset}); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.GetIncomingList", "store.sql_webhooks.get_incoming_by_user.app_error", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	return webhooks, nil

}

func (s SqlWebhookStore) GetIncomingByTeam(teamId string, offset, limit int) ([]*model.IncomingWebhook, *model.AppError) {
	var webhooks []*model.IncomingWebhook

	if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM IncomingWebhooks WHERE TeamId = :TeamId AND DeleteAt = 0 LIMIT :Limit OFFSET :Offset", map[string]interface{}{"TeamId": teamId, "Limit": limit, "Offset": offset}); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.GetIncomingByUser", "store.sql_webhooks.get_incoming_by_user.app_error", nil, "teamId="+teamId+", err="+err.Error(), http.StatusInternalServerError)
	}

	return webhooks, nil
}

func (s SqlWebhookStore) GetIncomingByChannel(channelId string) ([]*model.IncomingWebhook, *model.AppError) {
	var webhooks []*model.IncomingWebhook

	if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM IncomingWebhooks WHERE ChannelId = :ChannelId AND DeleteAt = 0", map[string]interface{}{"ChannelId": channelId}); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.GetIncomingByChannel", "store.sql_webhooks.get_incoming_by_channel.app_error", nil, "channelId="+channelId+", err="+err.Error(), http.StatusInternalServerError)
	}

	return webhooks, nil
}

func (s SqlWebhookStore) SaveOutgoing(webhook *model.OutgoingWebhook) (*model.OutgoingWebhook, *model.AppError) {
	if len(webhook.Id) > 0 {
		return nil, model.NewAppError("SqlWebhookStore.SaveOutgoing", "store.sql_webhooks.save_outgoing.override.app_error", nil, "id="+webhook.Id, http.StatusBadRequest)
	}

	webhook.PreSave()
	if err := webhook.IsValid(); err != nil {
		return nil, err
	}

	if err := s.GetMaster().Insert(webhook); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.SaveOutgoing", "store.sql_webhooks.save_outgoing.app_error", nil, "id="+webhook.Id+", "+err.Error(), http.StatusInternalServerError)
	}

	return webhook, nil
}

func (s SqlWebhookStore) GetOutgoing(id string) (*model.OutgoingWebhook, *model.AppError) {

	var webhook model.OutgoingWebhook

	if err := s.GetReplica().SelectOne(&webhook, "SELECT * FROM OutgoingWebhooks WHERE Id = :Id AND DeleteAt = 0", map[string]interface{}{"Id": id}); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.GetOutgoing", "store.sql_webhooks.get_outgoing.app_error", nil, "id="+id+", err="+err.Error(), http.StatusInternalServerError)
	}

	return &webhook, nil
}

func (s SqlWebhookStore) GetOutgoingList(offset, limit int) ([]*model.OutgoingWebhook, *model.AppError) {
	var webhooks []*model.OutgoingWebhook

	if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM OutgoingWebhooks WHERE DeleteAt = 0 LIMIT :Limit OFFSET :Offset", map[string]interface{}{"Offset": offset, "Limit": limit}); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.GetOutgoingList", "store.sql_webhooks.get_outgoing_by_channel.app_error", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	return webhooks, nil
}

func (s SqlWebhookStore) GetOutgoingByChannel(channelId string, offset, limit int) ([]*model.OutgoingWebhook, *model.AppError) {
	var webhooks []*model.OutgoingWebhook

	query := ""
	if limit < 0 || offset < 0 {
		query = "SELECT * FROM OutgoingWebhooks WHERE ChannelId = :ChannelId AND DeleteAt = 0"
	} else {
		query = "SELECT * FROM OutgoingWebhooks WHERE ChannelId = :ChannelId AND DeleteAt = 0 LIMIT :Limit OFFSET :Offset"
	}

	if _, err := s.GetReplica().Select(&webhooks, query, map[string]interface{}{"ChannelId": channelId, "Offset": offset, "Limit": limit}); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.GetOutgoingByChannel", "store.sql_webhooks.get_outgoing_by_channel.app_error", nil, "channelId="+channelId+", err="+err.Error(), http.StatusInternalServerError)
	}

	return webhooks, nil
}

func (s SqlWebhookStore) GetOutgoingByTeam(teamId string, offset, limit int) ([]*model.OutgoingWebhook, *model.AppError) {
	var webhooks []*model.OutgoingWebhook

	query := ""
	if limit < 0 || offset < 0 {
		query = "SELECT * FROM OutgoingWebhooks WHERE TeamId = :TeamId AND DeleteAt = 0"
	} else {
		query = "SELECT * FROM OutgoingWebhooks WHERE TeamId = :TeamId AND DeleteAt = 0 LIMIT :Limit OFFSET :Offset"
	}

	if _, err := s.GetReplica().Select(&webhooks, query, map[string]interface{}{"TeamId": teamId, "Offset": offset, "Limit": limit}); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.GetOutgoingByTeam", "store.sql_webhooks.get_outgoing_by_team.app_error", nil, "teamId="+teamId+", err="+err.Error(), http.StatusInternalServerError)
	}

	return webhooks, nil
}

func (s SqlWebhookStore) DeleteOutgoing(webhookId string, time int64) *model.AppError {
	_, err := s.GetMaster().Exec("Update OutgoingWebhooks SET DeleteAt = :DeleteAt, UpdateAt = :UpdateAt WHERE Id = :Id", map[string]interface{}{"DeleteAt": time, "UpdateAt": time, "Id": webhookId})
	if err != nil {
		return model.NewAppError("SqlWebhookStore.DeleteOutgoing", "store.sql_webhooks.delete_outgoing.app_error", nil, "id="+webhookId+", err="+err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (s SqlWebhookStore) PermanentDeleteOutgoingByUser(userId string) *model.AppError {
	_, err := s.GetMaster().Exec("DELETE FROM OutgoingWebhooks WHERE CreatorId = :UserId", map[string]interface{}{"UserId": userId})
	if err != nil {
		return model.NewAppError("SqlWebhookStore.DeleteOutgoingByUser", "store.sql_webhooks.permanent_delete_outgoing_by_user.app_error", nil, "id="+userId+", err="+err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (s SqlWebhookStore) PermanentDeleteOutgoingByChannel(channelId string) *model.AppError {
	_, err := s.GetMaster().Exec("DELETE FROM OutgoingWebhooks WHERE ChannelId = :ChannelId", map[string]interface{}{"ChannelId": channelId})
	if err != nil {
		return model.NewAppError("SqlWebhookStore.DeleteOutgoingByChannel", "store.sql_webhooks.permanent_delete_outgoing_by_channel.app_error", nil, "id="+channelId+", err="+err.Error(), http.StatusInternalServerError)
	}

	s.ClearCaches()

	return nil
}

func (s SqlWebhookStore) UpdateOutgoing(hook *model.OutgoingWebhook) (*model.OutgoingWebhook, *model.AppError) {
	hook.UpdateAt = model.GetMillis()

	if _, err := s.GetMaster().Update(hook); err != nil {
		return nil, model.NewAppError("SqlWebhookStore.UpdateOutgoing", "store.sql_webhooks.update_outgoing.app_error", nil, "id="+hook.Id+", "+err.Error(), http.StatusInternalServerError)
	}

	return hook, nil
}

func (s SqlWebhookStore) AnalyticsIncomingCount(teamId string) (int64, *model.AppError) {
	query :=
		`SELECT 
			COUNT(*)
		FROM
			IncomingWebhooks
		WHERE
			DeleteAt = 0`

	if len(teamId) > 0 {
		query += " AND TeamId = :TeamId"
	}

	v, err := s.GetReplica().SelectInt(query, map[string]interface{}{"TeamId": teamId})
	if err != nil {
		return 0, model.NewAppError("SqlWebhookStore.AnalyticsIncomingCount", "store.sql_webhooks.analytics_incoming_count.app_error", nil, "team_id="+teamId+", err="+err.Error(), http.StatusInternalServerError)
	}

	return v, nil
}

func (s SqlWebhookStore) AnalyticsOutgoingCount(teamId string) (int64, *model.AppError) {
	query :=
		`SELECT 
			COUNT(*)
		FROM
			OutgoingWebhooks
		WHERE
			DeleteAt = 0`

	if len(teamId) > 0 {
		query += " AND TeamId = :TeamId"
	}

	v, err := s.GetReplica().SelectInt(query, map[string]interface{}{"TeamId": teamId})
	if err != nil {
		return 0, model.NewAppError("SqlWebhookStore.AnalyticsOutgoingCount", "store.sql_webhooks.analytics_outgoing_count.app_error", nil, "team_id="+teamId+", err="+err.Error(), http.StatusInternalServerError)
	}

	return v, nil
}
