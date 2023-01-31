package main
 
import (
	"fmt"
	"log"
	"net/http"
	"strings" 
	"os"
	"io"
	"html/template"
	"path"
)





func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
 
	w.Write([]byte("Просто Отсканируйте Код"))
}








 

func showDescription(w http.ResponseWriter, r *http.Request) {

	


	t, err := template.ParseFiles("templates/description.html")

	t.ExecuteTemplate(w, "description.html", nil)

	if err != nil{

		fmt.Fprintf(w, err.Error())
	}



	id := r.URL.Query().Get("id") 

	strsplit := strings.Split(id,",") 

	for docnumber := range strsplit {
  	 	
		filename :=string("service/"+strsplit[docnumber]+".txt")

  	 	content, err := os.ReadFile(filename)

    	if err != nil {
        	w.Write([]byte("Услуга не найдена"))
        	fmt.Fprintf(w,"<br>")
   		 } else {
			w.Write([]byte(content))
			fmt.Fprintf(w,"<br>")
			fmt.Fprintf(w,"<br>")
		 }
 	 }

	



	

}
 
func uploadFile (w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		handleUpload (w,r)
		return
	}


	t, err := template.ParseFiles("templates/upload.html")

	if err != nil{

		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "upload.html", nil)

}


func handleUpload  (w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10<<20) //10MB

	file, fileHeader, err := r.FormFile("myfile")

	if err != nil {

		http.Error(w,"http.StatusBadRequest",http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := path.Base(fileHeader.Filename)

	dest, err := os.Create("service/"+filename)

	if err !=nil {

		http.Error(w,"http.StatusInternalServerError",http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	if _, err = io.Copy (dest,file); err !=nil {

		http.Error(w,"http.StatusInternalServerError",http.StatusInternalServerError)
		return

	}

	http.Redirect(w,r,"?success=true",http.StatusSeeOther)
}






 
func main() {


	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/description", showDescription)
	mux.HandleFunc("/upload/", uploadFile)
	
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080(locallhost)")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}