package shuffle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	//"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func executeCloudAction(action CloudSyncJob, apikey string) error {
	data, err := json.Marshal(action)
	if err != nil {
		log.Printf("Failed cloud webhook action marshalling: %s", err)
		return err
	}

	client := &http.Client{}
	syncUrl := fmt.Sprintf("https://shuffler.io/api/v1/cloud/sync/handle_action")
	req, err := http.NewRequest(
		"POST",
		syncUrl,
		bytes.NewBuffer(data),
	)

	req.Header.Add("Authorization", fmt.Sprintf(`Bearer %s`, apikey))
	newresp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(newresp.Body)
	if err != nil {
		return err
	}

	type Result struct {
		Success bool   `json:"success"`
		Reason  string `json:"reason"`
	}

	//log.Printf("Data: %s", string(respBody))
	responseData := Result{}
	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return err
	}

	if !responseData.Success {
		return errors.New(fmt.Sprintf("Cloud error from Shuffler: %s", responseData.Reason))
	}

	return nil
}

func HandleAlgoliaAppSearch(ctx context.Context, appname string) (string, error) {
	algoliaClient := os.Getenv("ALGOLIA_CLIENT")
	algoliaSecret := os.Getenv("ALGOLIA_SECRET")
	if len(algoliaClient) == 0 || len(algoliaSecret) == 0 {
		log.Printf("[WARNING] ALGOLIA_CLIENT or ALGOLIA_SECRET not defined")
		return "", errors.New("Algolia keys not defined")
	}

	algClient := search.NewClient(algoliaClient, algoliaSecret)
	algoliaIndex := algClient.InitIndex("appsearch")
	appname = strings.TrimSpace(strings.ToLower(strings.Replace(appname, "_", " ", -1)))
	res, err := algoliaIndex.Search(appname)
	if err != nil {
		log.Printf("[WARNING] Failed searching Algolia: %s", err)
		return "", err
	}

	var newRecords []AlgoliaSearchApp
	err = res.UnmarshalHits(&newRecords)
	if err != nil {
		log.Printf("[WARNING] Failed unmarshaling from Algolia: %s", err)
		return "", err
	}

	log.Printf("[INFO] Algolia hits for %s: %d", appname, len(newRecords))
	for _, newRecord := range newRecords {
		newApp := strings.TrimSpace(strings.ToLower(strings.Replace(newRecord.Name, "_", " ", -1)))
		if newApp == appname {
			return newRecord.ObjectID, nil
		}
	}

	return "", nil
}

func HandleAlgoliaWorkflowSearchByApp(ctx context.Context, appname string) ([]AlgoliaSearchWorkflow, error) {
	algoliaClient := os.Getenv("ALGOLIA_CLIENT")
	algoliaSecret := os.Getenv("ALGOLIA_SECRET")
	if len(algoliaClient) == 0 || len(algoliaSecret) == 0 {
		log.Printf("[WARNING] ALGOLIA_CLIENT or ALGOLIA_SECRET not defined")
		return []AlgoliaSearchWorkflow{}, errors.New("Algolia keys not defined")
	}

	algClient := search.NewClient(algoliaClient, algoliaSecret)
	algoliaIndex := algClient.InitIndex("workflows")

	appSearch := fmt.Sprintf("%s", appname)
	res, err := algoliaIndex.Search(appSearch)
	if err != nil {
		log.Printf("[WARNING] Failed app searching Algolia for creators: %s", err)
		return []AlgoliaSearchWorkflow{}, err
	}

	var newRecords []AlgoliaSearchWorkflow
	err = res.UnmarshalHits(&newRecords)
	if err != nil {
		log.Printf("[WARNING] Failed unmarshaling from Algolia with app creators: %s", err)
		return []AlgoliaSearchWorkflow{}, err
	}
	//log.Printf("[INFO] Algolia hits for %s: %d", appSearch, len(newRecords))

	allRecords := []AlgoliaSearchWorkflow{}
	for _, newRecord := range newRecords {
		allRecords = append(allRecords, newRecord)

	}

	return allRecords, nil
}

