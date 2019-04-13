package gocd

import (
	"context"
	"fmt"
)

type ScmService service

type ScmConfig struct {
	Links          *HALLinks      `json:"_links,omitempty"`
	Id             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	AutoUpdate     bool           `json:"auto_update,omitempty"`
	PluginMetadata PluginMetadata `json:"plugin_metadata,omitempty"`
	Configuration  []*Property    `json:"configuration,omitempty"`
	Version        string         `json:"version,omitempty"` // Version corresponds to the ETag header used when updating a pipeline config
}

type PluginMetadata struct {
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

func (ps *ScmService) Create(ctx context.Context, scm *ScmConfig) (scmRes *ScmConfig, resp *APIResponse, err error) {
	scmRes = &ScmConfig{}
	_, resp, err = ps.client.postAction(ctx, &APIClientRequest{
		Path:         "admin/scms",
		APIVersion:   apiV1,
		RequestBody:  scm,
		ResponseBody: scmRes,
	})
	return
}

// Update an scm object in the GoCD API.
func (ps *ScmService) Update(ctx context.Context, scm *ScmConfig) (ptr *ScmConfig, resp *APIResponse, err error) {
	ptr = &ScmConfig{}
	_, resp, err = ps.client.putAction(ctx, &APIClientRequest{
		Path:         "admin/scms/" + scm.Name,
		APIVersion:   apiV1,
		RequestBody:  scm,
		ResponseBody: ptr,
	})
	return
}

func (p *ScmConfig) SetVersion(version string) {
	p.Version = version
}

// GetVersion retrieves a version string for this pipeline
func (p *ScmConfig) GetVersion() (version string) {
	return p.Version
}
