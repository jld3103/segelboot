package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jld3103/segelboot/config"
	"github.com/jld3103/segelboot/efibootmgr"
	"github.com/jld3103/segelboot/lsblk"
)

//nolint:gocognit
func NewRootCmd() *cobra.Command {
	var configFile string
	var deleteEntries bool

	rootCmd := &cobra.Command{
		Use: "segelboot",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := config.ReadConfigFile(configFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			out, err := efibootmgr.ExecuteWithOutput(
				[]efibootmgr.Flag{},
				map[efibootmgr.Option]string{},
			)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o := efibootmgr.ParseOutput(out)

			if deleteEntries {
				for _, bootEntry := range o.BootEntries {
					if matchLabel(bootEntry) {
						err = deleteBootEntry(bootEntry)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
					}
				}
			} else {
				var blockDevices *[]lsblk.BlockDevice
				blockDevices, err = lsblk.Execute()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for _, entry := range c.Entries {
					var blockDevice *lsblk.BlockDevice
					var partitionIndex string
					blockDevice, partitionIndex, err = validateEntry(*blockDevices, entry)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					err = deleteBootEntriesForEntry(o.BootEntries, entry)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					err = createBootEntry(*blockDevice, partitionIndex, entry)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			}
		},
	}

	rootCmd.Flags().StringVarP(
		&configFile,
		"config",
		"c",
		"/etc/segelboot.conf",
		"Specify the location of the config file",
	)
	rootCmd.Flags().BoolVarP(
		&deleteEntries,
		"delete",
		"d",
		false,
		"Delete all entries created by segelboot",
	)

	return rootCmd
}

func deleteBootEntry(bootEntry efibootmgr.BootEntry) error {
	err := efibootmgr.Execute(
		[]efibootmgr.Flag{
			efibootmgr.FlagDeleteBootnum,
		},
		map[efibootmgr.Option]string{
			efibootmgr.OptionBootnum: bootEntry.Bootnum,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create delete entry: %w", err)
	}

	return nil
}

func deleteBootEntriesForEntry(bootEntries []efibootmgr.BootEntry, entry *config.Entry) error {
	for _, bootEntry := range bootEntries {
		if matchLabelWithEntry(bootEntry, entry) {
			fmt.Printf("Found existing boot entry %s for entry '%s'\n. Removing entry", bootEntry.Bootnum, entry.ID)
			err := deleteBootEntry(bootEntry)
			if err != nil {
				return err
			}

			// We don't exit the loop here in order to remove multiple entries in case we messed up when creating them
		}
	}

	return nil
}

func createBootEntry(blockDevice lsblk.BlockDevice, partitionIndex string, entry *config.Entry) error {
	err := efibootmgr.Execute(
		[]efibootmgr.Flag{
			efibootmgr.FlagCreate,
		},
		map[efibootmgr.Option]string{
			efibootmgr.OptionDisk:      fmt.Sprintf("/dev/%s", blockDevice.PkName),
			efibootmgr.OptionPartition: partitionIndex,
			efibootmgr.OptionLabel:     formatLabel(entry),
			efibootmgr.OptionLoader:    entry.Loader,
		},
		entry.CmdLine,
	)
	if err != nil {
		return fmt.Errorf("failed to create boot entry: %w", err)
	}

	return nil
}

func validateEntry(blockDevices []lsblk.BlockDevice, entry *config.Entry) (*lsblk.BlockDevice, string, error) {
	blockDevice, err := findBlockDeviceForEntry(blockDevices, entry)
	if err != nil {
		return nil, "", err
	}

	partitionIndex, err := findPartitionIndex(*blockDevice)
	if err != nil {
		return nil, "", err
	}

	err = validateLoader(*blockDevice, entry)
	if err != nil {
		return nil, "", err
	}

	return blockDevice, partitionIndex, nil
}

func findBlockDeviceForEntry(blockDevices []lsblk.BlockDevice, entry *config.Entry) (*lsblk.BlockDevice, error) {
	var blockDevice lsblk.BlockDevice
	found := false
	for _, b := range blockDevices {
		if b.PartUUID == entry.PartitionUUID {
			blockDevice = b
			found = true

			break
		}
	}
	if !found {
		//nolint:goerr113
		return nil, fmt.Errorf("could not find partition UUID '%s'", entry.PartitionUUID)
	}

	return &blockDevice, nil
}

func findPartitionIndex(blockDevice lsblk.BlockDevice) (string, error) {
	s := blockDevice.KName[len(blockDevice.PkName):]
	if strings.HasPrefix(s, "p") {
		s = s[1:]
	}
	regex := regexp.MustCompile("^\\d+$")
	if !regex.MatchString(s) {
		//nolint:goerr113
		return "", fmt.Errorf("failed to parse partition index '%s'", s)
	}

	return s, nil
}

func validateLoader(blockDevice lsblk.BlockDevice, entry *config.Entry) error {
	loaderPaths := make([]string, len(blockDevice.MountPoints))

	for i, mountPoint := range blockDevice.MountPoints {
		loaderPath := filepath.Join(mountPoint, entry.Loader)
		if _, err := os.Stat(loaderPath); err == nil {
			return nil
		}

		loaderPaths[i] = loaderPath
	}

	//nolint:goerr113
	return fmt.Errorf("could not find a linux kernel image at %s", strings.Join(loaderPaths, " or "))
}

func formatLabel(entry *config.Entry) string {
	return fmt.Sprintf("Segelboot: %s (%s)", entry.Name, entry.ID)
}

func matchLabel(bootEntry efibootmgr.BootEntry) bool {
	regex := regexp.MustCompile("^Segelboot: .* \\(.*\\)$")

	return regex.MatchString(bootEntry.Label)
}

func matchLabelWithEntry(bootEntry efibootmgr.BootEntry, entry *config.Entry) bool {
	regex := regexp.MustCompile(fmt.Sprintf("^Segelboot: .* \\(%s\\)$", entry.ID))

	return regex.MatchString(bootEntry.Label)
}
