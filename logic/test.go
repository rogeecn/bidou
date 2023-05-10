package logic

import "log"

func Fetch(url string) {
	resp, err := client.R().Get(url)
	if err != nil {
		log.Fatal("123", err)
	}

	if resp.IsErrorState() {
		log.Fatal(resp.ErrorResult().(error))
	}

	log.Println(resp.String())
}
