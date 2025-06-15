package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yek-j/filemanager/config"
	"github.com/yek-j/filemanager/plugins"
	"github.com/yek-j/filemanager/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./filemanager <config-file>")
		fmt.Println("Example: ./filemanager my-config.json")
		os.Exit(1)
	}

	configPath := os.Args[1]

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Config load failed: ", err)
	}
	fmt.Printf("Source path: %s\n", cfg.SourcePath)
	fmt.Printf("Work path: %s\n", cfg.WorkPath)
	fmt.Println("✅ Config loaded successfully")

	// ScanFiles
	fmt.Println("\n--- ScanFiles ---")
	scanReport, err := utils.ScanFiles(cfg)
	if err != nil {
		log.Fatal("scanFiles failed: ", err)
	}

	// 결과 출력
	fmt.Printf("Root exists: %v\n", scanReport.RootExists)
	fmt.Printf("Ready to process: %v\n", scanReport.ReadyToProcess)
	fmt.Printf("Total files: %d\n", scanReport.TotalFiles)

	// copyRootDir
	fmt.Println("\n--- copyRootDir ---")
	if scanReport.ReadyToProcess {
		err = utils.CopyRootDir(cfg)

		if err != nil {
			log.Fatal("CopyRootDir failed: ", err)
		}
		fmt.Println("✅ Copy completed successfully")

		// 플러그인 실행
		plugin, err := plugins.GetPlugin(cfg.Plugin)
		if err != nil {
			log.Fatal("Plugin not found: ", err)
		}

		fmt.Printf("Plugin: %s\n", plugin.GetName())
		err = plugin.Process(cfg)

		if err != nil {
			log.Fatal("Plugin process failed: ", err)
		}
		fmt.Println("✅ Plugin processing completed")
	}
}
