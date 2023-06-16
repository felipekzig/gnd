package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gnd/internal/domain"
	"github.com/spf13/cobra"
)

var All bool
var gndRootCmd = &cobra.Command{Use: "gnd"}

func defTimeValue(v time.Time) string {
	if v.IsZero() {
		return "-"
	}
	return v.Local().Format(time.DateOnly)
}

func Execute(ts domain.TaskService) {

	var dueDateFlag *string
	addCmd := &cobra.Command{
		Use:       "add",
		Short:     "Add a new task to the list",
		ValidArgs: []string{"task"},
		Aliases:   []string{"a"},
		Run: func(cmd *cobra.Command, args []string) {
			t := args[0]
			if len(args) > 1 {
				t = strings.Join(args, " ")
			}

			if dueDateFlag == nil {
				ts.Add(t, time.Time{})
			}

			dd, err := dateparse.ParseAny(*dueDateFlag)
			if err != nil {
				fmt.Println("Due date flag does not seem to be a valid date")
				return
			}

			if dd.Before(time.Now()) {
				fmt.Println("Due date can't be in the past")
				return
			}

			ts.Add(t, dd)
		},
	}
	dueDateFlag = addCmd.Flags().StringP("due-date", "d", "", "The due date for this task. Accepted formats: MM/DD/YYYY, MM-DD-YYYY")
	gndRootCmd.AddCommand(addCmd)

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ll"},
		Short:   "List all opened tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := ts.List(All)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Print("   |")
			fmt.Printf("%4s | %-50s | %12s", "ID", "Task", "Due Date")

			fmt.Print("\n")
			fmt.Printf("%s\n", strings.Repeat("-", 80))
			for _, t := range tasks {
				if t.Done {
					fmt.Printf(" \u2713 |")

				} else {
					fmt.Printf("   |")

				}
				fmt.Printf("%4d | %-50s | %12s", t.ID, t.Task, defTimeValue(t.DueDate))
				fmt.Print("\n")
			}
		},
	}
	listCmd.Flags().BoolVarP(&All, "all", "a", false, "Include all tasks")
	gndRootCmd.AddCommand(listCmd)

	gndRootCmd.AddCommand(&cobra.Command{
		Use:       "ping",
		Short:     "Ping a task to increase its relevance",
		ValidArgs: []string{"task_id"},
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("The task ID <%s> is not an integer!", args[0])
			}

			var t domain.Task
			t, err = ts.Ping(uint(id))
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}

			fmt.Printf("The task <%s> had its relevance increased to %d\n", t.Task, t.Pings)
		},
	})

	gndRootCmd.AddCommand(&cobra.Command{
		Use:       "done",
		Short:     "Mark task as completed",
		ValidArgs: []string{"task_id"},
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("The task ID <%s> is not an integer!\n", args[0])
			}

			var t domain.Task
			t, err = ts.Complete(uint(id))
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}

			fmt.Printf("Kudos! The task <%s> is completed!\n", t.Task)
		},
	})

	if err := gndRootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
