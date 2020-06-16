package main

import (
    // "bufio"
    "fmt"
    "strings"
    "log"
    "os"
    "encoding/csv"
    "strconv"
    // "net"
    
	"encoding/json"
	// "math/rand"
	"net/http"
    // "strconv"

    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
	"github.com/bitly/go-simplejson"
	
	"math"
	// "math/rand"
	"sort"
	"time"
)
var adults []Adult
type Adult struct {
	Id     int  `json:"id"`
	Fever     int  `json:"fever"`
    Tiredness int `json:"tiredness"`
    Dry_Cough   int `json:"dry_Cough"`
    Difficulty_in_Breathing int    `json:"difficulty_in_Breathing"`
    None_Sympton int    `json:"none_Sympton"`
    Sore_Throat int    `json:"sore_Throat"`
    Pains int  `json:"pains"`
    Nasal_Congestion int   `json:"nasal_Congestion"`
    Runny_Nose int `json:"runny_Nose"`
    Diarrhea int `json:"diarrhea"`
    None_Experiencing int  `json:"none_Experiencing"`
    Age float64  `json:"age"`
    Gender float64   `json:"gender"`
    Severity int   `json:"severity"`
    Contact float64   `json:"contact"`
    Country string   `json:"country"`
}

var testSet []Adult
var trainSet []Adult
var k int 
func main3() {


	k = 10
	fmt.Println("total ")
	fmt.Println(len(adults))
	for i := range adults {
		// if i < 2 {
		// 	testSet = append(testSet, adults[i])
		// } else {
		trainSet = append(trainSet, adults[i])
		// }
	}

	var predictions []int
	fmt.Println("test lenght")
	fmt.Println(len(testSet))
	start := time.Now()
	for x := range testSet {
		// neighbors := getNeighbors(trainSet, testSet[x], k)
		// result := getResponse(neighbors)
		result := testCase(trainSet,testSet[x],k)
		predictions = append(predictions, result[0].key)
		fmt.Printf("Predicted: %d, Actual: %d\n", result[0].key, testSet[x].Severity)
	}
    elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
	
	accuracy := getAccuracy(testSet, predictions)
	fmt.Printf("Accuracy con una diferencia +-1 nivel: %f%s\n", accuracy, "%")
}

func testCase(trainSetA []Adult,testSetObject Adult, k int) sortedClassVotes {
	fmt.Println(testSetObject)
	neighbors := getNeighbors2(trainSetA, testSetObject, k)
	result := getResponse(neighbors)
	return result
}


func getAccuracy(testSet []Adult, predictions []int) float64 {
	correct := 0

	for x := range testSet {
		difference := math.Abs(float64(testSet[x].Severity - predictions[x]))
		if testSet[x].Severity == predictions[x] {
			correct += 100
		} else if difference <= 33 {
			correct += 100
		} else if difference> 33  && difference<=66 {
			correct += 33
		}
	}

	return (float64(correct) / float64(len(testSet)*100)) * 100.00
}

type classVote struct {
	key   int
	value int
}

type sortedClassVotes []classVote

func (scv sortedClassVotes) Len() int           { return len(scv) }
func (scv sortedClassVotes) Less(i, j int) bool { return scv[i].value < scv[j].value }
func (scv sortedClassVotes) Swap(i, j int)      { scv[i], scv[j] = scv[j], scv[i] }

func getResponse(neighbors []Adult) sortedClassVotes {
	classVotes := make(map[int]int)

	for x := range neighbors {
		response := neighbors[x].Severity
		if contains(classVotes, response) {
			classVotes[response] += 1
		} else {
			classVotes[response] = 1
		}
	}

	scv := make(sortedClassVotes, len(classVotes))
	i := 0
	for k, v := range classVotes {
		scv[i] = classVote{k, v}
		i++
	}

	sort.Sort(sort.Reverse(scv))
	log.Printf("Distance len %d  %d", len(scv), scv)
	return scv
}

