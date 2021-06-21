package main

import (
	"os"
	
	"github.com/mailbadger/app/storage"
)

func main() {
	driver := os.Getenv("DATABASE_DRIVER")
	config := storage.MakeConfigFromEnv(driver)
	
	_ = storage.New(driver, config)
	
	// All of this should be a part of a transaction
	// ---------------------
	/*
		Part I - Fetch last processed id and fetch next chunk of events
	*/
	// ---------------------
	
	// ---------------------
	/*
		Part II - For each summary record fetch latest date record and add the sum if it is the same date update it,
		otherwise insert new one with the date from the summary record (keep in mind that this could be the first one)
	*/
	// ---------------------
	
	// ---------------------
	/*
		Part III - Add last processed id
	*/
	// ---------------------
}
