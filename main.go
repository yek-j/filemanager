package main

import (
	"fmt"
	"log"

	"github.com/yek-j/filemanager/config"
	"github.com/yek-j/filemanager/utils"
)

func main() {
	cfg, err := config.LoadConfig("window-test.json")
	if err != nil {
		log.Fatal("Config load failed: ", err)
	}
	fmt.Printf("Source path: %s\n", cfg.SourcePath)
	fmt.Printf("Work path: %s\n", cfg.WorkPath)
	fmt.Println("✅ Config loaded successfully")

	// 2. ScanFiles 테스트
	fmt.Println("\n--- Testing ScanFiles ---")
	scanReport, err := utils.ScanFiles(cfg)
	if err != nil {
		log.Fatal("scanFiles failed: ", err)
	}

	// 결과 출력
	fmt.Printf("Root exists: %v\n", scanReport.RootExists)
	fmt.Printf("Ready to process: %v\n", scanReport.ReadyToProcess)
	fmt.Printf("Total files: %d\n", scanReport.TotalFiles)

	// 3. copyRootDir 테스트
	fmt.Println("\n--- Testing copyRootDir ---")
	if scanReport.ReadyToProcess {
		err = utils.CopyRootDir(cfg)

		if err != nil {
			log.Fatal("CopyRootDir failed: ", err)
		}
		fmt.Println("✅ Copy completed successfully")
	}
}
