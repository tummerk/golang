package main

import (
	"CustomappTest/internal/application"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func main() {
	rtpStr := strings.TrimPrefix(os.Args[1], "rtp=")
	rtp, err := strconv.ParseFloat(rtpStr, 64)
	if err != nil {
		slog.Error("invalid rtp value: %v", slog.String("error", err.Error()))
	}
	port := ":64333"

	application.Run(port, rtp)
	if err != nil {
		slog.Error("run err: %v", err)
	}
}
