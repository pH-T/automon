package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var SysmonConfigs []string = []string{
	"https://raw.githubusercontent.com/NextronSystems/aurora-helpers/master/sysmon-config/aurora-sysmon-config.xml",
	"https://raw.githubusercontent.com/Neo23x0/sysmon-config/master/sysmonconfig-export.xml",
	"https://raw.githubusercontent.com/SwiftOnSecurity/sysmon-config/master/sysmonconfig-export.xml",
	"https://raw.githubusercontent.com/NextronSystems/evtx-baseline/master/sysmon-intense.xml",
	"https://raw.githubusercontent.com/OTRF/Blacksmith/master/resources/configs/sysmon/sysmon.xml",
}

var SysmonFolder string = "sysmon"
var SysmonZipFile string = "sysmon.zip"
var SysmonConfigFile string = "sysmon-config.xml"

var ListConfigsFlag bool
var SysmonDownloadOnlyFlag bool
var UninstallSysmonFlag bool
var UseConfigIndexFlag int
var SysmonURLFlag string
var SysmonArchFlag string
var ForceInstallFlag bool
var ConfigURLFlag string

func main() {
	flag.BoolVar(&ListConfigsFlag, "listconfigs", false, "Lists hardcoded Sysmon config URLs")
	flag.BoolVar(&ForceInstallFlag, "force", false, "Uninstalls Sysmon before installing")
	flag.BoolVar(&SysmonDownloadOnlyFlag, "sysmondownload", false, "Just downloads Sysmon")
	flag.BoolVar(&UninstallSysmonFlag, "uninstall", false, "Uninstall Sysmon")
	flag.IntVar(&UseConfigIndexFlag, "config", -1, "Which config should be used")
	flag.StringVar(&SysmonURLFlag, "sysmonURL", "https://download.sysinternals.com/files/Sysmon.zip", "URL to download Sysmon zip")
	flag.StringVar(&SysmonArchFlag, "arch", "64", "Which Sysmon version to use: 64 or 32")
	flag.StringVar(&ConfigURLFlag, "configURL", "", "URL to download config")
	flag.Parse()

	if ListConfigsFlag {
		fmt.Println("Available Sysmon configs:")
		for i, e := range SysmonConfigs {
			fmt.Printf("[%d] %s\n", i, e)
		}
		return
	}

	if UninstallSysmonFlag {
		uninstall()
		return
	}

	if SysmonDownloadOnlyFlag {
		fmt.Println("Downloading Sysmon zip...")
		if err := downloadFile(SysmonURLFlag, SysmonZipFile); err != nil {
			fmt.Printf("ERROR: was not able to download Sysmon from : %s : %v\n", SysmonURLFlag, err)
			return
		}

		fmt.Println("Unzipping Sysmon...")
		if err := unzip(SysmonZipFile, SysmonFolder); err != nil {
			fmt.Printf("ERROR: was not able to download Sysmon from : %s : %v\n", SysmonURLFlag, err)
			return
		}

		if err := os.Remove(SysmonZipFile); err != nil {
			fmt.Printf("NOTICE: was not able to remove: %s : %v\n", SysmonZipFile, err)
		}

		return
	}

	if UseConfigIndexFlag >= 0 || ConfigURLFlag != "" {
		configURL := ""
		if UseConfigIndexFlag >= 0 {
			configURL = SysmonConfigs[UseConfigIndexFlag]
		} else {
			configURL = ConfigURLFlag
		}

		fmt.Printf("Using config: %s\n", configURL)

		fmt.Println("Downloading config...")
		if err := downloadFile(configURL, SysmonConfigFile); err != nil {
			fmt.Printf("ERROR: was not able to download Sysmon config: %s : %v\n", configURL, err)
			return
		}

		fmt.Println("Downloading Sysmon zip...")
		if err := downloadFile(SysmonURLFlag, SysmonZipFile); err != nil {
			fmt.Printf("ERROR: was not able to download Sysmon from : %s : %v\n", SysmonURLFlag, err)
			return
		}

		fmt.Println("Unzipping Sysmon...")
		if err := unzip(SysmonZipFile, SysmonFolder); err != nil {
			fmt.Printf("ERROR: was not able to download Sysmon from : %s : %v\n", SysmonURLFlag, err)
			return
		}

		if err := os.Remove(SysmonZipFile); err != nil {
			fmt.Printf("NOTICE: was not able to remove: %s : %v\n", SysmonZipFile, err)
		}

		if ForceInstallFlag {
			uninstall()
		}

		install()

		return
	}

	flag.PrintDefaults()
}

func uninstall() {
	fmt.Println("Uninstalling Sysmon...")
	switch SysmonArchFlag {
	case "64":
		c := filepath.Join(SysmonFolder, "Sysmon64.exe")
		out, err := uninstallSysmon(c, SysmonConfigFile)
		if err != nil {
			fmt.Printf("ERROR: was not able to uninstall Sysmon 64: %v\n", err)
			fmt.Println(out)
			return
		}
		fmt.Println(out)
	case "32":
		c := filepath.Join(SysmonFolder, "Sysmon.exe")
		out, err := uninstallSysmon(c, SysmonConfigFile)
		if err != nil {
			fmt.Printf("ERROR: was not able to uninstall Sysmon 32: %v\n", err)
			fmt.Println(out)
			return
		}
		fmt.Println(out)
	default:
		fmt.Printf("ERROR: %s is unkown Sysmon architecture\n", SysmonArchFlag)
		return
	}
}

func install() {
	fmt.Println("Installing Sysmon...")
	switch SysmonArchFlag {
	case "64":
		c := filepath.Join(SysmonFolder, "Sysmon64.exe")
		out, err := installSysmon(c, SysmonConfigFile)
		if err != nil {
			fmt.Printf("ERROR: was not able to install Sysmon 64: %v\n", err)
			fmt.Println(out)
			return
		}
		fmt.Println(out)
	case "32":
		c := filepath.Join(SysmonFolder, "Sysmon.exe")
		out, err := installSysmon(c, SysmonConfigFile)
		if err != nil {
			fmt.Printf("ERROR: was not able to install Sysmon 32: %v\n", err)
			fmt.Println(out)
			return
		}
		fmt.Println(out)
	default:
		fmt.Printf("ERROR: %s is unkown Sysmon architecture\n", SysmonArchFlag)
		return
	}
}

func uninstallSysmon(sysmonFile, configFile string) (string, error) {
	cmd := exec.Command(sysmonFile, "-u")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func installSysmon(sysmonFile, configFile string) (string, error) {
	cmd := exec.Command(sysmonFile, "-accepteula", "-i", configFile)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func unzip(zipFile, dst string) error {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("was not able to zip.OpenReader(): %v", err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("Unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("Creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("was not able to os.MkdirAll(): %v", err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("was not able to os.OpenFile(): %v", err)
		}
		defer dstFile.Close()

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("was not able to f.Open(): %v", err)
		}
		defer fileInArchive.Close()

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("was not able to io.Copy(): %v", err)
		}
	}

	return nil

}

func downloadFile(url, dst string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("was not able to http.NewRequest(): %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("was not able to client.Do(): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statuscode != 200 - maybe new download url?")
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("was not able to os.Create(): %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("was not able to io.Copy(): %v", err)
	}

	return nil
}
