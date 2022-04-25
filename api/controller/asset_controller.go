package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/headdetect/its-a-twitter/api/utils"
)

func HandleAsset(writer http.ResponseWriter, request *http.Request) {
	requestedAsset, exists := GetPathValue(request, 0)

	if !exists {
		NotFoundResponse(writer)
		return
	}

	// Prevent a path-traversal attack //
	assetBase := filepath.Base(requestedAsset)
	mimeType := fmt.Sprintf("image/%s", strings.TrimLeft(filepath.Ext(assetBase), "."))

	mediaPath, _ := utils.GetStringOrDefault("MEDIA_PATH", "./assets/media")
	fullFilePath := fmt.Sprintf("%s/%s", mediaPath, assetBase)
	diskFile, err := os.ReadFile(fullFilePath)

	if err != nil {
		NotFoundResponse(writer)
		return
	}

	writer.Header().Add("Content-Type", mimeType)
	writer.Header().Add("Content-Disposition", "inline")
	writer.Header().Add("Cache-Control", "max-age=31536000")
	writer.WriteHeader(http.StatusOK)
	writer.Write(diskFile)
}
