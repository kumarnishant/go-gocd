package gocd

import (
	"context"
	"fmt"
)

type ScmService service

type ScmConfig struct {
	Links          *HALLinks     `json:"_links"`
	Id             string        `json:"id,omitempty"`
	Name           string        `json:"name,omitempty"`
	AutoUpdate     bool          `json:"auto_update,omitempty"`
	PluginMetadata PluginMetaDta `json:"plugin_metadata,omitempty"`
	configuration  []*Property   `json:"configuration,omitempty"`
}

type PluginMetaDta struct {
	Id      string `json:"id,omitempty"`
	Version string `json:"version,omitempty"`
}

type ScmListResp struct {
	Links    *HALLinks `json:"_links"`
	Embedded struct {
		Scms []*ScmConfig `json:"scms"`
	} `json:"_embedded"`
}

func (ps *ScmService) List(ctx context.Context) (*ScmListResp, *APIResponse, error) {
	pr := ScmListResp{}
	_, resp, err := ps.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/scms",
		ResponseBody: &pr,
		APIVersion:   apiV1,
	})

	return &pr, resp, err
}

func (ps *ScmService) Get(ctx context.Context, name string) (*ScmConfig, *APIResponse, error) {
	p := &ScmConfig{}
	_, resp, err := ps.client.getAction(ctx, &APIClientRequest{
		Path:         fmt.Sprintf("admin/scms/%s", name),
		ResponseBody: p,
		APIVersion:   apiV1,
	})
	return p, resp, err
}

func (ps *ScmService) Create(ctx context.Context,  scm *ScmConfig) (scmRes *ScmConfig, resp *APIResponse, err error) {
	scmRes = &ScmConfig{}
	_, resp, err = ps.client.postAction(ctx, &APIClientRequest{
		Path:       "admin/scms",
		APIVersion: apiV1,
		RequestBody: scm,
		ResponseBody: scmRes,
	})
	return
}