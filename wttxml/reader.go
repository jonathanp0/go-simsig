package wttxml

import (
	"archive/zip"
	"errors"
	"io/ioutil"
)

//Load a WTT file and return the contents of the SavedTimetable.xml file
func ReadSavedTimetable(path string) ([]byte, error) {
	wtt, err := zip.OpenReader(path)
	if err != nil {
		err = errors.New("cannot open WTT file" + err.Error())
		return nil, err
	}

	defer wtt.Close()

	for _, file := range wtt.File {
		if file.Name == "SavedTimetable.xml" {
			fileReader, err := file.Open()
			if err != nil {
				err = errors.New("error decompressing SavedTimetable.xml" + err.Error())
				return nil, err
			}
			return ioutil.ReadAll(fileReader)
		}
	}

	err = errors.New("Could not find SavedTimetable.xml in WTT file")

	return nil, err
}
