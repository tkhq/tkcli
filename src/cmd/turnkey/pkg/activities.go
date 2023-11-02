package pkg

import (
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/activities"
	"github.com/tkhq/go-sdk/pkg/api/models"
)

var activitiesListStatus []string

func init() {
	activitiesListCmd.Flags().StringSliceVar(&activitiesListStatus, "status", nil, "only include activities whose status is declared in this set")

	rootCmd.AddCommand(activitiesCmd)

	activitiesCmd.AddCommand(activitiesListCmd)
	activitiesCmd.AddCommand(activitiesGetCmd)
}

var activitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "Interact with the API activities",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
		LoadClient()
	},
}

var activitiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return the set of activities for an organization",
	Run: func(cmd *cobra.Command, args []string) {
		activitiesFilter := make([]models.ActivityStatus, len(activitiesListStatus))

		for i, s := range activitiesListStatus {
			if s == Help {
				Output(models.ActivityStatusEnum)

				return
			}

			if s == "all" {
				activitiesFilter = models.ActivityStatusEnum
				break
			}

			switch s {
			case "created":
				activitiesFilter[i] = models.ActivityStatusCreated
			case "pending":
				activitiesFilter[i] = models.ActivityStatusPending
			case "completed":
				activitiesFilter[i] = models.ActivityStatusCompleted
			case "failed":
				activitiesFilter[i] = models.ActivityStatusFailed
			case "consensus":
				activitiesFilter[i] = models.ActivityStatusConsensusNeeded
			case "consensus_needed":
				activitiesFilter[i] = models.ActivityStatusConsensusNeeded
			case "rejected":
				activitiesFilter[i] = models.ActivityStatusRejected
			default:
				activitiesFilter[i] = models.ActivityStatus(s)
			}
		}

		params := activities.NewGetActivitiesParams().WithDefaults().WithBody(&models.GetActivitiesRequest{
			FilterByStatus: activitiesFilter,
			OrganizationID: &Organization,
		})

		res, err := APIClient.V0().Activities.GetActivities(params, APIClient.Authenticator)
		if err != nil {
			OutputError(err)
		}

		Output(res.GetPayload().Activities)
	},
}

var activitiesGetCmd = &cobra.Command{
	Use:   "get <activity-id>",
	Short: "Return the details and status of a particular activity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		params := activities.NewGetActivityParams().WithDefaults().WithBody(&models.GetActivityRequest{
			ActivityID:     &id,
			OrganizationID: &Organization,
		})

		res, err := APIClient.V0().Activities.GetActivity(params, APIClient.Authenticator)
		if err != nil {
			OutputError(err)
		}

		Output(res.GetPayload().Activity)
	},
}
