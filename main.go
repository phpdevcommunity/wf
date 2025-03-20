//
//
// @Author: F.Michel

// @github: https://github.com/phpdevcommunity
package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var (
	workflows            []Workflow
	dockerComposeCommand = ""
)

func main() {

	currentDir := getCurrentDir()
	wfFiles, err := getWfFiles(currentDir)
	if err != nil {
		fmt.Println(err)
	}

	if len(wfFiles) == 0 {
		fmt.Println("No workflow files found")
	}

	workflows = make([]Workflow, 0)
	for _, wfFile := range wfFiles {

		wfs := ParseContentToWorkFlowStruct(wfFile)
		for _, wf := range wfs {
			if len(wf.Lines) == 0 {
				continue
			}
			workflows = append(workflows, wf)
		}
	}

	Commands := []*cli.Command{}
	values := InitDefaultVariables()
	for _, wf := range workflows {

		Commands = append(Commands, &cli.Command{
			Name:  wf.Name,
			Usage: wf.Comment,
			Action: func(c *cli.Context) error {
				executeWorkflow(wf, values)
				return nil
			},
		})
	}

	app := &cli.App{
		Commands: Commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
func executeWorkflow(wf Workflow, values *map[string]string) {
	for _, line := range wf.Lines {
		executeLine(line, values)
	}
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func Touch(filename string) (bool, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}
	fd.Close()
	return true, nil
}

func executeLine(lineOriginal string, values *map[string]string) {

	if lineOriginal == "" || strings.HasPrefix(lineOriginal, "#") {
		return
	}

	lineOriginal = ResolveVariables(*values, lineOriginal)
	actionOriginal := strings.Split(lineOriginal, " ")[0]
	action := strings.ToLower(actionOriginal)
	line := strings.Replace(lineOriginal, actionOriginal, action, 1)
	if action == "set" {
		line = strings.TrimPrefix(line, "set ")
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			pterm.Error.WithFatal().Println("Invalid set command")
		}
		if parts[0] == "" || parts[1] == "" {
			pterm.Error.WithFatal().Println("Invalid set command")
		}

		(*values)[parts[0]] = parts[1]
		pterm.Info.Printfln("Variable %s set to %s", parts[0], (*values)[parts[0]])

		return
	} else if strings.HasPrefix(line, "#") {
		return
	} else if action == "run" {
		line = strings.TrimPrefix(line, "run ")
		run(line, true)
		return
	} else if action == "echo" {
		line = strings.TrimPrefix(line, "echo ")
		fmt.Println(line)
		return
	} else if action == "exit" {
		os.Exit(0)
	} else if action == "touch" {
		line = strings.TrimPrefix(line, "touch ")
		fileToCreate := strings.TrimSpace(line)
		if fileToCreate == "" {
			pterm.Error.WithFatal().Println("Invalid touch command")
		}
		if FileExists(fileToCreate) {
			pterm.Info.Printfln("%s already exists, skipping...", fileToCreate)
			return
		}
		Touch(line)
		fmt.Println(pterm.FgGreen.Sprintf("%s created", line))
		return
	} else if action == "copy" || action == "cp" {
		line = strings.TrimPrefix(line, "copy ")
		parts := strings.Split(line, " ")

		if len(parts) < 2 {
			pterm.Error.WithFatal().Println("Invalid copy command")
		}
		fileOrFolder := parts[0]
		destination := parts[1]
		if !FileExists(fileOrFolder) {
			pterm.Error.WithFatal().Printfln("%s not found", fileOrFolder)
		}

		if FileExists(destination) {
			pterm.Info.Printfln("%s already exists, skipping...", destination)
			return
		}

		Copy(fileOrFolder, destination)
		fmt.Println(pterm.FgGreen.Sprintf("%s copied to %s", fileOrFolder, destination))
		return

	} else if action == "mkdir" {
		folderName := strings.Split(line, " ")[1]
		if folderName == "" {
			pterm.Error.WithFatal().Println("Invalid mkdir command")
		}

		if FileExists(folderName) {
			pterm.Info.Printfln("Folder %s already exists, skipping...", folderName)
			return
		}

		err := os.MkdirAll(folderName, os.ModePerm)
		if err != nil {
			pterm.Error.WithFatal().Println(err)
		}

		fmt.Println(pterm.FgGreen.Sprintf("Folder %s created", folderName))
		return
	} else if action == "set_permissions" {
		args := strings.Fields(line)
		if len(args) < 3 {
			pterm.Error.WithFatal().Printf("Invalid set_permissions command: %s", line)
		}

		folderName := args[1]
		permissions, err := strconv.ParseUint(args[2], 8, 32) // Use ParseUint instead of Atoi
		if err != nil {
			pterm.Error.WithFatal().Printf("Invalid permissions: %s", args[2])
		}

		err = os.Chmod(folderName, os.FileMode(permissions))
		if err != nil {
			pterm.Error.WithFatal().Printf("Failed to set permissions for folder/file %s: %s", folderName, err)
		}

		fmt.Println(pterm.FgGreen.Sprintf("Permissions 0%o set for folder/file %s", permissions, folderName))
		return
	} else if action == "sync_time" {
		fmt.Println("Checking time...")
		run("date", false)
		return
	} else if action == "docker_compose" {
		args := strings.Fields(line)
		if len(args) < 2 {
			pterm.Error.WithFatal().Printf("Invalid execute command: %s", line)
		}
		dockerComposeExec := args[1:]
		dockerCompose := GetDockerComposeCommand()
		run(fmt.Sprintf("%s %s", dockerCompose, strings.Join(dockerComposeExec, " ")), true)
		return
	} else if action == "wf" {
		args := strings.Fields(line)
		if len(args) < 2 {
			pterm.Error.WithFatal().Printf("Invalid wf command: %s", line)
		}
		WorkflowName := args[1]
		for _, wf := range workflows {
			if wf.Name == WorkflowName {
				pterm.Description.Printfln("Executing workflow: %s", WorkflowName)
				executeWorkflow(wf, values)
				break
			}
		}
		return
	}

	litsOfNotifiesOptions := []string{
		"notify_success",
		"notify_error",
		"notify_warning",
		"notify_info",
		"notify",
	}

	for _, option := range litsOfNotifiesOptions {
		if action == option {
			message := strings.TrimPrefix(line, option)
			message = strings.TrimSpace(message)

			message = strings.TrimPrefix(message, "\"")
			message = strings.TrimSuffix(message, "\"")

			if option == "notify" {
				fmt.Println(message)
			} else if option == "notify_success" {
				pterm.Success.Println(message)
			} else if option == "notify_error" {
				pterm.Error.Println(message)
			} else if option == "notify_warning" {
				pterm.Warning.Println(message)
			} else if option == "notify_info" {
				pterm.Info.Println(message)
			}

			return
			// pterm.DefaultBasicText.Println(message)
		}
	}

	pterm.Error.WithFatal().Println("Unknown command or invalid syntax : " + line)

}

func run(line string, printCommand bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	if printCommand {
		pterm.DefaultBasicText.Println(line)
	}
	args := strings.Split(line, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Run()
	if err != nil {
		pterm.Error.Println("Error :", err)
		os.Exit(1)
	}

}

func getWfFiles(currentDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(currentDir, "*.wf"))
}

func getCurrentDir() string {
	dir, _ := os.Getwd()
	return dir
}

func FileGetContents(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func Copy(source string, dest string) (bool, error) {
	fd1, err := os.Open(source)
	if err != nil {
		return false, err
	}
	defer fd1.Close()
	fd2, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer fd2.Close()
	_, e := io.Copy(fd2, fd1)
	if e != nil {
		return false, e
	}
	return true, nil
}

func ParseContentToWorkFlowStruct(filename string) map[string]Workflow {

	workflowsMap := map[string]Workflow{}
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentWorkflow := Workflow{
		Name:    "main",
		Comment: "",
	}
	workflowsMap[currentWorkflow.Name] = currentWorkflow

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		workflowName := currentWorkflow.Name
		comment := currentWorkflow.Comment

		if strings.HasPrefix(line, "[") {

			cursor := 0
			findWorkflowName := false
			for i := 0; i < len(line); i++ {
				if line[i] == ']' {
					workflowName = strings.TrimPrefix(line[:i], "[")
					cursor = i
					findWorkflowName = true
					break
				}
			}

			if !findWorkflowName {
				continue
			}

			afterWorkflowNameLine := line[cursor+1:]
			afterWorkflowNameLine = strings.TrimSpace(afterWorkflowNameLine)
			line = ""
			if strings.HasPrefix(afterWorkflowNameLine, "#") {
				comment = strings.TrimPrefix(afterWorkflowNameLine, "#")
				comment = strings.TrimSpace(comment)
			}

		}

		wf, ok := workflowsMap[workflowName]
		if !ok {
			wf = Workflow{
				Name:    workflowName,
				Comment: comment,
			}
			workflowsMap[workflowName] = wf
		}

		wf.Lines = append(wf.Lines, line)
		workflowsMap[workflowName] = wf
		currentWorkflow = wf

		currentWorkflow.Lines = append(currentWorkflow.Lines, line)

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return workflowsMap

}

func InitDefaultVariables() *map[string]string {

	values := map[string]string{}
	values["IP_LOCAL"] = GetLocalIP()
	values["CURRENT_PATH"] = getCurrentDir()

	return &values

}

func GetLocalIP() string {
	// Récupère toutes les interfaces réseau
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, i := range interfaces {
		// Récupère toutes les adresses de chaque interface
		addrs, err := i.Addrs()
		if err != nil {
			return ""
		}

		for _, addr := range addrs {
			var ip net.IP

			// Vérifie si l'adresse est de type IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Filtre pour obtenir une adresse IPv4 non-loopback
			if ip != nil && ip.IsLoopback() == false && ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return ""
}

func TokenGenerator(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func ResolveVariables(values map[string]string, line string) string {
	for key, value := range values {
		line = strings.ReplaceAll(line, "${"+key+"}", value)
	}
	line = strings.ReplaceAll(line, "${GENERATE_SECRET}", TokenGenerator(32))
	return line

}

func GetDockerComposeCommand() string {
	if dockerComposeCommand != "" {
		return dockerComposeCommand
	}

	cmd := exec.Command("docker", "compose", "--version")
	cmd.Env = os.Environ()
	_, err := cmd.Output()
	if err != nil {
		dockerComposeCommand = "docker-compose"
	} else {
		dockerComposeCommand = "docker compose"
	}
	return dockerComposeCommand
}

type Workflow struct {
	Name    string
	Comment string
	Lines   []string
}
