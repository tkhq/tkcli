package cmd

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
	Short: "activities interacts with the API activities",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
		LoadClient()
	},
}

var activitiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "list returns the set of activities for an organization",
	Run: func(cmd *cobra.Command, args []string) {
		activitiesFilter := make([]models.V1ActivityStatus, len(activitiesListStatus))

		for i, s := range activitiesListStatus {
			if s == Help {
				Output(models.ActivityStatuses())

				return
			}

			if s == "all" {
				activitiesFilter = models.ActivityStatuses()

				break
			}

			switch s {
			case "created":
				activitiesFilter[i] = models.V1ActivityStatusACTIVITYSTATUSCREATED
			case "pending":
				activitiesFilter[i] = models.V1ActivityStatusACTIVITYSTATUSPENDING
			case "completed":
				activitiesFilter[i] = models.V1ActivityStatusACTIVITYSTATUSCOMPLETED
			case "failed":
				activitiesFilter[i] = models.V1ActivityStatusACTIVITYSTATUSFAILED
			case "consensus":
				activitiesFilter[i] = models.V1ActivityStatusACTIVITYSTATUSCONSENSUSNEEDED
			case "consensus_needed":
				activitiesFilter[i] = models.V1ActivityStatusACTIVITYSTATUSCONSENSUSNEEDED
			case "rejected":
				activitiesFilter[i] = models.V1ActivityStatusACTIVITYSTATUSREJECTED
			default:
				activitiesFilter[i] = models.V1ActivityStatus(s)
			}
		}

		params := activities.NewPublicAPIServiceGetActivitiesParams().WithDefaults().WithBody(&models.V1GetActivitiesRequest{
			FilterByStatus: activitiesFilter,
			OrganizationID: &Organization,
		})

		res, err := APIClient.V0().Activities.PublicAPIServiceGetActivities(params, APIClient.Authenticator)
		if err != nil {
			OutputError(err)
		}

		Output(res.GetPayload().Activities)
	},
}

var activitiesGetCmd = &cobra.Command{
	Use:   "get <activity-id>",
	Short: "get returns the details and status of a particular activity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		params := activities.NewPublicAPIServiceGetActivityParams().WithDefaults().WithBody(&models.V1GetActivityRequest{
			ActivityID:     &id,
			OrganizationID: &Organization,
		})

		res, err := APIClient.V0().Activities.PublicAPIServiceGetActivity(params, APIClient.Authenticator)
		if err != nil {
			OutputError(err)
		}

		Output(res.GetPayload().Activity)
	},
}
