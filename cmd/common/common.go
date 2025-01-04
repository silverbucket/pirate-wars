package common

type Coordinates struct {
	X int
	Y int
}

const (
	LogFile          = "pirate-wars.log"
	WorldWidth       = 600
	WorldHeight      = 600
	TotalTowns       = 30
	ViewWidth        = 75
	ViewHeight       = 50
	MiniMapFactor    = 11
	TypeDeepWater    = 0
	TypeOpenWater    = 1
	TypeShallowWater = 2
	TypeBeach        = 3
	TypeLowland      = 4
	TypeHighland     = 5
	TypeRock         = 6
	TypePeak         = 7
	TypeTown         = 8
)

type ViewPort struct {
	width   int
	height  int
	topLeft int
}

//func CreateNewLogger(filename string, prefix string) *log.Logger {
//	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
//	if err != nil {
//		log.Fatalf("error opening file: %v", err)
//	}
//	defer f.Close()
//
//	logger := log.New(f, filename, log.LstdFlags)
//	logger.SetPrefix(fmt.Sprintf("%v: ", prefix))
//	logger.Println("TEST 1234")
//	return logger
//}

//func CreateLogger(filename string, namespace string) func(string) {
//	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	logger := log.New(f, namespace, log.LstdFlags)
//	return func(msg string) {
//		d := godump.Dumper{}
//		logger.Println(d.Sprint(msg))
//	}
//}
//
//func InitLogger(name string) {
//	CreateLogger(LogFile, name)
//}

//func Logger(name string) log {
//	var logger *os.File
//	if _, ok := os.LookupEnv("DEBUG"); ok {
//		var err error
//		logger, err = os.OpenFile(LogFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
//		if err != nil {
//			os.Exit(1)
//		}
//	}
//	spew.Fdump(m.dump, msg)
//	return logger
//}
