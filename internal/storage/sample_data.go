package storage

import (
	"context"

	"github.com/google/uuid"
	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/domain"
)

// SeedSampleData seeds the database with sample users and snippets
func (s *Storage) SeedSampleData(ctx context.Context) error {
	// Create sample users
	userIDs := make([]string, len(sampleUsers))
	for i, user := range sampleUsers {
		// Set a default password for all sample users
		passwordHash, err := auth.HashPassword("password123")
		if err != nil {
			return err
		}

		userCreation := domain.UserCreation{
			ID:           uuid.New().String(),
			Username:     user.Username,
			Email:        user.Email,
			PasswordHash: passwordHash,
		}

		createdUser, err := s.CreateUser(ctx, &userCreation)
		if err != nil {
			return err
		}
		userIDs[i] = createdUser.ID
	}

	// Create sample snippets
	for i, snippet := range sampleSnippetsData {
		// Assign snippets to users in a round-robin fashion
		authorID := userIDs[i%len(userIDs)]
		snippet.Author = &domain.User{
			ID: authorID,
		}

		if snippet.Language == "" {
			snippet.Language = "plainext" // Default language if not specified
		}

		if err := s.CreateSnippet(ctx, &snippet); err != nil {
			return err
		}
	}

	return nil
}

// SampleUsers contains predefined sample users
var sampleUsers = []domain.UserCreation{
	{
		Username:     "demo",
		Email:        "demo@example.com",
		PasswordHash: "", // Will be set during seeding
	},
	{
		Username:     "john_doe",
		Email:        "john@example.com",
		PasswordHash: "", // Will be set during seeding
	},
	{
		Username:     "alice_johnson",
		Email:        "alice@example.com",
		PasswordHash: "", // Will be set during seeding
	},
	{
		Username:     "bob_smith",
		Email:        "bob@example.com",
		PasswordHash: "", // Will be set during seeding
	},
}

