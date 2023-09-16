package main

import (
	"flag"
	"fmt"

	"github.com/thank243/zteOnu/app"
	"github.com/thank243/zteOnu/version"
)

func main() {
	version.Show()

	user := flag.String("u", "telecomadmin", "factory mode auth username")
	passwd := flag.String("p", "nE7jA%5m", "factory mode auth password")
	ip := flag.String("i", "192.168.1.1", "ONU ip address")
	port := flag.Int("port", 8080, "ONU http port")
	flag.Parse()

	fac := app.New(*user, *passwd, *ip, *port)
	fmt.Print("step 0 (reset factory): ")
	if err := fac.Reset(); err != nil {
		fmt.Printf("reset errors: %v\n", err)
		return
	}
	fmt.Println("ok")
	fmt.Print("step 1 (request factory mode): ")
	fac.ReqFactoryMode()
	fmt.Println("ok")

	fmt.Print("step 2 (send sq): ")
	ver, err := fac.SendSq()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ok")

	fmt.Print("step 3 (check login auth): ")
	switch ver {
	case 1:
		if err := fac.CheckLoginAuth(); err != nil {
			fmt.Println(err)
			return
		}
	case 2:
		if err := fac.SendInfo(); err != nil {
			fmt.Println(err)
			return
		}
		if err := fac.CheckLoginAuth(); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("ok")

	fmt.Print("step 4 (enter factory mode): ")
	tlUser, tlPass, err := fac.FactoryMode()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok")

	fmt.Printf("user: %s, passwd: %s", tlUser, tlPass)
}