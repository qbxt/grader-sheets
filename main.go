package main

import (
	"./functions"
	"./structures"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	gradersJson := ""
	configJson := ""
	submissionsZip := ""
	help := false

	flag.StringVar(&gradersJson, "g", "", "Location of graders.json.")
	flag.StringVar(&submissionsZip, "s", "", "Location of submissions.zip.")
	flag.StringVar(&configJson, "c", "", "Location of config.json.")
	flag.BoolVar(&help, "h", false, "Print help and exit")

	flag.Usage = func() {
		fmt.Printf("Usage: ./grader-sheets -c path/to/config.json -g path/to/graders.json -s path/to/submissions.zip\n")
	}

	flag.Parse()

	if help {
		fmt.Printf("Usage: ./grader-sheets -c path/to/config.json -g path/to/graders.json -s path/to/submissions.zip\n")
		return
	}

	if gradersJson == "" || submissionsZip == "" || configJson == "" {
		fmt.Printf("Usage: ./grader-sheets -c path/to/config.json -g path/to/graders.json -s path/to/submissions.zip\n")
		return
	}

	config, err := functions.UnmarshalConfig(configJson)
	if err != nil {
		panic(err)
	}
	config.SubmissionsZip = submissionsZip
	config.GradersJson = gradersJson

	graders, err := functions.UnmarshalGraders(gradersJson)
	if err != nil {
		panic(err)
	}

	graderAssignmentCount := functions.CountGraderAssignments(graders.G)
	if graderAssignmentCount == 0 {
		panic("0 grader assignments counted")
	}

	submissionCount, err := functions.CountSubmissions(config)
	if err != nil {
		panic(err)
	} else if submissionCount == 0 {
		panic("0 submissions detected")
	}

	fmt.Printf("Running all scripts in ./scripts\n")
	if err := functions.RunScripts(config.SubmissionsDirectory); err != nil {
		panic(err)
	}

	if submissionCount == graderAssignmentCount {
		fmt.Printf("Grader Assignment Count == Submission Count (%d == %d), continuing\n", graderAssignmentCount, submissionCount)
	} else if submissionCount > graderAssignmentCount {
		fmt.Printf("Warning: Grader Assignment Count less than Submission Count (%d < %d)\n", graderAssignmentCount, submissionCount)
		fmt.Printf("Some graders will (randomly) get more assignments than they requested. To prevent random assignment, increase the 'grade' field for graders in graders.json\n")
	} else if submissionCount < graderAssignmentCount {
		fmt.Printf("Warning: Grader Assignment Count greater than Submission Count (%d > %d)\n", graderAssignmentCount, submissionCount)
		fmt.Printf("Some graders will (randomly) get fewer assignments than they requested.\n")
	}

	graderList, err := functions.MakeGraderList(config, graders.G)

	printGraderEmails(graders)
	printGraderList(graderList)
	fmt.Printf("\nSuccessfully wrote grader info and grader sheet. Run ./print.sh to view the results.\n")
}

func printGraderList(input map[string]*[]string) {
	file, err := os.Create("./graderlist.txt")
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(file)

	_, err = fmt.Fprintf(w, "Name (LastFirst)|Student\n")
	if err != nil {
		panic(err)
	}
	for grader, gradees := range input {
		for _, g := range *gradees {
			_, err := fmt.Fprintf(w, "%s|%s\n", grader, strings.Split(g, "_")[0])
			if err != nil {
				panic(err)
			}
		}
	}

	_ = w.Flush()
	_ = file.Close()
}

func printGraderEmails(input *structures.Graders) {
	file, err := os.Create("./graderemails.txt")
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(file)

	_, err = fmt.Fprintf(w, "Name (LastFirst)|Email\n")
	if err != nil {
		panic(err)
	}
	for _, g := range *input.G {
		_, err := fmt.Fprintf(w, "%s%s|%s\n", g.Last, g.First, g.Email)
		if err != nil {
			panic(err)
		}
	}

	_ = w.Flush()
	_ = file.Close()
}
