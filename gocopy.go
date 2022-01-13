package main

import (
	"fmt"
	"gocopy/symbolic"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	cp "github.com/otiai10/copy"
	cli "github.com/urfave/cli/v2"
)

func StartCopy(src string, dest string, wg *sync.WaitGroup, c *cli.Context) {
	isSync := c.Bool("sync")
	isPreserve := c.Bool("preserve")

	copyOpt := cp.Options{
		OnSymlink: func(src string) cp.SymlinkAction {
			// Shallow creates new symlink to the dest of symlink.
			return cp.Shallow
		},
		OnDirExists: func(src, dest string) cp.DirExistsAction {
			if c.Bool("force") {
				// Replace deletes all contents under the dir and copy src files.
				return cp.Replace
			}
			// Merge preserves or overwrites existing files under the dir (default behavior).
			return cp.Merge
		},
		Sync:          isSync,
		PreserveOwner: isPreserve,
		PreserveTimes: isPreserve,
	}

	dir, err := os.Open(src)

	if err != nil {
		log.Fatal(err)
	}

	nameList, err := dir.Readdirnames(0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Copy " + src + " to " + dest)

	for _, name := range nameList {
		wg.Add(1)
		srcPath := src + "/" + name
		destPath := dest + "/" + name

		go func() {
			defer wg.Done()
			err := cp.Copy(srcPath, destPath, copyOpt)

			if err != nil {
				// symbolic check
				isSymbolink := symbolic.IsSameSymbolink(srcPath, destPath, isPreserve)

				if isSymbolink == symbolic.SAMESYM {
					fmt.Println("Same Symbolic Link : " + srcPath + " to " + destPath)
				} else if isSymbolink == symbolic.DIFFSYM {
					fmt.Println("Check Symbolic Link : " + srcPath + " and " + destPath)
				} else if isSymbolink == symbolic.NONESYM {
					log.Fatal(err)
				}
			}
		}()
	}

}

func main() {
	app := &cli.App{
		Name:      "gocopy",
		Usage:     "Copy directory using go",
		UsageText: "gocopy [Source Dir/File] [Destination Dir/File]",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "MJ Cheon",
				Email: "myungjae92@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name: "core", Aliases: []string{"c"}, Value: 1,
			},
			&cli.BoolFlag{
				Name: "preserve", Aliases: []string{"p"}, Value: false,
			},
			&cli.BoolFlag{
				Name: "force", Aliases: []string{"f"}, Value: false,
			},
			&cli.BoolFlag{
				Name: "sync", Aliases: []string{"s"}, Value: false,
			},
			&cli.BoolFlag{
				Name: "verbose", Aliases: []string{"v"}, Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			var wg sync.WaitGroup

			if c.NArg() == 2 {

				useCore := c.Int("core")
				isVerbose := c.Bool("verbose")

				if useCore <= runtime.NumCPU() {
					runtime.GOMAXPROCS(useCore)

					if isVerbose {
						fmt.Println("CPU Core (Used/Total): " + strconv.Itoa(useCore) + " / " + strconv.Itoa(runtime.NumCPU()))
						fmt.Println("")
					}
				}

				copySrc := c.Args().Get(0)
				copyDest := c.Args().Get(1)

				if isVerbose {
					fmt.Println("Start Time \t: " + time.Now().Format("2006-01-02 15:04:05.999") + "\n")
				}

				StartCopy(copySrc, copyDest, &wg, c)
				wg.Wait()

				if isVerbose {
					fmt.Println("\nEnd Time \t: " + time.Now().Format("2006-01-02 15:04:05.999"+"\n"))
				}

			} else {
				cli.ShowAppHelp(c)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
