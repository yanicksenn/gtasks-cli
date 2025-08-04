package gtasks

import (
	"fmt"
	"strings"
	"testing"

	"google.golang.org/api/tasks/v1"
)

func TestTasksLifecycle(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// 1. Create a task list to work with
	var taskListID string
	output := CaptureOutput(t, func() {
		opts := CreateTaskListOptions{Title: "Shopping List"}
		list, err := client.CreateTaskList(opts)
		if err != nil {
			t.Fatalf("CreateTaskList failed: %v", err)
		}
		taskListID = list.Id
	})

	// 2. Initial list of tasks should be empty
	output = CaptureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID}
		tasks, err := client.ListTasks(opts)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
		printTasks(tasks)
	})
	if !strings.Contains(output, "No tasks found.") {
		t.Errorf("expected empty list, got '%s'", output)
	}

	// 3. Create a task
	var taskID string
	output = CaptureOutput(t, func() {
		opts := CreateTaskOptions{TaskListID: taskListID, Title: "Buy Milk"}
		task, err := client.CreateTask(opts)
		if err != nil {
			t.Fatalf("CreateTask failed: %v", err)
		}
		taskID = task.Id
	})

	// 4. List should now contain the new task
	output = CaptureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID}
		tasks, err := client.ListTasks(opts)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
		printTasks(tasks)
	})
	if !strings.Contains(output, "[ ] Buy Milk") {
		t.Errorf("expected list to contain '[ ] Buy Milk', got '%s'", output)
	}

	// 5. Get the task by ID
	output = CaptureOutput(t, func() {
		opts := GetTaskOptions{TaskListID: taskListID, TaskID: taskID}
		task, err := client.GetTask(opts)
		if err != nil {
			t.Fatalf("GetTask failed: %v", err)
		}
		printTask(task)
	})
	if !strings.Contains(output, "Title:   Buy Milk") {
		t.Errorf("expected get to show 'Title:   Buy Milk', got '%s'", output)
	}

	// 6. Update the task
	output = CaptureOutput(t, func() {
		opts := UpdateTaskOptions{TaskListID: taskListID, TaskID: taskID, Title: "Buy Almond Milk"}
		_, err := client.UpdateTask(opts)
		if err != nil {
			t.Fatalf("UpdateTask failed: %v", err)
		}
	})

	// 7. Complete the task
	output = CaptureOutput(t, func() {
		opts := CompleteTaskOptions{TaskListID: taskListID, TaskID: taskID}
		_, err := client.CompleteTask(opts)
		if err != nil {
			t.Fatalf("CompleteTask failed: %v", err)
		}
	})

	// 8. List should show the task as completed
	output = CaptureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID, ShowCompleted: true}
		tasks, err := client.ListTasks(opts)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
		printTasks(tasks)
	})
	if !strings.Contains(output, "[x] Buy Almond Milk") {
		t.Errorf("expected list to contain '[x] Buy Almond Milk', got '%s'", output)
	}

	// 9. Uncomplete the task
	output = CaptureOutput(t, func() {
		opts := UncompleteTaskOptions{TaskListID: taskListID, TaskID: taskID}
		_, err := client.UncompleteTask(opts)
		if err != nil {
			t.Fatalf("UncompleteTask failed: %v", err)
		}
	})

	// 10. List should show the task as not completed
	output = CaptureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID, ShowCompleted: true}
		tasks, err := client.ListTasks(opts)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
		printTasks(tasks)
	})
	if !strings.Contains(output, "[ ] Buy Almond Milk") {
		t.Errorf("expected list to contain '[ ] Buy Almond Milk', got '%s'", output)
	}

	// 11. Delete the task
	output = CaptureOutput(t, func() {
		opts := DeleteTaskOptions{TaskListID: taskListID, TaskID: taskID}
		err := client.DeleteTask(opts)
		if err != nil {
			t.Fatalf("DeleteTask failed: %v", err)
		}
	})

	// 12. Final list of tasks should be empty again
	output = CaptureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID}
		tasks, err := client.ListTasks(opts)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
		printTasks(tasks)
	})
	if !strings.Contains(output, "No tasks found.") {
		t.Errorf("expected empty list, got '%s'", output)
	}
}

func TestTaskPrint(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// 1. Create a task list to work with
	var taskListID string
	CaptureOutput(t, func() {
		opts := CreateTaskListOptions{Title: "Shopping List"}
		list, err := client.CreateTaskList(opts)
		if err != nil {
			t.Fatalf("CreateTaskList failed: %v", err)
		}
		taskListID = list.Id
	})

	// 2. Create a task
	var taskID string
	CaptureOutput(t, func() {
		opts := CreateTaskOptions{TaskListID: taskListID, Title: "Buy Milk"}
		task, err := client.CreateTask(opts)
		if err != nil {
			t.Fatalf("CreateTask failed: %v", err)
		}
		taskID = task.Id
	})

	// 3. Print the title
	output := CaptureOutput(t, func() {
		task, err := client.GetTask(GetTaskOptions{TaskListID: taskListID, TaskID: taskID})
		if err != nil {
			t.Fatalf("GetTask failed: %v", err)
		}
		if err := PrintTaskProperty(task, "title", false); err != nil {
			t.Fatalf("PrintTaskProperty failed: %v", err)
		}
	})
	if !strings.Contains(output, "Buy Milk") {
		t.Errorf("expected output to contain 'Buy Milk', got '%s'", output)
	}
}

func printTasks(tasks *tasks.Tasks) {
	if len(tasks.Items) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("Tasks:")
	for _, item := range tasks.Items {
		status := " "
		if item.Status == "completed" {
			status = "x"
		}
		fmt.Printf("[%s] %s (%s)\n", status, item.Title, item.Id)
	}
}

func printTask(task *tasks.Task) {
	fmt.Printf("ID:      %s\n", task.Id)
	fmt.Printf("Title:   %s\n", task.Title)
	fmt.Printf("Status:  %s\n", task.Status)
	fmt.Printf("Notes:   %s\n", task.Notes)
	fmt.Printf("Due:     %s\n", task.Due)
	fmt.Printf("Self:    %s\n", task.SelfLink)
}