Package main



import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponsePayload struct {
	Sortedarrays [][]int `json:"sorted_arrays"`
	Time       int64   `json:"time_ns"`
}

func sortSequential(input [][]int) [][]int {
	result := make([][]int, len(input))

	for i, subarray := range input {
		sortedSubarray := make([]int, len(subarray))
		copy(sortedSubarray, subarray)
		sort.Ints(sortedSubarray)
		result[i] = sortedSubarray
	}

	return result
}

func sortConcurrent(input [][]int) [][]int {
	result := make([][]int, len(input))
	var wg sync.WaitGroup
	wg.Add(len(input))

	for i, subarray := range input {
		go func(i int, subarray []int) {
			defer wg.Done()
			sortedSubarray := make([]int, len(subarray))
			copy(sortedSubarray, subarray)
			sort.Ints(sortedSubarray)
			result[i] = sortedSubarray
		}(i, subArray)
	}

	wg.Wait()
	return result
}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortSequential(payload.ToSort)
	timeTaken := time.Since(startTime).Nanoseconds()

	response := ResponsePayload{
		Sortedarrays: sortedarrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortConcurrent(payload.ToSort)
	timeTaken := time.Since(startTime).Nanoseconds()

	response := ResponsePayload{
		Sortedarrays: sortedarrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/process-single", processSingleHandler)
	http.HandleFunc("/process-concurrent", processConcurrentHandler)

	fmt.Println("Server listening on :4545...")
	err := http.ListenAndServe(":4545", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
