package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if _, ok := pathsToUrls[r.URL.Path]; !ok {
			fallback.ServeHTTP(rw, r)
		}

		http.Redirect(rw, r, pathsToUrls[r.URL.Path], 301)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type pathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)

	if err != nil {
		panic(err)
	}

	return MapHandler(buildMap(parsedYaml), fallback), err
}

func parseYaml(yml []byte) ([]pathUrl, error) {
	result := []pathUrl{}
	err := yaml.Unmarshal(yml, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	url_map := make(map[string]string, len(pathUrls))

	for _, pathUrl := range pathUrls {
		url_map[pathUrl.Path] = pathUrl.Url
	}

	return url_map
}
