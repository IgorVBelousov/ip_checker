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

var old_ip string

func main() {
	single_proc := single.New("name")
	single_proc.Lock()
	defer single_proc.Unlock()

	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	icon, err := walk.NewIconFromResourceId(9)
	if err != nil {
		log.Fatal(err)
	}

	ni, err := walk.NewNotifyIcon()
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}

	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	go func() {
		first_loop := true
		for first_loop {
			ip, ip_err := get_ip()
			if ip_err == nil {
				first_loop = false
				old_ip = ip
				if err := ni.ShowCustom(
					"IP Checker",
					" Current IP is "+old_ip); err != nil {

					log.Fatal(err)
				}

				go func() {
					for {
						time.Sleep(time.Second * 60)
						ip, ip_err := get_ip()
						if ip_err == nil {
							if old_ip != ip {
								if err := ni.ShowCustom(
									"IP Checker",
									" Old IP - "+old_ip+" New IP - "+ip); err != nil {

									log.Fatal(err)
								}
								old_ip = ip
							}
						}
					}
				}()

			}
		}

	}()

	mw.Run()
}
