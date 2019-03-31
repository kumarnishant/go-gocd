package gocd

import (
	"context"
	"strconv"
)

type GetPipelineService service

type PipelineHistoryResponse struct {
	Pagination Pagination	`json:"pagination"`
	Pipelines []*PipelineInstanceResponse `json:"pipelines"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Total int `json:"total"`
	PageSize int `json:"page_size"`
}

// PipelineInstance describes a single pipeline run
// codebeat:disable[TOO_MANY_IVARS]
type PipelineInstanceResponse struct {
	BuildCause          BuildCause       `json:"build_cause"`
	Label               string           `json:"label"`
	Counter             int              `json:"counter"`
	PreparingToSchedule bool             `json:"preparing_to_schedule"`
	CanRun              bool             `json:"can_run"`
	Name                string           `json:"name"`
	NaturalOrder        float32          `json:"natural_order"`
	Comment             string           `json:"comment"`
	Stages              []*StageResponse `json:"stages"`
}

// Stage represents a GoCD Stage object.
// codebeat:disable[TOO_MANY_IVARS]
type StageResponse struct {
	Name              string `json:"name"`
	Jobs              []*Job `json:"jobs,omitempty"`
	Result            string `json:"result"`
	ApprovedBy        string `json:"approved_by"`
	ApprovalType      string `json:"approval_type"`
	CanRun            bool `json:"can_run"`
	Counter           string `json:"counter"`
	Scheduled         bool   `json:"scheduled"`
	ID                int    `json:"id"`
	OperatePermission bool   `json:"operate_permission"`
	RerunOfCounter    int    `json:"rerun_of_counter"`
	CancelledBy       string `json:"cancelled_by"`
}

type JobResponse struct {
	Name          string `json:"name"`
	ScheduledDate int    `json:"scheduled_date,omitempty"`
	Result        string `json:"result,omitempty"`
	State         string `json:"state,omitempty"`
	ID            int    `json:"id,omitempty"`
}

func (ps *GetPipelineService) Get(ctx context.Context, pipelineName string, instanceId int) (*PipelineInstanceResponse, *APIResponse, error) {
	pipeline := &PipelineInstanceResponse{}

	_, resp, err := ps.client.getAction(ctx, &APIClientRequest{
		Path:         "pipelines/" + pipelineName + "/instance/" + strconv.Itoa(instanceId),
		ResponseBody: pipeline,
	})
	return pipeline, resp, err
}

func (ps *GetPipelineService) GetHistory(ctx context.Context, pipelineName string, offset int) (*PipelineHistoryResponse, *APIResponse, error) {
	pipeline := &PipelineHistoryResponse{}

	_, resp, err := ps.client.getAction(ctx, &APIClientRequest{
		Path:         "pipelines/" + pipelineName + "/history/" + strconv.Itoa(offset),
		ResponseBody: pipeline,
	})
	return pipeline, resp, err
}