package gocd

import "context"

// PipelineGroupsService describes the HAL _link resource for the api response object for a pipeline group response.
type PipelineGroupsService service

// PipelineGroups represents a collection of pipeline groups
type PipelineGroups []*PipelineGroup

// PipelineGroup describes a pipeline group API response.
type PipelineGroup struct {
	Name      string      `json:"name"`
	Pipelines []*Pipeline `json:"pipelines"`
}

// List Pipeline groups
func (pgs *PipelineGroupsService) List(ctx context.Context, name string) (*PipelineGroups, *APIResponse, error) {

	pg := []*PipelineGroup{}
	_, resp, err := pgs.client.getAction(ctx, &APIClientRequest{
		Path:         "config/pipeline_groups",
		ResponseType: responseTypeJSON,
		ResponseBody: &pg,
	})

	filtered := PipelineGroups{}
	if name != "" && err == nil {
		for _, pipelineGroup := range pg {
			if pipelineGroup.Name == name {
				filtered = append(filtered, pipelineGroup)
			}
		}
	} else {
		filtered = pg
	}

	return &filtered, resp, err
}

func (pgs *PipelineGroupsService) Create(ctx context.Context, name string) (*PipelineGroup, *APIResponse, error) {
	pg := &PipelineGroup{}
	_, resp, err := pgs.client.postAction(ctx, &APIClientRequest{
		Path:         "admin/pipeline_groups",
		RequestBody:  PipelineGroup{Name: name},
		ResponseBody: pg,
		APIVersion:   apiV1,
		ResponseType: responseTypeJSON,
	})

	return pg, resp, err
}

/**
@author: vikram
it will fetch materials of the pipeline from any pipeline group
 */
func (pgs *PipelineGroupsService) ListPipelineMaterial(ctx context.Context, name string, pipelineName string) ([] Material, *APIResponse, error) {

	pg := []*PipelineGroup{}
	_, resp, err := pgs.client.getAction(ctx, &APIClientRequest{
		Path:         "config/pipeline_groups",
		ResponseType: responseTypeJSON,
		ResponseBody: &pg,
	})

	materials := [] Material{}
	if name != "" && err == nil {
		for _, pipelineGroup := range pg {
			if pipelineGroup.Name == name {
				for _,obj := range pipelineGroup.Pipelines {
					if obj.Name == pipelineName {
						materials = append(obj.Materials)
					}
				}
			}
		}
	}

	return materials, resp, err
}