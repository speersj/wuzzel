// wuzzel interacts with fuzzel to turn it into a window picker.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// SwayNode contains information we need to display / select windows
type SwayNode struct {
	ID            int        `json:"id"`
	Type          string     `json:"type"`
	Name          string     `json:"name"`
	Nodes         []SwayNode `json:"nodes"`
	FloatingNodes []SwayNode `json:"floating_nodes"`
}

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	fail := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format, args...)
		exitCode = 1
	}

	tree, err := swayTree()
	if err != nil {
		fail("Error getting window list: %v\n", err)
		return
	}

	windows := tree.swayWindows()
	names := ""
	for _, w := range windows {
		name := strings.TrimSpace(w.Name)
		if name != "" {
			names += name + "\n"
		}
	}

	outp, err := fuzzelExec(os.Args[1:], names)
	if err != nil {
		// typically this means the user did not selected anything
		return
	}

	// get the selected window name
	sel := fuzzelParseOutput(outp)
	if sel == "" {
		return
	}

	// find first window with matching name
	for _, w := range windows {
		if w.Name == sel {
			// focus the window
			cmd := exec.Command("swaymsg", fmt.Sprintf("[con_id=%d]", w.ID), "focus")
			if err := cmd.Run(); err != nil {
				fail("Error running swaymsg: %v\n", err)
				return
			}
		}
	}

}

// swayWindows equivalent to get_windows
func (n SwayNode) swayWindows() []SwayNode {
	var windows []SwayNode
	for _, node := range n.Nodes {
		if node.Name != "__i3" && node.Type == "output" {
			for _, ws := range node.Nodes {
				if ws.Type == "workspace" {
					windows = append(windows, ws.swayWorkspaceWindows()...)
				}
			}
		}
	}

	return windows
}

// swayWorkspaceWindows returns the windows in the given workspace
// equivalent to extract_nodes_iterative
func (n SwayNode) swayWorkspaceWindows() []SwayNode {
	var allNodes []SwayNode
	allNodes = append(allNodes, n.FloatingNodes...)

	for _, node := range n.Nodes {
		if len(node.Nodes) == 0 {
			allNodes = append(allNodes, node)
		} else {
			allNodes = append(allNodes, node.Nodes...)
		}
	}
	return allNodes
}

// swayTree returns the current window tree
func swayTree() (SwayNode, error) {
	cmd := exec.Command("swaymsg", "-t", "get_tree")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return SwayNode{}, err
	}

	var root SwayNode
	err = json.Unmarshal(output, &root)
	if err != nil {
		return SwayNode{}, err
	}
	return root, nil
}

// fuzzelParseOutput returns the last non-empty line from fuzzel output,
// which should be the name of the selected window or blank if none
func fuzzelParseOutput(fuzzelOut string) string {
	lines := strings.Split(fuzzelOut, "\n")
	if len(lines) < 2 {
		return ""
	}
	return lines[len(lines)-2]
}

// run fuzzel with given arguments.
// input should be a string with window names separated by newlines
func fuzzelExec(args []string, input string) (string, error) {
	cmd := exec.Command("fuzzel", append([]string{"-d"}, args...)...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(input))
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
