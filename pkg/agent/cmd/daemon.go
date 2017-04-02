//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"os"
	"os/signal"
	"syscall"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime"
)

func Agent(cmd *cli.Cmd) {

	var ctx = context.Get()
	var cfg = config.Get()

	cmd.Spec = "[-d]"

	var debug = cmd.Bool(cli.BoolOpt{Name: "d debug", Value: false, Desc: "Enable debug mode"})

	cmd.Before = func() {

		if *debug {
			cfg.Debug = *debug
		}

		if cfg.HttpServer.Port == 0 {
			cfg.HttpServer.Port = 2967
		}
	}

	cmd.Action = func() {

		var (
			sigs = make(chan os.Signal)
			done = make(chan bool, 1)
		)

		ctx.New(cfg)
		runtime.New(cfg.Runtime)

		// Handle SIGINT and SIGTERM.
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for {
				select {
				case <-sigs:
					done <- true
					return
				}
			}
		}()

		<-done

		logrus.Info("Handle SIGINT and SIGTERM.")
	}

}
