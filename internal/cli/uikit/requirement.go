package uikit

import "github.com/spf13/cobra"

func RequireCloudflared(cmd *cobra.Command) *cobra.Command {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[annotationNeedsCloudflared] = annotationEnabled
	return cmd
}

func RequiresCloudflared(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Annotations[annotationNeedsCloudflared] != "" {
			return true
		}
	}
	return false
}
