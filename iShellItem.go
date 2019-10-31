// +build windows

package cfd

import (
	"github.com/go-ole/go-ole"
	"syscall"
	"unsafe"
)

const (
	iidShellItemGUID = "{43826d1e-e718-42ee-bc55-a1e261c37bfe}"
)

var (
	shell32                         *syscall.LazyDLL
	procSHCreateItemFromParsingName *syscall.LazyProc
	iidShellItem                    *ole.GUID
)

func init() {
	shell32 = syscall.NewLazyDLL("Shell32.dll")
	procSHCreateItemFromParsingName = shell32.NewProc("SHCreateItemFromParsingName")
	iidShellItem, _ = ole.IIDFromString(iidShellItemGUID) // TODO handle error
}

type iShellItem struct {
	vtbl *iShellItemVtbl
}

type iShellItemVtbl struct {
	iUnknownVtbl
	BindToHandler  uintptr
	GetParent      uintptr
	GetDisplayName uintptr // func (sigdnName SIGDN, ppszName *LPWSTR) HRESULT
	GetAttributes  uintptr
	Compare        uintptr
}

func newIShellItem(path string) (*iShellItem, error) {
	var shellItem *iShellItem
	pathPtr := ole.SysAllocString(path)
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(pathPtr)))
	ret, _, _ := procSHCreateItemFromParsingName.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		0,
		uintptr(unsafe.Pointer(iidShellItem)),
		uintptr(unsafe.Pointer(&shellItem)))
	return shellItem, hresultToError(ret)
}

func (vtbl *iShellItemVtbl) getDisplayName(objPtr unsafe.Pointer) (string, error) {
	var ptr *uint16
	ret, _, _ := syscall.Syscall(vtbl.GetDisplayName,
		2,
		uintptr(objPtr),
		0x80058000, // SIGDN_FILESYSPATH
		uintptr(unsafe.Pointer(&ptr)))
	if err := hresultToError(ret); err != nil {
		return "", err
	}
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(ptr)))
	return ole.LpOleStrToString(ptr), nil
}
