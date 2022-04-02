package efibootmgr

type Flag string
type Option string

// See `efibootmgr --help` for the meaning of all flags and options

const (
	FlagActive   Flag = "active"
	FlagInactive Flag = "inactive"
)

const (
	FlagCreate     Flag = "create"
	FlagCreateOnly Flag = "create-only"
)

const (
	FlagDeleteBootnext Flag   = "delete-bootnext"
	OptionBootnext     Option = "bootnext"
)

const (
	FlagDeleteBootnum Flag   = "delete-bootnum"
	OptionBootnum     Option = "bootnum"
)

const (
	FlagDeleteBootorder Flag   = "delete-bootorder"
	OptionBootorder     Option = "bootorder"
)

const (
	FlagDeleteTimeout Flag   = "delete-timeout"
	OptionTimeout     Option = "timeout"
)

const (
	OptionDisk      Option = "disk"
	OptionPartition Option = "part"
)

const (
	OptionDevice Option = "device"
	OptionEdd    Option = "edd"
)

const (
	OptionMirrorAbove4G Option = "mirror-above-4g"
	OptionMirrorBelow4G Option = "mirror-below-4g"
)

const (
	FlagDriver  Flag = "driver"
	FlagSysprep Flag = "sysprep"
)

const (
	FlagRemoveDuplicates   Flag   = "remove-dups"
	FlagUnicode            Flag   = "unicode"
	OptionAppendBinaryArgs Option = "append-binary-args"
	OptionIface            Option = "iface"
	OptionLabel            Option = "label"
	OptionLoader           Option = "loader"
)

const (
	FlagGpt            Flag = "gpt"
	FlagWriteSignature Flag = "write-signature"
)

const (
	FlagHelp    Flag = "help"
	FlagVersion Flag = "version"
	FlagQuiet   Flag = "quiet"
	FlagVerbose Flag = "verbose"
)
