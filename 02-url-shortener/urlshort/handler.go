package urlshort

import (
	"net/http"
	"strings"

	"go.yaml.in/yaml/v4"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uri := strings.TrimSuffix(r.RequestURI, "/")
		if redirectURL, ok := pathsToUrls[uri]; ok {
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYAML)
	return MapHandler(pathMap, fallback), nil
}

type URLRedirect struct {
	Path string `yaml:path`
	URL  string `yaml:url`
}

func parseYAML(yml []byte) ([]URLRedirect, error) {
	parsedYAML := []URLRedirect{}
	err := yaml.Unmarshal(yml, &parsedYAML)
	if err != nil {
		return nil, err
	}
	return parsedYAML, nil
}

func buildMap(paths []URLRedirect) map[string]string {
	pathMap := make(map[string]string)
	for _, p := range paths {
		path := strings.TrimSuffix(p.Path, "/")
		pathMap[path] = p.URL
	}
	return pathMap
}
