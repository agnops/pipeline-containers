package main

import (
	"flag"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"log"
	"os"
	"time"
)

type UpgradeOptions struct {
	Namespace    string
	Timeout      time.Duration
	Wait         bool
	DisableHooks bool
	DryRun       bool
	ClientOnly   bool
	Force        bool
	ResetValues  bool
	SkipCRDs     bool
	ReuseValues  bool
	Recreate     bool
	MaxHistory   int
	Atomic       bool
	Description	 string
}

func UpgradeFromPath(chartPath string, releaseName string, valuesFile string, opts UpgradeOptions) (*release.Release, error) {
	settings := cli.New()

	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), opts.Namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	// Load the chart from the given path, this also ensures that
	// all chart dependencies are present
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	// Read and set values
	val, err := chartutil.ReadValuesFile(valuesFile)
	if err != nil && len(valuesFile) != 0 {
		//return nil, err
		log.Println(err.Error())
	}

	var res *release.Release

	install := action.NewInstall(actionConfig)
	installOptions(opts).configure(install, releaseName)
	res, err = install.Run(chartRequested, val.AsMap())

	if err != nil && err.Error() == "cannot re-use a name that is still in use" {
		upgrade := action.NewUpgrade(actionConfig)
		upgradeOptions(opts).configure(upgrade)
		res, err = upgrade.Run(releaseName, chartRequested, val.AsMap())
	}

	if err != nil {
		return nil, err
	}
	return res, err
}

type installOptions UpgradeOptions

func (opts installOptions) configure(action *action.Install, releaseName string) {
	action.Namespace = opts.Namespace
	action.ReleaseName = releaseName
	action.Atomic = opts.Atomic
	action.DisableHooks = opts.DisableHooks
	action.DryRun = opts.DryRun
	action.ClientOnly = opts.ClientOnly
	action.Timeout = opts.Timeout
	action.Wait = opts.Wait
	action.SkipCRDs = opts.SkipCRDs
	action.Description = opts.Description
}

type upgradeOptions UpgradeOptions

func (opts upgradeOptions) configure(action *action.Upgrade) {
	action.Namespace = opts.Namespace
	action.Atomic = opts.Atomic
	action.DisableHooks = opts.DisableHooks
	action.DryRun = opts.DryRun
	action.Force = opts.Force
	action.MaxHistory = opts.MaxHistory
	action.ResetValues = opts.ResetValues
	action.ReuseValues = opts.ReuseValues
	action.Timeout = opts.Timeout
	action.Wait = opts.Wait
	action.Description = opts.Description
}

func main() {
	var namespace string
	var releaseName string
	var chartPath string
	var valuesFile string
	var description string

	flag.StringVar(&namespace, "namespace", "default", "kubernetes namespace")
	flag.StringVar(&releaseName, "releaseName", "", "helm release")
	flag.StringVar(&chartPath, "chartPath", "", "helm chart full path")
	flag.StringVar(&valuesFile, "valuesFile", "", "helm chart values file")
	flag.StringVar(&description, "description", "", "helm chart description")
	flag.Parse()

	var opts UpgradeOptions

	opts.Wait = true

	opts.Namespace = namespace
	opts.Force = true

	rel, err := UpgradeFromPath(chartPath, releaseName, valuesFile, opts)

	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	} else {
		log.Println(rel.Manifest)
	}
}