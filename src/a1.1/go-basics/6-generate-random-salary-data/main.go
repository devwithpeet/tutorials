package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	data := generate()

	for _, v := range data {
		fmt.Println(strings.Join(v, ","))
	}
}

func getPositions() []string {
	return []string{"Software Developer", "Senior Software Developer", "DevOps Engineer", "UI/UX Designer", "Project Manager", "Data Scientist", "Quality Assurance Engineer", "Product Owner", "Frontend Developer", "Backend Developer"}
}

func getGermanNames() []string {
	return []string{"Alex", "Maria", "John", "Lena", "Simon", "Eva", "Felix", "Nina", "Oliver", "Sophie"}
}

func getPolishNames() []string {
	return []string{"Piotr", "Anna", "Marek", "Kasia", "Katarzyna", "Krzysztof", "Michal"}
}

func getGermanCities() []string {
	return []string{"Munich", "Berlin", "Frankfurt", "Hamburg", "Stuttgart", "Cologne", "Dresden", "Leipzig", "Bremen", "Nuremberg"}
}

func getPolishCities() []string {
	return []string{"Warsaw", "Wroclaw", "Krakow", "Gdansk"}
}

func generate() [][]string {
	positions := getPositions()
	germanNames := getGermanNames()
	polishNames := getPolishNames()
	germanCities := getGermanCities()
	polishCities := getPolishCities()

	var (
		name, city, country string
		salary              int
		data                [][]string
	)
	for i := 0; i < getCount(); i++ {
		if rand.Intn(10) == 8 {
			name = germanNames[rand.Intn(len(germanNames))]
			city = germanCities[rand.Intn(len(germanCities))]
			country = "Germany"
			salary = rand.Intn(40000) + 40000 // Random salary between 50000 and 100000
		} else {
			name = polishNames[rand.Intn(len(polishNames))]
			city = polishCities[rand.Intn(len(polishCities))]
			country = "Poland"
			salary = rand.Intn(50000) + 50000 // Random salary between 50000 and 100000
		}

		position := positions[rand.Intn(len(positions))]

		data = append(data, []string{name, position, city, country, strconv.Itoa(salary)})
	}

	return data
}

func getCount() int {
	if len(os.Args) < 2 {
		return 2000
	}

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	return num
}
