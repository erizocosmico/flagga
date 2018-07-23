package flagga

// Extractor extracts values from the sources to fill the flag value.
type Extractor interface {
	// Get checks the sources and tries to assign the flag value.
	Get(sources []Source, dst Value) (bool, error)
}

// Env returns an Extractor that will match environment variables with
// the given key.
func Env(key string) Extractor {
	return envExtractor(key)
}

type envExtractor string

func (e envExtractor) Get(sources []Source, dst Value) (bool, error) {
	for _, s := range sources {
		if _, ok := s.(envSource); !ok {
			continue
		}

		ok, err := s.Get(string(e), dst)
		if err != nil || !ok {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

type jsonExtractor string

// JSON returns an Extractor that will match the given key in a provided
// JSON file to set as value for the flag.
func JSON(key string) Extractor {
	return jsonExtractor(key)
}

func (e jsonExtractor) Get(sources []Source, dst Value) (bool, error) {
	for _, s := range sources {
		if _, ok := s.(*jsonSource); !ok {
			continue
		}

		ok, err := s.Get(string(e), dst)
		if err != nil {
			return false, err
		}

		if !ok {
			continue
		}

		return true, nil
	}

	return false, nil
}
