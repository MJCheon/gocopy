package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	cp "github.com/otiai10/copy"
	cli "github.com/urfave/cli/v2"
)

func mainCopy(src string, dest string, c *cli.Context) {
	startCopy(src, dest, c)
}

func startCopy(src string, dest string, c *cli.Context) {
	copyOpt := cp.Options{
		OnDirExists: func(src, dest string) cp.DirExistsAction {
			return 1
		},
		Sync:          false,
		PreserveOwner: c.Bool("preserve"),
		PreserveTimes: c.Bool("preserve"),
	}

	dir, err := os.Open(src)

	if err != nil {
		log.Fatal(err)
	}

	nameList, err := dir.Readdirnames(0)

	if err != nil {
		log.Fatal(err)
	}

	for _, name := range nameList {
		srcPath := src + "/" + name
		destPath := dest + "/" + name

		fmt.Println("Copy " + srcPath + " to " + destPath)
		err := cp.Copy(srcPath, destPath, copyOpt)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	app := &cli.App{
		Name:      "gocopy",
		Usage:     "Copy directory using go",
		UsageText: "gocopy [Source Dir] [Destination Dir]",
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
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 2 {

				useCore := c.Int("core")

				if useCore <= runtime.NumCPU() {
					runtime.GOMAXPROCS(useCore)
				}

				copySrc := c.Args().Get(0)
				copyDest := c.Args().Get(1)

				fmt.Println("Start Time : " + time.Now().Format(time.RFC3339))
				mainCopy(copySrc, copyDest, c)
				fmt.Println("End Time : " + time.Now().Format(time.RFC3339))
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