// SampleSnippetsData contains predefined sample snippets
var sampleSnippetsData = []domain.Snippet{
	{
		ID:       "1",
		Title:    "Vue 3 Composition API Example",
		Content:  "<template>\n\t<div>\n\t\t<p>Count: {{ count }}</p>\n\t\t<p>Double: {{ double }}</p>\n\t\t<button @click=\"increment\">Increment</button>\n\t</div>\n</template>\n\n<script setup>\nimport { ref, computed } from 'vue'\n\nconst count = ref(0)\nconst double = computed(() => count.value * 2)\n\nfunction increment() {\n\tcount.value++\n}\n</script>\n",
		Language: "vue",
	},
	{
		ID:       "2",
		Title:    "React Hooks: useEffect Example",
		Content:  "import React, { useState, useEffect } from 'react';\n\nfunction Example() {\n    const [count, setCount] = useState(0);\n\n    useEffect(() => {\n    document.title = `Count: ${count}`;\n    });\n\n    return (\n    <div>\n        <p>You clicked {count} times</p>\n        <button onClick={() => setCount(count + 1)}>Click me</button>\n    </div>\n    );\n}",
		Language: "jsx",
	},
	{
		ID:       "3",
		Title:    "Svelte Counter Example",
		Content:  "<script>\n\tlet count = 0;\n</script>\n\n<main>\n\t<p>The count is {count}</p>\n\t<button on:click={() => count++}>\n\t\tIncrement\n\t</button>\n</main>",
		Language: "svelte",
	},
	{
		ID:       "4",
		Title:    "Shell Script Example",
		Content:  "#!/bin/bash\n\necho \"Hello, World!\"",
		Language: "shell",
	},
	{
		ID:       "5",
		Title:    "Java Basic Hello World",
		Content:  "public class HelloWorld {\n    public static void main(String[] args) {\n        System.out.println(\"Hello, World!\");\n    }\n}",
		Language: "java",
	},
	{
		ID:       "6",
		Title:    "Go Simple HTTP Server",
		Content:  "package main\n\nimport (\n    \"fmt\"\n    \"net/http\"\n)\n\nfunc handler(w http.ResponseWriter, r *http.Request) {\n    fmt.Fprintf(w, \"Hello, World!\")\n}\n\nfunc main() {\n    http.HandleFunc(\"/\", handler)\n    http.ListenAndServe(\":8080\", nil)\n}",
		Language: "go",
	},
	{
		ID:       "7",
		Title:    "Python Data Analysis with Pandas",
		Content:  "import pandas as pd\n\ndata = {'Name': ['Tom', 'Jerry', 'Mickey'], 'Age': [20, 21, 19]}\ndf = pd.DataFrame(data)\nprint(df.describe())",
		Language: "python",
	},
	{
		ID:       "8",
		Title:    "TypeScript Generic Function",
		Content:  "function identity<T>(arg: T): T {\n    return arg;\n}\n\n// Usage\nlet output = identity<string>(\"myString\");\nlet output2 = identity<number>(42);",
		Language: "typescript",
	},
	{
		ID:       "9",
		Title:    "Rust Error Handling",
		Content:  "fn read_file(path: &str) -> Result<String, io::Error> {\n    let mut file = File::open(path)?;\n    let mut contents = String::new();\n    file.read_to_string(&mut contents)?;\n    Ok(contents)\n}",
		Language: "rust",
	},
	{
		ID:       "10",
		Title:    "Docker Compose Example",
		Content:  "version: '3'\nservices:\n  web:\n    build: .\n    ports:\n      - \"5000:5000\"\n  redis:\n    image: redis:alpine\n    ports:\n      - \"6379:6379\"",
		Language: "yaml",
	},
	{
		ID:       "11",
		Title:    "SQL Query with JOIN",
		Content:  "SELECT users.name, orders.order_date, products.name\nFROM users\nINNER JOIN orders ON users.id = orders.user_id\nINNER JOIN products ON orders.product_id = products.id\nWHERE orders.order_date > '2024-01-01'\nORDER BY orders.order_date DESC;",
		Language: "sql",
	},
	{
		ID:       "12",
		Title:    "Kotlin Coroutine Example",
		Content:  "suspend fun fetchUserData() = coroutineScope {\n    val user = async { userRepository.getUser() }\n    val posts = async { postRepository.getPosts() }\n    UserWithPosts(user.await(), posts.await())\n}",
		Language: "kotlin",
	},
	{
		ID:       "13",
		Title:    "GraphQL Query",
		Content:  "query {\n  user(id: \"123\") {\n    name\n    email\n    posts {\n      title\n      content\n      comments {\n        text\n        author {\n          name\n        }\n      }\n    }\n  }\n}",
		Language: "graphql",
	},
	{
		ID:       "14",
		Title:    "Swift UI View",
		Content:  "struct ContentView: View {\n    @State private var name = \"\"\n    \n    var body: some View {\n        VStack {\n            TextField(\"Enter your name\", text: $name)\n            Text(\"Hello, \\(name)!\")\n                .font(.title)\n        }\n        .padding()\n    }\n}",
		Language: "swift",
	},
	{
		ID:       "15",
		Title:    "Elixir Pattern Matching",
		Content:  "defmodule Math do\n  def factorial(0), do: 1\n  def factorial(n) when n > 0 do\n    n * factorial(n - 1)\n  end\nend",
		Language: "elixir",
	},
	{
		ID:       "16",
		Title:    "C# LINQ Example",
		Content:  "var result = numbers\n    .Where(n => n % 2 == 0)\n    .Select(n => n * n)\n    .OrderByDescending(n => n)\n    .Take(5)\n    .ToList();",
		Language: "c#",
	},
	{
		ID:       "17",
		Title:    "Haskell List Comprehension",
		Content:  "primes = [x | x <- [2..], all (\\y -> x `mod` y /= 0) [2..x-1]]\n\n-- Get first 10 prime numbers\ntake 10 primes",
		Language: "haskell",
	},
	{
		ID:       "18",
		Title:    "Scala Case Class",
		Content:  "case class Person(name: String, age: Int)\n\nval people = List(\n  Person(\"Alice\", 25),\n  Person(\"Bob\", 30),\n  Person(\"Charlie\", 35)\n)\n\nval youngPeople = people.filter(_.age < 30)",
		Language: "scala",
	},
	{
		ID:       "19",
		Title:    "PHP Laravel Route",
		Content:  "Route::middleware(['auth'])->group(function () {\n    Route::get('/profile', [ProfileController::class, 'show']);\n    Route::put('/profile', [ProfileController::class, 'update']);\n    Route::delete('/profile', [ProfileController::class, 'destroy']);\n});",
		Language: "php",
	},
	{
		ID:       "20",
		Title:    "Ruby on Rails Model",
		Content:  "class User < ApplicationRecord\n  has_many :posts, dependent: :destroy\n  has_many :comments\n  \n  validates :email, presence: true, uniqueness: true\n  validates :username, presence: true, length: { minimum: 3 }\n  \n  before_save :normalize_email\n  \n  private\n  \n  def normalize_email\n    self.email = email.downcase.strip\n  end\nend",
		Language: "ruby",
	},
	{
		ID:       "21",
		Title:    "Deno Web Server",
		Content:  "import { serve } from \"https://deno.land/std@0.168.0/http/server.ts\";\n\nconst handler = (req: Request): Response => {\n  const url = new URL(req.url);\n  return new Response(`Hello from ${url.pathname}!`);\n};\n\nserve(handler, { port: 8000 });",
		Language: "typescript",
	},
}