func HandleAlgoliaWorkflowSearchByUser(ctx context.Context, userId string) ([]AlgoliaSearchWorkflow, error) {
	algoliaClient := os.Getenv("ALGOLIA_CLIENT")
	algoliaSecret := os.Getenv("ALGOLIA_SECRET")
	if len(algoliaClient) == 0 || len(algoliaSecret) == 0 {
		log.Printf("[WARNING] ALGOLIA_CLIENT or ALGOLIA_SECRET not defined")
		return []AlgoliaSearchWorkflow{}, errors.New("Algolia keys not defined")
	}

	algClient := search.NewClient(algoliaClient, algoliaSecret)
	algoliaIndex := algClient.InitIndex("workflows")

	appSearch := fmt.Sprintf("%s", userId)
	res, err := algoliaIndex.Search(appSearch)
	if err != nil {
		log.Printf("[WARNING] Failed app searching Algolia for creators: %s", err)
		return []AlgoliaSearchWorkflow{}, err
	}

	var newRecords []AlgoliaSearchWorkflow
	err = res.UnmarshalHits(&newRecords)
	if err != nil {
		log.Printf("[WARNING] Failed unmarshaling from Algolia with app creators: %s", err)
		return []AlgoliaSearchWorkflow{}, err
	}
	//log.Printf("[INFO] Algolia hits for %s: %d", appSearch, len(newRecords))

	allRecords := []AlgoliaSearchWorkflow{}
	for _, newRecord := range newRecords {
		allRecords = append(allRecords, newRecord)

	}

	return allRecords, nil
}

func HandleAlgoliaAppSearchByUser(ctx context.Context, userId string) ([]AlgoliaSearchApp, error) {
	algoliaClient := os.Getenv("ALGOLIA_CLIENT")
	algoliaSecret := os.Getenv("ALGOLIA_SECRET")
	if len(algoliaClient) == 0 || len(algoliaSecret) == 0 {
		log.Printf("[WARNING] ALGOLIA_CLIENT or ALGOLIA_SECRET not defined")
		return []AlgoliaSearchApp{}, errors.New("Algolia keys not defined")
	}

	algClient := search.NewClient(algoliaClient, algoliaSecret)
	algoliaIndex := algClient.InitIndex("appsearch")

	appSearch := fmt.Sprintf("%s", userId)
	res, err := algoliaIndex.Search(appSearch)
	if err != nil {
		log.Printf("[WARNING] Failed app searching Algolia for creators: %s", err)
		return []AlgoliaSearchApp{}, err
	}

	var newRecords []AlgoliaSearchApp
	err = res.UnmarshalHits(&newRecords)
	if err != nil {
		log.Printf("[WARNING] Failed unmarshaling from Algolia with app creators: %s", err)
		return []AlgoliaSearchApp{}, err
	}
	//log.Printf("[INFO] Algolia hits for %s: %d", appSearch, len(newRecords))

	allRecords := []AlgoliaSearchApp{}
	for _, newRecord := range newRecords {
		newAppName := strings.TrimSpace(strings.Replace(newRecord.Name, "_", " ", -1))
		newRecord.Name = newAppName
		allRecords = append(allRecords, newRecord)

	}

	return allRecords, nil
}

