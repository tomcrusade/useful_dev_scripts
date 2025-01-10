package api_digitalocean

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
)

type DOListSnapshotsResponse struct {
	Snapshots []*entity.DigitaloceanSnapshot      `json:"snapshots,omitempty"`
	Links     *entity.DigitaloceanAPIResultCursor `json:"links,omitempty"`
	Meta      *entity.DigitaloceanAPIResultMeta   `json:"meta,omitempty"`
}

func (do *DigitaloceanAPI) ListSnapshots(limit int, pageNum int, resourceType string) (
	[]*entity.DigitaloceanSnapshot,
	*entity.DigitaloceanAPIResultCursor,
	error,
) {
	response, err := adapters.CallApi[DOListSnapshotsResponse](
		do.getAPICallParams(
			adapters.HttpMethodGet, "/v2/snapshots",
			struct {
				PerPage      int    `json:"per_page,omitempty"`
				Page         int    `json:"page,omitempty"`
				ResourceType string `json:"resource_type,omitempty"`
			}{limit, pageNum, resourceType},
		),
	)

	if err != nil {
		return []*entity.DigitaloceanSnapshot{}, &entity.DigitaloceanAPIResultCursor{}, err
	}

	return response.Snapshots, response.Links, nil
}
