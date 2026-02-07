package main

import (
	"flag"
	"log"
	"math/rand"

	"po-backend/configs"
	"po-backend/helper"
	"po-backend/models"
)

var names = []struct {
	Name     string
	Username string
	Bio      string
}{
	{"Alice Johnson", "alice", "Software engineer who loves hiking and coffee."},
	{"Bob Smith", "bob", "Music lover and weekend guitarist. Based in Portland."},
	{"Charlie Brown", "charlie", "Photographer capturing life one frame at a time."},
	{"Diana Lee", "diana", "Bookworm, tea enthusiast, and part-time baker."},
	{"Ethan Park", "ethan", "Fitness junkie and aspiring travel blogger."},
	{"Fiona Chen", "fiona", "UX designer by day, gamer by night."},
	{"George Miller", "george", "Dad jokes connoisseur. You've been warned."},
	{"Hannah Davis", "hannah", "Plant mom with way too many succulents."},
	{"Isaac Torres", "isaac", "Startup founder. Building the future, one bug at a time."},
	{"Julia Kim", "julia", "Film student and popcorn aficionado."},
}

var postContents = []string{
	"Just finished a 10k run! Personal best time. Feeling amazing right now.",
	"Anyone else think Monday mornings should be illegal?",
	"Made homemade pasta from scratch today. It actually turned out decent!",
	"Reading a fantastic book right now. Can't put it down.",
	"Hot take: pineapple absolutely belongs on pizza.",
	"Just adopted a rescue dog! Meet my new best friend.",
	"The sunset tonight was absolutely breathtaking. Nature never disappoints.",
	"Finally fixed that bug that's been haunting me for three days. Pure joy.",
	"Coffee is not a want, it's a need. Don't @ me.",
	"Started learning guitar last month. My fingers hurt but it's worth it.",
	"Who else is excited for the weekend? I have zero plans and I love it.",
	"Just watched the most mind-blowing documentary. Highly recommend it.",
	"Tried a new recipe tonight and my kitchen looks like a war zone.",
	"Morning walks with my dog are the best therapy.",
	"Unpopular opinion: winter is the best season. Fight me.",
	"Just got promoted at work! Hard work really does pay off.",
	"Why do I always think of the perfect comeback three hours later?",
	"Spent the entire day organizing my closet. Send help.",
	"Nothing beats a rainy afternoon with a good book and hot chocolate.",
	"Just finished my first marathon! Every step was worth it.",
}

var commentContents = []string{
	"This is so relatable!",
	"Congrats! That's amazing!",
	"Haha, couldn't agree more.",
	"Love this! Keep it up.",
	"Totally feel you on this one.",
	"That's awesome, well done!",
	"I needed to hear this today.",
	"Same here! You're not alone.",
	"This made my day, thanks for sharing.",
	"Wow, that's incredible!",
	"I've been saying this for years!",
	"So proud of you!",
	"This is gold.",
	"Can't stop laughing at this.",
	"You're an inspiration!",
	"Hard agree on this one.",
	"Tell me more! I'm intrigued.",
	"Best post I've seen all day.",
	"Living your best life!",
	"Facts. No debate needed.",
	"You always have the best takes.",
	"This deserves way more likes.",
	"Literally me every single day.",
	"Sending good vibes your way!",
	"I wish I could double like this.",
	"My thoughts exactly!",
	"This is wholesome content right here.",
	"Nailed it.",
	"Adding this to my to-do list.",
	"You never fail to make me smile.",
}

func main() {
	reset := flag.Bool("reset", false, "Drop all tables and re-seed the database")
	flag.Parse()

	cfg := configs.Envs
	if err := cfg.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if *reset {
		log.Println("Resetting database — dropping all tables...")
		if err := cfg.DB.Migrator().DropTable(
			&models.Notification{},
			&models.CommentLike{},
			&models.PostLike{},
			&models.Follow{},
			&models.Comment{},
			&models.Post{},
			&models.User{},
		); err != nil {
			log.Fatal("Failed to drop tables:", err)
		}
		log.Println("All tables dropped")
	}

	if err := cfg.InitializeDB(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	var users []models.User
	password, _ := helper.HashPassword("password")

	for _, n := range names {
		user := models.User{
			Name:     n.Name,
			Username: n.Username,
			Bio:      n.Bio,
			Password: password,
		}
		if err := cfg.DB.Create(&user).Error; err != nil {
			log.Printf("Error creating user %s: %v\n", n.Username, err)
			continue
		}
		users = append(users, user)
	}
	log.Printf("Created %d users\n", len(users))

	var posts []models.Post
	for i, content := range postContents {
		post := models.Post{
			Content: content,
			UserID:  users[i%len(users)].ID,
		}
		if err := cfg.DB.Create(&post).Error; err != nil {
			log.Printf("Error creating post: %v\n", err)
			continue
		}
		posts = append(posts, post)
	}
	log.Printf("Created %d posts\n", len(posts))

	commentCount := 0
	for i, content := range commentContents {
		comment := models.Comment{
			Content: content,
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[i%len(posts)].ID,
		}
		if err := cfg.DB.Create(&comment).Error; err != nil {
			log.Printf("Error creating comment: %v\n", err)
			continue
		}
		commentCount++
	}
	log.Printf("Created %d comments\n", commentCount)

	followCount := 0
	for i := range users {
		for j := 0; j < 3; j++ {
			target := rand.Intn(len(users))
			if target != i {
				follow := models.Follow{
					FollowerID:  users[i].ID,
					FollowingID: users[target].ID,
				}
				if err := cfg.DB.Create(&follow).Error; err == nil {
					followCount++
				}
			}
		}
	}
	log.Printf("Created %d follows\n", followCount)

	likeCount := 0
	for i := 0; i < 40; i++ {
		like := models.PostLike{
			PostID: posts[rand.Intn(len(posts))].ID,
			UserID: users[rand.Intn(len(users))].ID,
		}
		if err := cfg.DB.Create(&like).Error; err == nil {
			likeCount++
		}
	}
	log.Printf("Created %d post likes\n", likeCount)

	commentLikeCount := 0
	for i := 0; i < 20; i++ {
		like := models.CommentLike{
			CommentID: uint(rand.Intn(commentCount) + 1),
			UserID:    users[rand.Intn(len(users))].ID,
		}
		if err := cfg.DB.Create(&like).Error; err == nil {
			commentLikeCount++
		}
	}
	log.Printf("Created %d comment likes\n", commentLikeCount)

	log.Println("Seeding complete!")
}
