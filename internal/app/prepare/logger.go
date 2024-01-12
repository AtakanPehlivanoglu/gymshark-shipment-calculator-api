package prepare

import (
	"fmt"
	"github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/app/config"
	"log"
	"os"
)

func AppLogger(config *config.Config) *log.Logger {
	logger := log.New(os.Stdout, fmt.Sprintf("["+config.App.Name+"] "), log.Ldate|log.Ltime)

	return logger
}
