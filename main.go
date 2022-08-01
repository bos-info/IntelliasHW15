package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout     = "15:04:05"
	countOfRecords = 3
	fName          = "data.json"
	//Errors Description
	emptyDeptSt    = "empty departure station"
	badDeptSt      = "bad departure station input"
	emptyArrSt     = "empty arrival station"
	badArrSt       = "bad arrival station input"
	unsCriteria    = "unsupported criteria"
	invalidJSON    = "JSON is invalid"
	unmarshalError = "unmarshal error"
	parseArrTime   = "bad parse arrival time"
	parseDepTime   = "bad parse departure time"
)

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

func main() {
	var departureStation, arrivalStation, criteria string
	fmt.Println("----------- Інформаційна система УкрЗалізниця 3000 -----------")
	fmt.Println("Оберіть станцію відправлення")
	_, err := fmt.Scanf("%v", &departureStation)
	if err != nil {
		log.Fatal(emptyDeptSt)
	}
	departureStation = strings.Trim(departureStation, "\"")
	fmt.Println("Оберіть станцію прибуття")
	_, err = fmt.Scanf("%v", &arrivalStation)
	if err != nil {
		log.Fatal(emptyArrSt)
	}
	arrivalStation = strings.Trim(arrivalStation, "\"")
	fmt.Println("Оберіть критерій, по якому сортувати результат")
	_, err = fmt.Scanf("%v", &criteria)
	if err != nil {
		log.Fatal(unsCriteria)
	}
	criteria = strings.Trim(criteria, "\"")
	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Println("Х-Х-Х-Х-Х-Х-Х-Х-Х-Х-Х- ПОМИЛКА -Х-Х-Х-Х-Х-Х-Х-Х-Х-Х-Х")
		log.Fatal(err)
	}
	fmt.Println("---------------------- Результат пошуку ----------------------")
	if result == nil {
		fmt.Println("Відсутні потяги між вказанами станціями.")
	}
	for i := 0; i < len(result); i++ {
		fmt.Printf("TrainID: %d DepartureStationID: %d ArrivalStationID: %d  Price: %0.2f ArrivalTime: %v DepartureTime: %v\n",
			result[i].TrainID, result[i].DepartureStationID, result[i].ArrivalStationID, result[i].Price, result[i].ArrivalTime, result[i].DepartureTime)
	}
	fmt.Println("--------------------------------------------------------------")
}

//FindTrains шукаємо потяги що задовольняють введеним умовам користувача
func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	//читаємо json та записуємо вичитані дані в
	trainsSlice := readFile()
	var depStationInt, arrStationInt int
	//Перевіряємо валідність станції відправлення
	if len(departureStation) == 0 {
		return nil, errors.New(emptyDeptSt)
	}
	depStationInt, errDeptSt := strconv.Atoi(departureStation)
	if errDeptSt != nil || depStationInt < 0 {
		return nil, errors.New(badDeptSt)
	}
	//Перевіряємо валідність станції прибуття
	if len(arrivalStation) == 0 {
		return nil, errors.New(emptyArrSt)
	}
	arrStationInt, errArrSt := strconv.Atoi(arrivalStation)
	if errArrSt != nil || arrStationInt < 0 {
		return nil, errors.New(badArrSt)
	}
	var result Trains
	//Додаємо всі потяги, що задовольняють умові курсування між станціями
	for _, v := range trainsSlice {
		if v.DepartureStationID == depStationInt && v.ArrivalStationID == arrStationInt {
			result = append(result, v)
		}
	}
	//Перевіряємо правильність введення критерію сортування та сортуємо по ньому.
	switch criteria {
	case "price":
		sort.Slice(result, func(i, j int) bool { return result[i].Price < result[j].Price })
	case "arrival-time":
		sort.Slice(result, func(i, j int) bool { return result[i].ArrivalTime.Before(result[j].ArrivalTime) })
	case "departure-time":
		sort.Slice(result, func(i, j int) bool { return result[i].DepartureTime.Before(result[j].DepartureTime) })
	default:
		//якщо введено невалідний критерій повертаємо помилку
		return nil, errors.New(unsCriteria)
	}
	//якщо результатів більше ніж вимагається за умовою задачі
	if len(result) > countOfRecords {
		return result[:countOfRecords], nil // маєте повернути правильні значення
	}
	//якщо по вказаним критеріям потягів не знайдено повертаємо nil
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

//readFile вичитує дані з файлу data.json, парсить їх, та повертає підготовлені дані для пошуку.
func readFile() Trains {
	file, err := os.OpenFile(fName, os.O_RDONLY, os.FileMode(0600))
	if os.IsNotExist(err) {
		log.Fatal("File doesn't exist", err)
	}
	// оброюблюємо закриття файлу
	defer closeFile(file)
	byteValue, _ := ioutil.ReadAll(file)
	var uz Trains
	//Додаткова перевірка валідності структури JSON(не вимагалась за умовою задачі)
	isValid := json.Valid(byteValue)
	if !isValid {
		log.Fatal(invalidJSON)
	}
	//Обробка помилки при обробці даних з файлу
	err = json.Unmarshal(byteValue, &uz)
	if err != nil {
		log.Fatal(err)
	}
	return uz
}

//closeFile функція закриває файл
func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

//UnmarshalJSON описуємо метод анмаршалінгу JSON для Train
func (tr *Train) UnmarshalJSON(j []byte) error {
	var rawData map[string]any
	err := json.Unmarshal(j, &rawData)
	if err != nil {
		return errors.New(unmarshalError)
	}
	for k, v := range rawData {
		if k == "trainId" {
			tr.TrainID = int(v.(float64))
		}
		if k == "departureStationId" {
			tr.DepartureStationID = int(v.(float64))
		}
		if k == "arrivalStationId" {
			tr.ArrivalStationID = int(v.(float64))
		}
		if k == "price" {
			tr.Price = float32(v.(float64))
		}
		if k == "arrivalTime" {
			t, err := time.Parse(timeLayout, v.(string))
			if err != nil {
				return errors.New(parseArrTime)
			}
			tr.ArrivalTime = t
		}
		if k == "departureTime" {
			t, err := time.Parse(timeLayout, v.(string))
			if err != nil {
				return errors.New(parseDepTime)
			}
			tr.DepartureTime = t
		}
	}
	return nil
}
