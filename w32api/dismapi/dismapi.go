// Package dismapi calls Deployment Image Servicing and Management (DISM) API.
//
// https://learn.microsoft.com/en-us/windows-hardware/manufacture/desktop/dism/dism-api-reference
package dismapi

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/internal/utils"
	"github.com/Snshadow/winimg/w32api"
)

var (
	ErrDelete = errors.New("failed to release resource") // DismDelete failure
)

//sys	dismInitialize (logLevel DismLogLevel, logFilePath *uint16, scratchDirectory *uint16) (ret error) = dismapi.DismInitialize
//sys	dismShutdown() (ret error) = dismapi.DismShutdown
//sys	dismMountImage(imageFilePath *uint16, mountPath *uint16, imageIndex uint32, imageName *uint16, imageIdentifier DismImageIdentifier, flags uint32, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismMountImage
//sys	dismUnmountImage(mountPath *uint16, flags uint32, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismUnmountImage
//sys	dismOpenSession(imagePath *uint16, windowsDirectory *uint16, systemDrive *uint16, session *DismSession) (ret error) = dismapi.DismOpenSession
//sys	dismCloseSession(session DismSession) (ret error) = dismapi.DismCloseSession
//sys	dismGetLastErrorMessage(errorMessage **DismString) (ret error) = dismapi.DismGetLastErrorMessage
//sys	dismRemountImage(mountPath *uint16) (ret error) = dismapi.DismRemountImage
//sys	dismCommitImage(session DismSession, flags uint32, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismCommitImage
//sys	dismGetImageInfo(imageFilePath *uint16, imageInfo **DismImageInfo, count *uint32) (ret error) = dismapi.DismGetImageInfo
//sys	dismGetMountedImageInfo(mountedImageInfo **DismMountedImageInfo, count *uint32) (ret error) = dismapi.DismGetMountedImageInfo
//sys	dismCleanupMountPoints() (ret error) = dismapi.DismCleanupMountPoints
//sys	dismCheckImageHealth(session DismSession, scanImage bool, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer, imageHealth *DismImageHealthState) (ret error) = dismapi.DismCheckImageHealth
//sys	dismRestoreImageHealth(session DismSession, sourcePaths **uint16, sourcePathCount uint32, limitAccess bool, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismRestoreImageHealth
//sys	dismDelete(dismStructure unsafe.Pointer) (ret error) = dismapi.DismDelete
//sys	dismAddPackage(session DismSession, packagePath *uint16, ignoreCheck bool, preventPending bool, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismAddPackage
//sys	dismRemovePackage(session DismSession, identifier *uint16, packageIdentifier DismPackageIdentifier, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismRemovePackage
//sys	dismEnableFeature(session DismSession, featureName *uint16, identifier *uint16, packageIdentifier DismPackageIdentifier, limitAccess bool, sourcePaths **uint16, sourcePathCount uint32, enableAll bool, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismEnableFeature
//sys	dismDisableFeature(session DismSession, featureName *uint16, packageName *uint16, removePayload bool, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismDisableFeature
//sys	dismGetPackages(session DismSession, pPackage **DismPackage, count *uint32) (ret error) = dismapi.DismGetPackages
//sys	dismGetPackageInfo(session DismSession, identifier *uint16, packageIdentifier DismPackageIdentifier, packageInfo **DismPackageInfo) (ret error) = dismapi.DismGetPackageInfo
//sys	dismGetPackageInfoEx(session DismSession, identifier *uint16, packageIdentifier DismPackageIdentifier, packageInfoEx **DismPackageInfoEx) (ret error) = dismapi.DismGetPackageInfoEx
//sys	dismGetFeatures(session DismSession, identifier *uint16, packageIdentifier DismPackageIdentifier, feature **DismFeature, count *uint32) (ret error) = dismapi.DismGetFeatures
//sys	dismGetFeaturesEx(session DismSession, identifier *uint16, packageIdentifier DismPackageIdentifier, flags uint32, feature **DismFeatureEx, count *uint32) (ret error) = dismapi._DismGetFeaturesEx
//sys	dismGetFeatureInfo(session DismSession, featureName *uint16, identifier *uint16, packageIdentifier DismPackageIdentifier, featureInfo **DismFeatureInfo) (ret error) = dismapi.DismGetFeatureInfo
//sys	dismGetFeatureParent(session DismSession, featureName *uint16, identifier *uint16, packageIdentifier DismPackageIdentifier, feature **DismFeature, count *uint32) (ret error) = dismapi.DismGetFeatureParent
//sys	dismApplyUnattend(session DismSession, unattendFile *uint16, singleSession bool) (ret error) = dismapi.DismApplyUnattend
//sys	dismAddDriver(session DismSession, driverPath *uint16, forceUnsigned bool) (ret error) = dismapi.DismAddDriver
//sys	dismRemoveDriver(session DismSession, driverPath *uint16) (ret error) = dismapi.DismRemoveDriver
//sys	dismGetDrivers(session DismSession, allDrivers bool, driverPackage **DismDriverPackage, count *uint32) (ret error) = dismapi.DismGetDrivers
//sys	dismGetDriverInfo(session DismSession, driverPath *uint16, driver **DismDriver, count *uint32, driverPackage **DismDriverPackage) (ret error) = dismapi.DismGetDriverInfo
//sys	dismGetCapabilities(session DismSession, capability **DismCapability, count *uint32) (ret error) = dismapi.DismGetCapabilities
//sys	dismGetCapabilityInfo(session DismSession, name *uint16, info **DismCapabilityInfo) (ret error) = dismapi.DismGetCapabilityInfo
//sys	dismAddCapability(session DismSession, name *uint16, limitAccess bool, sourcePaths **uint16, sourcePathCount uint32, cancelEvent windows.Handle, progress uintptr, UserData unsafe.Pointer) (ret error) = dismapi.DismAddCapability
//sys	dismRemoveCapability(session DismSession, name *uint16, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismRemoveCapability
//sys	dismCleanImage(session DismSession, cleanupType DismCleanupType, flags uint32, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi._DismCleanImage
//sys	dismGetReservedStorageState(session DismSession, state *uint32) (ret error) = dismapi.DismGetReservedStorageState
//sys	dismSetReservedStorageState(session DismSession, state uint32) (ret error) = dismapi.DismSetReservedStorageState
//sys	dismGetProvisionedAppxPackages(session DismSession, pPackage **DismAppxPackage, count *uint32) (ret error) = dismapi.DismGetProvisionedAppxPackages
//sys	dismAddProvisionedAppxPackage(session DismSession, appPath *uint16, dependencyPackages **uint16, dependencyPackageCount uint32, optionalPackages **uint16, optionalPackagesCount uint32, licensePaths **uint16, licensePathsCount uint32, skipLicense bool, customDataPath *uint16, region *uint16, stubPackageOption DismStubPackageOption) (ret error) = dismapi.DismAddProvisionedAppxPackage
//sys	dismRemoveProvisionedAppxPackage(session DismSession, packageName *uint16) (ret error) = dismapi.DismRemoveProvisionedAppxPackage
//sys	dismAddLanguage(session DismSession, languageName *uint16, preventPending bool, limitAccess bool, sourcePaths **uint16, sourcePathCount uint32, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismAddLanguage
//sys	dismRemoveLanguage(session DismSession, languageName *uint16, cancelEvent windows.Handle, progress uintptr, userData unsafe.Pointer) (ret error) = dismapi.DismRemoveLanguage

