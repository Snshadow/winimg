// Package winimg provides functions to access and manupulate Windows image(.wim, .esd) files.
package winimg

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/w32api"
	"github.com/Snshadow/winimg/w32api/dismapi"
	"github.com/Snshadow/winimg/w32api/wimgapi"
	"github.com/deckarep/golang-set/v2"
	"github.com/puzpuzpuz/xsync/v3"
)

var (
	ErrNotMounted          = errors.New("the path is not a mount path")
	ErrAlreadyMounted      = errors.New("the image is already mounted")
	ErrDismSessionNotExist = errors.New("the specified DismSession does not exist")
)

var (
	// do not run DismInitialize when it is already initialized for a process
	dismInitialized bool
	dismMutex       sync.Mutex

	// default to information level
	dismLogLevel = dismapi.DismLogErrorsWarningsInfo
	// use working directory by default
	dismLogPath = ".\\dism.log"
	// use DISM API default path
	dismScratchDir = ""

	// track current use of DISM API
	curDismImg = mapset.NewSet[*DismImageFile]()
)

// ConfigureDism configures global settings for DISM API for process.
func ConfigureDism(logLevel dismapi.DismLogLevel, logPath, scratchDir string) {
	dismLogLevel = logLevel
	dismLogPath = logPath
	dismScratchDir = scratchDir
}

// initDism initalizes DISM API with global settings.
func initDism() error {
	dismMutex.Lock()
	defer dismMutex.Unlock()

	if !dismInitialized {
		err := dismapi.DismInitialize(dismLogLevel, dismLogPath, dismScratchDir)
		// returns DISMAPI_S_RELOAD_IMAGE_SESSION_REQUIRED if already initialized
		if err != nil && err.(*w32api.WinimgInternalErr).Errno() !=
			dismapi.DISMAPI_S_RELOAD_IMAGE_SESSION_REQUIRED {
			return err
		}
		dismInitialized = true
	}

	return nil
}

// GetDismMountedImages returns information of mounted images as
// an array of [dismapi.GoDismMountedImageInfo].
func GetDismMountedImages() ([]dismapi.GoDismMountedImageInfo, error) {
	return dismapi.DismGetMountedImageInfo()
}

// CleanupDismMountPoints calls [dismapi.DismCleanupMountPoints].
func CleanupDismMountPoints() error {
	return dismapi.DismCleanupMountPoints()
}

// RemountDismImage calls [dismapi.DismRemountImage].
//
// Call [DismImageFile].OpenSession with mountPath after
// it is remounted.
func RemountDismImage(mountPath string) error {
	return dismapi.DismRemountImage(mountPath)
}

// DismImageSession handles operations related with DismSession.
type DismImageSession struct {
	session   dismapi.DismSession
	mountPath string

	sesMapRef mapset.Set[*DismImageSession]
}

// GetMountPath returns mount path being currently used.
func (d *DismImageSession) GetMountPath() string {
	return d.mountPath
}

// ApplyUnattend applies unattend answer file to
// the associated image with current session.
func (d *DismImageSession) ApplyUnattend(unattendFile string, singleSession bool) error {
	err := dismapi.DismApplyUnattend(d.session, unattendFile, singleSession)
	if err != nil {
		return err
	}

	return nil
}

// GetCapabilties gets capabilites from an image.
func (d *DismImageSession) GetCapabilities() ([]dismapi.GoDismCapability, error) {
	return dismapi.DismGetCapabilities(d.session)
}

// GetDrivers gets drivers in an image.
//
// If allDrivers is true, it gets all drivers including
// drivers from the image itself, otherwise only get
// drivers not from the image.
func (d *DismImageSession) GetDrivers(allDrivers bool) ([]dismapi.GoDismDriverPackage, error) {
	return dismapi.DismGetDrivers(d.session, allDrivers)
}

// GetFeatures gets feature in an image, regardless of
// whether the features are enabled or disabled.
//
// An optional parameter identifier can be an absolute
// path to a .cab file or the package name.
func (d *DismImageSession) GetFeatures(identifier string) ([]dismapi.GoDismFeature, error) {
	pkgIdentifier := dismapi.DismPackageNone
	if identifier != "" {
		if strings.HasSuffix(strings.ToLower(identifier), ".cab") {
			pkgIdentifier = dismapi.DismPackagePath
		} else {
			pkgIdentifier = dismapi.DismPackageName
		}
	}

	return dismapi.DismGetFeatures(d.session, identifier, pkgIdentifier)
}

