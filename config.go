package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func GetConfig(configFile string) {

	//	var split, header, format string
	var split, header string

	//	var roscop = codeRoscopFromCsv("code_roscop.csv")

	err := gcfg.ReadFileInto(&cfg, configFile)
	if err == nil {
		split = cfg.Ctd.Split
		header = cfg.Ctd.Header
		//		format = cfg.Ctd.Format
		//		cruisePrefix = cfg.Ctd.CruisePrefix
		//		stationPrefixLength = cfg.Ctd.StationPrefixLength
		// TODOS: complete
		nc.Attributes["cycle_mesure"] = cfg.Cruise.CycleMesure
		nc.Attributes["plateforme"] = cfg.Cruise.Plateforme
		nc.Attributes["institute"] = cfg.Cruise.Institute
		nc.Attributes["pi"] = cfg.Cruise.Pi
		nc.Attributes["timezone"] = cfg.Cruise.Timezone
		nc.Attributes["type"] = cfg.Ctd.Type
		nc.Attributes["sn"] = cfg.Ctd.Sn
	} else {
		log.Fatal(err)
	}
	fields := strings.Split(split, ",")
	for i := 0; i < len(fields); i += 2 {
		if v, err := strconv.Atoi(fields[i+1]); err == nil {
			map_var[fields[i]] = v - 1
		}
	}
	// split header
	hdr = strings.Fields(header)

	// fill map_format from code_roscop
	for _, key := range hdr {
		map_format[key] = nc.Roscop[key].format
	}
	if *optDebug {
		fmt.Println(header)
	}
}
