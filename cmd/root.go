/*
Copyright © 2022 Cole Arendt <dev@colearendt.com>

*/
package cmd

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "traefik-plugin-init",
	Short: "Download traefik plugins to a directory",
	Long: `A small CLI that is useful for downloading traefik plugins to a directory.

This is especially useful as an init container in Kubernetes.
    `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Initializing traefik plugins")

		// get environment prefix
		pref, ok := os.LookupEnv("TRAEFIK_PLUGIN_PREFIX")
		if !ok {
			pref = "TRAEFIK_PLUGIN_REPO_"
		}
		log.Printf("Using TRAEFIK_PLUGIN_PREFIX=%v\n", pref)

		// TODO: handle possible panic...
		prefReg := regexp.MustCompile("^" + pref)

		// get the directory to clone into
		dir, ok := os.LookupEnv("TRAEFIK_PLUGIN_PATH")
		if !ok {
			dir = "/plugin-storage"
		}
		log.Printf("Using TRAEFIK_PLUGIN_PATH=%v\n", dir)

		var arrPlugins []string

		// find env vars matching prefix
		var allEnvVars []string = os.Environ()
		for _, e := range allEnvVars {
			if prefReg.MatchString(e) {
				arrPlugins = append(arrPlugins, e)
			}
		}

		if len(arrPlugins) > 0 {
			for _, plugin := range arrPlugins {
				err := clonePlugin(plugin, dir)
				if err != nil {
					log.Printf("Error initializing %v: %v\n", plugin, err)
				}
			}
		} else {
			log.Print("No plugin variables found. Exiting")
			return
		}
		log.Printf("Done cloning plugins")
	},
}

func clonePlugin(env string, dir string) (err error) {
	envVarRegex := regexp.MustCompile("^([^=]+)=([^@]+)@?([^@]*)?")

	res := envVarRegex.FindStringSubmatch(env)

	repoRegex := regexp.MustCompile("^.*/([^/]+)/([^/]+)")

	// set up data from input
	var repoUrl string
	var repoOwner string
	var repoName string
	if len(res) >= 2 {
		repoUrl = res[2]
		if len(repoUrl) == 0 {
			return errors.New("input is empty")
		}

		repoDirs := repoRegex.FindStringSubmatch(repoUrl)
		if len(repoDirs) < 3 {
			return errors.New("input does not look like a repository (owner/repository)")
		}

		repoOwner = repoDirs[1]
		repoName = repoDirs[2]

	} else {
		return errors.New("input did not match regex")
	}

	ref := res[3]
	if len(ref) == 0 {
		ref = "main"
	}

	masterRef := plumbing.NewBranchReferenceName("master")
	mainRef := plumbing.NewBranchReferenceName("main")

	// clone or otherwise open the repository
	// TODO: a way to be less redundant here?
	cloneOptsMaster := git.CloneOptions{Depth: 0, URL: repoUrl, SingleBranch: true, ReferenceName: masterRef}
	cloneOptsMain := git.CloneOptions{Depth: 0, URL: repoUrl, SingleBranch: true, ReferenceName: mainRef}

	var repoLocal *git.Repository
	// TODO: a way to build a path more cleanly (i.e. multiple slashes, etc.)
	fullDir := dir + "/" + repoOwner + "/" + repoName
	if !exists(fullDir) {
		log.Printf("Cloning repository '%v' to '%v'\n", repoUrl, fullDir)
		repoLocal, err = git.PlainClone(fullDir, false, &cloneOptsMaster)
		if err != nil {
			log.Print(err)
			repoLocal, err = git.PlainClone(fullDir, false, &cloneOptsMain)
			if err != nil {
				log.Print(err)
				return errors.New("error cloning repository")
			}
		}
	} else {
		log.Printf("Found repository '%v' at '%v'\n", repoUrl, fullDir)
		repoLocal, err = git.PlainOpen(fullDir)
		if err != nil {
			log.Print(err)
			return errors.New("error opening repository")
		}
	}

	// check out the proper ref
	worktreeLocal, err := repoLocal.Worktree()
	if err != nil {
		log.Print(err)
		return errors.New("error getting repository worktree")
	}
	// TODO: pull to get new updates to the ref??
	gitRef := plumbing.NewReferenceFromStrings("local", ref)
	err = worktreeLocal.Checkout(&git.CheckoutOptions{Hash: gitRef.Hash()})
	if err != nil {
		log.Print(err)
		return errors.New("error checking out reference")
	}

	log.Printf("Cloning repository '%s' complete", repoUrl)
	return err
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.traefik-plugin-init.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
