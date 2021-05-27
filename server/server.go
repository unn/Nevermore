// Copyright 2021 Nevermore.

// Origination Copyright:
// Copyright 2015 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package main

import (
	"github.com/ArcCS/Nevermore/comms"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/jarvoral"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/stats"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	logFile, err := os.OpenFile("log_"+ time.Now().Format("01-02-2006") +".txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	stats.Start()
	// Lets set some settings
	config.Server.Motd, _ = data.LoadSetting("motd")

	go jarvoral.StartJarvoral()
	objects.Load()
	comms.Listen(config.Server.Host, config.Server.Port)
}
