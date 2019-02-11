package mempool

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"testing"
)

const (
	txErrorVarPat = "^[ 	]+(Err[A-Z][A-Za-z]+)"
)

func Test_faultMapCoverage(t *testing.T) {
	errNames, err := openAndExtract("../types/errors.go", txErrorVarPat)
	if err != nil {
		t.Errorf("failed to extract file: %v", err)
	}

	for k, v := range faultMap {
		if k.Error() != v.Error() {
			t.Errorf("Invalid Error wrapper %v, want %v", v.Error(), k.Error())
		}
	}

	if len(errNames) != len(faultMap) {
		sort.Strings(errNames)
		t.Errorf("txPenalties does not cover all put tx errors %v, want %v \n%v", len(faultMap),
			len(errNames), errNames )
	}

}

func openAndExtract(filepath, pattern string) ([]string, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	result := make([]string, 0, 10)
	for scanner.Scan() {
		found := reg.FindStringSubmatch(scanner.Text())
		if found == nil {
			continue
		}
		result = append(result, found[1])
	}

	return result, nil
}


func getKeys(m map[error]error) []error {
	keys := make([]error, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
