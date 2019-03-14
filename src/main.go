package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ImageItem struct {
	ID  string `json:"id"`
	Isu string `json:"isu"`
	Itg int    `json:"itg"`
	Ity string `json:"ity"`
	Oh  int    `json:"oh"`
	Ou  string `json:"ou"`
	Ow  int    `json:"ow"`
	Pt  string `json:"pt"`
	Rh  string `json:"rh"`
	Rid string `json:"rid"`
	Rt  int    `json:"rt"`
	Ru  string `json:"ru"`
	S   string `json:"s"`
	Sc  int    `json:"sc"`
	St  string `json:"st"`
	Th  int    `json:"th"`
	Tu  string `json:"tu"`
	Tw  int    `json:"tw"`
}

func main() {

	imageFolderName := flag.String("folder", "img", "a string")
	searchTerm := flag.String("searchTerm", "apple seed", "a string")
	count := flag.Int("count", 150, "a int")
	flag.Parse()

	println("Folder:" + *imageFolderName)
	println("Search term:" + *searchTerm)
	println("Count:" + strconv.Itoa(*count))

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	folderPath := dir + "\\" + *imageFolderName
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, os.ModePerm)
	}

	//Preparing term for the search
	term := strings.Replace(*searchTerm, " ", "+", -1)
	var index int

	for index <= *count {

		imageLinks := GetImageLinks(term, index)
		println("GOT ITEMS: ", len(imageLinks))

		for _, imageLink := range imageLinks {

			img := GetImageItemFromJson(imageLink)
			fullName := GetFileFullName(img, folderPath)
			index += 1
			println(strconv.Itoa(index) + "." + img.Ou)
			println(strconv.Itoa(index) + "." + fullName)

			if err := DownloadImage(fullName, img.Ou); err != nil {
				//panic(err)
				println("<==============ERROR======================>")
				println(err.Error())
				println(imageLink)
				println(img.Ou)
				println("<==============ERROR======================>")
			}

			if index >= *count {
				println(strconv.Itoa(index))
				break
			}
		}
	}
}

func GetImageLinks(term string, index int) []string {

	url := "https://www.google.com/search?q=" + term + "&oq=" + term

	if index == 0 {
		url += "&biw=1536&bih=723&tbm=isch&sa=1&ei=6qqGXM_oDenYjwSw1b-oAw"
	} else {
		url += "&ijn=" + strconv.Itoa(index/100) + "&start=" + strconv.Itoa(index) +
			"&biw=1536&bih=723&tbm=isch&sa=1&ei=6qqGXM_oDenYjwSw1b-oAw&yv=3&as_st=y&tbm=isch&asearch=ichunk&async=_id:rg_s,_pms:s,_fmt:pc"
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	page := string(body)

	r := regexp.MustCompile("<div class=\"rg_meta notranslate\">([\\s\\S]*?)</div>")
	imageLinks := r.FindAllStringSubmatch(page, -1)

	var result []string

	for _, ImageItem := range imageLinks {
		result = append(result, ImageItem[1])
	}

	return result
}

func GetFileFullName(img ImageItem, folderPath string) string {
	url := img.Ou
	fileName := img.ID[0 : len(img.ID)-1]

	if len(img.Ity) != 0 {
		fileName += "_" + url[strings.LastIndex(img.Ou, "/")+1:strings.LastIndex(img.Ou, ".")] + "." + img.Ity
	} else {
		fileName += ".jpeg"
	}

	fileName = folderPath + "\\" + CleanFileName(fileName)
	//println("<===================================================>")
	// println("-----------------------------------------------------")
	// println("Source:", img.Ou)
	//println(fileName)
	// println(img.Ou)
	//println("<===================================================>")

	return fileName
}

func CleanFileName(fileName string) string {
	symbols := [6]string{"*", "?", "%", "\\", "/"}
	for _, symbol := range symbols {
		fileName = strings.Replace(fileName, symbol, "", -1)
	}
	return fileName
}

func GetImageItemFromJson(jsonString string) ImageItem {
	img := ImageItem{}

	err := json.Unmarshal([]byte(jsonString), &img)
	if err != nil {
		panic(err)
	}

	return img
}

func DownloadImage(filePath string, url string) error {

	// Get the data
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("authority", "cdn-images-1.medium.com")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
