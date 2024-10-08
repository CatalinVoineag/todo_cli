package main

import (
	"fmt"
	"os"
	"flag"
	"bufio"
	"io"
	"strings"

	"github.com/catalinVoineag/todo_cli"
)

var todoFileName = ".todo.json"

func main() {
  flag.Usage = func() {
    fmt.Fprintf(flag.CommandLine.Output(),
    "%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
    fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
    fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
    flag.PrintDefaults()
  }

  add := flag.Bool("add", false, "Add task to the Todo list")
  list := flag.Bool("list", false, "List all tasks")
  complete := flag.Int("complete", 0, "Item to be completed")
  flag.Parse()

  if os.Getenv("TODO_FILENAME") != "" {
    todoFileName = os.Getenv("TODO_FILENAME")
  }

  l := &todo_cli.List{}

  if err := l.Get(todoFileName); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  switch {
  case *list:
    fmt.Print(l)
  case *complete > 0:
    if err := l.Complete(*complete); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }

    if err := l.Save(todoFileName); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
  case *add:
    // When any arguments (excluding flags) are provided, they will be
    // used as the new task
    t, err := getTask(os.Stdin, flag.Args()...)

    if err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
    l.Add(t)

    if err := l.Save(todoFileName); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
  default:
    fmt.Fprintln(os.Stderr, "Invalid option")
    os.Exit(1)
  }
}

// getTask function decides where to get the description for a new
// task from: arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
  if len(args) > 0 {
    return strings.Join(args, " "), nil
  }

  s := bufio.NewScanner(r)
  s.Scan()

  if err := s.Err(); err != nil {
    return "", err
  }

  if len(s.Text()) == 0 {
    return "", fmt.Errorf("Task cannot be blank")
  }

  return s.Text(), nil
}
