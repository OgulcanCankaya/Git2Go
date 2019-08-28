package main

import (
	"bufio"
	"fmt"
	"github.com/libgit2/git2go"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose option...\n1.Status\n2.Fetch\n3.Merge\n4.Pull\n5.Push\n6.Clone\n7.Commit\n8.Custom command")

	/*put some logging mechanisms*/

	text, _ := reader.ReadString('\n')
	for text != "q\n" {
		if text == "1\n" {
			println("Status Funct")
			fmt.Println("Enter path of repo:")
			path, _ := reader.ReadString('\n')
			path = path[:len(path)-1]
			gitstat(path)
		}
		if text == "2\n" {
			println("Fetch Funct")
			fmt.Println("Please enter the git repository path")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\n")
			repo, _ := git.OpenRepository(path)
			gitFetch(repo)
		}
		if text == "3\n" {
			println("Merge Funct")
			fmt.Println("Please enter the git repository path")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\n")
			repo, _ := git.OpenRepository(path)
			err := gitMerge(repo, &signature)
			checkErr(err)
		}
		if text == "4\n" {
			println("Pull Funct")
			fmt.Println("Please enter the git repository path")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\n")
			repo, _ := git.OpenRepository(path)
			err := gitPull2(repo)
			if err != nil {
				fmt.Println("we couldnt pull :| ")
				log.Println(err)
			}
		}
		if text == "5\n" {
			println("Push Funct")
			fmt.Println("Please enter the git repository path")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\n")
			repo, err := git.OpenRepository(path)
			checkErr(err)
			gitPush(repo)
		}
		/*git clone*/
		if text == "6\n" {
			println("Clone Funct")
			fmt.Println("Please enter the repo address details in the following format.")
			fmt.Println("git@github.com:{username}/{repository-name}.git")
			url, _ := reader.ReadString('\n')
			url=strings.TrimSuffix(url,"\n")
			fmt.Println("Please enter the path")
			path, _ := reader.ReadString('\n')
			path = path[:len(path)-1]
			go gitClone(url, path,&cloneOpt)
		}
		if text == "7\n" {
			fmt.Println("Please enter the git repository path")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\n")
			repo, _ := git.OpenRepository(path)
			fmt.Println("Please enter message to commit: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSuffix(path, "\n")
			fmt.Println("Commiting...")
			//gitCommit (repo , message , &signature)
			err := GitAddCommit(&signature, repo, message)
			if err != nil {
				log.Println(err)
			}
		}
		if text == "8\n" {
			fmt.Println("Please enter the git repository path")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\n")
			/*check if a .git file exists*/
			/*put a loop for that job*/ /*TO DO*/
			fmt.Println("Type the command in following format:\n{command} {optional-- side-command} -{params} {args} == git commit -m \"initial commit\"")
			commandStr, _ := reader.ReadString('\n')
			commandArg := strings.Split(commandStr, " ")
			argLen := len(commandArg)
			for i := range commandArg {
				commandArg[i] = strings.TrimSuffix(commandArg[i], "\n")
				commandArg[i] = strings.TrimSuffix(commandArg[i], " ")
			}
			cmdName := commandArg[0]
			commandArg[argLen-1] = strings.TrimSuffix(commandArg[argLen-1], "\n")
			fmt.Println("\nRunning custom command on : ", path)
			runCom(cmdName, commandArg, path)
		}
		time.Sleep( 3/2 * time.Second)
		println("Wanna do more?  \nAnswer y\\n ?:")
		text, _ = reader.ReadString('\n')
		if text == "y\n" {
			fmt.Println("Choose option...\n1.Status\n2.Fetch\n3.Merge\n4.Pull\n5.Push\n6.Clone\n7.Commit\n8.Custom command")
			text, _ = reader.ReadString('\n')
		}
		if text == "n\n" {
			text = "q\n"
		}
	}
	fmt.Println("\nGoodBye :)")
}