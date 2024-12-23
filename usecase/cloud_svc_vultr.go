package usecase

import (
	"dev_scripts/entity"
	"dev_scripts/repository/api_vultr"
	"errors"
	"fmt"
	"strings"
	"time"
)

var instanceCheckMaxCount = 120

type CloudSvcVultr struct {
	env      *entity.EnvCloudServer
	vultrAPI *api_vultr.VultrAPI
}

func NewCloudSvcVultr(env *entity.EnvCloudServer) *CloudSvcVultr {
	vultrRepo := api_vultr.NewVultrAPI(env)
	return &CloudSvcVultr{env, vultrRepo}
}

// --

func (uc *CloudSvcVultr) StartInstance() (*entity.VultrInstance, error) {
	var chosenCloudVM *entity.VultrInstance
	var err error

	if chosenCloudVM, err = uc.getActiveVM(); err != nil {
		return nil, err
	}

	if chosenCloudVM == nil || chosenCloudVM.MainIP == "" {
		config := entity.VultrAPIRequestCreateInstanceConfig{
			Region:  strings.TrimSpace(uc.env.VmRegion),
			Plan:    strings.TrimSpace(uc.env.VmVultrPlan),
			Backups: "disabled",
			Label:   strings.TrimSpace(uc.env.VmLabel),
		}

		chosenISO, err := uc.getISoID()
		if err != nil {
			return nil, err
		}
		config.IsoID = chosenISO.Id

		if chosenSshKey, err := uc.getSSHKeyID(); err != nil {
			return nil, err
		} else if chosenSshKey != nil {
			config.SshKeyIDs = []string{chosenSshKey.Id}
		}

		if activeSnapshot, err := uc.getLatestActiveSnapshot(); err != nil {
			return nil, err
		} else if activeSnapshot != nil {
			config.SnapshotID = activeSnapshot.ID
		}

		if chosenFirewallGroup, err := uc.getFirewallGroup(); err != nil {
			return nil, err
		} else if chosenFirewallGroup != nil {
			config.FirewallGroupID = chosenFirewallGroup.Id
		}

		//reservedIps, _, err0 := vultrApi.GetReservedIps()
		//if err0 != nil {
		//	return nil, fmt.Errorf("failed to get reserved ips because: %w", err0)
		//} else if len(reservedIps) > 0 {
		//	config.ReservedIPV4 = reservedIps[0].Id
		//}

		if chosenCloudVM, err = uc.createInstance(config); err != nil {
			return nil, err
		}
	}

	if err = uc.attachExternalStorage(chosenCloudVM); err != nil {
		return nil, err
	}
	return chosenCloudVM, err
}

func (uc *CloudSvcVultr) getActiveVM() (*entity.VultrInstance, error) {
	instances, _, err2 := uc.vultrAPI.ListInstances()
	if err2 != nil {
		return nil, fmt.Errorf("ERROR: Failed to get instance list because: %v", err2)
	}
	if len(instances) > 0 {
		instance := &instances[0]
		fmt.Printf(
			"Instance found : %s \n",
			entity.ConvertToJSON[*entity.VultrInstance](instance),
		)
		return instance, nil
	} else {
		return nil, nil
	}
}

func (uc *CloudSvcVultr) getLatestActiveSnapshot() (*entity.VultrSnapshot, error) {
	snapshots, _, err := uc.vultrAPI.ListSnapshots(api_vultr.ListSnapshotSortByDateCreatedDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest active snapshot because: %w", err)
	}
	return uc.findActiveSnapshot(snapshots, 0), nil
}
func (uc *CloudSvcVultr) findActiveSnapshot(snapshots []entity.VultrSnapshot, currentIdx int) *entity.VultrSnapshot {
	if currentIdx >= len(snapshots) {
		return nil
	}
	snapshot := snapshots[currentIdx]
	if snapshot.Status == "complete" && snapshot.Description == strings.TrimSpace(uc.env.VmLabel) {
		fmt.Printf(
			"Active snapshot found : %s \n",
			entity.ConvertToJSON[entity.VultrSnapshot](snapshot),
		)
		return &snapshot
	}
	return uc.findActiveSnapshot(snapshots, currentIdx+1)
}

func (uc *CloudSvcVultr) getISoID() (*entity.VultrISO, error) {
	listOfISO, _, err1 := uc.vultrAPI.GetISOs()
	if err1 != nil {
		return nil, fmt.Errorf("failed to get iso list because: %w", err1)
	}
	for _, iso := range listOfISO {
		if iso.Filename == strings.TrimSpace(uc.env.VmISoFilename) {
			return &iso, nil
		}
	}
	return nil, fmt.Errorf(
		"failed to find iso key with filename: %s. ISO: %v",
		uc.env.VmISoFilename,
		entity.ConvertToJSON(listOfISO),
	)
}

