package download

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type download struct {
	url           string
	targetPath    string
	totalSections int
}

// New returns a new download instance.
func New(url, targetPath string, totalSections int) *download {
	return &download{
		url:           url,
		targetPath:    targetPath,
		totalSections: totalSections,
	}
}

// Do performs the download.
func (d download) Do() error {
	fmt.Println("Connecting...")
	req, err := d.newRequest("HEAD")
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("got %v\n", res.StatusCode)
	if res.StatusCode > 299 {
		return fmt.Errorf("can't process, response is %v", res.StatusCode)
	}

	size, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	fmt.Printf("file size is %v\n", size)

	var sections = make([][2]int, d.totalSections)
	eachSize := size / d.totalSections
	fmt.Printf("each section size is %v\n", eachSize)

	// Example: if total size is 100 bytes and total sections is 10,
	// each section size will be 10 bytes. So the sections will be:
	// [[0, 10], [11, 21], [22, 32], [33, 43], [44, 54], [55, 65], [66, 76],
	// [77, 87], [88, 98], [99, 99]]
	for i := range sections {
		if i == 0 {
			// starting byte of first section
			sections[i][0] = 0
		} else {
			// starting byte of other sections
			sections[i][0] = sections[i-1][1] + 1
		}

		if i < d.totalSections-1 {
			// ending byte of other sections
			sections[i][1] = sections[i][0] + eachSize
		} else {
			// ending byte of last section
			sections[i][1] = size - 1
		}
	}
	fmt.Println(sections)
	for i, s := range sections {
		err = d.downloadSection(i, s)
		if err != nil {
			return err
		}
		fmt.Printf("section %v completed\n", i)
	}
	err = d.mergeFiles(sections)
	if err != nil {
		return err
	}
	return nil
}

func (d download) newRequest(method string) (*http.Request, error) {
	req, err := http.NewRequest(method, d.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Gopher Download Manager v1.0")
	return req, nil
}

func (d download) downloadSection(i int, s [2]int) error {
	req, err := d.newRequest("GET")
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", s[0], s[1]))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("downloaded %v bytes for section %v? %v\n", res.ContentLength, i, s)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("temp/section-%v.tmp", i), b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (d download) mergeFiles(sections [][2]int) error {
	f, err := os.OpenFile(d.targetPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := range sections {
		b, err := ioutil.ReadFile(fmt.Sprintf("temp/section-%v.tmp", i))
		if err != nil {
			return err
		}
		n, err := f.Write(b)
		if err != nil {
			return err
		}
		fmt.Printf("%v bytes merged\n", n)
	}
	return nil
}
