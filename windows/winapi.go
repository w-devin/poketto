//go:build windows

package windows

import (
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	modntdll    = windows.NewLazySystemDLL("ntdll.dll")

	procIsProcessInJob    = modkernel32.NewProc("IsProcessInJob")
	procIsProcessInJobErr = procIsProcessInJob.Find()

	procNtQueryObject               = modntdll.NewProc("NtQueryObject")
	procNtQueryObjectErr            = procNtQueryObject.Find()
	procNtQuerySystemInformation    = modntdll.NewProc("NtQuerySystemInformation")
	procNtQuerySystemInformationErr = procNtQuerySystemInformation.Find()
)

const (
	// ObjectInformationClass values used to call NtQueryObject (https://docs.microsoft.com/en-us/windows/win32/api/winternl/nf-winternl-ntqueryobject)
	ObjectNameInformationClass = 0x1
	ObjectTypeInformationClass = 0x2

	// Includes all processes in the system in the snapshot. (https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot)
	Th32csSnapProcess uint32 = 0x00000002
)

type API interface {
	// IsProcessInJob determines whether the process is running in the specified job.
	IsProcessInJob(procHandle windows.Handle, jobHandle windows.Handle, result *bool) error

	// GetObjectType gets the object type of the given handle
	GetObjectType(handle windows.Handle) (string, error)

	// GetObjectName gets the object name of the given handle
	GetObjectName(handle windows.Handle) (string, error)

	// QuerySystemExtendedHandleInformation retrieves Extended handle system information.
	QuerySystemExtendedHandleInformation() ([]SystemHandleInformationExItem, error)

	// CurrentProcess returns the handle for the current process.
	// It is a pseudo handle that does not need to be closed.
	CurrentProcess() windows.Handle

	// CloseHandle closes an open object handle.
	CloseHandle(h windows.Handle) error

	// OpenProcess returns an open handle
	OpenProcess(desiredAccess uint32, inheritHandle bool, pID uint32) (windows.Handle, error)

	// DuplicateHandle duplicates an object handle.
	DuplicateHandle(hSourceProcessHandle windows.Handle, hSourceHandle windows.Handle, hTargetProcessHandle windows.Handle, lpTargetHandle *windows.Handle, dwDesiredAccess uint32, bInheritHandle bool, dwOptions uint32) error

	// CreateToolhelp32Snapshot takes a snapshot of the specified processes, as well as the heaps, modules, and threads used by these processes.
	CreateToolhelp32Snapshot(flags uint32, pID uint32) (windows.Handle, error)

	// Process32First retrieves information about the first process encountered in a system snapshot.
	Process32First(snapshot windows.Handle, procEntry *windows.ProcessEntry32) error

	// Process32Next retrieves information about the next process recorded in a system snapshot.
	Process32Next(snapshot windows.Handle, procEntry *windows.ProcessEntry32) error
}

type WindowsApi struct {
}

func (a *WindowsApi) IsProcessInJob(procHandle windows.Handle, jobHandle windows.Handle, result *bool) error {
	if procIsProcessInJobErr != nil {
		return procIsProcessInJobErr
	}
	r1, _, e1 := syscall.SyscallN(procIsProcessInJob.Addr(), uintptr(procHandle), uintptr(jobHandle), uintptr(unsafe.Pointer(result)))
	if r1 == 0 {
		if e1 != 0 {
			return e1
		}
		return syscall.EINVAL
	}
	return nil
}

// GetObjectType gets the object type of the given handle
func (a *WindowsApi) GetObjectType(handle windows.Handle) (string, error) {
	buffer := make([]byte, 1024*10)
	length := uint32(0)

	status := a.NtQueryObject(handle, ObjectTypeInformationClass,
		&buffer[0], uint32(len(buffer)), &length)
	if status != windows.STATUS_SUCCESS {
		return "", status
	}

	return (*ObjectTypeInformation)(unsafe.Pointer(&buffer[0])).TypeName.String(), nil
}