func HandleAlgoliaCreatorSearch(ctx context.Context, username string) (AlgoliaSearchCreator, error) {
	cacheKey := fmt.Sprintf("algolia_creator_%s", username)
	searchCreator := AlgoliaSearchCreator{}
	cache, err := GetCache(ctx, cacheKey)
	if err == nil {
		cacheData := []byte(cache.([]uint8))
		//log.Printf("CACHE: %d", len(cacheData))
		//log.Printf("CACHEDATA: %#v", cacheData)
		err = json.Unmarshal(cacheData, &searchCreator)
		if err == nil {
			return searchCreator, nil
		}
	}

	algoliaClient := os.Getenv("ALGOLIA_CLIENT")
	algoliaSecret := os.Getenv("ALGOLIA_SECRET")
	if len(algoliaClient) == 0 || len(algoliaSecret) == 0 {
		log.Printf("[WARNING] ALGOLIA_CLIENT or ALGOLIA_SECRET not defined")
		return searchCreator, errors.New("Algolia keys not defined")
	}

	algClient := search.NewClient(algoliaClient, algoliaSecret)
	algoliaIndex := algClient.InitIndex("creators")
	res, err := algoliaIndex.Search(username)
	if err != nil {
		log.Printf("[WARNING] Failed searching Algolia creators: %s", err)
		return searchCreator, err
	}

	var newRecords []AlgoliaSearchCreator
	err = res.UnmarshalHits(&newRecords)
	if err != nil {
		log.Printf("[WARNING] Failed unmarshaling from Algolia creators: %s", err)
		return searchCreator, err
	}

	//log.Printf("RECORDS: %d", len(newRecords))
	foundUser := AlgoliaSearchCreator{}
	for _, newRecord := range newRecords {
		if strings.ToLower(newRecord.Username) == strings.ToLower(username) || newRecord.ObjectID == username || ArrayContainsLower(newRecord.Synonyms, username) {
			foundUser = newRecord
			break
		}
	}

	// Handling search within a workflow, and in the future, within apps
	if len(foundUser.ObjectID) == 0 {
		if len(username) == 36 {
			// Check workflows
			algoliaIndex := algClient.InitIndex("workflows")
			res, err := algoliaIndex.Search(username)
			if err != nil {
				log.Printf("[WARNING] Failed searching Algolia creator workflow: %s", err)
				return searchCreator, err
			}

			var newRecords []AlgoliaSearchWorkflow
			err = res.UnmarshalHits(&newRecords)
			if err != nil {
				log.Printf("[WARNING] Failed unmarshaling from Algolia creator workflow: %s", err)

				if len(newRecords) > 0 && len(newRecords[0].ObjectID) > 0 {
					log.Printf("[INFO] Workflow search ID: %#v", newRecords[0].ObjectID)
				} else {
					return searchCreator, err
				}
			}

			//log.Printf("[DEBUG] Got %d records for workflow sub", len(newRecords))
			if len(newRecords) == 1 {
				if len(newRecords[0].Creator) > 0 && username != newRecords[0].Creator {
					foundCreator, err := HandleAlgoliaCreatorSearch(ctx, newRecords[0].Creator)
					if err != nil {
						return searchCreator, err
					}

					foundUser = foundCreator
				} else {
					return searchCreator, errors.New("User not found")
				}
			} else {
				return searchCreator, errors.New("User not found")
			}
		} else {
			return searchCreator, errors.New("User not found")
		}
	}

	if project.CacheDb {
		data, err := json.Marshal(foundUser)
		if err != nil {
			return foundUser, nil
		}

		err = SetCache(ctx, cacheKey, data)
		if err != nil {
			log.Printf("[WARNING] Failed updating algolia username cache: %s", err)
		}
	}

	return foundUser, nil
}

