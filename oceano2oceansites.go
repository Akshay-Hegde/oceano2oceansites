package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const PROG_NAME string = "oceano2oceansites"
const PROG_VERSION string = "0.2.2"
const PROG_DATE string = "05/20/2016"

// use for echo mode
// Discard is an io.Writer on which all Write calls succeed
// without doing anything. Discard = devNull(0) = int(0)
// when echo is define in program argument list, echo = os.Stdout
var echo io.Writer = ioutil.Discard

// use for debug mode
var debug io.Writer = ioutil.Discard

// usefull shortcut macros
var p = fmt.Println
var f = fmt.Printf

// default configuration file
var cfgname string = "oceano2oceansites.ini"

// default physical parameters file definition is embeded in code_roscop.go
var code_roscop string = "roscop/code_roscop.csv"

// file prefix for --all option: "-all" for all parameters, "" empty by default
var prefixAll = ""

// default output directory
var outputDir = "out"

// Create an empty map.
var map_var = map[string]int{}
var map_format = map[string]string{}
var data = make(map[string]interface{})

var hdr []string
var cfg Config

// matrix used to store profils
type Data_2D struct {
	data [][]float64
}

// map used a matrix for each parameters
type AllData_2D map[string]Data_2D

// the representation in memory of a data set is similar to
// that of a netcdf file
type Nc struct {
	// store dimensions
	Dimensions map[string]int

	// store one dimension variables (eg: TIME, LATITUDE, ...)
	Variables_1D map[string]interface{}
	// store two dimensions variables (eg: PRES, DEPTH, TEMP, ...)
	Variables_2D AllData_2D
	// store global attributes
	Attributes map[string]string

	// used to store max of profiles value
	Extras_f map[string]float64
	// used to store max of profiles type
	Extras_s map[string]string
	// give access to physical parameters
	Roscop roscop

	// store header
	//hdr []string
}

// interface common for all data sets (profile, trajectory and time-series
// and instruments
type Process interface {
	Read([]string)
	GetConfig(string)
	//	WriteHeader(map[string]string, []string)
	WriteAscii(map[string]string, []string)
	WriteNetcdf(InstrumentType)
}

// nc implement interface Process
var nc Process

// define new receiver type based on netcdf equivalent structure
type Ctd struct{ Nc }
type Btl struct{ Nc }

func NewCtd() *Ctd { return &Ctd{} }
func NewBtl() *Btl { return &Btl{} }

// main body
func main() {

	// slice of filename to read and extract data
	var files []string

	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if os.Getenv("OCEANO2OCEANSITES_INI") != "" {
		cfgname = os.Getenv("OCEANO2OCEANSITES_INI")
	}
	if os.Getenv("ROSCOP_CSV") != "" {
		code_roscop = os.Getenv("ROSCOP_CSV")
	}

	// get options parse args list and return all given files to read
	// and configuration file name
	files, optCfgfile := GetOptions()

	// test if output directories exists and create them if not
	mkOutputDir()

	// read the first file and try to find the instrument type, return a bit mask
	typeInstrument = AnalyseFirstFile(files)

	// following the instrument type, allocate the rigth receiver based on
	// Process interface
	switch typeInstrument {
	case CTD:
		//nc = &Ctd{}
		nc = NewCtd()
	case BTL:
		//nc = &Btl{}
		nc = NewBtl()
	default:
		f("main: invalide option typeInstrument -> %d\n", typeInstrument)
		p("Exiting...")
		os.Exit(0)
	}

	// read configuration file, by default, optCfgfile = cfgname
	nc.GetConfig(optCfgfile)
	// debug
	fmt.Fprintln(debug, map_format)

	// read and process all data files
	nc.Read(files)

	// write ASCII file
	nc.WriteAscii(map_format, hdr)

	// write netcdf file
	nc.WriteNetcdf(typeInstrument)
}
