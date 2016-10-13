package main

import (
	"bufio"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func isMounted(targetDir string) bool {
	findmntCmd := exec.Command("findmnt", "-n", targetDir)
	findmntStdout, err := findmntCmd.StdoutPipe()
	if err != nil {
		fail("Could not get findmount stdout pipe: ", err.Error())
	}

	if err = findmntCmd.Start(); err != nil {
		fail("findmnt failed to start: ", err.Error())
	}

	findmntScanner := bufio.NewScanner(findmntStdout)
	findmntScanner.Split(bufio.ScanWords)
	findmntScanner.Scan()
	if findmntScanner.Err() != nil {
		fail("Couldn't read findnmnt output: ", findmntScanner.Err().Error())
	}

	findmntText := findmntScanner.Text()
	if err = findmntCmd.Wait(); err != nil {
		_, isExitError := err.(*exec.ExitError)
		if !isExitError {
			fail("findmnt failed: ", err.Error())
		}
	}

	return findmntText == targetDir
}

func mount(targetDir string, device string, opt JsonOptions) {
	if opt.FsType == "" {
		fail("No filesystem type specified")
	}

	var res unix.Stat_t
	if err := unix.Stat(device, &res); err != nil {
		fail("Could not stat ", device, ": ", err.Error())
	}

	if res.Mode&unix.S_IFMT != unix.S_IFBLK {
		fail("Not a block device: ", device)
	}

	if isMounted(targetDir) {
		succeed()
	}

	mkfsCmd := exec.Command("mkfs", "-t", opt.FsType, device)
	if mkfsOut, err := mkfsCmd.CombinedOutput(); err != nil {
		fail("Could not mkfs: ", err.Error(), " Output: ", string(mkfsOut))
	}

	if err := os.MkdirAll(targetDir, 0777); err != nil {
		fail("Could not create directory: ", err.Error())
	}

	mountCmd := exec.Command("mount", device, targetDir)
	if mountOut, err := mountCmd.CombinedOutput(); err != nil {
		fail("Could not mount: ", err.Error(), " Output: ", string(mountOut))
	}

	succeed()
}

func unmount(targetDir string) {
	if !isMounted(targetDir) {
		succeed()
	}

	umountCmd := exec.Command("umount", targetDir)
	if umountOut, err := umountCmd.CombinedOutput(); err != nil {
		fail("Could not umount: ", err.Error(), " Output: ", string(umountOut))
	}

	succeed()
}