// GetObjectName gets the object name of the given handle
func (a *WindowsApi) GetObjectName(handle windows.Handle) (string, error) {
	buffer := make([]byte, 1024*2)
	var length uint32

	MapingHandle, err := windows.CreateFileMapping(handle, nil, windows.PAGE_READONLY, 0, 1024, windows.StringToUTF16Ptr("testFileMapping"))
	defer windows.CloseHandle(MapingHandle)
	if err != nil {
		return "", err
	}

	status := a.NtQueryObject(handle, ObjectNameInformationClass,
		&buffer[0], uint32(len(buffer)), &length)
	if status != windows.STATUS_SUCCESS {
		return "", status
	}

	return (*UnicodeString)(unsafe.Pointer(&buffer[0])).String(), nil
}

func (a *WindowsApi) QuerySystemExtendedHandleInformation() ([]SystemHandleInformationExItem, error) {
	buffer := make([]byte, 1024)
	var retLen uint32
	var status windows.NTStatus

	for {
		status = a.NtQuerySystemInformation(
			windows.SystemExtendedHandleInformation,
			unsafe.Pointer(&buffer[0]),
			uint32(len(buffer)),
			&retLen,
		)

		if status == windows.STATUS_BUFFER_OVERFLOW ||
			status == windows.STATUS_BUFFER_TOO_SMALL ||
			status == windows.STATUS_INFO_LENGTH_MISMATCH {
			if int(retLen) <= cap(buffer) {
				buffer = unsafe.Slice(&buffer[0], int(retLen))
			} else {
				buffer = make([]byte, int(retLen))
			}
			continue
		}
		// if no error
		break
	}

	if status>>30 != 3 {
		buffer = (buffer)[:int(retLen)]

		handlesList := (*SystemExtendedHandleInformation)(unsafe.Pointer(&buffer[0]))
		handles := unsafe.Slice(&handlesList.Handles[0], int(handlesList.NumberOfHandles))

		return handles, nil
	}

	return nil, status
}

func (a *WindowsApi) QuerySystemHandleInformation() ([]SystemHandleInformationItem, error) {
	buffer := make([]byte, 1024)
	var retLen uint32
	var status windows.NTStatus

	for {
		status = a.NtQuerySystemInformation(
			windows.SystemHandleInformation,
			unsafe.Pointer(&buffer[0]),
			uint32(len(buffer)),
			&retLen,
		)

		if status == windows.STATUS_BUFFER_OVERFLOW ||
			status == windows.STATUS_BUFFER_TOO_SMALL ||
			status == windows.STATUS_INFO_LENGTH_MISMATCH {
			if int(retLen) <= cap(buffer) {
				buffer = unsafe.Slice(&buffer[0], int(retLen))
			} else {
				buffer = make([]byte, int(retLen))
			}
			continue
		}
		// if no error
		break
	}

	if status>>30 != 3 {
		buffer = (buffer)[:int(retLen)]

		handlesList := (*SystemHandleInformation)(unsafe.Pointer(&buffer[0]))
		handles := unsafe.Slice(&handlesList.Handles[0], int(handlesList.NumberOfHandles))

		return handles, nil
	}

	return nil, status
}

func (a *WindowsApi) OpenProcess(desiredAccess uint32, inheritHandle bool, pID uint32) (windows.Handle, error) {
	return windows.OpenProcess(desiredAccess, inheritHandle, pID)
}

func (a *WindowsApi) CloseHandle(h windows.Handle) error {
	return windows.CloseHandle(h)
}

// CurrentProcess returns the handle for the current process.
// It is a pseudo handle that does not need to be closed.
func (a *WindowsApi) CurrentProcess() windows.Handle {
	return windows.CurrentProcess()
}

func (a *WindowsApi) DuplicateHandle(hSourceProcessHandle windows.Handle, hSourceHandle windows.Handle, hTargetProcessHandle windows.Handle, lpTargetHandle *windows.Handle, dwDesiredAccess uint32, bInheritHandle bool, dwOptions uint32) error {
	return windows.DuplicateHandle(hSourceProcessHandle, hSourceHandle, hTargetProcessHandle, lpTargetHandle, dwDesiredAccess, bInheritHandle, dwOptions)
}

