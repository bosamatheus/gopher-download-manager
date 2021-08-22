package download

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type download struct {
	url           string
	targetPath    string
	totalSections int
}

// New returns a new download instance.
func New(url, filename string, threads string) (*download, error) {
	totalSections, err := stringToInt(threads)
	if err != nil {
		return nil, fmt.Errorf("invalid number of threads: %v", threads)
	}
	return &download{
		url:           url,
		targetPath:    "data/" + filename,
		totalSections: totalSections,
	}, nil
}

// Do performs the download.
func (d download) Do() error {
	size, err := d.getDownloadSize()
	if err != nil {
		return err
	}

	sections := d.makeSections(size)
	d.downloadAllSequentially(sections)

	err = d.merge(sections)
	if err != nil {
		return err
	}
	return nil
}

func (d download) getDownloadSize() (int, error) {
	fmt.Println("getting download size")
	req, err := newRequest("HEAD", d.url)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	if res.StatusCode > 299 {
		return 0, fmt.Errorf("can't process, response is %v", res.StatusCode)
	}

	size, err := stringToInt(res.Header.Get("Content-Length"))
	if err != nil {
		return 0, err
	}
	fmt.Printf("download size is %v\n", size)
	return size, nil
}

// makeSections makes a list of sections to download.
// For example, if total size is 100 bytes and total sections is 10,
// each section size will be 10 bytes. So the sections will be:
// [[0, 10], [11, 21], [22, 32], [33, 43], [44, 54], [55, 65], [66, 76],
// [77, 87], [88, 98], [99, 99]]
func (d download) makeSections(size int) [][2]int {
	sectionSize := size / d.totalSections
	fmt.Printf("each section size is %v\n", sectionSize)

	var sections = make([][2]int, d.totalSections)
	for i := range sections {
		if i == 0 {
			// starting byte of first section
			sections[i][0] = 0
		} else {
			// starting byte of other sections
			sections[i][0] = sections[i-1][1] + 1
		}

		if i == d.totalSections-1 {
			// ending byte of last section
			sections[i][1] = size - 1
		} else {
			// ending byte of other sections
			sections[i][1] = sections[i][0] + sectionSize
		}
	}
	fmt.Printf("sections are:\n%v\n", sections)
	return sections
}

func (d download) downloadAllSequentially(sections [][2]int) error {
	for i, s := range sections {
		err := d.downloadSection(i, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d download) downloadSection(i int, s [2]int) error {
	req, err := newRequest("GET", d.url)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", s[0], s[1]))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("downloaded %v bytes for section %v: %v\n", res.ContentLength, i, s)

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

func (d download) merge(sections [][2]int) error {
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
