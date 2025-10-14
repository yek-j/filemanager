package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
		workStartTime := time.Now()
		fmt.Printf("Starting file processing at %s\n", workStartTime.Format("15:04:05"))

		copyStartTime := time.Now()
		err = utils.CopyRootDir(cfg)

		if err != nil {
			log.Fatal("CopyRootDir failed: ", err)
		}
		fmt.Println("✅ Copy completed successfully")

		copyDuration := time.Since(copyStartTime)
		fmt.Printf("✅ Backup completed in %v\n", copyDuration)

		// 플러그인 실행 - 순서대로
		processStartTime := time.Now()
		for _, pluginCfg := range cfg.Plugin {
			pluginStartTime := time.Now()
			plugin, err := plugins.GetPlugin(&pluginCfg)
			if err != nil {
				log.Fatal("Plugin not found: ", err)
			}

			fmt.Printf("Plugin: %s\n", plugin.GetName())
			err = plugin.Process(cfg)

			if err != nil {
				log.Fatal("Plugin process failed: ", err)
			}

			fmt.Printf("⭕ %s Plugin processing completed\n", pluginCfg.Name)
			pluginDuration := time.Since(pluginStartTime)
			fmt.Printf("%s Plugin work time: %v\n", pluginCfg.Name, pluginDuration)
		}

		fmt.Println("✅ Plugin processing completed")
		processDuration := time.Since(processStartTime)
		fmt.Printf("✅ File processing completed in %v\n", processDuration)

		// 전체 작업 시간
		totalWorkTime := time.Since(workStartTime)
		fmt.Printf("Total work time: %v (Copy: %v, Process: %v)\n", totalWorkTime, copyDuration, processDuration)
	} else {
		fmt.Println("CHECK: System not ready for processing")
	}
}
