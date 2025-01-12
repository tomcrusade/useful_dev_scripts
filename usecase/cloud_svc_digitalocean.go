package usecase

import (
	"dev_scripts/entity"
	"dev_scripts/repository/api_digitalocean"
	"errors"
	"fmt"
	"strings"
	"time"
)

var instanceCheckMaxCount = 120

type CloudSvcDigitalocean struct {
	env *entity.EnvCloudServer
	api *api_digitalocean.DigitaloceanAPI
}

func NewCloudSvcDigitalocean(env *entity.EnvCloudServer, tokenEnv *entity.EnvResourceToken) *CloudSvcDigitalocean {
	return &CloudSvcDigitalocean{env, api_digitalocean.NewDigitaloceanAPI(env, tokenEnv)}
}

// --

func (uc *CloudSvcDigitalocean) StartInstance() (*entity.DigitaloceanDroplet, error) {
	var chosenCloudVM *entity.DigitaloceanDroplet
	var err error

	if chosenCloudVM, err = uc.getRelatedVM(); err != nil {
		return nil, err
	}

	if chosenCloudVM == nil {
		var backupPolicy *api_digitalocean.DOCreateDropletRequestBackupPolicy
		hasBackupPlan := uc.env.VmBackupPlan != nil
		if hasBackupPlan {
			backupPolicy = &api_digitalocean.DOCreateDropletRequestBackupPolicy{
				Plan:    strings.TrimSpace(uc.env.VmBackupPlan.Plan),
				Weekday: strings.TrimSpace(uc.env.VmBackupPlan.Weekday),
				Hour:    uc.env.VmBackupPlan.Hour,
			}
		}
		config := api_digitalocean.DOCreateDropletRequest{
			Name:         strings.TrimSpace(uc.env.VmLabel),
			Region:       strings.TrimSpace(uc.env.VmRegion),
			Size:         strings.TrimSpace(uc.env.VmResourcePlan),
			Backups:      hasBackupPlan,
			BackupPolicy: backupPolicy,
			Monitoring:   true,
			SSHKeys:      []string{uc.env.SSHKey},
			Tags:         []string{strings.TrimSpace(uc.env.VmLabel)},
		}
		if imageID, err := uc.getImageID(); err != nil {
			return nil, err
		} else {
			config.Image = imageID
		}

		if chosenCloudVM, err = uc.createInstance(config); err != nil {
			return nil, err
		}
	}

	return chosenCloudVM, err
}

func (uc *CloudSvcDigitalocean) getRelatedVM() (*entity.DigitaloceanDroplet, error) {
	instances, _, _, err2 := uc.api.ListDroplets(
		api_digitalocean.DOListDropletAPIRequest{
			Name: uc.env.VmLabel,
		},
	)
	if err2 != nil {
		return nil, fmt.Errorf(
			"ERROR: Failed to get droplet by name %s because: %v",
			uc.env.VmLabel,
			err2,
		)
	}
	if len(instances) > 0 {
		instance := instances[0]
		fmt.Printf(
			"Instance found : %s \n",
			entity.ConvertToJSON[*entity.DigitaloceanDroplet](instance),
		)
		return instance, nil
	} else {
		return nil, nil
	}
}

func (uc *CloudSvcDigitalocean) getImageID() (string, error) {
	if uc.env.VmChooseSnapshotOverISO {
		snapshots, _, err := uc.api.ListSnapshots(20, 1, "droplet")
		if err != nil {
			return "", err
		}
		for _, snapshot := range snapshots {
			for _, tag := range snapshot.Tags {
				if tag == uc.env.VmLabel {
					return snapshot.ID, nil
				}
			}
		}
	}
	return uc.env.VmISO, nil
}

func (uc *CloudSvcDigitalocean) createInstance(config api_digitalocean.DOCreateDropletRequest) (
	*entity.DigitaloceanDroplet,
	error,
) {
	instance, actions, err := uc.api.CreateDroplets(config)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create instance from config: %v because: %w",
			entity.ConvertToJSON(config),
			err,
		)
	}

	if err = uc.api.AttachDropletToFirewall(uc.env.VmFirewall, []int{instance.ID}); err != nil {
		return nil, fmt.Errorf(
			"failed to attach firewall to instance ID %d because: %w",
			instance.ID,
			err,
		)
	}

	if err = uc.waitUntilRunning(instance.ID, actions[0].ID, 0); err != nil {
		return nil, err
	}

	return instance, nil
}
func (uc *CloudSvcDigitalocean) waitUntilRunning(dropletID int, dropletActionID int, currentCount int) error {
	fmt.Print("waiting...")
	time.Sleep(30 * time.Second)
	fmt.Printf("instance running check #%v \n", currentCount)

	action, err := uc.api.GetDropletAction(dropletID, dropletActionID)
	if err != nil {
		return fmt.Errorf(
			"failed to get action from instance id %d with action id %d because: %w",
			dropletID,
			dropletActionID,
			err,
		)
	}

	if currentCount == instanceCheckMaxCount && action.Status != "completed" {
		return errors.New(
			fmt.Sprintf(
				"instance id %d with action id %d not completed after %v checks",
				dropletID,
				dropletActionID,
				currentCount,
			),
		)
	}
	if action.Status == "completed" {
		fmt.Printf(
			"instance id %d with action id %d completed  \n",
			dropletID,
			dropletActionID,
		)
		return nil
	}
	return uc.waitUntilRunning(dropletID, dropletActionID, currentCount+1)
}