// DismInitialize initializes DISM API. DismInitialize must be called
// once per process, before calling any other DISM API functions.
func DismInitialize(
	logLevel DismLogLevel,
	logFilePath string,
	scratchDirectory string,
) error {
	var u16LogFilePath, u16ScratchDir *uint16
	var err error

	if logFilePath != "" {
		if u16LogFilePath, err = windows.UTF16PtrFromString(logFilePath); err != nil {
			return err
		}
	}

	if scratchDirectory != "" {
		if u16ScratchDir, err = windows.UTF16PtrFromString(scratchDirectory); err != nil {
			return err
		}
	}

	if err = dismInitialize(
		logLevel,
		u16LogFilePath,
		u16ScratchDir,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismShutdown shuts down DISM API. DismShutdown must be called once
// per process. Other DISM API function calls will fail after
// DismShutdown has been called.
func DismShutdown() error {
	if err := dismShutdown(); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismMountImage mounts a WIM or VHD image file to a
// specified location.
func DismMountImage(
	imageFilePath string,
	mountPath string,
	imageIndex uint32,
	imageName string,
	imageIdentifier DismImageIdentifier,
	flags uint32,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	var u16ImgFilePath, u16MntPath, u16ImgName *uint16
	var err error

	if u16ImgFilePath, err = windows.UTF16PtrFromString(imageFilePath); err != nil {
		return err
	}

	if u16MntPath, err = windows.UTF16PtrFromString(mountPath); err != nil {
		return err
	}

	if imageName != "" {
		if u16ImgName, err = windows.UTF16PtrFromString(imageName); err != nil {
			return err
		}
	}

	if err = dismMountImage(
		u16ImgFilePath,
		u16MntPath,
		imageIndex,
		u16ImgName,
		imageIdentifier,
		flags,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismUnmountImage unmounts a Windows image from
// a specified location.
func DismUnmountImage(
	mountPath string,
	flags uint32,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	u16MntPath, err := windows.UTF16PtrFromString(mountPath)
	if err != nil {
		return err
	}

	if err = dismUnmountImage(
		u16MntPath,
		flags,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismOpenSession associates an offline or online
// Windows image with a [DismSession].
func DismOpenSession(
	imagePath string,
	windowsDirectory string,
	systemDrive string,
) (DismSession, error) {
	var u16ImgPath, u16WinDir, u16SysDrv *uint16
	var err error

	if u16ImgPath, err = windows.UTF16PtrFromString(imagePath); err != nil {
		return 0, err
	}

	if windowsDirectory != "" {
		if u16WinDir, err = windows.UTF16PtrFromString(windowsDirectory); err != nil {
			return 0, err
		}
	}

	if systemDrive != "" {
		if u16SysDrv, err = windows.UTF16PtrFromString(systemDrive); err != nil {
			return 0, err
		}
	}

	var ses DismSession

	if err = dismOpenSession(
		u16ImgPath,
		u16WinDir,
		u16SysDrv,
		&ses,
	); err != nil {
		return 0, w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return ses, nil
}

// DismCloseSession closes a [DismSession] opened
// with [DismOpenSession].
func DismCloseSession(session DismSession) error {
	if err := dismCloseSession(session); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismGetLastErrorMessage gets the error message in the current
// thread, immediately after a failure.
//
// DismGetLastErrorMessage does not apply to the [DismShutdown] function,
// [DismDelete] function, or the DismGetLastErrorMessage function.
func DismGetLastErrorMessage() (string, error) {
	var errorMsg *DismString

	if err := dismGetLastErrorMessage(&errorMsg); err != nil {
		return "", w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	lastErrMsg := windows.UTF16PtrToString(errorMsg.Value)
	n := len(lastErrMsg)
	for ; n > 0 && lastErrMsg[n-1] == '\n' || lastErrMsg[n-1] == '\r'; n-- {
	}

	if delErr := DismDelete(unsafe.Pointer(errorMsg)); delErr != nil {
		return lastErrMsg[:n], errors.Join(ErrDelete, delErr)
	}

	return lastErrMsg[:n], nil
}

// DismRemountImage remounts a previously mounted Windows image
// from the .wim or .vhd file at the path specified by MountPath.
// Use [DismOpenSession] to associate the image with a
// [DismSession] after it is remounted.
func DismRemountImage(mountPath string) error {

	u16MntPath, err := windows.UTF16PtrFromString(mountPath)
	if err != nil {
		return err
	}

	if err = dismRemountImage(u16MntPath); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismCommitImage commit the changes made to a Windows image
// in a mounted .wim or .vhd file. The image must be mounted
// using the [DismMountImage].
func DismCommitImage(
	session DismSession,
	flags uint32,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	if err := dismCommitImage(
		session,
		flags,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismGetImageInfo returns an array of [GoDismImageInfo]
// structures that describe the images in a .wim or
// .vhd file.
func DismGetImageInfo(imageFilePath string) (imageInfo []GoDismImageInfo, err error) {
	u16ImgFilePath, err := windows.UTF16PtrFromString(imageFilePath)
	if err != nil {
		return nil, err
	}

	var infoPtr *DismImageInfo
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(infoPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetImageInfo(
		u16ImgFilePath,
		&infoPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
		return
	}

	if infoPtr != nil && count != 0 {
		bufSize := GetPackedSize(*infoPtr)
		stPtr := (*byte)(unsafe.Pointer(infoPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismImageInfo](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismImageInfo
			goSt.fill(&unpacked)

			imageInfo = append(imageInfo, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismMountedImageinfo returns an array of
// [GoDismMountedImageInfo] structures describing
// currently mounted images.
func DismGetMountedImageInfo() (mountedImageInfo []GoDismMountedImageInfo, err error) {
	var infoPtr *DismMountedImageInfo
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(infoPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetMountedImageInfo(
		&infoPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
		return
	}

	if infoPtr != nil && count != 0 {
		bufSize := GetPackedSize(*infoPtr)
		stPtr := (*byte)(unsafe.Pointer(infoPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismMountedImageInfo](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismMountedImageInfo
			goSt.fill(&unpacked)

			mountedImageInfo = append(mountedImageInfo, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismCleanupMountPoints removes files and releases
// resources associated with corrupted or invalid
// mount paths.
func DismCleanupMountPoints() error {
	err := dismCleanupMountPoints()
	if err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismCheckImageHealth checks whether the image
// can be serviced or is corrupted.
func DismCheckImageHealth(
	session DismSession,
	scanImage bool,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) (
	imageHealth DismImageHealthState,
	err error,
) {
	if err = dismCheckImageHealth(
		session,
		scanImage,
		cancelEvent,
		progress,
		userData,
		&imageHealth,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return
}

// DismRestoreImagehealth Repairs a corrupted
// image that has been identified as repairable
// by [DismCheckImageHealth].
func DismRestoreImageHealth(
	session DismSession,
	sourcePaths []string,
	limitAccess bool,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	u16SrcPaths, err := utils.StrSliceToUtf16PtrArr(sourcePaths)
	if err != nil {
		return err
	}

	var srcPtr **uint16

	if len(u16SrcPaths) != 0 {
		srcPtr = &u16SrcPaths[0]
	}

	if err = dismRestoreImageHealth(
		session,
		srcPtr,
		uint32(len(u16SrcPaths)),
		limitAccess,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismDelete releases resources held by a structure
// or an array of structures returned by other DISM
// API Functions.
func DismDelete(dismStructure unsafe.Pointer) error {
	if err := dismDelete(dismStructure); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismAddPackage adds a single .cab or .msu file
// to a Windows image.
func DismAddPackage(
	session DismSession,
	packagePath string,
	ignoreCheck bool,
	preventPending bool,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	u16PkgPath, err := windows.UTF16PtrFromString(packagePath)
	if err != nil {
		return err
	}

	if err = dismAddPackage(
		session,
		u16PkgPath,
		ignoreCheck,
		preventPending,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismRemovePackage removes a package from an image.
func DismRemovePackage(
	session DismSession,
	identifer string,
	packageIdentifer DismPackageIdentifier,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	var u16Id *uint16
	var err error

	if identifer != "" {
		if u16Id, err = windows.UTF16PtrFromString(identifer); err != nil {
			return err
		}
	}

	if err = dismRemovePackage(
		session,
		u16Id,
		packageIdentifer,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismEnableFeature enables a feature in an image.
// Features are identified by a name and can optionally
// be tied to a package.
func DismEnableFeature(
	session DismSession,
	featureName string,
	identifier string,
	packageIdentifer DismPackageIdentifier,
	limitAccess bool,
	sourcePaths []string,
	enableAll bool,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	var u16Feat, u16Id *uint16
	var err error

	if u16Feat, err = windows.UTF16PtrFromString(featureName); err != nil {
		return err
	}

	if identifier != "" {
		if u16Id, err = windows.UTF16PtrFromString(identifier); err != nil {
			return err
		}
	}

	u16SrcPaths, err := utils.StrSliceToUtf16PtrArr(sourcePaths)
	if err != nil {
		return err
	}

	var srcPtr **uint16
	if len(u16SrcPaths) != 0 {
		srcPtr = &u16SrcPaths[0]
	}

	if err = dismEnableFeature(
		session,
		u16Feat,
		u16Id,
		packageIdentifer,
		limitAccess,
		srcPtr,
		uint32(len(u16SrcPaths)),
		enableAll,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismDisableFeature disables a feature in the current image.
func DismDisableFeature(
	session DismSession,
	featureName string,
	packageName string,
	removePayload bool,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	var u16FeatName, u16PkgName *uint16
	var err error

	if u16FeatName, err = windows.UTF16PtrFromString(featureName); err != nil {
		return err
	}

	if packageName != "" {
		if u16PkgName, err = windows.UTF16PtrFromString(packageName); err != nil {
			return err
		}
	}

	if err = dismDisableFeature(
		session,
		u16FeatName,
		u16PkgName,
		removePayload,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismGetPackages gets an array of [GoDismPackage]
// from an image.
func DismGetPackages(session DismSession) (packages []GoDismPackage, err error) {
	var pkgPtr *DismPackage
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(pkgPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetPackages(
		session,
		&pkgPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
		return
	}

	if pkgPtr != nil && count != 0 {
		bufSize := GetPackedSize(*pkgPtr)
		stPtr := (*byte)(unsafe.Pointer(pkgPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismPackage](unsafe.Slice(stPtr, bufSize))
			var goSt GoDismPackage
			goSt.fill(&unpacked)

			packages = append(packages, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismGetPackageInfo gets standard package properties
// as [GoDismPackageInfo].
func DismGetPackageInfo(
	session DismSession,
	identifier string,
	packageIdentifier DismPackageIdentifier,
) (packageInfo GoDismPackageInfo, err error) {
	u16Id, err := windows.UTF16PtrFromString(identifier)
	if err != nil {
		return
	}

	var pkgInfoPtr *DismPackageInfo

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(pkgInfoPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetPackageInfo(
		session,
		u16Id,
		packageIdentifier,
		&pkgInfoPtr,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
		return
	}

	if pkgInfoPtr != nil {
		buf := unsafe.Slice((*byte)(unsafe.Pointer(pkgInfoPtr)), GetPackedSize(*pkgInfoPtr))

		unpacked, _ := ToStruct[DismPackageInfo](buf)

		packageInfo.fill(&unpacked)
	}

	return
}

// DismGetPackageInfoEx gets standard packages properties as
// [DismGetPackageInfo] and CapabilityId in [GoDismPackageInfoEx].
func DismGetPackageInfoEx(
	session DismSession,
	identifier string,
	packageIdentifier DismPackageIdentifier,
) (packageInfoEx GoDismPackageInfoEx, err error) {
	u16Id, err := windows.UTF16PtrFromString(identifier)
	if err != nil {
		return
	}

	var pkgInfoPtrEx *DismPackageInfoEx

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(pkgInfoPtrEx)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetPackageInfoEx(
		session,
		u16Id,
		packageIdentifier,
		&pkgInfoPtrEx,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
		return
	}

	if pkgInfoPtrEx != nil {
		buf := unsafe.Slice((*byte)(unsafe.Pointer(pkgInfoPtrEx)), GetPackedSize(*pkgInfoPtrEx))

		unpacked, _ := ToStruct[DismPackageInfoEx](buf)

		packageInfoEx.fill(&unpacked)
	}

	return
}

// DismGetFeatures gets all the features in an image as an
// array of [GoDismFeature] regardless of whether the features
// are enabled or disabled.
func DismGetFeatures(
	session DismSession,
	identifier string,
	packageIdentifier DismPackageIdentifier,
) (feature []GoDismFeature, err error) {
	u16Id, err := windows.UTF16PtrFromString(identifier)
	if err != nil {
		return
	}

	var featPtr *DismFeature
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(featPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetFeatures(
		session,
		u16Id,
		packageIdentifier,
		&featPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))

		return
	}

	if featPtr != nil && count != 0 {
		bufSize := GetPackedSize(*featPtr)
		stPtr := (*byte)(unsafe.Pointer(featPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismFeature](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismFeature
			goSt.fill(&unpacked)

			feature = append(feature, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismGetFeatures gets all the features in an image as an array
// of [GoDismFeatureEx] regardless of whether the features are
// enabled or disabled.
func DismGetFeaturesEx(
	session DismSession,
	identifier string,
	packageIdentifier DismPackageIdentifier,
	flags uint32,
) (feature []GoDismFeatureEx, err error) {
	u16Id, err := windows.UTF16PtrFromString(identifier)
	if err != nil {
		return
	}

	var featExPtr *DismFeatureEx
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(featExPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetFeaturesEx(
		session,
		u16Id,
		packageIdentifier,
		flags,
		&featExPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))

		return
	}

	if featExPtr != nil && count != 0 {
		bufSize := GetPackedSize(*featExPtr)
		stPtr := (*byte)(unsafe.Pointer(featExPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismFeatureEx](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismFeatureEx
			goSt.fill(&unpacked)

			feature = append(feature, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismGetFeatureInfo gets detailed info for the
// specified feature as [GoDismFeatureInfo].
func DismGetFeatureInfo(
	session DismSession,
	featureName string,
	identifier string,
	packageIdentifier DismPackageIdentifier,
) (featureInfo GoDismFeatureInfo, err error) {
	var u16featName, u16Id *uint16
	if u16featName, err = windows.UTF16PtrFromString(featureName); err != nil {
		return
	}

	if identifier != "" {
		if u16Id, err = windows.UTF16PtrFromString(identifier); err != nil {
			return
		}
	}

	var featInfoPtr *DismFeatureInfo

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(featInfoPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetFeatureInfo(
		session,
		u16featName,
		u16Id,
		packageIdentifier,
		&featInfoPtr,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	if featInfoPtr != nil {
		buf := unsafe.Slice((*byte)(unsafe.Pointer(featInfoPtr)), GetPackedSize(*featInfoPtr))

		unpacked, _ := ToStruct[DismFeatureInfo](buf)
		featureInfo.fill(&unpacked)
	}

	return
}

// DismGetFeatureParent gets the parent features
// of a specified feature as [GoDismFeature].
func DismGetFeatureParent(
	session DismSession,
	featureName string,
	identifier string,
	packageIdentifier DismPackageIdentifier,
) (feature []GoDismFeature, err error) {
	var u16FeatName, u16Id *uint16
	if u16FeatName, err = windows.UTF16PtrFromString(featureName); err != nil {
		return
	}

	if identifier != "" {
		if u16Id, err = windows.UTF16PtrFromString(identifier); err != nil {
			return
		}
	}

	var featPtr *DismFeature
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(featPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetFeatureParent(
		session,
		u16FeatName,
		u16Id,
		packageIdentifier,
		&featPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
		return
	}

	if featPtr != nil && count != 0 {
		bufSize := GetPackedSize(*featPtr)
		stPtr := (*byte)(unsafe.Pointer(featPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismFeature](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismFeature
			goSt.fill(&unpacked)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismApplyUnattend applies an unattended
// answer file to a Windows image.
func DismApplyUnattend(
	session DismSession,
	unattendFile string,
	singleSession bool,
) error {
	u16unattend, err := windows.UTF16PtrFromString(unattendFile)
	if err != nil {
		return err
	}

	if err = dismApplyUnattend(
		session,
		u16unattend,
		singleSession,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismAddDriver adds a third party driver (.inf)
// to an offline Windows image.
func DismAddDriver(
	session DismSession,
	driverPath string,
	forceUnsigned bool,
) error {
	u16DrvPath, err := windows.UTF16PtrFromString(driverPath)
	if err != nil {
		return err
	}

	if err = dismAddDriver(
		session,
		u16DrvPath,
		forceUnsigned,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismRemoveDriver removes a third party driver (.inf)
// in an offline Windows image.
func DismRemoveDriver(
	session DismSession,
	driverPath string,
) error {
	u16DrvPath, err := windows.UTF16PtrFromString(driverPath)
	if err != nil {
		return err
	}

	if err = dismRemoveDriver(
		session,
		u16DrvPath,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismGetDriver lists the drivers in an image
// as an array of [GoDismDriverPackage].
func DismGetDrivers(
	session DismSession,
	allDrivers bool,
) (driverPackage []GoDismDriverPackage, err error) {
	var drvPkgPtr *DismDriverPackage
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(drvPkgPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetDrivers(
		session,
		allDrivers,
		&drvPkgPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
		return
	}

	if drvPkgPtr != nil && count != 0 {
		bufSize := GetPackedSize(*drvPkgPtr)
		stPtr := (*byte)(unsafe.Pointer(drvPkgPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismDriverPackage](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismDriverPackage
			goSt.fill(&unpacked)

			driverPackage = append(driverPackage, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismGetDriverInfo gets information about an .inf file
// in a specified image as an array of [GoDismDriver] and
// [GoDismDriverPackage].
func DismGetDriverInfo(
	session DismSession,
	driverPath string,
	getPackage bool,
) (
	driver []GoDismDriver,
	driverPackage *GoDismDriverPackage,
	err error,
) {
	u16DrvPath, err := windows.UTF16PtrFromString(driverPath)
	if err != nil {
		return
	}

	var packageParam **DismDriverPackage // optional
	if getPackage {
		var packagePtr *DismDriverPackage
		packageParam = &packagePtr
	}

	var drvPtr *DismDriver
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(drvPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
		if packageParam != nil {
			if delErr := DismDelete(unsafe.Pointer(*packageParam)); delErr != nil {
				err = errors.Join(err, ErrDelete, delErr)
			}
		}
	}()

	if err = dismGetDriverInfo(
		session,
		u16DrvPath,
		&drvPtr,
		&count,
		packageParam,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
		return
	}

	if packageParam != nil && *packageParam != nil {
		unpacked, _ := ToStruct[DismDriverPackage](unsafe.Slice((*byte)(unsafe.Pointer(*packageParam)), GetPackedSize(**packageParam)))

		driverPackage = &GoDismDriverPackage{}
		driverPackage.fill(&unpacked)
	}

	if drvPtr != nil && count != 0 {
		bufSize := GetPackedSize(*drvPtr)
		stPtr := (*byte)(unsafe.Pointer(drvPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismDriver](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismDriver
			goSt.fill(&unpacked)

			driver = append(driver, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismGetCapabilties gets DISM capabilities as an array
// of [GoDismCapability].
func DismGetCapabilities(session DismSession) (capability []GoDismCapability, err error) {
	var capPtr *DismCapability
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(capPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetCapabilities(
		session,
		&capPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))

		return
	}

	if capPtr != nil && count != 0 {
		bufSize := GetPackedSize(*capPtr)
		stPtr := (*byte)(unsafe.Pointer(capPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismCapability](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismCapability
			goSt.fill(&unpacked)

			capability = append(capability, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismGetCapabilityInfo gets information of
// DISM capability as [GoDismCapabilityInfo].
func DismGetCapabilityInfo(
	session DismSession,
	name string,
) (info GoDismCapabilityInfo, err error) {
	u16Name, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return
	}

	var infoPtr *DismCapabilityInfo

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(infoPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if err = dismGetCapabilityInfo(
		session,
		u16Name,
		&infoPtr,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))

		return
	}

	if infoPtr != nil {
		unpacked, _ := ToStruct[DismCapabilityInfo](unsafe.Slice((*byte)(unsafe.Pointer(infoPtr)), GetPackedSize(*infoPtr)))

		info.fill(&unpacked)
	}

	return
}

// DismAddCapability adds a capability to an image.
func DismAddCapability(
	session DismSession,
	name string,
	limitAccess bool,
	sourcePaths []string,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	u16Name, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	u16SrcPaths, err := utils.StrSliceToUtf16PtrArr(sourcePaths)
	if err != nil {
		return err
	}

	var srcPtr **uint16
	if len(u16SrcPaths) != 0 {
		srcPtr = &u16SrcPaths[0]
	}

	if err = dismAddCapability(
		session,
		u16Name,
		limitAccess,
		srcPtr,
		uint32(len(u16SrcPaths)),
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismRemoveCapability removes a capability from an image.
func DismRemoveCapability(
	session DismSession,
	name string,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	u16Name, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	if err = dismRemoveCapability(
		session,
		u16Name,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismCleanImage cleans up superseded elements in an image.
func DismCleanImage(
	session DismSession,
	cleanupType DismCleanupType,
	flags uint32,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	if err := dismCleanImage(
		session,
		cleanupType,
		flags,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismGetReservedStorageState gets state of
// reserved storage from an image.
func DismGetReservedStorageState(session DismSession) (state uint32, err error) {
	if err = dismGetReservedStorageState(
		session,
		&state,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return
}

// DismSetReservedStorageState sets state of
// reserved storage to an image.
func DismSetReservedStorageState(
	session DismSession,
	state uint32,
) error {
	if err := dismSetReservedStorageState(
		session,
		state,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismGetProvisionedAppxPackage gets provisioned
// appx packages from an image as an array of
// [GoDismAppxPackage].
func DismGetProvisionedAppxPackages(session DismSession) (appxPackage []GoDismAppxPackage, err error) {
	var pkgPtr *DismAppxPackage
	var count uint32

	defer func() {
		if delErr := DismDelete(unsafe.Pointer(pkgPtr)); delErr != nil {
			err = errors.Join(err, ErrDelete, delErr)
		}
	}()

	if dismGetProvisionedAppxPackages(
		session,
		&pkgPtr,
		&count,
	); err != nil {
		err = w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))

		return
	}

	if pkgPtr != nil && count != 0 {
		bufSize := GetPackedSize(*pkgPtr)
		stPtr := (*byte)(unsafe.Pointer(pkgPtr))

		for i := uint32(0); i < count; i++ {
			unpacked, _ := ToStruct[DismAppxPackage](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismAppxPackage
			goSt.fill(&unpacked)

			appxPackage = append(appxPackage, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}

	return
}

// DismAddProvisionedAppxPackage adds provisioned
// appx package to an image.
func DismAddProvisionedAppxPackage(
	session DismSession,
	appPath string,
	dependencyPackages []string,
	optionalPackages []string,
	licensePaths []string,
	skipLicense bool,
	customDataPath string,
	region string,
	stubPackageOption DismStubPackageOption,
) error {
	var (
		u16AppPath, u16CustomDataPath, u16Region *uint16
		depPtr, optPtr, licensePtr               **uint16
		depLen, optLen, licenseLen               uint32
		err                                      error
	)

	if u16AppPath, err = windows.UTF16PtrFromString(appPath); err != nil {
		return err
	}

	if customDataPath != "" {
		if u16CustomDataPath, err = windows.UTF16PtrFromString(customDataPath); err != nil {
			return err
		}
	}

	if region != "" {
		if u16Region, err = windows.UTF16PtrFromString(region); err != nil {
			return err
		}
	}

	if len(dependencyPackages) != 0 {
		u16Dep, err := utils.StrSliceToUtf16PtrArr(dependencyPackages)
		if err != nil {
			return err
		}

		depPtr = &u16Dep[0]
		depLen = uint32(len(u16Dep))
	}

	if len(optionalPackages) != 0 {
		u16Opt, err := utils.StrSliceToUtf16PtrArr(optionalPackages)
		if err != nil {
			return err
		}

		optPtr = &u16Opt[0]
		optLen = uint32(len(u16Opt))
	}

	if len(licensePaths) != 0 {
		u16License, err := utils.StrSliceToUtf16PtrArr(licensePaths)
		if err != nil {
			return err
		}

		licensePtr = &u16License[0]
		licenseLen = uint32(len(u16License))
	}

	if err = dismAddProvisionedAppxPackage(
		session,
		u16AppPath,
		depPtr,
		depLen,
		optPtr,
		optLen,
		licensePtr,
		licenseLen,
		skipLicense,
		u16CustomDataPath,
		u16Region,
		stubPackageOption,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismRemoveProvisionedAppxPackage removes provisioned
// appx package from an image.
func DismRemoveProvisionedAppxPackage(
	session DismSession,
	packageName string,
) error {
	var u16PkgName *uint16
	var err error

	if packageName != "" {
		if u16PkgName, err = windows.UTF16PtrFromString(packageName); err != nil {
			return err
		}
	}

	if err = dismRemoveProvisionedAppxPackage(
		session,
		u16PkgName,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismAddLanguage adds language to an image.
func DismAddLanguage(
	session DismSession,
	languageName string,
	preventPending bool,
	limitAccess bool,
	sourcePaths []string,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	var u16LangName *uint16
	var err error

	if languageName != "" {
		if u16LangName, err = windows.UTF16PtrFromString(languageName); err != nil {
			return err
		}
	}

	u16SrcPaths, err := utils.StrSliceToUtf16PtrArr(sourcePaths)
	if err != nil {
		return err
	}

	if err = dismAddLanguage(
		session,
		u16LangName,
		preventPending,
		limitAccess,
		&u16SrcPaths[0],
		uint32(len(u16SrcPaths)),
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}

// DismRemoveLanguage removes language from an image.
func DismRemoveLanguage(
	session DismSession,
	languageName string,
	cancelEvent windows.Handle,
	progress uintptr,
	userData unsafe.Pointer,
) error {
	var u16LangName *uint16
	var err error

	if languageName != "" {
		if u16LangName, err = windows.UTF16PtrFromString(languageName); err != nil {
			return err
		}
	}

	if err = dismRemoveLanguage(
		session,
		u16LangName,
		cancelEvent,
		progress,
		userData,
	); err != nil {
		return w32api.WrapInternalErr(utils.HresultToError(err), moddismapi.Handle(), getErrMsg(err))
	}

	return nil
}
