package main

import (
	"awesomeProject11/Book"
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

func workerFB2Remove(id int, TitleWork string, wg *sync.WaitGroup) {
	FileName := fmt.Sprintf("%d.fb2", id)
	Dir := path.Join("Novel", TitleWork, "fb2", FileName)
	err := os.Remove(Dir)
	if err != nil {
		log.Fatalf("id worker %d err %v", id, err)
	}
	wg.Done()
}

func workerTXTRemove(id int, TitleWork string, wg *sync.WaitGroup) {
	FileName := fmt.Sprintf("%d.txt", id)
	Dir := path.Join("Novel", TitleWork, "txt", FileName)
	err := os.Remove(Dir)
	if err != nil {
		log.Fatalf("id worker %d err %v", id, err)
	}
	wg.Done()
}

func workerFB2Create(id int, TitleWork string, wg *sync.WaitGroup) {
	FileName := fmt.Sprintf("%d.txt", id)
	Dir := path.Join("Novel", TitleWork, "txt", FileName)
	f, err := os.Open(Dir)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return
	}
	Chapter := fmt.Sprintf("Chapter %d", id)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//fmt.Println(lines)
	fb2 := &Book.FB2{
		Description: Book.Description{
			TitleInfo: Book.TitleInfo{
				Genre: Book.Genre{Text: "sf"},
				Author: Book.Author{
					FirstName: Book.FirstName{Text: "John"},
					LastName:  Book.LastName{Text: "Doe"},
				},
				BookTitle: Book.BookTitle{Text: "Example Book"},
			},
			DocumentInfo: Book.DocumentInfo{
				Author: Book.Author{
					FirstName: Book.FirstName{Text: "John"},
					LastName:  Book.LastName{Text: "Doe"},
				},
				Date: Book.Date{Text: "2023-11-15"},
			},
		},

		Body: Book.Body{
			Section: Book.Section{
				Title:     Book.Title{Text: Chapter},
				Paragraph: Book.Paragraph{Text: lines},
			},
		},
	}

	out, _ := xml.MarshalIndent(fb2, " ", " ")

	FileNameOut := fmt.Sprintf("%d.fb2", id)
	DirOut := path.Join("Novel", TitleWork, "fb2", FileNameOut)
	err = os.WriteFile(DirOut, out, 0644)
	if err != nil {
		log.Fatalf("worker %d, err %v", id, err)
	}
	wg.Done()
}

func MergentFB2(count int, TitleWork string) {
	//files := []string{"book1.fb2", "book2.fb2", "book3.fb2"}

	// Create a new FB2 file
	merged := Book.FB2{
		Description: Book.Description{
			TitleInfo: Book.TitleInfo{
				Genre: Book.Genre{Text: "sf"},
				Author: Book.Author{
					FirstName: Book.FirstName{Text: "John"},
					LastName:  Book.LastName{Text: "Doe"},
				},
				BookTitle: Book.BookTitle{Text: "Merged Book"},
			},
			DocumentInfo: Book.DocumentInfo{
				Author: Book.Author{
					FirstName: Book.FirstName{Text: "John"},
					LastName:  Book.LastName{Text: "Doe"},
				},
				Date: Book.Date{Text: "2023-11-15"},
			},
		},
		Body: Book.Body{
			Section: Book.Section{
				Title: Book.Title{Text: "Merged Book"},
			},
		},
	}

	// Read each file and append its content to the merged FB2 file

	for i := 1; i < count; i++ {
		FileName := fmt.Sprintf("%d.fb2", i)
		Dir := path.Join("Novel", TitleWork, "fb2", FileName)
		data, err := os.ReadFile(Dir)
		if err != nil {
			fmt.Printf("Error reading file: %v", err)
			return
		}

		var book Book.FB2
		err = xml.Unmarshal(data, &book)
		if err != nil {
			fmt.Printf("Error unmarshalling file: %v", err)
			return
		}

		merged.Body.Section.Paragraph.Text = append(merged.Body.Section.Paragraph.Text, book.Body.Section.Paragraph.Text...)
	}

	// Write the merged FB2 file
	out, _ := xml.MarshalIndent(merged, " ", " ")
	fileName := "merged.fb2"
	Dir := path.Join("Novel", TitleWork, "fb2", fileName)
	err := os.WriteFile(Dir, out, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v", err)
		return
	}
}

func MergentTXT(count int, TitleWork string) {
	fileName := "merged.txt"
	Dir := path.Join("Novel", TitleWork, "txt", fileName)
	outFile, err := os.Create(Dir)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outFile.Close()
	for i := 1; i < count; i++ {
		FileName := fmt.Sprintf("%d.txt", i)
		Dir := path.Join("Novel", TitleWork, "txt", FileName)
		inFile, err := os.Open(Dir)
		if err != nil {
			fmt.Printf("Error reading file: %v", err)
			return
		}
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			log.Fatalf("Error copying file: %v", err)
		}
		inFile.Close()
	}
}

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the file name from the URL
	TitleWork := r.URL.Query().Get("TitleWork")
	fileName := r.URL.Query().Get("fileName")
	T := r.URL.Query().Get("type")

	// Join to get the full file path
	filePath := filepath.Join("Novel", TitleWork, T, fileName)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the content-disposition header to force a download
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))

	// Serve the file
	http.ServeFile(w, r, filePath)
}

func main() {
	//var name string
	//name = "a-will-eternal"
	//var wg sync.WaitGroup
	//Dir := path.Join("Novel", name, "txt")
	//files, err := os.ReadDir(Dir)
	//if err != nil {
	//	panic(err)
	//}
	//var command uint8
	//count := len(files)
	//fmt.Printf("введите 0 или Ctrl+C для выхода\nвведите 1 для создания fb2 из txt\nвведите 2 для удаления txt\nвведите 3 для удаления fb2\nвведите 4 для обьединения fb2\nвведите 5 для обьединения txt\n")
	//fmt.Scanln(&command)
	//switch command {
	//case 0:
	//	os.Exit(1)
	//case 1:
	//	for i := 1; i < count; i++ {
	//		wg.Add(1)
	//		go workerFB2Create(i, name, &wg)
	//	}
	//	wg.Wait()
	//case 2:
	//	for i := 1; i < count; i++ {
	//		wg.Add(1)
	//		go workerTXTRemove(i, name, &wg)
	//	}
	//	wg.Wait()
	//case 3:
	//	for i := 1; i < count; i++ {
	//		wg.Add(1)
	//		go workerFB2Remove(i, name, &wg)
	//	}
	//	wg.Wait()
	//case 4:
	//	MergentFB2(count, name)
	//	wg.Wait()
	//case 5:
	//	MergentTXT(count, name)
	//	wg.Wait()
	//}

	var fe []string
	Dir := "Novel"
	fies, _ := os.ReadDir(Dir)
	for _, file := range fies {
		if file.IsDir() {
			file.Type()
			fe = append(fe, file.Name())
		}
	}
	fmt.Println(fe)

	//Serve static files from the "static1" directory
	http.Handle("/Novel/", http.StripPrefix("/Novel/", http.FileServer(http.Dir("Novel"))))

	http.HandleFunc("/download", downloadFileHandler)
	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
