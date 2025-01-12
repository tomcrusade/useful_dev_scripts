package api_vultr

import (
	"dev_scripts/adapters"
	entity2 "dev_scripts/entity"
	"fmt"
	"sort"
)

type ListSnapshotSortBy string

const (
	ListSnapshotSortByDateCreatedDesc ListSnapshotSortBy = "-date_created"
	ListSnapshotSortByDateCreatedAsc  ListSnapshotSortBy = "+date_created"
)

func (vultr *VultrAPI) ListSnapshots(sortBy ListSnapshotSortBy) (
	[]entity2.VultrSnapshot,
	entity2.VultrAPIMeta,
	error,
) {
	response, err := adapters.CallApi[struct {
		Snapshots []entity2.VultrSnapshot `json:"snapshots"`
		Meta      entity2.VultrAPIMeta    `json:"meta"`
	}](vultr.getAPICallParams(adapters.HttpMethodGet, "/snapshots", struct{}{}))
	return sortSnapshots(response.Snapshots, sortBy), response.Meta, err
}

func sortSnapshots(
	snapshots []entity2.VultrSnapshot,
	sortBy ListSnapshotSortBy,
) []entity2.VultrSnapshot {
	newSnapshots := snapshots
	sort.SliceStable(
		newSnapshots,
		func(i, j int) bool {
			if sortBy == ListSnapshotSortByDateCreatedDesc {
				prevTime, _ := entity2.ConvertVultrTimeToDate(snapshots[i].DateCreated)
				currentTime, _ := entity2.ConvertVultrTimeToDate(snapshots[j].DateCreated)
				return prevTime.After(currentTime)
			} else if sortBy == ListSnapshotSortByDateCreatedAsc {
				prevTime, _ := entity2.ConvertVultrTimeToDate(snapshots[i].DateCreated)
				currentTime, _ := entity2.ConvertVultrTimeToDate(snapshots[j].DateCreated)
				return prevTime.Before(currentTime)
			}
			return true
		},
	)
	return newSnapshots
}

func (vultr *VultrAPI) GetSnapshot(snapshotId string) (entity2.VultrSnapshot, error) {
	response, err := adapters.CallApi[struct {
		Snapshot entity2.VultrSnapshot `json:"snapshot"`
	}](
		vultr.getAPICallParams(
			adapters.HttpMethodGet,
			fmt.Sprintf("/snapshots/%v", snapshotId),
			struct{}{},
		),
	)
	return response.Snapshot, err
}

func (vultr *VultrAPI) RemoveSnapshot(snapshotId string) error {
	_, err := adapters.CallApi[struct{}](
		vultr.getAPICallParams(
			adapters.HttpMethodDelete,
			fmt.Sprintf("/snapshots/%v", snapshotId),
			struct{}{},
		),
	)
	return err
}

func (vultr *VultrAPI) CreateSnapshot(instanceId string, description string) (
	entity2.VultrSnapshot,
	error,
) {
	var param = struct {
		InstanceId  string `json:"instance_id"`
		Description string `json:"description"`
	}{InstanceId: instanceId, Description: description}
	response, err := adapters.CallApi[struct {
		Snapshot entity2.VultrSnapshot `json:"snapshot"`
	}](vultr.getAPICallParams(adapters.HttpMethodPost, "/snapshots", param))
	return response.Snapshot, err
}
