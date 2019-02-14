package gocd

import "context"

type Property struct {
	Key            string `json:"key"`
	Value          string `json:"value,omitempty"`
	EncryptedValue string `json:"encrypted_value,omitempty"`
}
type ArtifactStore struct {
	Links      *HALLinks   `json:"_links,omitempty"`
	PluginId   string      `json:"plugin_id"`
	Properties []*Property `json:"properties"`
	ID         string      `json:"id"`
}

type ArtifactStoreService service

func (crs *ArtifactStoreService) Create(ctx context.Context, store *ArtifactStore) (out *ArtifactStore, resp *APIResponse, err error) {
	out = &ArtifactStore{}
	_, resp, err = crs.client.postAction(ctx, &APIClientRequest{
		Path:         "admin/artifact_stores",
		RequestBody:  store,
		ResponseBody: out,
		APIVersion:   apiV1,
	})
	return
}