type distancePair struct {
	record   Adult
	distance float64
}

type distancePairs []distancePair

func (slice distancePairs) Len() int           { return len(slice) }
func (slice distancePairs) Less(i, j int) bool { return slice[i].distance < slice[j].distance }
func (slice distancePairs) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

func getNeighbors(trainingSet []Adult, testRecord Adult, k int) []Adult {
	var distances distancePairs
	for i := range trainingSet {
		dist := euclidianDistance(testRecord, trainingSet[i])
		distances = append(distances, distancePair{trainingSet[i], dist})
	}

	sort.Sort(distances)
	log.Printf("Distance len %d", len(distances))

	var neighbors []Adult

	for x := 0; x < k; x++ {
		neighbors = append(neighbors, distances[x].record)
	}

	return neighbors
}

func getNeighbors2(trainingSet []Adult, testRecord Adult, k int) []Adult {
	var distances distancePairs
	distancesAux  :=  [][]distancePair{[]distancePair{}}
	numproc := 10
	// Process 0 esta con len 0 siempre?
    end := make(chan bool)
    for i := 0; i < numproc; i++ {
		distancesAux = append(distancesAux,[]distancePair{})
        go func(id int) {
            for f := id; f < len(trainingSet); f += numproc {
				dist := euclidianDistance(testRecord, trainingSet[f])
				distancesAux[id] = append(distancesAux[id], distancePair{trainingSet[f], dist})
            }
            end <- true
        }(i)
	}
	for i := 0; i < numproc; i++ {
		<-end
	}
	for i := 0; i < numproc; i++ {
	log.Printf("Len por process %d %d", i,len(distancesAux[i]))
	distances = append(distances,distancesAux[i]...)
	}

	sort.Sort(distances)
	log.Printf("Distance len %d", len(distances))

	var neighbors []Adult

	for x := 0; x < k; x++ {
		neighbors = append(neighbors, distances[x].record)
	}

	return neighbors
}


func concurrent(){
	
		// 	go func(instanceOne Adult, instanceTwo Adult){
	// 		var distance float64

	// distance += math.Pow(float64((instanceOne.Fnlwgt - instanceTwo.Fnlwgt)), 2)
	// distance += math.Pow(float64((instanceOne.EducationNum - instanceTwo.EducationNum)), 2)
	// distance += math.Pow(float64((instanceOne.CapitalGain - instanceTwo.CapitalGain)), 2)
	// distance += math.Pow(float64((instanceOne.CapitalLoss - instanceTwo.CapitalLoss)), 2)
	// distance += math.Pow(float64((instanceOne.Hours - instanceTwo.Hours)), 2)
	// distance += math.Pow(float64((instanceOne.Sex - instanceTwo.Hours)), 2)
	
	// dist := math.Sqrt(distance)

	// distances = append(distances, distancePair{instanceTwo, dist})

	// 	}(testRecord,trainingSet[i])
}

func euclidianDistance(instanceOne Adult, instanceTwo Adult) float64 {
	var distance float64

	distance += math.Pow(float64((instanceOne.Fever - instanceTwo.Fever)), 2)
	distance += math.Pow(float64((instanceOne.Tiredness - instanceTwo.Tiredness)), 2)
	distance += math.Pow(float64((instanceOne.Dry_Cough - instanceTwo.Dry_Cough)), 2)
	distance += math.Pow(float64((instanceOne.Difficulty_in_Breathing - instanceTwo.Difficulty_in_Breathing)), 2)
	distance += math.Pow(float64((instanceOne.None_Sympton - instanceTwo.None_Sympton)), 2)
	distance += math.Pow(float64((instanceOne.Sore_Throat - instanceTwo.Sore_Throat)), 2)
	distance += math.Pow(float64((instanceOne.Pains - instanceTwo.Pains)), 2)
	distance += math.Pow(float64((instanceOne.Nasal_Congestion - instanceTwo.Nasal_Congestion)), 2)
	distance += math.Pow(float64((instanceOne.Runny_Nose - instanceTwo.Runny_Nose)), 2)
	distance += math.Pow(float64((instanceOne.Diarrhea - instanceTwo.Diarrhea)), 2)
	distance += math.Pow(float64((instanceOne.None_Experiencing - instanceTwo.None_Experiencing)), 2)
	distance += math.Pow(float64((instanceOne.Age - instanceTwo.Age)), 2)
	distance += math.Pow(float64((instanceOne.Gender - instanceTwo.Gender)), 2)
	distance += math.Pow(float64((instanceOne.Contact - instanceTwo.Contact)), 2)
	// distance += math.Pow(float64((instanceOne.Severity - instanceTwo.Severity)), 2)

	return math.Sqrt(distance)
}


func errHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func contains(votesMap map[int]int, name int) bool {
	for s, _ := range votesMap {
		if s == name {
			return true
		}
	}

	return false
}

func lineToStruc(lines [][]string){
		// Loop through lines & turn into object
		i := 0
        for _, line := range lines {
			Fever,_ := strconv.Atoi(strings.TrimSpace(line[0]))
			Tiredness,_ := strconv.Atoi(strings.TrimSpace(line[1]))
			Dry_Cough,_ := strconv.Atoi(strings.TrimSpace(line[2]))
			Difficulty_in_Breathing,_ := strconv.Atoi(strings.TrimSpace(line[3]))
			None_Sympton,_ := strconv.Atoi(strings.TrimSpace(line[4]))
			Sore_Throat,_ := strconv.Atoi(strings.TrimSpace(line[5]))
			Pains,_ := strconv.Atoi(strings.TrimSpace(line[6]))
			Nasal_Congestion,_ := strconv.Atoi(strings.TrimSpace(line[7]))
			Runny_Nose,_ := strconv.Atoi(strings.TrimSpace(line[8]))
			Diarrhea,_ := strconv.Atoi(strings.TrimSpace(line[9]))
			None_Experiencing,_ := strconv.Atoi(strings.TrimSpace(line[10]))
			Age_0_9,_ := strconv.Atoi(strings.TrimSpace(line[11]))
			Age_10_19,_ := strconv.Atoi(strings.TrimSpace(line[12]))
			Age_20_24,_ := strconv.Atoi(strings.TrimSpace(line[13]))
			Age_25_59,_ := strconv.Atoi(strings.TrimSpace(line[14]))
			Age_60,_ := strconv.Atoi(strings.TrimSpace(line[15]))
			Gender_Female,_ := strconv.Atoi(strings.TrimSpace(line[16]))
			Gender_Male,_ := strconv.Atoi(strings.TrimSpace(line[17]))
			Gender_Transgender,_ := strconv.Atoi(strings.TrimSpace(line[18]))
			Severity_Mild,_ := strconv.Atoi(strings.TrimSpace(line[19]))
			Severity_Moderate,_ := strconv.Atoi(strings.TrimSpace(line[20]))
			Severity_None,_ := strconv.Atoi(strings.TrimSpace(line[21]))
			Severity_Severe,_ := strconv.Atoi(strings.TrimSpace(line[22]))
			Contact_Dont_Know,_ := strconv.Atoi(strings.TrimSpace(line[23]))
			Contact_No,_ := strconv.Atoi(strings.TrimSpace(line[24]))
			Contact_Yes,_ := strconv.Atoi(strings.TrimSpace(line[25]))
			Country := strings.TrimSpace(line[26])

			Age:= 0.0
			if(Age_0_9 == 1){
				Age = 0
			}
			if(Age_10_19 == 1){
				Age = 0.25
			}
			if(Age_20_24 == 1){
				Age = 0.5
			}
			if(Age_25_59 == 1){
				Age = 0.75
			}
			if(Age_60 == 1){
				Age = 1
			}
			
			

			Gender := 0.0
			if(Gender_Transgender == 1){
				Gender = 0
			}
			if(Gender_Female == 1){
				Gender = 0.5
			}
			if(Gender_Male == 1){
				Gender = 1
			}

		

			// TESTING
			Severity:= 0
			if(Severity_Mild == 1){
				Severity = 33
			}
			if(Severity_Moderate == 1){
				Severity = 66
			}
			if(Severity_None == 1){
				Severity = 0
			}
			if(Severity_Severe == 1){
				Severity = 100
			}

		

			Contact:= 0.0
			if(Contact_Dont_Know == 1){
				Contact = 0.5
			}
			if(Contact_No == 1){
				Contact = 0
			}
			if(Contact_Yes == 1){
				Contact = 1
			}

			
            adults = append(adults,Adult{
				Id:i,
				Fever:Fever,
				Tiredness:Tiredness,
				Dry_Cough:Dry_Cough,
				Difficulty_in_Breathing:Difficulty_in_Breathing,
				None_Sympton:None_Sympton,
				Sore_Throat:Sore_Throat,
				Pains:Pains,
				Nasal_Congestion:Nasal_Congestion,
				Runny_Nose:Runny_Nose,
				Diarrhea:Diarrhea,
				None_Experiencing:None_Experiencing,
				Age: Age,
				Gender: Gender,
				Severity:Severity,
				Contact:Contact,
				Country:Country,
			})
			
			i++
        }
}