// --
// --
// --
//
//func (uc *CloudSvcDigitalocean) StopInstance() error {
//	var existingInstance *entity.VultrInstance
//	var newSnapshot *entity.VultrSnapshot
//	var err error
//
//	if existingInstance, err = uc.getRelatedVM(); err != nil {
//		return err
//	} else if existingInstance == nil {
//		return fmt.Errorf("there is no active instance available")
//	}
//
//	if newSnapshot, err = uc.newSnapshot(existingInstance); err != nil {
//		return err
//	}
//
//	if err = uc.deleteVM(existingInstance); err != nil {
//		return err
//	}
//
//	if err = uc.deleteOldestSnapshot(existingInstance, *newSnapshot); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (uc *CloudSvcDigitalocean) newSnapshot(instance *entity.VultrInstance) (
//	*entity.VultrSnapshot,
//	error,
//) {
//	var createdSnapshot entity.VultrSnapshot
//	var err error
//
//	if createdSnapshot, err = uc.api.CreateSnapshot(
//		instance.ID,
//		strings.TrimSpace(uc.env.VmLabel),
//	); err != nil {
//		return nil, fmt.Errorf(
//			"failed to create new snapshot from instance %v because: %w",
//			instance.ID,
//			err,
//		)
//	}
//	if err = uc.waitUntilAvailable(&createdSnapshot, 0); err != nil {
//		return &createdSnapshot, err
//	}
//
//	fmt.Printf(
//		"New snapshot created and activated: %s \n",
//		entity.ConvertToJSON[entity.VultrSnapshot](createdSnapshot),
//	)
//	return &createdSnapshot, nil
//}
//func (uc *CloudSvcDigitalocean) waitUntilAvailable(snapshot *entity.VultrSnapshot, currentCount int) error {
//	fmt.Print("waiting...")
//	time.Sleep(30 * time.Second)
//	fmt.Printf("snapshot availability check #%v \n", currentCount)
//
//	updatedSnapshot, err := uc.api.GetSnapshot(snapshot.ID)
//	if err != nil {
//		return fmt.Errorf("failed to get snapshot %v because: %w", snapshot.ID, err)
//	}
//	if currentCount == instanceCheckMaxCount && updatedSnapshot.Status != "complete" {
//		return errors.New(
//			fmt.Sprintf(
//				"snapshot not active yet after %v checks, this is the latest status: %v",
//				currentCount,
//				updatedSnapshot.Status,
//			),
//		)
//	}
//	if updatedSnapshot.Status == "complete" {
//		return nil
//	}
//	return uc.waitUntilAvailable(snapshot, currentCount+1)
//}
//
//func (uc *CloudSvcDigitalocean) deleteOldestSnapshot(instance *entity.VultrInstance, snapshot entity.VultrSnapshot) error {
//	getSnapshotsResult, _, err := uc.api.ListSnapshots(api_vultr.ListSnapshotSortByDateCreatedAsc)
//	if err != nil {
//		return fmt.Errorf(
//			"failed to get oldest snapshot from instance %v because: %w",
//			instance.ID,
//			err,
//		)
//	}
//
//	oldSnapshot := uc.findActiveSnapshot(getSnapshotsResult, 0)
//	if oldSnapshot == nil || oldSnapshot.ID == snapshot.ID {
//		return nil
//	}
//	if err2 := uc.api.RemoveSnapshot(oldSnapshot.ID); err2 != nil {
//		return fmt.Errorf("failed to remove snapshot %v because: %w", oldSnapshot.ID, err2)
//	}
//
//	fmt.Printf("Oldest snapshot deleted \n")
//	return nil
//}
//
//func (uc *CloudSvcDigitalocean) deleteVM(vm *entity.VultrInstance) error {
//	if err := uc.api.HaltDroplet([]string{vm.ID}); err != nil && !errors.Is(
//		err,
//		errors.New("EOF"),
//	) {
//		return fmt.Errorf("failed to halt instance %v because: %v", vm.ID, err)
//	}
//
//	if err := uc.api.RemoveDroplet(vm.ID); err != nil && !errors.Is(err, errors.New("EOF")) {
//		return fmt.Errorf("failed to remove instance %v because: %w", vm.ID, err)
//	}
//
//	if err := uc.waitUntilTerminated(vm, 0); err != nil {
//		return fmt.Errorf("failed to validate instance %v removal because: %w", vm.ID, err)
//	}
//
//	fmt.Printf(
//		"Instance deleted : %s \n",
//		entity.ConvertToJSON[*entity.VultrInstance](vm),
//	)
//
//	return nil
//}
//func (uc *CloudSvcDigitalocean) waitUntilTerminated(vm *entity.VultrInstance, currentCount int) error {
//	fmt.Print("waiting...")
//	time.Sleep(30 * time.Second)
//	fmt.Printf("snapshot deletion check #%v \n", currentCount)
//
//	deletedInstance, err3 := uc.api.GetDroplet(vm.ID)
//	if err3 != nil {
//		return nil
//	}
//	if currentCount == instanceCheckMaxCount && deletedInstance.ID != "" {
//		return errors.New(
//			fmt.Sprintf(
//				"instance id %v not deleted after %v checks",
//				vm.ID,
//				currentCount,
//			),
//		)
//	}
//	if deletedInstance.ID == "" {
//		return nil
//	}
//	return uc.waitUntilTerminated(vm, currentCount+1)
//}
