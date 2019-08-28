package main
import (
	git "github.com/libgit2/git2go"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)
var cloneOpt = git.CloneOptions{
	CheckoutOpts:         nil,
	FetchOptions:         &git.FetchOptions{

		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:          credentialsCallback,
			CertificateCheckCallback:     certificateCheckCallback,
		},

		Prune:           0,
		UpdateFetchhead: false,
		DownloadTags:    0,
		Headers:         nil,
	},
	Bare:                 false,
	CheckoutBranch:       "",
	RemoteCreateCallback: nil,
}
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

func AddAll(repo git.Repository) error {
	idx, err := repo.Index()
	if err != nil {
		return err
	}
	err = idx.AddAll([]string{}, git.IndexAddDefault, nil)
	if err != nil {
		return err
	}
	err = idx.Write()
	return err
}

func gitPull2 ( repo *git.Repository  ) error {
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}
	called := false
	err = remote.Fetch([]string{}, &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: func(url string, username_from_url string, allowed_types git.CredType) (git.ErrorCode, *git.Cred) {
				if called {
					return git.ErrUser, nil
				}
				called = true
				ret, cred := git.NewCredSshKey("git", "/home/ogulcan/.ssh/id_rsa.pub", "/home/ogulcan/.ssh/id_rsa", "")
				return git.ErrorCode(ret), &cred
			},
		},
	}, "")
	if err != nil {
		return err
	}
	remoteBranch, err := repo.References.Lookup("refs/remotes/" + "origin" + "/" + "master")
	if err != nil {
		return err
	}
	mergeRemoteHead, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		return err
	}
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = mergeRemoteHead
	if err = repo.Merge(mergeHeads, nil, nil); err != nil {
		return err
	}
	// Check if the index has conflicts after the merge
	idx, err := repo.Index()
	if err != nil {
		return err
	}
	currentBranch, err := repo.Head()
	if err != nil {
		return err
	}
	localCommit, err := repo.LookupCommit(currentBranch.Target())
	if err != nil {
		return err
	}
	// If index has conflicts, read old tree into index and
	// return an error.
	if idx.HasConflicts() {
		_ = repo.ResetToCommit(localCommit, git.ResetHard, &git.CheckoutOpts{})
		_ = repo.StateCleanup()
		return errors.New("conflict")
	}
	// If everything looks fine, create a commit with the two parents
	treeID, err := idx.WriteTree()
	if err != nil {
		return err
	}
	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return err
	}
	remoteCommit, err := repo.LookupCommit(remoteBranch.Target())
	if err != nil {
		return err
	}
	sig := &git.Signature{Name: "OgulcanCankaya", Email: "email", When: time.Now()}
	_, err = repo.CreateCommit("HEAD", sig, sig, "merged", tree, localCommit, remoteCommit)
	if err != nil {
		return err
	}
	_ = repo.StateCleanup()
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
	branchName := "master"
	if err := remote.Fetch([]string{}, &git.FetchOptions{}, ""); err != nil {
		return err
	}
	// Merge
	remoteRef, err := repo.References.Lookup("refs/remotes/origin/" + branchName)
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
	err := repo.Remotes.AddFetch("master","+refs/heads/*:refs/remotes/origin/*")
	checkErr(err)
}

func gitPush(repo *git.Repository) error {
	remote, err := repo.Remotes.Lookup("master")
	if err != nil {
		return err
	}
	called := false
	err = remote.Push([]string{"refs/heads/" + "master"}, &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: func(url string, username_from_url string, allowed_types git.CredType) (git.ErrorCode, *git.Cred) {
				if called {
					return git.ErrUser, nil
				}
				called = true
				ret, cred := git.NewCredSshKey("git", "/home/ogulcan/.ssh/id_rsa.pub", "/home/ogulcan/.ssh/id_rsa", "")
				return git.ErrorCode(ret), &cred
			},
		},
	})
	return err
}

func gitPush2(repo *git.Repository) error {
	remote, err := repo.Remotes.Lookup("master")
	if err != nil {
		return err
	}
	called := false
	err = remote.Push([]string{"refs/heads/" + "origin/master"}, &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: func(url string, username_from_url string, allowed_types git.CredType) (git.ErrorCode, *git.Cred) {
				if called {
					return git.ErrUser, nil
				}
				called = true
				ret, cred := git.NewCredSshKey("git", "/home/ogulcan/.ssh/id_rsa.pub", "/home/ogulcan/.ssh/id_rsa", "")
				return git.ErrorCode(ret), &cred
			},
		},
	})
	return err
}

func gitstat(path string){
	cmd := exec.Command("git" , "status" , path)
	cmd.Dir=path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	ret, cred := git.NewCredSshKey("git", "/home/ogulcan/.ssh/id_rsa.pub", "/home/ogulcan/.ssh/id_rsa", "")
	return git.ErrorCode(ret), &cred
}

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}

func gitClone( url string, path string, cloneOpt *git.CloneOptions){
	repo, err := git.Clone(url, path, cloneOpt)
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
	gitCommit(repo, "merged origin/master to master",signature)
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
	if err := remote.Fetch([]string{}, &git.FetchOptions{}, ""); err != nil {
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