func (uc *CloudSvcVultr) getFirewallGroup() (*entity.VultrFirewallGroup, error) {
	listOfFirewallGroup, _, err1 := uc.vultrAPI.ListFirewallGroups()
	if err1 != nil {
		return nil, fmt.Errorf("failed to get firewall group list because: %w", err1)
	}
	for _, firewallGroup := range listOfFirewallGroup {
		if firewallGroup.Description == strings.TrimSpace(uc.env.VmFirewallLabel) {
			return &firewallGroup, nil
		}
	}
	return nil, nil
}

func (uc *CloudSvcVultr) getSSHKeyID() (*entity.VultrSSHKey, error) {
	sshKeys, _, err1 := uc.vultrAPI.GetSSHKeys()
	if err1 != nil {
		return nil, fmt.Errorf("failed to get ssh keys because: %w", err1)
	}
	for _, sshKey := range sshKeys {
		if sshKey.Name == strings.TrimSpace(uc.env.SSHKeyLabel) {
			return &sshKey, nil
		}
	}
	return nil, fmt.Errorf("failed to find ssh key with label: %s", uc.env.SSHKeyLabel)
}

func (uc *CloudSvcVultr) createInstance(config entity.VultrAPIRequestCreateInstanceConfig) (
	*entity.VultrInstance,
	error,
) {
	instance, _, err := uc.vultrAPI.CreateInstance(config)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create instance from config: %v because: %w",
			entity.ConvertToJSON(config),
			err,
		)
	}

	return uc.waitUntilRunning(&instance, 0)
}
func (uc *CloudSvcVultr) waitUntilRunning(vm *entity.VultrInstance, currentCount int) (
	*entity.VultrInstance,
	error,
) {
	fmt.Print("waiting...")
	time.Sleep(30 * time.Second)
	fmt.Printf("instance running check #%v \n", currentCount)

	updatedInstance, err3 := uc.vultrAPI.GetInstance(vm.ID)
	if err3 != nil {
		return updatedInstance, fmt.Errorf("failed to get instance %v because: %w", vm.ID, err3)
	}

	if currentCount == instanceCheckMaxCount && updatedInstance.PowerStatus != "running" && updatedInstance.MainIP != "0.0.0.0" {
		return updatedInstance, errors.New(
			fmt.Sprintf(
				"instance not started after %v checks, this is the latest power status: %v, status: %v, ip: %v",
				currentCount,
				updatedInstance.PowerStatus,
				updatedInstance.Status,
				updatedInstance.MainIP,
			),
		)
	}
	if updatedInstance.PowerStatus == "running" {
		fmt.Printf(
			"Instance created and running : %s \n",
			entity.ConvertToJSON[*entity.VultrInstance](vm),
		)
		return updatedInstance, nil
	}
	return uc.waitUntilRunning(vm, currentCount+1)
}

func (uc *CloudSvcVultr) attachExternalStorage(vm *entity.VultrInstance) error {
	if len(uc.env.VmBlockStoragesLabel) > 0 {
		obtainedBlockStorages, _, err := uc.vultrAPI.GetBlockStorages()
		if err != nil {
			return fmt.Errorf("failed to get block storages because: %w", err)
		}

		var totalMatchedBlockStorages int
		labelToBlockStorage := map[string]*entity.VultrBlockStorage{}
		for _, storageLabel := range uc.env.VmBlockStoragesLabel {
			labelToBlockStorage[strings.TrimSpace(storageLabel)] = nil
		}
		for _, blockStorageData := range obtainedBlockStorages {
			if _, ok := labelToBlockStorage[blockStorageData.Label]; ok {
				labelToBlockStorage[blockStorageData.Label] = &blockStorageData
				totalMatchedBlockStorages += 1
			}
			if totalMatchedBlockStorages >= len(labelToBlockStorage) {
				break
			}
		}

		for label, blockStorage := range labelToBlockStorage {
			if blockStorage == nil {
				return fmt.Errorf("failed to find block storage with label %s", label)
			} else {
				if err := uc.vultrAPI.AttachBlockStorage(blockStorage.Id, vm.ID); err != nil {
					return fmt.Errorf(
						"failed to attach storage ID %s into instance ID %s because: %w",
						blockStorage.Id,
						vm.ID,
						err,
					)
				}
			}
		}
		fmt.Printf(
			"Instance block storage %s attached \n",
			entity.ConvertToJSON(labelToBlockStorage),
		)
	}
	return nil
}

// --
// --
// --

