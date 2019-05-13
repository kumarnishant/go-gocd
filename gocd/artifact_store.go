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

func (crs *ArtifactStoreService) Update(ctx context.Context, store *ArtifactStore, eTag string) (out *ArtifactStore, resp *APIResponse, err error) {
	out = &ArtifactStore{}
	request := &APIClientRequest{
		Path:         "admin/artifact_stores/" + store.ID,
		APIVersion:   apiV1,
		RequestBody:  store,
		ResponseBody: out,
		Headers:      map[string]string{"If-Match": eTag},
	}
	_, resp, err = crs.client.putAction(ctx, request)
	return
}

func (crs *ArtifactStoreService) Get(ctx context.Context, storeId string) (store *ArtifactStore, resp *APIResponse, err error) {
	//var contentResponse string
	store = &ArtifactStore{}
	_, resp, err = crs.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/artifact_stores/" + storeId,
		ResponseBody: store,
		//ResponseType: responseTypeText,
		APIVersion: apiV1,
	})

	return
}

func (crs *ArtifactStoreService) GetAll(ctx context.Context) (stores []ArtifactStore, resp *APIResponse, err error) {
	//var contentResponse string
	//var stores []ArtifactStore
	_, resp, err = crs.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/artifact_stores",
		ResponseBody: &stores,
		ResponseType: responseTypeText,
		APIVersion:   apiV1,
	})

	return
}
