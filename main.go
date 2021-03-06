package main

import (
    "html/template"
    "log"
    //"fmt"
    "math"
    "net/http"
    "os"
    "time"
    "net/url"
    "strconv"
    "bytes"
    
    "github.com/joho/godotenv"
    "github.com/freshman-tech/news-demo-starter-files/news"
)



type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
	Domains string
}


func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}


func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}


func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}


func indexHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	    
        tpl := template.Must(template.ParseFiles("index.html"))
        
        u, err := url.Parse(r.URL.String())
        
       
        params := u.Query()
            
        
        page := params.Get("page")
        domain := params.Get("domain")
        
        if page == "" {
        	page = "1"
        }
        
        
        results, err := newsapi.FetchAllNews(page, domain)
        
        
        if err != nil {
        	http.Error(w, err.Error(), http.StatusInternalServerError)
        	return
        }
        
        nextPage, err := strconv.Atoi(page)
        if err != nil {
        	http.Error(w, err.Error(), http.StatusInternalServerError)
        	return
        }
        
        content := &Search{
            Query: "",
        	Domains: news.IndexDomains,
        	NextPage:   nextPage,
        	TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
        	Results:    results,
        }
        
        if ok := !content.IsLastPage(); ok {
        	content.NextPage++
        }
        
        buf := &bytes.Buffer{}
        err = tpl.Execute(buf, content)
        if err != nil {
        	http.Error(w, err.Error(), http.StatusInternalServerError)
        	return
        }
        
        buf.WriteTo(w)
	}
}


func searchHandler(newsapi *news.Client) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
	    
	    tpl := template.Must(template.ParseFiles("search.html"))
	
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := newsapi.FetchByQuery(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			Domains: news.IndexDomains,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:    results,
		}
		
		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}
		
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}


func main() {
    err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    
    fs := http.FileServer(http.Dir("assets"))
    mux := http.NewServeMux()
    
    apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apiKey, 20)

    mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
    mux.HandleFunc("/", indexHandler(newsapi))
    mux.HandleFunc("/search", searchHandler(newsapi))
    http.ListenAndServe(":"+port, mux)
}