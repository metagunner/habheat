package main

import (
	"context"
	"errors"

	"github.com/jesseduffield/gocui"
	"github.com/metagunner/habheath/pkg/database"
	"github.com/metagunner/habheath/pkg/gui"
)

func main() {
	db := database.NewDB("./test.db")
	if err := db.Open(); err != nil {
		panic(err)
	}
	database.SeedTestData(context.Background(), db, 2024, 7)

	gui := gui.NewGui(db)
	err := gui.Run()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			panic(err)
		}
	}
}

//func lookup(tab []string, val string) (int, string, error) {
//	for i, v := range tab {
//		if len(val) >= len(v) && match(val[0:len(v)], v) {
//			return i, val[len(v):], nil
//		}
//	}
//	return -1, val, errBad
//}

//256-colors escape codes
//		for i := 0; i < 256; i++ {
//			str := fmt.Sprintf("\x1b[48;5;%dm\x1b[30m%3d\x1b[0m ", i, i)
//			str += fmt.Sprintf("\x1b[38;5;%dm%3d\x1b[0m ", i, i)
//
//			if (i+1)%10 == 0 {
//				str += "\n"
//			}
//
//			fmt.Fprint(v, str)
//		}
//
//		fmt.Fprint(v, "\n\n")
//
//		// 8-colors escape codes
//		ctr := 0
//		for i := 0; i <= 7; i++ {
//			for _, j := range []int{1, 4, 7} {
//				str := fmt.Sprintf("\x1b[3%d;%dm%d:%d\x1b[0m ", i, j, i, j)
//				if (ctr+1)%20 == 0 {
//					str += "\n"
//				}
//
//				fmt.Fprint(v, str)
//
//				ctr++
//			}
//		}
//
