package gocd

import "context"

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
	apiVersion, err := ps.client.getAPIVersion(ctx, "admin/scms")
	if err != nil {
		return nil, nil, err
	}
	pr := ScmListResp{}
	_, resp, err := ps.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/scms",
		ResponseBody: &pr,
		APIVersion:   apiVersion,
	})

	return &pr, resp, err
}
