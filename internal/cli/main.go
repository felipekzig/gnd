package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gnd/internal/domain"
	"github.com/spf13/cobra"
)

var All bool
var gndRootCmd = &cobra.Command{Use: "gnd"}

func Execute(ts domain.TaskService) {

	gndRootCmd.AddCommand(&cobra.Command{
		Use:       "add",
		Short:     "Add a new task to the list",
		ValidArgs: []string{"task", "priority"},
		Args:      cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var p int
			if len(args) > 1 {
				p, _ = strconv.Atoi(args[1])
			}

			ts.Add(args[0], uint16(p))
		},
	})

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all opened tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := ts.List(All)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Printf("%4s | %-50s | Priority", "ID", "Task")
			if All {
				fmt.Printf(" | Done")
			}
			fmt.Print("\n")
			for _, t := range tasks {
				fmt.Printf("%4X | %-50s | %8d", t.ID, t.Task, t.Priority)
				if All {
					fmt.Printf(" | %4t", t.Done)
				}
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
