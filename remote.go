package gitrmt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Remote struct {
	handler  Handler
	lazyWork []func() (string, error)
}

func NewRemote(handler Handler) *Remote {
	log.Printf("$GIT_DIR=%v\n", os.Getenv("GIT_DIR"))

	return &Remote{
		handler:  handler,
		lazyWork: make([]func() (string, error), 0),
	}
}

func (r *Remote) Run(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)

loop:
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		command = strings.Trim(command, "\n")

		switch {
		case command == "capabilities":
			r.output(
				out,
				"%s\n",
				strings.Join(r.handler.Capabilities(), "\n"),
			)
		case strings.HasPrefix(command, "list"):
			list, err := r.handler.List(strings.HasPrefix(command, "list for-push"))
			if err != nil {
				return err
			}

			for _, e := range list {
				r.output(out, "%s\n", e)
			}

			_, _ = fmt.Fprint(out, "\n")
		case strings.HasPrefix(command, "push "):
			refs := strings.Split(command[5:], ":")
			isForce := strings.HasPrefix(refs[0], "+")
			r.push(refs, isForce)
		case strings.HasPrefix(command, "fetch "):
			parts := strings.Split(command, " ")
			if parts[1] != "0000000000000000000000000000000000000000" {
				r.fetch(parts[1], parts[2])
			}
		case command == "":
			fallthrough
		case command == "\n":
			log.Println("doing work...")

			for _, task := range r.lazyWork {
				resp, err := task()
				if err != nil {
					return fmt.Errorf("error processing task: %w", err)
				}
				r.output(out, "%s", resp)
			}
			r.output(out, "%s", "\n")
			r.lazyWork = nil
			break loop
		default:
			return fmt.Errorf("received unknown command %q", command)
		}
	}

	return r.handler.Finish()
}

func (r *Remote) output(out io.Writer, format, resp string) {
	if _, err := fmt.Fprintf(out, format, resp); err != nil {
		log.Printf("error outputting %q: %v", resp, err)
	}
}

func (r *Remote) fetch(sha string, ref string) {
	r.lazyWork = append(r.lazyWork, func() (string, error) {
		return "", nil
	})
}

func (r *Remote) push(refs []string, force bool) {
	src, dst := refs[0], refs[1]

	r.lazyWork = append(r.lazyWork, func() (string, error) {
		done, err := r.handler.Push(src, dst, force)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("ok %s\n", done), nil
	})
}
