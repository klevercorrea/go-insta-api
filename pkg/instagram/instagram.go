package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/klevercorrea/go-insta-api/pkg/redis"
	"github.com/klevercorrea/go-insta-api/pkg/utils"

	"github.com/Davincible/goinsta/v3"
)

// Login logs in to an Instagram account and returns the session as a base64-encoded string
func Login(ctx context.Context, username string, password string) (string, error) {
	key := "session:" + username

	// Check if session exists in Redis
	sessionString, err := redis.GetFromRedis(key)
	if err == nil && sessionString != "" {
		return sessionString, nil
	}

	insta := goinsta.New(username, password)
	err = insta.Login()
	if err != nil {
		return "", err
	}
	sessionString, err = utils.ExportAsBase64String(insta)
	if err != nil {
		return "", err
	}

	// Save session to Redis
	if err := redis.SaveToRedis(key, sessionString); err != nil {
		return "", fmt.Errorf("error saving session to Redis: %w", err)
	}

	return sessionString, nil
}

func GetProfileData(ctx context.Context, insta *goinsta.Instagram, username string) (*goinsta.User, error) {
	profile, err := insta.Profiles.ByName(username)
	if err != nil {
		return nil, fmt.Errorf("error getting profile data: %w", err)
	}

	key := "profile_data"
	val, _ := redis.GetFromRedis(key)

	if val == "" {
		profileData := map[string]int{
			"followers": profile.FollowerCount,
			"following": profile.FollowingCount,
		}
		profileDataJSON, _ := json.Marshal(profileData)
		if err := redis.SaveToRedis(key, string(profileDataJSON)); err != nil {
			return nil, fmt.Errorf("error saving profile data to Redis: %w", err)
		}
		log.Println("Saved to Redis:")
	} else {
		log.Println("Data from Redis:")
	}

	log.Println("Followers:", profile.FollowerCount)
	log.Println("Following:", profile.FollowingCount)

	return profile, nil
}

func GetFollowers(ctx context.Context, insta *goinsta.Instagram, totalFollowers int) error {
	key := "followers"
	var allFollowers []string
	totalProcessed := 0
	maxRetries := 3

	// Check if followers exist in Redis
	followersExist, err := redis.KeyExists(key)
	if err != nil {
		return fmt.Errorf("error checking followers in Redis: %w", err)
	}
	if followersExist {
		// Get followers from Redis
		followersJSON, err := redis.GetFromRedis(key)
		if err != nil {
			return fmt.Errorf("error getting followers from Redis: %w", err)
		}
		if err := json.Unmarshal([]byte(followersJSON), &allFollowers); err != nil {
			return fmt.Errorf("error unmarshalling followers from Redis: %w", err)
		}
		totalProcessed = len(allFollowers)
		fmt.Printf("Found %d followers in Redis.\n", totalProcessed)
	}

	// Get followers from API if not found in Redis
	if totalProcessed < totalFollowers {
		followers := insta.Account.Followers("")
		for retries := 0; retries < maxRetries; retries++ {
			for followers.Next() {
				for _, user := range followers.Users {
					if totalProcessed >= len(allFollowers) {
						fmt.Printf("Follower %d of %d: %s\n", totalProcessed+1, totalFollowers, user.Username)
						allFollowers = append(allFollowers, user.Username)
					}
					totalProcessed++
					if totalProcessed == totalFollowers {
						break
					}
				}
				if totalProcessed == totalFollowers {
					break
				}
				// Save the current page of followers to Redis
				followersJSON, _ := json.Marshal(allFollowers)
				if err := redis.SaveToRedis(key, string(followersJSON)); err != nil {
					return fmt.Errorf("error saving followers to Redis: %w", err)
				}
				allFollowers = nil

				// Sleep for a while
				sleepDuration := utils.RandomTime(30, 100)
				utils.Countdown(int(sleepDuration.Seconds()))
			}
			if followers.Error() == nil {
				break
			}
			fmt.Printf("Error getting followers: %s\n", followers.Error())
			if retries == maxRetries-1 {
				return fmt.Errorf("error getting followers: %s", followers.Error())
			}
			fmt.Printf("Retrying in %d seconds...\n", 2*retries)
			time.Sleep(2 * time.Duration(retries) * time.Second)
		}
	}

	// Save the final page of followers to Redis
	followersJSON, _ := json.Marshal(allFollowers)
	if err := redis.SaveToRedis(key, string(followersJSON)); err != nil {
		return fmt.Errorf("error saving followers to Redis: %w", err)
	}

	if len(allFollowers) != totalFollowers {
		return fmt.Errorf("expected to find %d followers, but found %d instead", totalFollowers, len(allFollowers))
	}

	return nil
}
