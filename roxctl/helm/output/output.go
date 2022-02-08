package output

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/image"
	"github.com/stackrox/rox/pkg/buildinfo"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/helm/charts"
	"github.com/stackrox/rox/pkg/images/defaults"
	"github.com/stackrox/rox/roxctl/common/environment"
	"github.com/stackrox/rox/roxctl/common/flags"
	"github.com/stackrox/rox/roxctl/helm/internal/common"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func handleRhacsWarnings(rhacs, imageFlavorProvided bool, logger environment.Logger) {
	if rhacs {
		logger.WarnfLn("'--rhacs' is deprecated, please use '--%s=%s' instead", flags.ImageDefaultsFlagName, defaults.ImageFlavorNameRHACSRelease)
	} else if !imageFlavorProvided {
		logger.WarnfLn("Default image registries have changed. Images will be taken from 'registry.redhat.io'. Specify '--%s=%s' command line argument to use images from 'stackrox.io' registries.", flags.ImageDefaultsFlagName, defaults.ImageFlavorNameStackRoxIORelease)
	}
}

func getMetaValues(flavorName string, imageFlavorProvided, rhacs, release bool, logger environment.Logger) (*charts.MetaValues, error) {
	handleRhacsWarnings(rhacs, imageFlavorProvided, logger)
	if rhacs {
		if imageFlavorProvided {
			return nil, fmt.Errorf("%w: flag '--rhacs' is deprecated and must not be used together with '--%s'. Remove '--rhacs' flag and specify only '--%s'", errorhelpers.ErrInvalidArgs, flags.ImageDefaultsFlagName, flags.ImageDefaultsFlagName)
		}
		flavorName = defaults.ImageFlavorNameRHACSRelease
	}
	imageFlavor, err := defaults.GetImageFlavorByName(flavorName, release)
	if err != nil {
		return nil, fmt.Errorf("%w: '--%s': %v", errorhelpers.ErrInvalidArgs, flags.ImageDefaultsFlagName, err)
	}
	return charts.GetMetaValuesForFlavor(imageFlavor), nil
}

func outputHelmChart(chartName string, outputDir string, removeOutputDir bool, imageFlavor string, flavorProvided, rhacs bool, logger environment.Logger) error {
	// Lookup chart template prefix.
	chartTemplatePathPrefix := common.ChartTemplates[chartName]
	if chartTemplatePathPrefix == "" {
		return errors.New("unknown chart, see --help for list of supported chart names")
	}

	metaVals, err := getMetaValues(imageFlavor, flavorProvided, rhacs, buildinfo.ReleaseBuild, logger)
	if err != nil {
		return err
	}

	if outputDir == "" {
		outputDir = fmt.Sprintf("./stackrox-%s-chart", chartName)
		fmt.Fprintf(os.Stderr, "No output directory specified, using default directory %q.\n", outputDir)
	}

	if _, err := os.Stat(outputDir); err == nil {
		if removeOutputDir {
			if err := os.RemoveAll(outputDir); err != nil {
				return errors.Wrapf(err, "failed to remove output dir %s", outputDir)
			}
			fmt.Fprintf(os.Stderr, "Removed output directory %s\n", outputDir)
		} else {
			fmt.Fprintf(os.Stderr, "Directory %q already exists, use --remove or select a different directory with --output-dir.\n", outputDir)
			return fmt.Errorf("directory %q already exists", outputDir)
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "failed to check if directory %q exists", outputDir)
	}

	// load image with templates
	templateImage := image.GetDefaultImage()
	if flags.IsDebug() {
		templateImage = flags.GetDebugHelmImage()
	}

	// Load and render template files.
	renderedChartFiles, err := templateImage.LoadAndInstantiateChartTemplate(chartTemplatePathPrefix, metaVals)
	if err != nil {
		return errors.Wrapf(err, "loading and instantiating %s helmtpl", chartName)
	}

	// Write rendered files to output directory.
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return errors.Wrapf(err, "unable to create folder %q", outputDir)
	}
	for _, f := range renderedChartFiles {
		if err := writeFile(f, outputDir); err != nil {
			return errors.Wrapf(err, "error writing file %q", f.Name)
		}
	}
	fmt.Fprintf(os.Stderr, "Written Helm chart %s to directory %q.\n", chartName, outputDir)

	return nil
}

// Command for writing Helm Chart.
func Command(cliEnvironment environment.Environment) *cobra.Command {
	var outputDir string
	var removeOutputDir bool
	var rhacs bool
	var imageFlavor string

	c := &cobra.Command{
		Use: fmt.Sprintf("output <%s>", common.PrettyChartNameList),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("incorrect number of arguments, see --help for usage information")
			}
			chartName := args[0]
			flavorProvided := cmd.Flags().Changed(flags.ImageDefaultsFlagName)
			return outputHelmChart(chartName, outputDir, removeOutputDir, imageFlavor, flavorProvided, rhacs, cliEnvironment.Logger())
		},
	}
	c.PersistentFlags().StringVar(&outputDir, "output-dir", "", "path to the output directory for Helm chart (default: './stackrox-<chart name>-chart')")
	c.PersistentFlags().BoolVar(&removeOutputDir, "remove", false, "remove the output directory if it already exists")
	c.PersistentFlags().BoolVar(&rhacs, "rhacs", false,
		fmt.Sprintf("render RHACS chart flavor (deprecated: use '--%s=%s' instead)", flags.ImageDefaultsFlagName, defaults.ImageFlavorNameRHACSRelease))

	if !buildinfo.ReleaseBuild {
		flags.AddHelmChartDebugSetting(c)
	}
	flags.AddImageDefaults(c.PersistentFlags(), &imageFlavor)
	return c
}

func writeFile(file *loader.BufferedFile, destDir string) error {
	outputPath := filepath.Join(destDir, filepath.FromSlash(file.Name))
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return errors.Wrapf(err, "creating directory for file %q", file.Name)
	}

	perms := os.FileMode(0644)
	if filepath.Ext(file.Name) == ".sh" {
		perms = os.FileMode(0755)
	}
	return os.WriteFile(outputPath, file.Data, perms)
}
