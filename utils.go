package microutils

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	prefixB  = "Bytes"
	prefixKB = "KBytes"
	prefixMB = "MBytes"
	prefixGB = "GBytes"
	prefixTB = "TBytes"
)

type FileInfo struct {
	Name string
	Path string
	Data []byte
}

func (f *FileInfo) Size() int {
	if f == nil {
		return 0
	}
	return len(f.Data)
}

func BytesToSizeString(s int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	size := float64(s)

	switch {
	case s < KB:
		return fmt.Sprintf("%.4f B", size)
	case s < MB:
		return fmt.Sprintf("%.4f KB", size/KB)
	case s < GB:
		return fmt.Sprintf("%.4f MB", size/MB)
	case s < TB:
		return fmt.Sprintf("%.4f GB", size/GB)
	default:
		return fmt.Sprintf("%.4f TB", size/TB)
	}
}

func PrintFatalErr(err error) {
	fmt.Printf("ERROR: %v\n", err)
	os.Exit(1)
}

func PrintErr(err error) {
	fmt.Printf("ERROR: %v", err)
}

func PrintJSON(pretty bool, v any) error {
	if pretty {
		return jsonPrintPretty(v)
	}
	return jsonPrint(v)
}

func jsonPrint(v any) error {
	data, err := json.Marshal(v)
	if err == nil {
		fmt.Println(string(data))
	}
	return err
}

func jsonPrintPretty(v any) error {
	data, err := json.MarshalIndent(v, " ", "  ")
	if err == nil {
		fmt.Println(string(data))
	}
	return err
}

func IsInputFromPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) == 0
}
