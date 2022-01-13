package symbolic

import (
	"io/fs"
	"log"
	"os"
	"syscall"
)

// Symbolic link compare value
const (
	SAMESYM = iota
	DIFFSYM
	NONESYM
)

type SymFileStat struct {
	SymbolicPath string
	Realpath     string
	Uid          uint32
	Gid          uint32
	Size         int64
	Blksize      int64
}

func getSymbolicStat(symFileInfo fs.FileInfo, symPath string) SymFileStat {
	var symFileStat SymFileStat

	realPath, _ := os.Readlink(symPath)
	fileStat := symFileInfo.Sys().(*syscall.Stat_t)

	symFileStat.SymbolicPath = symFileInfo.Name()
	symFileStat.Realpath = realPath
	symFileStat.Uid = fileStat.Uid
	symFileStat.Gid = fileStat.Gid
	symFileStat.Size = fileStat.Size
	symFileStat.Blksize = fileStat.Blksize

	return symFileStat
}

func IsSameSymbolink(src string, dest string, isPreserve bool) int {
	srcSymFileInfo, err := os.Lstat(src)

	if err != nil {
		log.Fatal(err)
	}

	destSymFileInfo, err := os.Lstat(src)

	if err != nil {
		log.Fatal(err)
	}

	if srcSymFileInfo.Mode()&os.ModeSymlink != 0 {
		if destSymFileInfo.Mode()&os.ModeSymlink != 0 {
			srcSysStat := getSymbolicStat(srcSymFileInfo, src)
			destSysStat := getSymbolicStat(destSymFileInfo, dest)

			if srcSysStat == destSysStat {
				return SAMESYM
			} else {
				return DIFFSYM
			}
		}
	}

	return NONESYM
}
