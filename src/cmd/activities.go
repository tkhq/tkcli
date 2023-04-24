package cmd

import (
	"github.com/pkg/errors"
	"github.com/tkhq/go-sdk/pkg/api/client/activities"
	"github.com/tkhq/go-sdk/pkg/api/models"

	"github.com/spf13/cobra"
)

var activitiesListStatus []string

func init() {
	activitiesListCmd.Flags().StringSliceVar(&activitiesListStatus, "status", nil, "only include activities whose status is declared in this set")

	rootCmd.AddCommand(activitiesCmd)

	activitiesCmd.AddCommand(activitiesListCmd)
}

var activitiesCmd = &cobra.Command{
	Use:   "activities interacts with the API activities",
	Short: "activities interacts with the API activities",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		LoadKeypair("")
		LoadClient()
	},
}

var activitiesListCmd = &cobra.Command{
	Use:   "list returns the set of activities for an organization",
	Short: "list returns the set of activities for an organization",
	Run: func(cmd *cobra.Command, args []string) {
		activitiesFilter := make([]models.V1ActivityStatus, len(activitiesListStatus))

		for i, s := range activitiesListStatus {
			switch s {
			case "all":
				activitiesFilter = models.ActivityStatuses()

				break
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

		if res.Error() != "" {
			OutputError(errors.Errorf("request failed: %s", res.Error()))
		}

		Output(res.GetPayload().Activities)
	},
}
