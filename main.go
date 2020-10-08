package main

func main() {
	h := newHub()
	server := getApplicationServer(*h)
	server.ListenAndServe()
}
