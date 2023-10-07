package server

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"math"
	"peargram/bookmarks"
	"peargram/likes"
	"peargram/models"
	"strings"
	"time"

	"github.com/gin-contrib/multitemplate"
)

const PROFILE_COLS = 3

func Renderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	funcs := template.FuncMap{
		"hasLiked": func(username string, postID int) bool {
			return likes.HasLiked(username, postID)
		},
		"hasBookmarked": func(username string, postID int) bool {
			return bookmarks.HasBookmarked(username, postID)
		},
		"isLastPost": func(post models.Post, posts []models.Post) bool {
			if posts[len(posts)-1] == post {
				return true
			}
			return false
		},
		"isEmpty": func(content string) bool {
			return content == ""
		},
		"orderedPosts": func(posts []models.Post) [][3]models.Post {
			var rows [][3]models.Post
			rowAmount := int(math.Ceil(float64(len(posts)) / PROFILE_COLS))

			for i := 0; i < rowAmount; i++ {
				var post1, post2, post3 models.Post

				post1 = posts[i*PROFILE_COLS+0]

				if i*PROFILE_COLS+1 < len(posts) {
					post2 = posts[i*PROFILE_COLS+1]
				}

				if i*PROFILE_COLS+2 < len(posts) {
					post3 = posts[i*PROFILE_COLS+2]
				}

				row := [3]models.Post{post1, post2, post3}
				rows = append(rows, row)
			}
			return rows
		},
		"emptyPostAmount": func(posts []models.Post) int {
			rowAmount := int(math.Ceil(float64(len(posts)) / PROFILE_COLS))
			remainingPosts := len(posts) - rowAmount*PROFILE_COLS

			return remainingPosts
		},
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"renderImage": func(content string) string {
			base64Img := base64.StdEncoding.EncodeToString([]byte(content))

			return base64Img
		},
		"getTags": func(tagStr string) []string {
			tags := strings.Split(tagStr, ",")
			return tags
		},
		"getPageButtons": func(pageCount int, page int) []int {
			// TODO: Compatibility with N pages
			var pageArr []int
			for i := 1; i <= pageCount; i++ {
				pageArr = append(pageArr, i)
			}

			return pageArr
		},
		"inc": func(views int) int {
			return views + 1
		},
		"renderViews": func(views int) string {
			if views == 1 {
				return "1 view"
			} else {
				return fmt.Sprintf("%d views", views)
			}
		},
		"renderDate": func(unixDateInt int) string {
			var visualAmount int
			var pluralS string

			if unixDateInt == 0 {
				return "never"
			}

			now := time.Now()

			unixDate := time.Unix(int64(unixDateInt), 0)

			diffSeconds := int(math.Floor(now.Sub(unixDate).Seconds()))

			if diffSeconds < 60 {
				visualAmount = diffSeconds
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d second%s ago", visualAmount, pluralS)
			} else if diffSeconds < 3600 {
				visualAmount = int(math.Floor(float64(diffSeconds) / float64(60)))
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d minute%s ago", visualAmount, pluralS)
			} else if diffSeconds < 86400 {
				visualAmount = int(math.Floor(float64(diffSeconds) / float64(3600)))
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d hour%s ago", visualAmount, pluralS)
			} else if diffSeconds < 604800 {
				visualAmount = int(math.Floor(float64(diffSeconds) / float64(86400)))
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d day%s ago", visualAmount, pluralS)
			} else if diffSeconds < 2419200 {
				visualAmount = int(math.Floor(float64(diffSeconds) / float64(604800)))
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d week%s ago, at %d:%d", visualAmount, pluralS, unixDate.Hour(), unixDate.Minute())
			} else {
				return fmt.Sprintf("%s %d, %d", unixDate.Month(), unixDate.Day(), unixDate.Year())
			}
		},
		"groupNotifications": func(notis []models.Notification) []NotificationGroup {
			var group []NotificationGroup
			var addedDays []string

			for _, notification := range notis {
				parsedUnix := time.Unix(int64(notification.Date), 0)
				finalDate := fmt.Sprintf("%02d/%02d/%d", parsedUnix.Day(), parsedUnix.Month(), parsedUnix.Year())

				added := false
				for i := 0; i < len(addedDays); i++ {
					if addedDays[i] == finalDate {
						added = true
						break
					}
				}

				if added {
					for i := 0; i < len(group); i++ {
						if group[i].DayDate == finalDate {
							group[i].Notifications = append(group[i].Notifications, notification)
							break
						}
					}

				} else {
					newGroup := NotificationGroup{DayDate: finalDate, Notifications: []models.Notification{notification}}
					group = append(group, newGroup)

					addedDays = append(addedDays, finalDate)
				}

			}
			return group
		}, "renderDay": func(dayDate string) string { // DayDate format: DD/MM/YYYY : string
			parsedTime, err := time.Parse("02/01/2006", dayDate)
			if err != nil {
				fmt.Println(err)
				return "ERROR_DAY_PARSE"
			}
			diffDays := int(time.Now().Sub(parsedTime).Hours() / 24)

			if diffDays == 0 {
				return "Today"
			}
			if diffDays == 1 {
				return "Yesterday"
			}

			return fmt.Sprintf("%d days ago", diffDays)
			// TODO: Group like this:
			//	- Today
			//	- Yesterday
			//	- This week
			//	- This month
			//	- Earlier

		}, "notificationButtonTarget": func(notification models.Notification) string {
			if notification.Type == models.NOTIFICATION_FOLLOW || notification.Type == models.NOTIFICATION_FOLLOW_REQUEST {
				return "profile"
			}
			return "post"

		}, "notificationButtonDetail": func(notification models.Notification) string {
			if notification.Type == models.NOTIFICATION_FOLLOW || notification.Type == models.NOTIFICATION_FOLLOW_REQUEST {
				return notification.Actor

			}
			return fmt.Sprintf("%d", *notification.Post)
		}, "notificationMessage": func(notificationType int) string {
			switch notificationType {
			case models.NOTIFICATION_POST_LIKE:
				return "has liked your post."
			case models.NOTIFICATION_POST_COMMENT:
				return "has commented in your post."
			case models.NOTIFICATION_COMMENT_REPLY:
				return "has replied your comment."
			case models.NOTIFICATION_COMMENT_LIKE:
				return "has liked your comment."
			case models.NOTIFICATION_MENTION:
				return "has mentioned you."
			case models.NOTIFICATION_FOLLOW:
				return "has started following you."
			case models.NOTIFICATION_FOLLOW_REQUEST:
				return "wants to follow you."
			}
			return "ERR_NOTIFICATION_MESSAGE"
		}, "lastMessageContent": func(chat models.Chat) string {
			return chat.Messages[0].Content
		}, "lastMessageDate": func(chat models.Chat) int {
			return chat.Messages[0].Date
		}, "renderDateShort": func(unixDateInt int) string {
			const (
				ONE_MINUTE = 60 // seconds
				ONE_HOUR   = 3600
				ONE_DAY    = 86400
				ONE_WEEK   = 604800
				ONE_MONTH  = 2419200
				ONE_YEAR   = 31536000
			)

			if unixDateInt == 0 {
				return "-1y"
			}

			now := time.Now()

			unixDate := time.Unix(int64(unixDateInt), 0)

			diffSeconds := int(math.Floor(now.Sub(unixDate).Seconds()))

			if diffSeconds < 60 {
				return fmt.Sprintf("%ds", diffSeconds)
			} else if diffSeconds < ONE_HOUR {
				return fmt.Sprintf("%dm", int(math.Floor(float64(diffSeconds)/ONE_MINUTE)))
			} else if diffSeconds < ONE_DAY {
				return fmt.Sprintf("%dh", int(math.Floor(float64(diffSeconds)/ONE_HOUR)))
			} else if diffSeconds < ONE_WEEK {
				return fmt.Sprintf("%dd", int(math.Floor(float64(diffSeconds)/ONE_DAY)))
			} else if diffSeconds < ONE_MONTH {
				return fmt.Sprintf("%dw", int(math.Floor(float64(diffSeconds)/ONE_WEEK)))
			} else if diffSeconds < ONE_YEAR {
				return fmt.Sprintf("%dm", int(math.Floor(float64(diffSeconds)/ONE_MONTH)))
			} else {
				return fmt.Sprintf("%dy", int(math.Floor(float64(diffSeconds)/ONE_YEAR)))
			}
		},
	}

	r.AddFromFilesFuncs("base", funcs, "web/templates/base.html", "web/templates/sidenav.html",
		"web/templates/index.html", "web/templates/search.html",
		"web/templates/messages.html", "web/templates/notifications.html",
		"web/templates/profile.html", "web/templates/post.html",
		"web/templates/error.html", "web/templates/settings.html",
		"web/templates/activity.html", "web/templates/saved.html",
	)
	r.AddFromFilesFuncs("index", funcs, "web/templates/index.html")
	r.AddFromFilesFuncs("search", funcs, "web/templates/search.html")
	r.AddFromFilesFuncs("messages", funcs, "web/templates/messages.html")
	r.AddFromFilesFuncs("notifications", funcs, "web/templates/notifications.html")
	r.AddFromFilesFuncs("profile", funcs, "web/templates/profile.html")

	r.AddFromFilesFuncs("error", funcs, "web/templates/error.html")

	r.AddFromFilesFuncs("post", funcs, "web/templates/post.html")

	r.AddFromFilesFuncs("settings", funcs, "web/templates/settings.html")
	r.AddFromFilesFuncs("activity", funcs, "web/templates/activity.html")
	r.AddFromFilesFuncs("saved", funcs, "web/templates/saved.html")

	r.AddFromFilesFuncs("signin", funcs, "web/templates/signin.html")
	r.AddFromFilesFuncs("signup", funcs, "web/templates/signup.html")
	r.AddFromFilesFuncs("endSignup", funcs, "web/templates/endSignup.html")
	return r
}

type NotificationGroup struct {
	DayDate       string
	Notifications []models.Notification
}
