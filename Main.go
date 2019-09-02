package main
import (
	"bufio"
	"errors"
	"fmt"
	"github.com/libgit2/git2go"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var signature = git.Signature{
	Name:  "Ogulcan Cankaya",
	Email: "ogulcan985@outlook.com",
	When:  time.Now(),
}
/*func to take all repo directories*/
func checkErr(err error){
	if err != nil {
		log.Println(err)
	}
}

func gitPull2 ( repo *git.Repository  ) error {

	// Locate remote
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}
	// Fetch changes from remote
	if err := remote.Fetch([]string{}, &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:          credentialsCallback,
			CertificateCheckCallback:     certificateCheckCallback,
		},}, ""); err != nil {
		return err
	}
	// Get remote master
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/"+"master")
	if err != nil {
		return err
	}
	remoteBranchID := remoteBranch.Target()
	// Get annotated commit
	annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		return err
	}
	// Do the merge analysis
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = annotatedCommit
	analysis, _, err := repo.MergeAnalysis(mergeHeads)
	if err != nil {
		return err
	}
	// Get repo head
	head, err := repo.Head()
	if err != nil {
		return err
	}
	if analysis & git.MergeAnalysisUpToDate != 0 {
		return nil
	}  else if analysis & git.MergeAnalysisNormal != 0 {
		// Just merge changes
		if err := repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
			return err
		}
		// Check for conflicts
		index, err := repo.Index()
		if err != nil {
			return err
		}
		if index.HasConflicts() {
			return errors.New("Conflicts encountered. Please resolve them.")
		}
		// Make the merge commit
		sig, err := repo.DefaultSignature()
		if err != nil {
			return err
		}
		// Get Write Tree
		treeId, err := index.WriteTree()
		if err != nil {
			return err
		}
		tree, err := repo.LookupTree(treeId)
		if err != nil {
			return err
		}
		localCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}
		remoteCommit, err := repo.LookupCommit(remoteBranchID)
		if err != nil {
			return err
		}
		repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)
		// Clean up
		repo.StateCleanup()
	} else if analysis & git.MergeAnalysisFastForward != 0 {
		// Fast-forward changes
		// Get remote tree
		remoteTree, err := repo.LookupTree(remoteBranchID)
		if err != nil {
			return err
		}
		// Checkout
		if err := repo.CheckoutTree(remoteTree, nil); err != nil {
			return err
		}
		branchRef, err := repo.References.Lookup("refs/heads/"+"master")
		if err != nil {
			return err
		}
		// Point branch to the object
		branchRef.SetTarget(remoteBranchID, "")
		if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unexpected merge analysis result %d", analysis)
	}
	return nil
}
func   gitPull( repo *git.Repository ) error  {
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		remote, err = repo.Remotes.Create("origin", repo.Path())
		if err != nil {
			return err
		}
	}
	if err := remote.Fetch([]string{}, &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:          credentialsCallback,
			CertificateCheckCallback:     certificateCheckCallback,
		},}, ""); err != nil {
		return err
	}
	remoteRef, err := repo.References.Lookup("refs/remotes/origin/" + "master")
	if err != nil {
		return err
	}
	mergeRemoteHead, err := repo.AnnotatedCommitFromRef(remoteRef)
	if err != nil {
		return err
	}
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = mergeRemoteHead
	if err = repo.Merge(mergeHeads, nil, nil); err != nil {
		return err
	}
	return nil
}
func gitFetch(repo *git.Repository)  {
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		remote, err = repo.Remotes.Create("origin", repo.Path())
		if err != nil {
			checkErr(err)
		}
	}
	if err := remote.Fetch([]string{}, &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:          credentialsCallback,
			CertificateCheckCallback:     certificateCheckCallback,
		},}, ""); err != nil {
		checkErr(err)
	}
}
func gitPush(repo *git.Repository)  {
	remote, err := repo.Remotes.Lookup("master")
	if err != nil {
		checkErr(err)
	}
	err = remote.Push([]string{"refs/heads/" + "/origin"}, nil)
}
func gitPush2(repo *git.Repository) error {
	remote, err := repo.Remotes.Lookup("master")
	if err != nil {
		return err
	}
	err = remote.Push([]string{"refs/heads/" + "origin/master"}, &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: credentialsCallback,
			},
		},
	)
	return err
}
func gitStat(path string){
	cmd := exec.Command("git" , "status" , path)
	cmd.Dir=path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
func credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	ret, cred := git.NewCredSshKey("git", "/home/ogulcan/.ssh/id_ed25519.pub", "/home/ogulcan/.ssh/id_ed25519", "")
	return git.ErrorCode(ret), &cred
}
func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}
func gitClone( url string, path string){
	cloneOpt := git.CloneOptions{
		FetchOptions:         &git.FetchOptions{
			RemoteCallbacks: git.RemoteCallbacks{
				CredentialsCallback:          credentialsCallback,
				CertificateCheckCallback:     certificateCheckCallback,
			},
		},
	}
	repo, err := git.Clone(url, path, &cloneOpt)
	fmt.Println(repo)
	if err != nil {
		fmt.Println("www.errors.com\n", err)
		//log.Panic(err)
	}
	log.Print(repo)
}
func GitAddCommit(sig *git.Signature, repo *git.Repository, message string) error {
	// Retrieve index
	index, err := repo.Index()
	if err != nil {
		log.Println("Index - ", err)
	}
	// See if we had conflicts before we added everything to the index
	indexHadConflicts := index.HasConflicts()
	// Add everything to the index
	err = index.AddAll([]string{}, git.IndexAddDefault, nil)
	if err != nil {
		log.Println("AddAll - ", err)
	}
	// Write the index
	err = index.Write()
	if err != nil {
		log.Println("Write - ", err)
	}
	// Write the current index tree to the repo
	treeId, err := index.WriteTreeTo(repo)
	if err != nil {
		log.Println("WriteTreeTo - ", err)
	}
	//Retrieve the tree we just wrote
	tree, err := repo.LookupTree(treeId)
	if err != nil {
		log.Println("LookupTree - ", err)
	}
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/master")
	if err != nil {
		log.Println("Fetch 2 - ", err)
	}
	currentBranch, err := repo.Head()
	if err != nil {
		log.Println("Head - ", err)
	}
	// Retrieve current commit
	currentTip, err := repo.LookupCommit(currentBranch.Target())
	if err != nil {
		log.Println("LookupCommit - ", err)
	}
	// Create a new one on top
	currentCommit, err := repo.CreateCommit("HEAD", sig, sig, message, tree, currentTip)
	if err != nil {
		log.Println("CreateCommit - ", err)
	}
	//  If there were conflicts, do the merge commit
	if indexHadConflicts == true {
		localCommit, err := repo.LookupCommit(currentCommit)
		if err != nil {
			log.Println("Fetch 11 - ", err)
		}
		remoteCommit, err := repo.LookupCommit(remoteBranch.Target())
		if err != nil {
			log.Println("Fetch 12 - ", err)
		}
		// Create a new one
		commitId, _ := repo.CreateCommit("HEAD", sig, sig, "Merge commit", tree, localCommit, remoteCommit)
		println(commitId)
		// Clean up
		_ = repo.StateCleanup()
	}
	return err
}
func gitMerge(repo *git.Repository, signature *git.Signature) error {
	/*commit  "merge operation" for double merge protection  */
	/*ahead behind check ???*/
	_ = GitAddCommit(signature, repo, "merged origin/master to master")
	sourceBranchName := "origin/master"
	destinationBranchName := "master"
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		remote, err = repo.Remotes.Create("origin", repo.Path())
		if err != nil{
			log.Print("Failed lookup remote branch " + sourceBranchName)
			return err
		}
	}
	if err := remote.Fetch([]string{}, &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:          credentialsCallback,
			CertificateCheckCallback:     certificateCheckCallback,
		},
	}, ""); err != nil {
		return err
	}
	// Merge
	remoteRef, err := repo.References.Lookup("refs/remotes/origin/" + destinationBranchName)
	checkErr(err)
	mergeRemoteHead, err := repo.AnnotatedCommitFromRef(remoteRef)
	checkErr(err)
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = mergeRemoteHead
	if err = repo.Merge(mergeHeads, nil, nil); err != nil {
		return err
	}
	return nil
}
func gitCommit (repo *git.Repository, message string, signature *git.Signature){
	/*Adding directory*/
	idx, err := repo.Index()
	if err != nil {
		panic(err)
	}
	var path []string
	pathAll := append(path,".")
	log.Println(pathAll)
	err = idx.AddAll(pathAll, git.IndexAddDefault, nil)
	if err != nil {
		log.Println(err)
		return
	}
	err = idx.Write()
	treeId, err := idx.WriteTreeTo(repo)
	if err != nil {
		panic(err)
	}
	/*added directory*/
	/*trying to create commit*/
	head, err := repo.Head()
	if err != nil {
		log.Println(err)
	}
	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		log.Println(err)
	}
	tree, err := repo.LookupTree(treeId)
	if err != nil {
		log.Println("LookUpTree error - ",err)
	}
	commitId, err := repo.CreateCommit("HEAD", signature, signature, message, tree, headCommit)
	if err != nil {
		panic(err)
	}
	log.Println(commitId)
	indexHadConflicts := idx.HasConflicts()
	if indexHadConflicts == true {
		localCommit, err := repo.LookupCommit(commitId)
		if err != nil {
			log.Println(err)
		}
		commitId,err = repo.CreateCommit("HEAD", signature, signature, "Merge commit", tree, localCommit, headCommit)
	}
}
func runCom(cmdName string, commandArg []string, path string) {
	/* Try to execute the commands on every .git*/
	fmt.Println("Command output is:")
	cmd := exec.Command("ls")
	if len(commandArg) > 1 {
		cmd = exec.Command(cmdName,commandArg[1:]...)
	}
	if len(commandArg) == 1 {
		cmd = exec.Command(cmdName)
	}
	cmd.Dir=path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
} /**needs maintenance i think*/
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
			gitStat(path)
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
				fmt.Println("we couldnt pull2 :| ")
				log.Println(err)
			}
			fmt.Println("after pull2")
			err = gitPull(repo)
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
			gitPush2(repo)
		}
		/*git clone*/
		if text == "6\n" {
			println("Clone Funct")
			fmt.Println("Please enter the repo address details in the following format.")
			fmt.Println("{username}/{repository-name}.git")
			url, _ := reader.ReadString('\n')
			url=strings.TrimSuffix(url,"\n")
			url = "git@github.com:"+url
			fmt.Println(url)
			fmt.Println("Please enter the path")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix( path,"\n")
			go gitClone(url, path)
		}
		if text == "7\n" {
			fmt.Println("Please enter the git repository path")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\n")
			repo, _ := git.OpenRepository(path)
			fmt.Println("Please enter message to commit: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSuffix(message, "\n")
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