func (uc *CloudSvcVultr) StopInstance() error {
	var existingInstance *entity.VultrInstance
	var newSnapshot *entity.VultrSnapshot
	var err error

	if existingInstance, err = uc.getActiveVM(); err != nil {
		return err
	} else if existingInstance == nil {
		return fmt.Errorf("there is no active instance available")
	}

	if newSnapshot, err = uc.newSnapshot(existingInstance); err != nil {
		return err
	}

	if err = uc.deleteVM(existingInstance); err != nil {
		return err
	}

	if err = uc.deleteOldestSnapshot(existingInstance, *newSnapshot); err != nil {
		return err
	}

	return nil
}

func (uc *CloudSvcVultr) newSnapshot(instance *entity.VultrInstance) (
	*entity.VultrSnapshot,
	error,
) {
	var createdSnapshot entity.VultrSnapshot
	var err error

	if createdSnapshot, err = uc.vultrAPI.CreateSnapshot(
		instance.ID,
		strings.TrimSpace(uc.env.VmLabel),
	); err != nil {
		return nil, fmt.Errorf(
			"failed to create new snapshot from instance %v because: %w",
			instance.ID,
			err,
		)
	}
	if err = uc.waitUntilAvailable(&createdSnapshot, 0); err != nil {
		return &createdSnapshot, err
	}

	fmt.Printf(
		"New snapshot created and activated: %s \n",
		entity.ConvertToJSON[entity.VultrSnapshot](createdSnapshot),
	)
	return &createdSnapshot, nil
}
func (uc *CloudSvcVultr) waitUntilAvailable(snapshot *entity.VultrSnapshot, currentCount int) error {
	fmt.Print("waiting...")
	time.Sleep(30 * time.Second)
	fmt.Printf("snapshot availability check #%v \n", currentCount)

	updatedSnapshot, err := uc.vultrAPI.GetSnapshot(snapshot.ID)
	if err != nil {
		return fmt.Errorf("failed to get snapshot %v because: %w", snapshot.ID, err)
	}
	if currentCount == instanceCheckMaxCount && updatedSnapshot.Status != "complete" {
		return errors.New(
			fmt.Sprintf(
				"snapshot not active yet after %v checks, this is the latest status: %v",
				currentCount,
				updatedSnapshot.Status,
			),
		)
	}
	if updatedSnapshot.Status == "complete" {
		return nil
	}
	return uc.waitUntilAvailable(snapshot, currentCount+1)
}

func (uc *CloudSvcVultr) deleteOldestSnapshot(instance *entity.VultrInstance, snapshot entity.VultrSnapshot) error {
	getSnapshotsResult, _, err := uc.vultrAPI.ListSnapshots(api_vultr.ListSnapshotSortByDateCreatedAsc)
	if err != nil {
		return fmt.Errorf(
			"failed to get oldest snapshot from instance %v because: %w",
			instance.ID,
			err,
		)
	}

	oldSnapshot := uc.findActiveSnapshot(getSnapshotsResult, 0)
	if oldSnapshot == nil || oldSnapshot.ID == snapshot.ID {
		return nil
	}
	if err2 := uc.vultrAPI.RemoveSnapshot(oldSnapshot.ID); err2 != nil {
		return fmt.Errorf("failed to remove snapshot %v because: %w", oldSnapshot.ID, err2)
	}

	fmt.Printf("Oldest snapshot deleted \n")
	return nil
}

func (uc *CloudSvcVultr) deleteVM(vm *entity.VultrInstance) error {
	if err := uc.vultrAPI.HaltInstances([]string{vm.ID}); err != nil && !errors.Is(
		err,
		errors.New("EOF"),
	) {
		return fmt.Errorf("failed to halt instance %v because: %v", vm.ID, err)
	}

	if err := uc.vultrAPI.RemoveInstance(vm.ID); err != nil && !errors.Is(err, errors.New("EOF")) {
		return fmt.Errorf("failed to remove instance %v because: %w", vm.ID, err)
	}

	if err := uc.waitUntilTerminated(vm, 0); err != nil {
		return fmt.Errorf("failed to validate instance %v removal because: %w", vm.ID, err)
	}

	fmt.Printf(
		"Instance deleted : %s \n",
		entity.ConvertToJSON[*entity.VultrInstance](vm),
	)

	return nil
}
func (uc *CloudSvcVultr) waitUntilTerminated(vm *entity.VultrInstance, currentCount int) error {
	fmt.Print("waiting...")
	time.Sleep(30 * time.Second)
	fmt.Printf("snapshot deletion check #%v \n", currentCount)

	deletedInstance, err3 := uc.vultrAPI.GetInstance(vm.ID)
	if err3 != nil {
		return nil
	}
	if currentCount == instanceCheckMaxCount && deletedInstance.ID != "" {
		return errors.New(
			fmt.Sprintf(
				"instance id %v not deleted after %v checks",
				vm.ID,
				currentCount,
			),
		)
	}
	if deletedInstance.ID == "" {
		return nil
	}
	return uc.waitUntilTerminated(vm, currentCount+1)
}