func readFile(filePath string) ([][]string, error) {

 // Open CSV file
 f, err := os.Open(filePath)
 if err != nil {
     return [][]string{}, err
 }
 defer f.Close()

 // Read File into a Variable
 lines, err := csv.NewReader(f).ReadAll()
 if err != nil {
     return [][]string{}, err
 }

 return lines, nil
}

// Adult struct (Model)



// Get all adults
func getAdults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adults)
}

// Get single adult
func getAdult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through adults and find one with the id from the params
	for _, item := range adults {
        id,_ := strconv.Atoi(params["id"])
		if item.Id == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Adult{})
}

// Get single adult
func getCategory(w http.ResponseWriter, r *http.Request) {
    // w.Header().Set("Content-Type", "text/html; charset=utf-8")
    // w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
 

	w.Header().Set("Content-Type", "application/json")
    // params := mux.Vars(r) // Gets params


    var adult Adult
    _ = json.NewDecoder(r.Body).Decode(&adult)
    
    k := 8
    fmt.Println(k)
    result := testCase(adults,adult,k)
    fmt.Printf("Predicted: %d, Actual: %d\n", result[0].key, adult.Severity)
    
	json := simplejson.New()
	json.Set("knn", result[0].key)
	json.Set("actual", adult.Severity)
    json.Set("predicted", result[0].key == adult.Severity)
    if(adult.Severity >= 0){
		fmt.Println("added?")
		adults = append(adults, adult)
	}
    payload, _ := json.MarshalJSON()
    w.Write(payload)
}


// Add new adult
func createAdult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var adult Adult
	_ = json.NewDecoder(r.Body).Decode(&adult)
	adults = append(adults, adult)
	json.NewEncoder(w).Encode(adult)
}


func main() {
    lines, err := readFile("data/Cleaned-Data.csv")
    if err != nil {
        panic(err)
    }
    fmt.Println("Leyo archivos")
    lineToStruc(lines)
    fmt.Println("Parseo Archivos")

	r := mux.NewRouter()

	r.HandleFunc("/adults", getAdults).Methods("GET")
	r.HandleFunc("/adults/{id}", getAdult).Methods("GET")
	r.HandleFunc("/adults", createAdult).Methods("POST")
    r.HandleFunc("/knn", getCategory).Methods("POST")
    
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	//router.HandleFunc("/", RootEndpointGET).Methods("GET")
	//router.HandleFunc("/", RootEndpointPOST).Methods("POST")

    // Start server
    port := ":8000"
    fmt.Println("Escuchando en " + port )
    main3()
    log.Fatal(http.ListenAndServe(port, handlers.CORS(headers, methods, origins)(r)))

}
