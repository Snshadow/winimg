package winimg

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/w32api/dismapi"
	"github.com/Snshadow/winimg/w32api/wimgapi"
)

// NewDismProgress creates callback function to be
// used get information of DISM operation progress.
func NewDismProgress(cb dismapi.DismProgressCallback) uintptr {
	return windows.NewCallback(cb)
}

// DismProgressOpts contains optional parameters to
// cancel or get information while DISM operation
// is in progress.
type DismProgressOpts struct {
	// create with [windows.CreateEvent] and
	// trigger with [windows.SetEvent]
	//
	// call [windows.ResetEvent] to use again
	// for other DISM operation
	//
	// close with [windows.CloseHandle]
	CancelEvent windows.Handle
	// callback function returned from [NewDismProgress]
	Progress uintptr
	// userData passed to [dismapi.DismProgressCallback]
	UserData unsafe.Pointer
}

// DismMountOpts contains options used for mounting DISM image.
type DismMountOpts struct {
	// if not empty, used instead of image index
	ImageName                          string
	ReadOnly, Optimize, CheckIntegrity bool
	DismProgressOpts
}

// DismUnmountOpts contains options used for unmounting DISM image.
type DismUnmountOpts struct {
	// unmount even if its not mounted from [DismImageFile]
	Force bool
	// Append, GenerateIntegrity or SupportEa are ignored if
	// Commit is false
	Commit, Append, GenerateIntegrity, SupportEa bool
	DismProgressOpts
}

// DismCommitOpts contains options used for saving
// changes in DISM image.
type DismCommitOpts struct {
	Append, GenerateIntegrity, SupportEa bool
	DismProgressOpts
}

// DismAddCapabilityOpts contains options used for
// adding capability to DISM image.
type DismAddCapabilityOpts struct {
	Name string
	// do not use WU/WSUS Update
	LimitAccess bool
	SourcePaths []string
	DismProgressOpts
}

// DismAddPackageOpts contains options used for adding
// package(.cab, .msu, .mum file) to DISM image.
type DismAddPackageOpts struct {
	IgnoreCheck    bool
	PreventPending bool
	DismProgressOpts
}

// DismDisableFeatureOpts contains options used for
// disabling or removing feature in a DISM image.
type DismDisableFeatureOpts struct {
	PackageName   string
	RemovePayload bool
	DismProgressOpts
}

// DismEnableFeatureOpts contains options used
// for enabling feature in a DISM image.
type DismEnableFeatureOpts struct {
	// optional, path to .cab file or package name
	Identifier string
	// do not use WU
	LimitAccess bool
	SourcePaths []string
	EnableAll   bool
	DismProgressOpts
}

// DismRestoreImageHealthOpts contains options
// used for repairing a corrupt image.
type DismRestoreImageHealthOpts struct {
	SourcePaths []string
	LimitAccess bool
	DismProgressOpts
}

// NewWimMessageCallback creates callback function to
// be used with [wimgapi.WIMRegisterMessageCallback] or
// [wimgapi.WIMUnregisterMessageCallback].
func NewWimMessageCallback(cb wimgapi.WIMMessageCallback) uintptr {
	return windows.NewCallback(cb)
}

// WimApplyOpts contains options used for applying
// volume image to a directory.
type WimApplyOpts struct {
	// verifies that files match original data
	Verify bool
	// specify that the image is be sequentially read for
	// caching and performance
	Index bool
	// apply without physically creating directories or files,
	// used for enumerating files and directories
	NoApply bool
	// disable restoring security information for directories
	NoDirAcl bool
	// disable restoring srcurity information for files
	NoFileAcl bool
	// disable automatic path fixups for junctions or symbolic
	// links
	NoReparseFix bool
	// send [wimgapi.WIM_MSG_FILEINFO] during apply operation
	FileInfo bool
	// validate the image for Trusted Desktop
	ConfirmTrustedFile bool
	// format image to install on WIMBoot
	WIMBoot bool
	// compress operating system files
	Compact bool
	// apply image with extended attributes(EA)
	SupportEa bool
}

// WimExportOpts contains options used for
// exporting volume image to other wim file.
type WimExportOpts struct {
	// export image even if it is stored in destination
	// .wim file
	AllowDuplicate bool
	// do not export image resources or XML information
	OnlyResources bool
	// export only resources and XML information
	OnlyMetadata bool
	// verify source image file
	VerifySource bool
	// verify destination image file
	VerifyDestination bool
}

// WimImageMountOpts contains options used for
// mounting volume image to a directory.
type WimImageMountOpts struct {
	// mount the image without the ability to save changes
	ReadOnly bool
	// verifies that files match original data
	Verify bool
	// disable automatic path fixups for junctions or symbolic
	// links
	NoReparseFix bool
	// disable restoring security information for directories
	NoDirAcl bool
	// disable restoring srcurity information for files
	NoFileAcl bool
}

type WimCompression uint32

const (
	// no compression
	None WimCompression = iota
	// xpress
	Fast
	// lzx
	Maximum
	LZMS
)

type WimCreateMode uint32

const (
	// open only if the file exists
	OpenExisting WimCreateMode = iota
	// open the file, create it first if it does not exist
	OpenOrCreate
	// create only if the file does not exist
	CreateIfNotExist
	// always create a new file, truncating the existing file
	AlwaysCreate
)

// WimCreateFileOpts contains optsions used for
// creating or opening .wim file.
type WimCreateFileOpts struct {
	// do not read, write or mount
	NoRead, NoWrite, NoMount bool
	// creation behavior, default to [OpenExisting]
	CreateMode WimCreateMode
	// compression level
	Compression WimCompression
	// temporary path used for image operation
	TempPath string
}

// WimCaptureOpts contains options used for
// capturing directory to a .wim file as a
// volume image.
type WimCaptureOpts struct {
	// verify single-instance files byte by byte
	Verify bool
	// disable automatic path fixups for junctions or symbolic
	// links
	NoReparseFix bool
	// disable capturing security information for directories
	NoDirAcl bool
	// disable capturing security information for files
	NoFileAcl bool
	// format image to install on WIMBoot
	WIMBoot bool
	// capture image with extended attributes(EA)
	SupportEa bool
}

// WimRemountOpts contains options used for
// remounting previously mounted volume image
// at a directory.
//
// Defined for future use.
type WimRemountOpts struct {
}