func HandleAlgoliaCreatorUpload(ctx context.Context, user User, overwrite bool) (string, error) {
	algoliaClient := os.Getenv("ALGOLIA_CLIENT")
	algoliaSecret := os.Getenv("ALGOLIA_SECRET")
	if len(algoliaClient) == 0 || len(algoliaSecret) == 0 {
		log.Printf("[WARNING] ALGOLIA_CLIENT or ALGOLIA_SECRET not defined")
		return "", errors.New("Algolia keys not defined")
	}

	algClient := search.NewClient(algoliaClient, algoliaSecret)
	algoliaIndex := algClient.InitIndex("creators")
	res, err := algoliaIndex.Search(user.Id)
	if err != nil {
		log.Printf("[WARNING] Failed searching Algolia creators: %s", err)
		return "", err
	}

	var newRecords []AlgoliaSearchCreator
	err = res.UnmarshalHits(&newRecords)
	if err != nil {
		log.Printf("[WARNING] Failed unmarshaling from Algolia creators: %s", err)
		return "", err
	}

	//log.Printf("RECORDS: %d", len(newRecords))
	for _, newRecord := range newRecords {
		if newRecord.ObjectID == user.Id {
			log.Printf("[INFO] Object %s already exists in Algolia", user.Id)

			if overwrite {
				break
			} else {
				return user.Id, errors.New("User ID already exists!")
			}
		}
	}

	timeNow := int64(time.Now().Unix())
	records := []AlgoliaSearchCreator{
		AlgoliaSearchCreator{
			ObjectID:   user.Id,
			TimeEdited: timeNow,
			Image:      user.PublicProfile.GithubAvatar,
			Username:   user.PublicProfile.GithubUsername,
		},
	}

	_, err = algoliaIndex.SaveObjects(records)
	if err != nil {
		log.Printf("[WARNING] Algolia Object put err: %s", err)
		return "", err
	}

	log.Printf("[INFO] SUCCESSFULLY UPLOADED creator %s with ID %s TO ALGOLIA!", user.Username, user.Id)
	return user.Id, nil
}

// Shitty temorary system
// Adding schedule to run over with another algorithm
// as well as this one, as to increase priority based on popularity:
// searches, clicks & conversions (CTR)
func GetWorkflowPriority(workflow Workflow) int {
	prio := 0
	if len(workflow.Tags) > 2 {
		prio += 1
	}

	if len(workflow.Name) > 5 {
		prio += 1
	}

	if len(workflow.Description) > 100 {
		prio += 1
	}

	if len(workflow.WorkflowType) > 0 {
		prio += 1
	}

	if len(workflow.UsecaseIds) > 0 {
		prio += 3
	}

	if len(workflow.Comments) >= 2 {
		prio += 2
	}

	return prio
}