func (a *WindowsApi) CreateToolhelp32Snapshot(flags uint32, pID uint32) (windows.Handle, error) {
	return windows.CreateToolhelp32Snapshot(flags, pID)
}

func (a *WindowsApi) Process32First(snapshot windows.Handle, procEntry *windows.ProcessEntry32) error {
	return windows.Process32First(snapshot, procEntry)
}

func (a *WindowsApi) Process32Next(snapshot windows.Handle, procEntry *windows.ProcessEntry32) error {
	return windows.Process32Next(snapshot, procEntry)
}

// System handle extended information item, returned by NtQuerySystemInformation (https://docs.microsoft.com/en-us/windows/win32/api/winternl/nf-winternl-ntquerysysteminformation)
type SystemHandleInformationExItem struct {
	Object                uintptr
	UniqueProcessID       uintptr
	HandleValue           uintptr
	GrantedAccess         uint32
	CreatorBackTraceIndex uint16
	ObjectTypeIndex       uint16
	HandleAttributes      uint32
	Reserved              uint32
}

// System extended handle information summary, returned by NtQuerySystemInformation (https://docs.microsoft.com/en-us/windows/win32/api/winternl/nf-winternl-ntquerysysteminformation)
type SystemExtendedHandleInformation struct {
	NumberOfHandles uintptr
	Reserved        uintptr
	Handles         [1]SystemHandleInformationExItem
}

type SystemHandleInformationItem struct {
	ProcessId        uint32
	ObjectTypeNumber byte
	Flags            byte
	HandleValue      uint16
	ObjectPointer    uintptr
	GrantedAccess    uint32
}

type SystemHandleInformation struct {
	NumberOfHandles uintptr
	Handles         [1]SystemHandleInformationItem
}

// Object type returned by calling NtQueryObject function
type ObjectTypeInformation struct {
	TypeName               UnicodeString
	TotalNumberOfObjects   uint32
	TotalNumberOfHandles   uint32
	TotalPagedPoolUsage    uint32
	TotalNonPagedPoolUsage uint32
}

// Unicode string returned by NtQueryObject calls (https://docs.microsoft.com/en-us/windows/win32/api/subauth/ns-subauth-unicode_string)
type UnicodeString struct {
	Length        uint16
	AllocatedSize uint16
	WString       *byte
}

func (u UnicodeString) String() string {
	defer func() {
		// TODO: may we recover?
		_ = recover()
	}()

	data := unsafe.Slice((*uint16)(unsafe.Pointer(u.WString)), int(u.Length*2))

	return windows.UTF16ToString(data)
}

func (a *WindowsApi) NtQueryObject(handle windows.Handle, objectInformationClass uint32, objectInformation *byte, objectInformationLength uint32, returnLength *uint32) (ntStatus windows.NTStatus) {
	if procNtQueryObjectErr != nil {
		return windows.STATUS_PROCEDURE_NOT_FOUND
	}
	r0, _, _ := syscall.SyscallN(procNtQueryObject.Addr(), uintptr(handle), uintptr(objectInformationClass), uintptr(unsafe.Pointer(objectInformation)), uintptr(objectInformationLength), uintptr(unsafe.Pointer(returnLength)), 0)
	if r0 != 0 {
		ntStatus = windows.NTStatus(r0)
	}
	return
}

func (a *WindowsApi) NtQuerySystemInformation(sysInfoClass int32, sysInfo unsafe.Pointer, sysInfoLen uint32, retLen *uint32) (ntstatus windows.NTStatus) {
	if procNtQuerySystemInformationErr != nil {
		return windows.STATUS_PROCEDURE_NOT_FOUND
	}
	r0, _, _ := syscall.SyscallN(procNtQuerySystemInformation.Addr(), uintptr(sysInfoClass), uintptr(sysInfo), uintptr(sysInfoLen), uintptr(unsafe.Pointer(retLen)), 0, 0)
	if r0 != 0 {
		ntstatus = windows.NTStatus(r0)
	}

	return
}

func (a *WindowsApi) GetKernelPath(userPath string) string {
	if strings.HasPrefix(userPath, "\\\\") {
		return userPath
	}

	return "\\\\?\\" + userPath
}
