package utils

import (
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procFormatMessageW = modkernel32.NewProc("FormatMessageW")
)

func formatMessage(flags uint32, source uintptr, messageId, languageId uint32, buffer *uint16, size uint32, args *byte) (uint32, error) {
	r, _, err := procFormatMessageW.Call(
		uintptr(flags),
		source,
		uintptr(messageId),
		uintptr(languageId),
		uintptr(unsafe.Pointer(buffer)),
		uintptr(size),
		uintptr(unsafe.Pointer(args)),
	)
	if r == 0 {
		return 0, err
	}

	return uint32(r), nil
}

// GetErrorMessage returns message using the errno and module(dll) handle.
func GetErrorMessage(errno uint32, module uintptr) string {
	var flags uint32 = windows.FORMAT_MESSAGE_FROM_HMODULE |
		windows.FORMAT_MESSAGE_FROM_SYSTEM |
		windows.FORMAT_MESSAGE_ARGUMENT_ARRAY |
		windows.FORMAT_MESSAGE_IGNORE_INSERTS
	buf := make([]uint16, 300)
	b := &buf[0]
	n, err := formatMessage(flags, module, errno, 0, b, uint32(len(buf)), nil)
	if err == nil {
		for ; n > 0 && buf[n-1] == '\n' || buf[n-1] == '\r'; n-- {
		}
		return windows.UTF16ToString(buf[:n])
	} else if err != windows.ERROR_INSUFFICIENT_BUFFER {
		return "winapi error: 0x" + strconv.FormatUint(uint64(errno), 16)
	}

	// retry with heap allocation
	n, err = formatMessage(flags|windows.FORMAT_MESSAGE_ALLOCATE_BUFFER, module, errno, 0, (*uint16)(unsafe.Pointer(&b)), 0, nil)
	if err != nil {
		return "winapi error: 0x" + strconv.FormatUint(uint64(errno), 16)
	}

	defer windows.LocalFree(windows.Handle(uintptr(unsafe.Pointer(b))))

	buf = unsafe.Slice(b, n)
	for ; n > 0 && buf[n-1] == '\n' || buf[n-1] == '\r'; n-- {
	}
	return windows.UTF16ToString(buf[:n])
}

// HresultToError parses HRESULT to error value for comparision
//
// see https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-erref/0642cb2f-2075-4469-918c-4441e69c548a
func HresultToError(err error) error {
	errno, ok := err.(windows.Errno)
	if !ok {
		return err
	}

	hresVal := int32(errno)
	if windows.Errno(hresVal) != errno {
		return err
	}

	if hresVal&0x10000000 != 0 { // this is an NTStatus value
		return windows.NTStatus(hresVal)
	}

	if hresVal&0x20000000 == 0 && hresVal&0x07ff0000 == windows.FACILITY_WIN32 {
		// this is an undecorated error code
		return windows.Errno(hresVal & 0xffff)
	}

	return err
}

func StrSliceToUtf16PtrArr(strSlice []string) ([]*uint16, error) {
	var u16Arr []*uint16

	for _, str := range strSlice {
		// add non-empty strings
		if str != "" {
			u16Ptr, err := windows.UTF16PtrFromString(str)
			if err != nil {
				return nil, err
			}

			u16Arr = append(u16Arr, u16Ptr)
		}
	}

	return u16Arr, nil
}

func PZZWSTRToStrings(pzzwstr **uint16) []string {
	result := make([]string, 0)
	if pzzwstr == nil || *pzzwstr == nil {
		return result
	}
	bufPtr := *pzzwstr

	for *bufPtr != 0 {
		length := 0
		for ptr := bufPtr; *ptr != 0; ptr = (*uint16)(unsafe.Add(unsafe.Pointer(ptr), 2)) {
			length++
		}

		result = append(result, windows.UTF16PtrToString(bufPtr))

		bufPtr = (*uint16)(unsafe.Add(unsafe.Pointer(bufPtr), (length+1)*2))
	}

	return result
}
