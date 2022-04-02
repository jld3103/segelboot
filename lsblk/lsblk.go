package lsblk

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Output struct {
	BlockDevices []BlockDevice
}

type BlockDevice struct {
	PartUUID    string
	PkName      string
	KName       string
	MountPoints []string
}

func Execute() (*[]BlockDevice, error) {
	cmd := exec.Command("lsblk", "--json", "--output", "PARTUUID,PKNAME,KNAME,MOUNTPOINTS")
	d, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsblk: %w", err)
	}

	o := &Output{}
	err = json.Unmarshal(d, o)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output of lsblk: %w", err)
	}

	return &o.BlockDevices, nil
}
