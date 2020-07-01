package main

import (
	"fmt"
	"log"
	"math"
	// "math/rand"
	"sort"
	"time"
	// "strconv"
)
var testSet []Adult
var trainSet []Adult
var k int 
func main3() {


	k = 8
	// fmt.Println("total ")
	fmt.Println(len(adults))
	for i := range adults {
		trainSet = append(trainSet, adults[i])
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
	log.Printf("Time %s", elapsed)
	
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
	// log.Printf("Distance len %d", len(distances))

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
			// log.Printf("Si peudo %d",id)
			// log.Printf("len %d",len(distancesAux[id]))

            for f := id; f < len(trainingSet); f += numproc {
				dist := euclidianDistance(testRecord, trainingSet[f])
	
				distancesAux[id] = append(distancesAux[id], distancePair{trainingSet[f], dist})
				if(id == 0){
					// log.Printf("len %d",len(distancesAux[id]))
				}
            }
            end <- true
        }(i)
    }
    for i := 0; i < numproc; i++ {
        <-end
    }
    for i := 0; i < numproc; i++ {
		// log.Printf("len %d",len(distancesAux[i]))
		distances = append(distances,distancesAux[i]...)
    }

	sort.Sort(distances)

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