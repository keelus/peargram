package models

type PendingSignup struct {
	GoogleID    string `db:"googleID"`
	GoogleEmail string `db:"googleEmail"`
}

type User struct {
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type UserDetails struct {
	Username        string `db:"username"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	Avatar          string `db:"avatar"`
	PostAmount      int
	FollowerAmount  int
	FollowingAmount int
}

type Post struct {
	ID            int    `db:"id"`
	Username      string `db:"username"`
	Content       string `db:"content"`
	Likes         uint
	Comments      []Comment
	CommentAmount uint
	Date          int `db:"date"`
	UserDetails   UserDetails
}

type Comment struct {
	ID       int    `db:"id"`
	ParentID int    `db:"parentID"`
	PostID   int    `db:"postID"`
	Actor    string `db:"actor"`
	Content  string `db:"content"`
	Date     int    `db:"date"`
}

type CommentBlock struct {
	ParentComment    Comment
	ChildrenComments []Comment
}

type Notification struct {
	ID     int    `db:"id" json:"ID"`
	Actor  string `db:"actor" json:"Actor"`
	Target string `db:"target" json:"Target"`
	Type   int    `db:"type" json:"Type"`
	Date   int    `db:"date" json:"Date"`
	Post   *int   `db:"post" json:"Post"`
}

type Message struct {
	ID        int     `db:"id" json:"ID"`
	Actor     string  `db:"actor" json:"Actor"`
	Target    string  `db:"target" json:"Target"`
	Content   string  `db:"content" json:"Content"`
	Reactions *string `db:"reactions" json:"Reactions"`
	Date      int     `db:"date" json:"Date"`
}

type Chat struct {
	Participants      [2]string // Participants[0] = First user/current logged user
	Messages          []Message
	MessageTotalCount uint
}

const (
	NOTIFICATION_POST_LIKE = iota
	NOTIFICATION_POST_COMMENT
	NOTIFICATION_COMMENT_REPLY
	NOTIFICATION_COMMENT_LIKE
	NOTIFICATION_MENTION
	NOTIFICATION_FOLLOW
	NOTIFICATION_FOLLOW_REQUEST
)