func handleAlgoliaWorkflowUpdate(ctx context.Context, workflow Workflow) (string, error) {
	log.Printf("[INFO] Should try to UPLOAD the Workflow to Algolia")

	algoliaClient := os.Getenv("ALGOLIA_CLIENT")
	algoliaSecret := os.Getenv("ALGOLIA_SECRET")
	if len(algoliaClient) == 0 || len(algoliaSecret) == 0 {
		log.Printf("[WARNING] ALGOLIA_CLIENT or ALGOLIA_SECRET not defined")
		return "", errors.New("Algolia keys not defined")
	}

	algClient := search.NewClient(algoliaClient, algoliaSecret)
	algoliaIndex := algClient.InitIndex("workflows")

	//res, err := algoliaIndex.Search("%s", api.ID)
	res, err := algoliaIndex.Search(workflow.ID)
	if err != nil {
		log.Printf("[WARNING] Failed searching Algolia: %s", err)
		return "", err
	}

	var newRecords []AlgoliaSearchWorkflow
	err = res.UnmarshalHits(&newRecords)
	if err != nil {
		log.Printf("[WARNING] Failed unmarshaling from Algolia workflow upload: %s", err)
		return "", err
	}

	found := false
	record := AlgoliaSearchWorkflow{}
	for _, newRecord := range newRecords {
		if newRecord.ObjectID == workflow.ID {
			log.Printf("[INFO] Workflow Object %s already exists in Algolia", workflow.ID)
			record = newRecord
			found = true
			break
		}
	}

	if !found {
		return "", errors.New(fmt.Sprintf("Couldn't find public workflow for ID %s", workflow.ID))
	}

	record.TimeEdited = int64(time.Now().Unix())
	categories := []string{}
	actions := []string{}
	triggers := []string{}
	actionRefs := []ActionReference{}
	for _, action := range workflow.Actions {
		if !ArrayContains(actions, action.AppName) {
			// Using this API as the original is kinda stupid
			foundApps, err := HandleAlgoliaAppSearchByUser(ctx, action.AppName)
			if err == nil && len(foundApps) > 0 {
				actionRefs = append(actionRefs, ActionReference{
					Name:       foundApps[0].Name,
					Id:         foundApps[0].ObjectID,
					ImageUrl:   foundApps[0].ImageUrl,
					ActionName: []string{action.Name},
				})
			}

			actions = append(actions, action.AppName)
		} else {
			for refIndex, ref := range actionRefs {
				if ref.Name == action.AppName {
					if !ArrayContains(ref.ActionName, action.Name) {
						actionRefs[refIndex].ActionName = append(actionRefs[refIndex].ActionName, action.Name)
					}
				}
			}
		}
	}

	for _, trigger := range workflow.Triggers {
		if !ArrayContains(triggers, trigger.TriggerType) {
			triggers = append(triggers, trigger.TriggerType)
		}
	}

	if workflow.WorkflowType != "" {
		record.Type = workflow.WorkflowType
	}

	record.Name = workflow.Name
	record.Description = workflow.Description
	record.UsecaseIds = workflow.UsecaseIds
	record.Triggers = triggers
	record.Actions = actions
	record.TriggerAmount = len(triggers)
	record.ActionAmount = len(actions)
	record.Tags = workflow.Tags
	record.Categories = categories
	record.ActionReferences = actionRefs

	record.Priority = GetWorkflowPriority(workflow)

	records := []AlgoliaSearchWorkflow{
		record,
	}

	//log.Printf("[WARNING] Returning before upload with data %#v", records)
	//return records[0].ObjectID, nil
	//return "", errors.New("Not prepared yet!")

	_, err = algoliaIndex.SaveObjects(records)
	if err != nil {
		log.Printf("[WARNING] Algolia Object update err: %s", err)
		return "", err
	}

	return workflow.ID, nil
}

// Returns an error if the users' org is over quota
func ValidateExecutionUsage(ctx context.Context, orgId string) (*Org, error) {
	log.Printf("[DEBUG] Validating usage of org %#v", orgId)

	org, err := GetOrg(ctx, orgId)
	if err != nil {
		return org, errors.New(fmt.Sprintf("Failed getting the organization %s: %s", orgId, err))
	}

	//if len(org.ChildOrgs) > 0 || len(org.ManagerOrgs) > 0 {
	//	log.Printf("[DEBUG] Execution for %s is allowed due to being a child-or parent org. This is only accessible to customers. We're not force-stopping them.")
	//	return org, nil
	//}
	//parentOrg.ChildOrgs = append(parentOrg.ChildOrgs, OrgMini{

	info, err := GetOrgStatistics(ctx, orgId)
	if err == nil {
		org.SyncFeatures.AppExecutions.Usage = info.MonthlyAppExecutions
		if org.SyncFeatures.AppExecutions.Limit <= 10000 {
			org.SyncFeatures.AppExecutions.Limit = 5000
		} else {
			// FIXME: Not strictly enforcing other limits yet
			// Should just warn our team about them going over
			org.SyncFeatures.AppExecutions.Limit = 15000000000
		}

		// FIXME: When inside this, check if usage should be sent to the user
		if org.SyncFeatures.AppExecutions.Usage > org.SyncFeatures.AppExecutions.Limit {
			return org, errors.New(fmt.Sprintf("You are above your limited usage of app executions this month (%d / %d) when running with triggers. Contact support@shuffler.io or the live chat to extend this.", org.SyncFeatures.AppExecutions.Usage, org.SyncFeatures.AppExecutions.Limit))
		}

		return org, nil
	} else {
		log.Printf("[WARNING] Failed finding executions for org %s (%s)", org.Name, org.Id)
	}

	return org, nil
}
