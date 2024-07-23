package apputil

import (
	"os"
	"regexp"
	"strings"
)

import (
	"bytes"
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"encoding/json"
	"fmt"
	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"path"
)

var hmrRegex = regexp.MustCompile(`(?m)const socketUrl = getSocketUrl\(\((.*)\)\)`)

const filePermission = os.FileMode(0o666)

// ShouldReplace returns true if the file being served should have its contents scanned for tokens
func ShouldReplace(path string, fileList []string) bool {
	for _, file := range fileList {
		if strings.HasSuffix(path, file) {
			return true
		}
	}
	return false
}

// ReplaceTokens replaces tokens in static files and caches the file based on the content hash in order to avoid
// scanning the same files over and over. Note that the cache will be lost on restart
func ReplaceTokens(c echo.Context, content []byte) []byte {
	appName := c.Get(apptypes.AppNameKey).(string)
	matches := hmrRegex.FindSubmatch(content)
	if len(matches) > 1 {
		sum := md5.Sum(content) //nolint:gosec
		hash := hex.EncodeToString(sum[:])
		if err := os.MkdirAll(path.Join(os.TempDir(), "gateway"), filePermission); err != nil {
			log.Warn().Err(err).Msgf("failed to create cache directory")
		}

		cachePath := path.Join(os.TempDir(), "gateway", hash)
		if _, err := os.Stat(cachePath); err == nil {
			cnt, err := os.ReadFile(cachePath)
			if err != nil {
				log.Warn().Err(err).Msgf("failed to load cached files")
			} else {
				return cnt
			}
		}

		val := matches[1]
		var asMap map[string]any
		if err := json.Unmarshal(val, &asMap); err != nil {
			log.Warn().Err(err).Msgf("fauled to unmarshal hmr socket url")
		} else {
			asMap["port"] = ""
			asMap["path"] = fmt.Sprintf("/app/%s%s", appName, asMap["path"])
			b, _ := json.Marshal(asMap)
			content = bytes.ReplaceAll(content, val, b)
			if err := os.WriteFile(cachePath, content, filePermission); err != nil {
				log.Warn().Err(err).Msgf("failed to write modified file to cache")
			}
		}
	}

	return content
}
