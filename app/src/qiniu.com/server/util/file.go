package util

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"

	"qiniu.com/server/typo"
)

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func GetNodes(fname string) (nodes []typo.NodeConfigItem, err error) {
	var data []byte
	data, err = ioutil.ReadFile(fname)

	if err != nil {
		Logger.Errorf("read node list failed, error: %s", err)
		return
	}

	err = json.Unmarshal(data, &nodes)
	if err != nil {
		Logger.Errorf("parse node list failed, error: %s", err)
		return
	}
	return
}

func GetGroups(fname string) (groups []typo.GroupConfigItem, err error) {
	var data []byte
	data, err = ioutil.ReadFile(fname)

	if err != nil {
		Logger.Errorf("read group list failed, error: %s", err)
		return
	}

	err = json.Unmarshal(data, &groups)
	if err != nil {
		Logger.Errorf("parse group list failed, error: %s", err)
		return
	}

	return
}
