package countlog

import (
	"os"
	"path/filepath"
	"time"
)

type fileLogOutput struct {
	windowSize          int64
	logFile             string
	rotateAfter         int64
	openedFile          *os.File
	openedFileArchiveTo string
}

func (output *fileLogOutput) Close() {
	if output.openedFile != nil {
		output.openedFile.Close()
	}
}

func (output *fileLogOutput) openLogFile(logFile string) {
	var err error
	output.openedFile, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		os.Stderr.Write([]byte("failed to open log file: " +
			logFile + ", " + err.Error() + "\n"))
		os.Stderr.Sync()
	}
	output.openedFileArchiveTo = logFile + "." + time.Now().Format("200601021504")
}

func (output *fileLogOutput) archiveLogFile(logFile string) {
	output.openedFile.Close()
	output.openedFile = nil
	err := os.Rename(logFile, output.openedFileArchiveTo)
	if err != nil {
		os.Stderr.Write([]byte("failed to rename to archived log file: " +
			output.openedFileArchiveTo + ", " + err.Error() + "\n"))
		os.Stderr.Sync()
	}
}

func (output *fileLogOutput) OutputLog(timestamp int64, formattedEvent []byte) {
	if timestamp > output.rotateAfter {
		now := time.Now()
		output.rotateAfter = (int64(now.UnixNano()/output.windowSize) + 1) * output.windowSize
		output.archiveLogFile(output.logFile)
		output.openLogFile(output.logFile)
	}
	if output.openedFile != nil {
		output.openedFile.Write(formattedEvent) // silently ignore error
	}
}

type osFileLogOutput struct {
	logFile *os.File
}

func (output *osFileLogOutput) Close() {
}

func (output *osFileLogOutput) OutputLog(timestamp int64, formattedEvent []byte) {
	output.logFile.Write(formattedEvent)
}

func NewFileLogOutput(logFile string) LogOutput {
	switch logFile {
	case "STDOUT":
		return &osFileLogOutput{os.Stdout}
	case "STDERR":
		return &osFileLogOutput{os.Stderr}
	default:
		output := &fileLogOutput{
			windowSize: int64(time.Hour),
		}
		err := os.MkdirAll(filepath.Dir(logFile), 0755)
		if err != nil {
			os.Stderr.Write([]byte("failed to create dir for log file: " +
				filepath.Dir(logFile) + ", " + err.Error() + "\n"))
			os.Stderr.Sync()
		}
		output.openLogFile(logFile)
		output.rotateAfter = (int64(time.Now().UnixNano()/output.windowSize) + 1) * output.windowSize
		return output
	}
}
