package main

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/marcsauter/single"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func get_ip() (string, error) {
	resp, err := http.Get("https://4.ifcfg.me/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := fmt.Sprintf("%s", body)
	return ip, err
}

func err_fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var old_ip string
var old_t string

func main() {
	single_proc := single.New("IP Checker")
	single_proc.Lock()
	defer single_proc.Unlock()

	mw, err := walk.NewMainWindow()
	err_fatal(err)

	icon, err := walk.NewIconFromResourceId(9)
	err_fatal(err)

	ni, err := walk.NewNotifyIcon()
	err_fatal(err)
	defer ni.Dispose()

	err_fatal(ni.SetIcon(icon))

	exitAction := walk.NewAction()
	err_fatal(exitAction.SetText("E&xit"))

	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	err_fatal(ni.ContextMenu().Actions().Add(exitAction))

	err_fatal(ni.SetVisible(true))

	go func() {
		first_loop := true
		for first_loop {
			ip, ip_err := get_ip()
			if ip_err == nil {
				first_loop = false
				old_ip = ip
				old_t = fmt.Sprint(time.Now().Format("_2 15:04:05"))
				err_fatal(ni.ShowCustom(
					"IP Checker",
					old_t+" Current IP is "+old_ip))
				err_fatal(ni.SetToolTip(old_t + " Current IP is " + old_ip))

				go func() {
					for {
						time.Sleep(time.Second * 60)
						ip, ip_err := get_ip()
						if ip_err == nil {
							if old_ip != ip {
								t := fmt.Sprint(time.Now().Format("_2 15:04:05"))
								err_fatal(ni.ShowCustom(
									"IP Checker",
									t+" - "+ip+"\n"+old_t+" - "+old_ip))
								err_fatal(ni.SetToolTip(t + " - " + ip + "\n" + old_t + " - " + old_ip))

								old_ip = ip
								old_t = t
							}
						}
					}
				}()

			}
		}

	}()

	mw.Run()
}