// GetPackages gets packages in an image.
func (d *DismImageSession) GetPackages() ([]dismapi.GoDismPackage, error) {
	return dismapi.DismGetPackages(d.session)
}

// CheckImageHealth checks whether the image can be serviced
// or corrupted, with or without scanning(to get previous status)
// the image.
func (d *DismImageSession) CheckImageHealth(scanImage bool, opts *DismProgressOpts) (dismapi.DismImageHealthState, error) {
	var o DismProgressOpts
	if opts != nil {
		o = *opts
	}

	return dismapi.DismCheckImageHealth(
		d.session,
		scanImage,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// AddCapability adds a capability to an image.
func (d *DismImageSession) AddCapability(name string, opts *DismAddCapabilityOpts) error {
	var o DismAddCapabilityOpts
	if opts != nil {
		o = *opts
	}

	return dismapi.DismAddCapability(
		d.session,
		name,
		o.LimitAccess,
		o.SourcePaths,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// AddDriver adds a third party driver (.inf) to an image.
func (d *DismImageSession) AddDriver(driverPath string, forceUnsigned bool) error {
	return dismapi.DismAddDriver(d.session, driverPath, forceUnsigned)
}

// AddPackage adds a package file(.cab, .msu, .mum) to an image.
func (d *DismImageSession) AddPackage(packagePath string, opts *DismAddPackageOpts) error {
	var o DismAddPackageOpts
	if opts != nil {
		o = *opts
	}

	return dismapi.DismAddPackage(
		d.session,
		packagePath,
		o.IgnoreCheck,
		o.PreventPending,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// DisableFeature disables a feature in an image.
func (d *DismImageSession) DisableFeature(featureName string, opts *DismDisableFeatureOpts) error {
	var o DismDisableFeatureOpts
	if opts != nil {
		o = *opts
	}

	return dismapi.DismDisableFeature(
		d.session,
		featureName,
		o.PackageName,
		o.RemovePayload,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// EnableFeature enables a feature in an image,
// featureName can have multiple features, separated
// with semicolons.
func (d *DismImageSession) EnableFeature(featureName string, opts *DismEnableFeatureOpts) error {
	var o DismEnableFeatureOpts
	if opts != nil {
		o = *opts
	}

	pkgIdentifier := dismapi.DismPackageNone
	if o.Identifier != "" {
		if strings.HasSuffix(strings.ToLower(o.Identifier), ".cab") {
			pkgIdentifier = dismapi.DismPackagePath
		} else {
			pkgIdentifier = dismapi.DismPackageName
		}
	}

	return dismapi.DismEnableFeature(
		d.session,
		featureName,
		o.Identifier,
		pkgIdentifier,
		o.LimitAccess,
		o.SourcePaths,
		o.EnableAll,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// RemoveCapability removes a capability in an image.
func (d *DismImageSession) RemoveCapability(name string, opts *DismProgressOpts) error {
	var o DismProgressOpts
	if opts != nil {
		o = *opts
	}

	return dismapi.DismRemoveCapability(
		d.session,
		name,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// RemoveDriver removes an out-of-box(installed) driver from
// an image with .inf file path.
func (d *DismImageSession) RemoveDriver(driverPath string) error {
	return dismapi.DismRemoveDriver(d.session, driverPath)
}

// RemovePackage removes a package from an image.
//
// identifier can be an absolute path to a .cab file or package name.
func (d *DismImageSession) RemovePackage(identifier string, opts *DismProgressOpts) error {
	var o DismProgressOpts
	if opts != nil {
		o = *opts
	}

	pkgIdentifier := dismapi.DismPackageNone
	if identifier != "" {
		if strings.HasSuffix(strings.ToLower(identifier), ".cab") {
			pkgIdentifier = dismapi.DismPackagePath
		} else {
			pkgIdentifier = dismapi.DismPackageName
		}
	}

	return dismapi.DismRemovePackage(
		d.session,
		identifier,
		pkgIdentifier,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// RestoreImageHealth repairs an image which has been
// identified as repairable of CheckImageHealth.
func (d *DismImageSession) RestoreImageHealth(opts *DismRestoreImageHealthOpts) error {
	var o DismRestoreImageHealthOpts
	if opts != nil {
		o = *opts
	}

	return dismapi.DismRestoreImageHealth(
		d.session,
		o.SourcePaths,
		o.LimitAccess,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// CommitImage commits changes from mount point to an image.
func (d *DismImageSession) CommitImage(opts *DismCommitOpts) error {
	var o DismCommitOpts
	if opts != nil {
		o = *opts
	}

	var flags uint32
	if o.Append {
		flags |= dismapi.DISM_COMMIT_APPEND
	}
	if o.GenerateIntegrity {
		flags |= dismapi.DISM_COMMIT_GENERATE_INTEGRITY
	}
	if o.SupportEa {
		flags |= dismapi.DISM_COMMIT_SUPPORT_EA
	}

	return dismapi.DismCommitImage(
		d.session,
		flags,
		o.CancelEvent,
		o.Progress,
		o.UserData,
	)
}

// Close closes opened DismSession for [DismImageSession].
func (d *DismImageSession) Close() error {
	if err := dismapi.DismCloseSession(d.session); err != nil {
		return err
	}

	d.sesMapRef.Remove(d)

	return nil
}

// DismImageFile stores mount points and DismSession
// for Windows image file.
type DismImageFile struct {
	// path of .wim or .vhd(x) file
	imageFilePath string
	// mounted paths used in [dismapi.DismMountImage]
	mountPoints mapset.Set[string]
	// DismSessions from [dismapi.DismOpenSession], mapped with mounted paths
	sessions mapset.Set[*DismImageSession]
}

// NewDismImageFile creates structure for using DISM API, it
// runs [dismapi.DismInitialize] if DISM API is not initialized.
func NewDismImageFile(imageFilePath string) (*DismImageFile, error) {
	if err := initDism(); err != nil {
		return nil, err
	}

	dismImg := &DismImageFile{
		imageFilePath: imageFilePath,
		mountPoints:   mapset.NewSet[string](),
		sessions:      mapset.NewSet[*DismImageSession](),
	}

	curDismImg.Add(dismImg)

	return dismImg, nil
}

// ImageFilePath returns a path of .wim of .vhd file.
func (d *DismImageFile) ImageFilePath() string {
	return d.imageFilePath
}

// GetMountPoints returns currently mounted paths.
func (d *DismImageFile) GetMountPoints() []string {
	return d.mountPoints.ToSlice()
}

// GetSessions returns currently opened sessions from this image.
func (d *DismImageFile) GetSessions() []*DismImageSession {
	return d.sessions.ToSlice()
}

// Mount mounts an image file with imageIndex or with imageName if specified.
func (d *DismImageFile) Mount(mountPath string, imageIndex uint32, opts *DismMountOpts) error {
	var o DismMountOpts
	if opts != nil {
		o = *opts
	}

	identifier := dismapi.DismImageIndex
	if o.ImageName != "" {
		identifier = dismapi.DismImageName
	}

	var flags uint32 = dismapi.DISM_MOUNT_READWRITE
	if o.ReadOnly {
		flags = dismapi.DISM_MOUNT_READONLY
	}
	if o.Optimize {
		flags |= dismapi.DISM_MOUNT_OPTIMIZE
	}
	if o.CheckIntegrity {
		flags |= dismapi.DISM_MOUNT_CHECK_INTEGRITY
	}

	if err := dismapi.DismMountImage(d.imageFilePath, mountPath, imageIndex,
		o.ImageName, identifier, flags, o.CancelEvent, o.Progress,
		o.UserData); err != nil {
		return err
	}

	d.mountPoints.Add(mountPath)

	return nil
}

// Unmount unmounts the specified mount path.
func (d *DismImageFile) Unmount(mountPath string, opts *DismUnmountOpts) error {
	if !opts.Force && !d.mountPoints.ContainsOne(mountPath) {
		return ErrNotMounted
	}

	var o DismUnmountOpts
	if opts != nil {
		o = *opts
	}

	var flags uint32 = dismapi.DISM_DISCARD_IMAGE
	if o.Commit {
		flags = dismapi.DISM_COMMIT_IMAGE
		if o.Append {
			flags |= dismapi.DISM_COMMIT_APPEND
		}
		if o.GenerateIntegrity {
			flags |= dismapi.DISM_COMMIT_GENERATE_INTEGRITY
		}
		if o.SupportEa {
			flags |= dismapi.DISM_COMMIT_SUPPORT_EA
		}
	}

	if err := dismapi.DismUnmountImage(mountPath, flags, o.CancelEvent,
		o.Progress, o.UserData); err != nil {
		return err
	}

	d.mountPoints.Remove(mountPath)

	return nil
}

// OpenSession creates DismImageSession from a mount path, if windowsPath or
// systemDrive is empty, default value will be used.
func (d *DismImageFile) OpenSession(mountPath, windowsPath, systemDrive string) (*DismImageSession, error) {
	if !d.mountPoints.ContainsOne(mountPath) {
		return nil, ErrNotMounted
	}

	ses, err := dismapi.DismOpenSession(mountPath, windowsPath, systemDrive)
	if err != nil {
		return nil, err
	}

	newSession := &DismImageSession{
		session:   ses,
		mountPath: mountPath,
		sesMapRef: d.sessions,
	}

	d.sessions.Add(newSession)

	return newSession, nil
}

// GetImageInfo returns information of images in .wim or .vhd(x) file.
func (d *DismImageFile) GetImageInfo() ([]dismapi.GoDismImageInfo, error) {
	return dismapi.DismGetImageInfo(d.imageFilePath)
}

// Close cleans up all DismSessions and mount points associated with
// the image file, note that this will discard changes in mount points.
func (d *DismImageFile) Close() error {
	var err error

	for session := range d.sessions.Iter() {
		closeErr := session.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}

	for mntPath := range d.mountPoints.Iter() {
		unmountErr := dismapi.DismUnmountImage(mntPath, dismapi.DISM_DISCARD_IMAGE, 0, 0, nil)
		if unmountErr != nil {
			err = errors.Join(err, unmountErr)
		} else {
			d.mountPoints.Remove(mntPath)
		}
	}

	curDismImg.Remove(d)

	dismMutex.Lock()
	if dismInitialized && curDismImg.IsEmpty() {
		// shutdown DISM API if it is no longer used
		shutdownErr := dismapi.DismShutdown()
		if shutdownErr != nil && shutdownErr.(*w32api.WinimgInternalErr).Errno() !=
			dismapi.DISMAPI_E_DISMAPI_NOT_INITIALIZED {
			err = errors.Join(err, shutdownErr)
		} else {
			dismInitialized = false
		}
	}
	dismMutex.Unlock()

	return err
}

// RegisterWimLog registers logFile for logging wimgapi
// operations, flags is reserved and always 0.
func RegisterWimLog(logFile string, flags uint32) error {
	err := wimgapi.WIMRegisterLogFile(logFile, flags)
	if err != nil {
		return err
	}

	return nil
}

// UnregisterWimLog unregisters logfile for logging
// wimgapi operations used in [RegisterWimLog].
func UnregisterWimLog(logFile string) error {
	err := wimgapi.WIMUnregisterLogFile(logFile)
	if err != nil {
		return err
	}

	return nil
}

// WimVolumeImage stores handle and mount path for
// a volume image in .wim file.
type WimVolumeImage struct {
	handle windows.Handle
	// if not empty, the image handle is mounted
	mountPath string

	imageIndex uint32
	wimRef     *WimImageFile
}

// GetHandle returns handle of a volume image.
func (w *WimVolumeImage) GetHandle() windows.Handle {
	return w.handle
}

// GetMountPath returns mounted path of a volume image.
func (w *WimVolumeImage) GetMountPath() string {
	return w.mountPath
}

// GetImageInfo returns xml info of a volume image in bytes.
func (w *WimVolumeImage) GetImageInfo() ([]byte, error) {
	return wimgapi.WIMGetImageInformation(w.handle)
}

// SetImageInfo sets xml info of a image with bytes.
func (w *WimVolumeImage) SetImageInfo(buf []byte) error {
	return wimgapi.WIMSetImageInformation(w.handle, buf)
}

// Apply applies a volume image in .wim file to specified path.
func (w *WimVolumeImage) Apply(applyPath string, opts *WimApplyOpts) error {
	if !w.wimRef.tempPathSet {
		if err := w.wimRef.setTempPath(); err != nil {
			return err
		}
	}

	var o WimApplyOpts
	if opts != nil {
		o = *opts
	}

	var flags uint32
	if o.Verify {
		flags |= wimgapi.WIM_FLAG_VERIFY
	}
	if o.Index {
		flags |= wimgapi.WIM_FLAG_INDEX
	}
	if o.NoApply {
		flags |= wimgapi.WIM_FLAG_NO_APPLY
	}
	if o.NoDirAcl {
		flags |= wimgapi.WIM_FLAG_NO_DIRACL
	}
	if o.NoFileAcl {
		flags |= wimgapi.WIM_FLAG_NO_FILEACL
	}
	if o.NoReparseFix {
		flags |= wimgapi.WIM_FLAG_NO_RP_FIX
	}
	if o.FileInfo {
		flags |= wimgapi.WIM_FLAG_FILEINFO
	}
	if o.ConfirmTrustedFile {
		flags |= wimgapi.WIM_FLAG_APPLY_CI_EA
	}
	if o.WIMBoot {
		flags |= wimgapi.WIM_FLAG_WIM_BOOT
	}
	if o.Compact {
		flags |= wimgapi.WIM_FLAG_APPLY_COMPACT
	}
	if o.SupportEa {
		flags |= wimgapi.WIM_FLAG_SUPPORT_EA
	}

	if err := wimgapi.WIMApplyImage(w.handle, applyPath, flags); err != nil {
		return err
	}

	return nil
}

// Export transfers a volume image in .wim file to another .wim file.
func (w *WimVolumeImage) Export(wim *WimImageFile, opts *WimExportOpts) error {
	if wim == nil {
		return errors.New("destination file is nil")
	}

	if !w.wimRef.tempPathSet {
		if err := w.wimRef.setTempPath(); err != nil {
			return err
		}
	}
	if !wim.tempPathSet {
		if err := wim.setTempPath(); err != nil {
			return err
		}
	}

	var o WimExportOpts
	if opts != nil {
		o = *opts
	}

	var flags uint32
	if o.AllowDuplicate {
		flags |= wimgapi.WIM_EXPORT_ALLOW_DUPLICATES
	}
	if o.OnlyResources {
		flags |= wimgapi.WIM_EXPORT_ONLY_RESOURCES
	}
	if o.OnlyMetadata {
		flags |= wimgapi.WIM_EXPORT_ONLY_METADATA
	}
	if o.VerifySource {
		flags |= wimgapi.WIM_EXPORT_VERIFY_SOURCE
	}
	if o.VerifyDestination {
		flags |= wimgapi.WIM_EXPORT_VERIFY_DESTINATION
	}

	if err := wimgapi.WIMExportImage(w.handle, wim.handle, flags); err != nil {
		return err
	}

	return nil
}

// Mount mounts a volume image in .wim file to mountPath.
func (w *WimVolumeImage) Mount(mountPath string, opts *WimImageMountOpts) error {
	if w.mountPath != "" {
		return ErrAlreadyMounted
	}

	var o WimImageMountOpts
	if opts != nil {
		o = *opts
	}

	var flags uint32
	if o.ReadOnly {
		flags |= wimgapi.WIM_FLAG_MOUNT_READONLY
	}
	if o.Verify {
		flags |= wimgapi.WIM_FLAG_VERIFY
	}
	if o.NoReparseFix {
		flags |= wimgapi.WIM_FLAG_NO_RP_FIX
	}
	if o.NoDirAcl {
		flags |= wimgapi.WIM_FLAG_NO_DIRACL
	}
	if o.NoFileAcl {
		flags |= wimgapi.WIM_FLAG_NO_FILEACL
	}

	if err := wimgapi.WIMMountImageHandle(w.handle, mountPath, flags); err != nil {
		return err
	}

	w.mountPath = mountPath

	return nil
}

// Unmount unmounts a volume image in .wim file from mounted path.
func (w *WimVolumeImage) Unmount() error {
	if w.mountPath == "" {
		return ErrNotMounted
	}

	if err := wimgapi.WIMUnmountImageHandle(w.handle, 0); err != nil {
		return err
	}

	w.mountPath = ""

	return nil
}

// Close closes a handle of volume image in .wim, unmounts it if required.
func (w *WimVolumeImage) Close() error {
	var err error

	if w.mountPath != "" {
		err = errors.Join(err, w.Unmount())
	}

	closeErr := wimgapi.WIMCloseHandle(w.handle)
	if closeErr != nil {
		err = errors.Join(err, closeErr)
	} else {
		w.wimRef.imageHandles.Delete(w.imageIndex)
	}

	return err
}

// WimImageFile contains handle for .wim file and associated mount points.
type WimImageFile struct {
	// handle from [wimgapi.WIMCreateFile] with .wim file
	handle windows.Handle
	// handles from [wimgapi.WIMLoadImage] and [wimgapi.WIMCaptureImage],
	// mapped with image index
	imageHandles *xsync.MapOf[uint32, *WimVolumeImage]
	// path of .wim file
	imageFilePath string
	// temporary path for capture and apply
	tempPath string
	// image count in .wim file
	imageCount uint32
	// is temporary path set or created?
	tempPathSet, tempCreated bool
	// this .wim file is newly created
	isCreated bool
}

// NewWIMImageFile opens or creates image file, if opts.TempPath is empty,
// temporary directory will be created with random name if needed.
func NewWIMImageFile(imageFilePath string, opts *WimCreateFileOpts) (*WimImageFile, error) {
	var o WimCreateFileOpts
	if opts != nil {
		o = *opts
	}

	var access, createMode, flags uint32
	if !o.NoRead {
		access |= wimgapi.WIM_GENERIC_READ
	}
	if !o.NoWrite {
		access |= wimgapi.WIM_GENERIC_WRITE
	}
	if !o.NoMount {
		access |= wimgapi.WIM_GENERIC_MOUNT
	}

	switch o.CreateMode {
	case OpenExisting:
		createMode = wimgapi.WIM_OPEN_EXISTING
	case OpenOrCreate:
		createMode = wimgapi.WIM_OPEN_ALWAYS
	case CreateIfNotExist:
		createMode = wimgapi.WIM_CREATE_NEW
	case AlwaysCreate:
		createMode = wimgapi.WIM_CREATE_ALWAYS
	}

	wimHandle, created, err := wimgapi.WIMCreateFile(imageFilePath, access, createMode,
		flags, wimgapi.WimCompressionType(o.Compression), true)
	if err != nil {
		return nil, err
	}

	return &WimImageFile{
		handle:        wimHandle,
		imageHandles:  xsync.NewMapOf[uint32, *WimVolumeImage](),
		imageFilePath: imageFilePath,
		imageCount:    wimgapi.WIMGetImageCount(wimHandle),
		tempPath:      o.TempPath,
		isCreated:     created,
	}, nil
}

// setTempPath sets temporary path used for WIM operation.
func (w *WimImageFile) setTempPath() error {
	if w.tempPath == "" {
		w.tempPath, _ = os.MkdirTemp("", filepath.Base(w.imageFilePath))
		w.tempCreated = true
	}

	err := wimgapi.WIMSetTemporaryPath(w.handle, w.tempPath)
	if err != nil {
		return err
	}

	w.tempPathSet = true

	return nil
}

// IsCreated returns if this .wim file is newly created.
func (w *WimImageFile) IsCreated() bool {
	return w.isCreated
}

// GetHandle returns handle of an opened .wim file.
func (w *WimImageFile) GetHandle() windows.Handle {
	return w.handle
}

// GetImageCount updates and returns the number of
// volume images stored in .wim file.
func (w *WimImageFile) GetImageCount() uint32 {
	w.imageCount = wimgapi.WIMGetImageCount(w.handle)

	return w.imageCount
}

// GetImageFilePath returns path of the image(.wim) file.
func (w *WimImageFile) GetImageFilePath() string {
	return w.imageFilePath
}

// Capture captures a directory path and stores it in an .wim
// file which is added as a volume image and returns it.
func (w *WimImageFile) Capture(capturePath string, opts *WimCaptureOpts) (*WimVolumeImage, error) {
	if !w.tempPathSet {
		if err := w.setTempPath(); err != nil {
			return nil, err
		}
	}

	var o WimCaptureOpts
	if opts != nil {
		o = *opts
	}

	var flags uint32
	if o.Verify {
		flags |= wimgapi.WIM_FLAG_VERIFY
	}
	if o.NoReparseFix {
		flags |= wimgapi.WIM_FLAG_NO_RP_FIX
	}
	if o.NoDirAcl {
		flags |= wimgapi.WIM_FLAG_NO_DIRACL
	}
	if o.NoFileAcl {
		flags |= wimgapi.WIM_FLAG_NO_FILEACL
	}
	if o.WIMBoot {
		flags |= wimgapi.WIM_FLAG_WIM_BOOT
	}
	if o.SupportEa {
		flags |= wimgapi.WIM_FLAG_SUPPORT_EA
	}

	volHnd, err := wimgapi.WIMCaptureImage(w.handle, capturePath, flags)
	if err != nil {
		return nil, err
	}

	w.imageCount++

	newVolume := &WimVolumeImage{
		handle:     volHnd,
		imageIndex: w.imageCount,
		wimRef:     w,
	}

	w.imageHandles.Store(w.imageCount, newVolume)

	return newVolume, nil
}

// GetImageInfo returns xml info of .wim file in bytes.
func (w *WimImageFile) GetImageInfo() ([]byte, error) {
	return wimgapi.WIMGetImageInformation(w.handle)
}

func (w *WimImageFile) GetAttributes() (wimgapi.GoWimInfo, error) {
	return wimgapi.WIMGetAttributes(w.handle)
}

// LoadImage loads volume image from .wim file with image index.
func (w *WimImageFile) LoadImage(imageIndex uint32) (*WimVolumeImage, error) {
	if !w.tempPathSet {
		if err := w.setTempPath(); err != nil {
			return nil, err
		}
	}

	hnd, err := wimgapi.WIMLoadImage(w.handle, imageIndex)
	if err != nil {
		return nil, err
	}

	imgHandle := &WimVolumeImage{
		handle:     hnd,
		imageIndex: imageIndex,
		wimRef:     w,
	}

	w.imageHandles.Store(imageIndex, imgHandle)

	return imgHandle, err
}

// SetImageInfo sets information of a .wim file with xml info
// as bytes.
func (w *WimImageFile) SetImageInfo(buf []byte) error {
	return wimgapi.WIMSetImageInformation(w.handle, buf)
}

// RegisterMessageCallback registers callback to be used for imaging operation.
func (w *WimImageFile) RegisterMessageCallback(callback uintptr, userData unsafe.Pointer) error {
	_, err := wimgapi.WIMRegisterMessageCallback(w.handle, callback, userData)

	return err
}

// UnregisterMessageCallback unregisters callback previously registered
// with RegisterMessageCallback, using 0 unregisters all callbacks.
func (w *WimImageFile) UnregisterMessageCallback(callback uintptr) error {
	return wimgapi.WIMUnregisterMessageCallback(w.handle, callback)
}

// Close cleans up all mount points and volume image handles and
// unregisters all callbacks, note that this will discard changes
// made in mount points created with Mount().
func (w *WimImageFile) Close() error {
	// unregister all callback functions
	err := wimgapi.WIMUnregisterMessageCallback(w.handle, 0)

	w.imageHandles.Range(func(index uint32, imgHandle *WimVolumeImage) bool {
		if imgHandle.mountPath != "" {
			unmountErr := wimgapi.WIMUnmountImageHandle(imgHandle.handle, 0)
			if unmountErr != nil {
				err = errors.Join(err, unmountErr)
			} else {
				imgHandle.mountPath = ""
			}
		}

		closeErr := imgHandle.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}

		return true
	})

	if w.handle != 0 {
		err = errors.Join(err, wimgapi.WIMCloseHandle(w.handle))
	}

	// remove temporary directory if it was created
	if w.tempCreated {
		err = errors.Join(err, os.RemoveAll(w.tempPath))
	}

	return err
}

// MountWimImage mounts an image with mount path, .wim file path and image
// index, if tempPath is empty, the image will not be mounted for edits.
func MountWimImage(mountPath, imageFilePath string, imageIndex uint32, tempPath string) error {
	err := wimgapi.WIMMountImage(mountPath, imageFilePath, imageIndex, tempPath)
	if err != nil {
		return err
	}

	return nil
}

// UnmountWimImage unmounts an image with mountPath and imageIndex, if
// commitChanges is true, changes in mounted directory will be saved to
// .wim file only if tempPath was specified in [MountWimImage].
func UnmountWimImage(mountPath, imageFilePath string, imageIndex uint32, commitChanges bool) error {
	err := wimgapi.WIMUnmountImage(mountPath, imageFilePath, imageIndex, commitChanges)
	if err != nil {
		return err
	}

	return nil
}

// RemountWimImage remounts an image mount to mountPath, it maps
// contents of mounted image volume to the directory. The opts
// is currently unused.
func RemountWimImage(mountPath string, opts *WimRemountOpts) error {
	var flags uint32

	if err := wimgapi.WIMRemountImage(mountPath, flags); err != nil {
		return err
	}

	return nil
}

type WinImage struct {
	filePath string

	DismImage *DismImageFile
	WimImage  *WimImageFile
}

func (w *WinImage) Path() string {
	return w.filePath
}
