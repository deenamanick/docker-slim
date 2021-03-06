package xray

import (
	"fmt"

	"github.com/docker-slim/docker-slim/internal/app/master/commands"

	"github.com/urfave/cli"
)

const (
	Name  = "xray"
	Usage = "Shows what's inside of your container image and reverse engineers its Dockerfile"
	Alias = "x"
)

var CLI = cli.Command{
	Name:    Name,
	Aliases: []string{Alias},
	Usage:   Usage,
	Flags: []cli.Flag{
		commands.Cflag(commands.FlagTarget),
		commands.Cflag(commands.FlagPull),
		commands.Cflag(commands.FlagShowPullLogs),
		cflag(FlagChanges),
		cflag(FlagChangesOutput),
		cflag(FlagLayer),
		cflag(FlagAddImageManifest),
		cflag(FlagAddImageConfig),
		cflag(FlagLayerChangesMax),
		cflag(FlagAllChangesMax),
		cflag(FlagAddChangesMax),
		cflag(FlagModifyChangesMax),
		cflag(FlagDeleteChangesMax),
		cflag(FlagChangePath),
		cflag(FlagChangeData),
		commands.Cflag(commands.FlagRemoveFileArtifacts),
	},
	Action: func(ctx *cli.Context) error {
		commands.ShowCommunityInfo()
		targetRef := ctx.String(commands.FlagTarget)

		if targetRef == "" {
			if len(ctx.Args()) < 1 {
				fmt.Printf("docker-slim[%s]: missing image ID/name...\n\n", Name)
				cli.ShowCommandHelp(ctx, Name)
				return nil
			} else {
				targetRef = ctx.Args().First()
			}
		}

		gcvalues, err := commands.GlobalCommandFlagValues(ctx)
		if err != nil {
			return err
		}

		doPull := ctx.Bool(commands.FlagPull)
		doShowPullLogs := ctx.Bool(commands.FlagShowPullLogs)

		changes, err := parseChangeTypes(ctx.StringSlice(FlagChanges))
		if err != nil {
			fmt.Printf("docker-slim[%s]: invalid change types: %v\n", Name, err)
			return err
		}

		changesOutputs, err := parseChangeOutputTypes(ctx.StringSlice(FlagChangesOutput))
		if err != nil {
			fmt.Printf("docker-slim[%s]: invalid change output types: %v\n", Name, err)
			return err
		}

		layers, err := commands.ParseTokenSet(ctx.StringSlice(FlagLayer))
		if err != nil {
			fmt.Printf("docker-slim[%s]: invalid layer selectors: %v\n", Name, err)
			return err
		}

		layerChangesMax := ctx.Int(FlagLayerChangesMax)
		allChangesMax := ctx.Int(FlagAllChangesMax)
		addChangesMax := ctx.Int(FlagAddChangesMax)
		modifyChangesMax := ctx.Int(FlagModifyChangesMax)
		deleteChangesMax := ctx.Int(FlagDeleteChangesMax)

		changePaths := ctx.StringSlice(FlagChangePath)
		changeDataPatterns := ctx.StringSlice(FlagChangeData)

		doAddImageManifest := ctx.Bool(FlagAddImageManifest)
		doAddImageConfig := ctx.Bool(FlagAddImageConfig)
		doRmFileArtifacts := ctx.Bool(commands.FlagRemoveFileArtifacts)

		xc := commands.NewExecutionContext(Name)

		OnCommand(
			xc,
			gcvalues,
			targetRef,
			doPull,
			doShowPullLogs,
			changes,
			changesOutputs,
			layers,
			layerChangesMax,
			allChangesMax,
			addChangesMax,
			modifyChangesMax,
			deleteChangesMax,
			changePaths,
			changeDataPatterns,
			doAddImageManifest,
			doAddImageConfig,
			doRmFileArtifacts)

		commands.ShowCommunityInfo()
		return nil
	},
}

func parseChangeTypes(values []string) (map[string]struct{}, error) {
	changes := map[string]struct{}{}
	if len(values) == 0 {
		values = append(values, "all")
	}

	for _, item := range values {
		switch item {
		case "none":
			return nil, nil
		case "all":
			changes["delete"] = struct{}{}
			changes["modify"] = struct{}{}
			changes["add"] = struct{}{}
		case "delete":
			changes["delete"] = struct{}{}
		case "modify":
			changes["modify"] = struct{}{}
		case "add":
			changes["add"] = struct{}{}
		}
	}

	return changes, nil
}

func parseChangeOutputTypes(values []string) (map[string]struct{}, error) {
	outputs := map[string]struct{}{}
	if len(values) == 0 {
		values = append(values, "all")
	}

	for _, item := range values {
		switch item {
		case "all":
			outputs["report"] = struct{}{}
			outputs["console"] = struct{}{}
		case "report":
			outputs["report"] = struct{}{}
		case "console":
			outputs["console"] = struct{}{}
		}
	}

	return outputs, nil
}
