package main

import "github.com/weaviate/weaviate/entities/models"

// first 5 objects contain something about programming without word 'programming' while other are random topics.
var data []*models.Object = []*models.Object{
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Debugging Java applications can be challenging."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Writing code in Python is both fun and efficient."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "HTML and CSS are the building blocks of web development."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Using version control systems like Git is essential for collaborative coding."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Understanding algorithms is crucial for effective software development."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Swimming is a fun and beneficial form of exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Playing musical instruments can be both challenging and fulfilling."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Playing musical instruments can be both challenging and fulfilling."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Learning a new language opens up a world of opportunities."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Playing musical instruments can be both challenging and fulfilling."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Meditation helps in reducing stress and improving focus."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Playing musical instruments can be both challenging and fulfilling."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Gardening is a peaceful and rewarding hobby."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Running is a simple and effective way to stay fit."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Playing musical instruments can be both challenging and fulfilling."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "The art of cooking is a blend of flavors and techniques."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Watching movies is a popular form of entertainment."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Practicing yoga promotes physical and mental well-being."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Reading books is a journey through knowledge and imagination."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "The art of cooking is a blend of flavors and techniques."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Running is a simple and effective way to stay fit."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Practicing yoga promotes physical and mental well-being."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Photography captures moments and memories."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Gardening is a peaceful and rewarding hobby."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Swimming is a fun and beneficial form of exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Reading books is a journey through knowledge and imagination."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Painting allows for creative expression and relaxation."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "The art of cooking is a blend of flavors and techniques."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Photography captures moments and memories."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Reading books is a journey through knowledge and imagination."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Running is a simple and effective way to stay fit."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Playing musical instruments can be both challenging and fulfilling."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Learning a new language opens up a world of opportunities."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Painting allows for creative expression and relaxation."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "The art of cooking is a blend of flavors and techniques."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Meditation helps in reducing stress and improving focus."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Practicing yoga promotes physical and mental well-being."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Swimming is a fun and beneficial form of exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Photography captures moments and memories."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Practicing yoga promotes physical and mental well-being."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Running is a simple and effective way to stay fit."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Swimming is a fun and beneficial form of exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Running is a simple and effective way to stay fit."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Meditation helps in reducing stress and improving focus."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Photography captures moments and memories."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Running is a simple and effective way to stay fit."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Practicing yoga promotes physical and mental well-being."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Practicing yoga promotes physical and mental well-being."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Reading books is a journey through knowledge and imagination."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Watching movies is a popular form of entertainment."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Playing musical instruments can be both challenging and fulfilling."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Painting allows for creative expression and relaxation."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Hiking in nature is both invigorating and relaxing."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Meditation helps in reducing stress and improving focus."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Watching movies is a popular form of entertainment."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Photography captures moments and memories."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Learning a new language opens up a world of opportunities."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Painting allows for creative expression and relaxation."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Photography captures moments and memories."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Practicing yoga promotes physical and mental well-being."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Practicing yoga promotes physical and mental well-being."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "The art of cooking is a blend of flavors and techniques."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Running is a simple and effective way to stay fit."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Running is a simple and effective way to stay fit."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Photography captures moments and memories."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "The art of cooking is a blend of flavors and techniques."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Photography captures moments and memories."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Playing musical instruments can be both challenging and fulfilling."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "The art of cooking is a blend of flavors and techniques."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Swimming is a fun and beneficial form of exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Cycling is both an eco-friendly transport and a great way to exercise."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Traveling to new places is an adventure worth having."},
	},
	{
		Class:      "Records",
		Properties: map[string]string{"content": "Painting allows for creative expression and relaxation."},
	},
}
