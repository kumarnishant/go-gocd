package gocd

import (
	"context"
)

type DashboardService service

type Dashboard struct {
	Links             *HALLinks          `json:"_links"`
	Personalization   string             `json:"_personalization"`
	EmbeddedDashboard *EmbeddedDashboard `json:"_embedded"`
}

type EmbeddedDashboard struct {
	PipelineGroups []*EmbeddedPipelineGroup `json:"pipeline_groups"`
	Environments   []*EmbeddedEnvironment   `json:"environments"`
	Pipelines      []*EmbeddedPipeline      `json:"pipelines"`
}

type EmbeddedPipelineGroup struct {
	Links         *HALLinks `json:"_links,omitempty"`
	Name          string    `json:"name"`
	Pipelines     []string  `json:"pipelines"`
	CanAdminister bool      `json:"can_administer"`
}

type EmbeddedEnvironment struct {
	Links         *HALLinks `json:"_links,omitempty"`
	Name          string    `json:"name"`
	Pipelines     []string  `json:"pipelines"`
	CanAdminister bool      `json:"can_administer"`
}

type EmbeddedPipeline struct {
	Links               *HALLinks         `json:"_links,omitempty"`
	Name                string            `json:"name"`
	LastUpdateTimeStamp int64             `json:"last_update_time_stamp"`
	Locked              bool              `json:"locked"`
	PauseInfo           PauseInfo         `json:"pause_info"`
	CanOperate          bool              `json:"can_operate"`
	CanAdminister       bool              `json:"can_administer"`
	CanUnlock           bool              `json:"can_unlock"`
	CanPause            bool              `json:"can_pause"`
	FromConfigRepo      bool              `json:"from_config_repo"`
	EmbeddedInstances   *EmbeddedInstances `json:"_embedded"`
}

type EmbeddedInstances struct {
	Instances []*EmbeddedInstance `json:"instances"`
}

type EmbeddedInstance struct {
	Links          *HALLinks      `json:"_links,omitempty"`
	Label          string         `json:"label"`
	Counter        int            `json:"counter"`
	TriggeredBy    string         `json:"triggered_by"`
	ScheduledAt    string         `json:"scheduled_at"`
	EmbeddedStages *EmbeddedStages `json:"_embedded"`
}

type EmbeddedStages struct {
	Stages []*EmbeddedStage `json:"stages"`
}

type EmbeddedStage struct {
	Links       *HALLinks `json:"_links,omitempty"`
	Name        string    `json:"name"`
	Counter     int       `json:"counter"`
	ScheduledAt string    `json:"scheduled_at"`
	Status      string    `json:"status"`
	ApprovedBy  string    `json:"approved_by"`
}

type PauseInfo struct {
	Paused      bool   `json:"paused"`
	PausedBy    string `json:"paused_by"`
	PauseReason string `json:"pause_reason"`
}

func (ds *DashboardService) Get(ctx context.Context) (*Dashboard, *APIResponse, error) {
	var dashboard Dashboard

	_, resp, err := ds.client.getAction(ctx, &APIClientRequest{
		Path:         "dashboard?viewName=Default",
		APIVersion:   apiV3,
		ResponseBody: &dashboard,
	})
	return &dashboard, resp, err
}
