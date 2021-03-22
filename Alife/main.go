package main

import (
	"Alife/lib"
	"Alife/model"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

//goland:noinspection GoBoolExpressions
func main() {
	fmt.Printf("Running model\n")
	if len(os.Args) != 5 {
		log.Fatal(errors.New("arguements should be: ExperimentName string, WorldDynamics string, Bonds []int, DSImode string"))
	}
	// number of runs
	n := 2

	// experiment name
	directoryName := os.Args[1]

//----------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------
//----------------------------------VARIABLES TESTED IN THE EXPERIMENT--------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------
	// variables testes in the experiment
	worldDynamics := os.Args[2]
	bondedAgents := bonds(os.Args[3])
	DSImode := os.Args[4]

	// for now always 6 agents
	numberOfAgents := 6
	// for now always Neutral, so 0.5 for every agent
	cortisolThresholdCondition := "Neutral"
//----------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

	// create directory for this experiment's data
	name := "data/" + directoryName
	if _, errDir := os.Stat(name); os.IsNotExist(errDir) {
		_ = os.Mkdir(name, 0700)
	}
	// create a txt file with experiment variables
	paramsFile, errFile := os.Create("data/"+directoryName+"/params.txt")
	check(errFile)
	paramsString := fmt.Sprintf("World Dynamics: %v, \nNumber Of Agents: %v, \nBonded Agents: %v, \nDSImode: %v, \nCortistol Threshold: %v, \n",
								worldDynamics, numberOfAgents, bondedAgents, DSImode, cortisolThresholdCondition)
	_, err2 := paramsFile.WriteString(paramsString)
	check(err2)

	for i:=0;i<n;i++ {
		if ((i+1) % 10) == 0 {
			fmt.Printf("Iteration: %v\n",i+1)
		}
		rand.Seed(time.Now().UnixNano())

		iterations := 15000
		filename := directoryName + "/run_" + strconv.Itoa(i+1) + ".csv"

		// world setup
		visualisation := false
		w, h := 99, 99
		visionLength := 20
		// vision angle is both to the right and left so actually twice this amount
		// vision in degrees and must be smaller than 90 (overall smaller than 180)
		visionAngle := 40

		a := lib.NewSimulation()
		grid2D := model.NewWorld(w, h, numberOfAgents, visionLength, visionAngle)
		a.SetWorld(grid2D)

		// channels for communication with UI and Writer
		chGrid := make(chan [][]interface{})
		chVar := make(chan []float64, numberOfAgents)

		//check if correct mode
		errDSI := checkDSImode(DSImode)
		if errDSI != nil {
			log.Fatal(errDSI)
		}

		// initialise agents from 1 to 6
		for i := 1; i < numberOfAgents+1; i++ {
			x, y := randomFloat(float64(w)), randomFloat(float64(h))
			addAgent(x, y, i, i, numberOfAgents, a, grid2D, chVar, false, cortisolThresholdCondition, DSImode)
		}

		// set up bonds between agents
		errBond := initialiseBonds(bondedAgents, numberOfAgents, a)
		if errBond != nil {
			log.Fatal(errBond)
		}

		// pick world settings
		setupWorld(a, grid2D, worldDynamics)

		a.LimitIterations(iterations)

		// reporting function, does something each iteration
		// in this case updates the UI
		if visualisation {
			a.SetReportFunc(func(a *lib.ABM) {
				chGrid <- grid2D.Dump(func(a lib.Agent) int {
					time.Sleep(100 * time.Nanosecond)
					if a == nil {
						return 0
					}
					return a.ID()
				})
			})
		}

		// get current time for datafile name
		//t := time.Now()

		//chan to make sure we finish writing before finishing
		finished := make(chan bool)

		// start writer go routine
		//writer := lib.NewWriter(finished, chVar, t.Format(time.Stamp)+".csv", numberOfAgents)
		writer := lib.NewWriter(finished, chVar, filename, numberOfAgents)
		go writer.Loop()

		go func() {
			a.StartSimulation()
			chVar <- nil
			close(chGrid)
		}()

		if visualisation {
			ui := lib.NewUI(w, h)
			ui.AddGrid(chGrid)
			ui.Loop()
			ui.Stop()
		}

		//wait for writer to finish
		<-finished
		close(chVar)
		close(finished)
	}
}
//______________________________________________________________________________________________________________________
//______________________________________________________________________________________________________________________
//______________________________________________________________________________________________________________________
//______________________________________________________________________________________________________________________
//______________________________________________________________________________________________________________________

func addAgent(x, y float64, id, rank, numOfAgents int, a *lib.ABM, grid2D *model.Grid, ch chan []float64,
				trail bool, CortisolThresholdCondition, DSImode string) {
	cell, err := model.NewAgent(a, id, rank, numOfAgents, x, y, ch, trail, CortisolThresholdCondition, DSImode)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}

func addFood(x, y float64, a *lib.ABM, grid2D *model.Grid) {
	cell, err := model.NewFood(a, x, y)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}

func setupWorld(a *lib.ABM, grid2D *model.Grid, condition string) {
	err := grid2D.SetWorldDynamics(condition)
	if err != nil {
		log.Fatal(err)
	}
	switch condition{
	case "Four":
		addFood(9, 9, a, grid2D)
		addFood(89, 89, a, grid2D)
		addFood(9, 89, a, grid2D)
		addFood(89, 9, a, grid2D)
	case "Seasonal":
		addFood(9, 9, a, grid2D)
		addFood(89, 89, a, grid2D)
		addFood(9, 89, a, grid2D)
		addFood(89, 9, a, grid2D)
	case "Extreme":
		addFood(9, 9, a, grid2D)
		addFood(89, 89, a, grid2D)
		addFood(9, 89, a, grid2D)
		addFood(89, 9, a, grid2D)
	}
}

func initialiseBonds(bondedAgents []int, numberOfAgents int, a *lib.ABM) error {
	for i:= 0; i < len(bondedAgents); i++ {
		if bondedAgents[i] < 1 || bondedAgents[i] > numberOfAgents { return errors.New("invalid agent id") }
		for j:=0; j < len(bondedAgents); j++ {
			if i != j {
				if bondedAgents[i] == bondedAgents[j]{
					return errors.New("agent bond duplicate; agent cannot bond with itself")
				}
			}
		}
	}
	for i:= 0; i < len(bondedAgents); i++ {
		for _,agent := range a.Agents(){
			if agent.ID() == bondedAgents[i] {
				var bonds []int
				for j:=0; j < len(bondedAgents); j++ {
					if i != j {
						bonds = append(bonds, bondedAgents[j])
					}
				}
				agent.(*model.Agent).SetBonds(bonds)
			}
		}
	}
	return nil
}

func bonds(arg string) []int{

	temp := strings.Replace(arg, "[", "", -1)
	temp2 := strings.Replace(temp, "]", "", -1)

	t := strings.Split(temp2, ",")

	if t[0] == "" {
		return nil
	}

	var t2 []int

	for _, i := range t {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		t2 = append(t2, j)
	}
	return t2
}

func checkDSImode(DSImode string) error {
	if DSImode == "Fixed" || DSImode == "Variable" {
		return nil
	}
	return errors.New("DSImode must be one of: Fixed, Variable")
}

// needed to make sure it's never 0
func randomFloat(max float64) float64 {
	var res float64
	for {
		res = rand.Float64()
		if res != 0 {
			break
		}
	}
	return res*max
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}