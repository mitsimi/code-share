<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <section class="bg-background relative overflow-hidden">
      <div class="container mx-auto px-4 pb-16 sm:py-16">
        <div class="mx-auto max-w-4xl text-center">
          <div class="mb-8">
            <Badge variant="secondary" class="mb-4">
              <CodeIcon class="mr-2 h-3 w-3" />
              <span class="flex space-x-2">
                <span> Share.</span> <span>Discover.</span> <span>Code.</span>
              </span>
            </Badge>
          </div>

          <h1
            class="text-foreground mb-6 text-4xl font-bold tracking-tight sm:text-6xl lg:text-7xl"
          >
            Share code snippets
            <span class="from-primary to-primary/60 bg-gradient-to-r bg-clip-text text-transparent">
              effortlessly
            </span>
          </h1>

          <p class="text-muted-foreground mx-auto mb-8 max-w-2xl text-lg sm:text-xl">
            Discover useful code snippets from the developer community or share your own coding
            knowledge. Make coding collaboration simple and fun.
          </p>

          <div class="flex flex-col gap-4 sm:flex-row sm:justify-center">
            <Button size="lg" class="text-base" @click="handleGetStarted">
              <FileCodeIcon class="mr-2 h-4 w-4" />
              Start Sharing Code
            </Button>
            <Button variant="outline" size="lg" class="text-base" @click="handleBrowseSnippets">
              <SearchIcon class="mr-2 h-4 w-4" />
              Browse Snippets
            </Button>
          </div>

          <div class="text-muted-foreground mt-12 flex items-center justify-center gap-8 text-sm">
            <div class="flex items-center gap-2">
              <CheckIcon class="text-primary h-4 w-4" />
              Free to use
            </div>
            <div class="flex items-center gap-2">
              <CheckIcon class="text-primary h-4 w-4" />
              No signup required to browse
            </div>
            <div class="flex items-center gap-2">
              <CheckIcon class="text-primary h-4 w-4" />
              Open source friendly
            </div>
          </div>
        </div>

        <!-- Hero Code Preview -->
        <div class="mt-16 sm:mt-24">
          <div class="relative mx-auto max-w-5xl">
            <div class="bg-card overflow-hidden rounded-lg border shadow">
              <div class="bg-muted/30 flex items-center justify-between border-b px-4 py-3">
                <!-- Window Management Buttons -->
                <div class="flex items-center gap-2">
                  <div class="h-3 w-3 rounded-full bg-red-500"></div>
                  <div class="h-3 w-3 rounded-full bg-yellow-500"></div>
                  <div class="h-3 w-3 rounded-full bg-green-500"></div>
                </div>
                <div class="text-muted-foreground font-mono text-sm">snippet.js</div>
              </div>
              <div class="bg-card p-6">
                <div class="space-y-2 font-mono text-sm">
                  <div class="text-muted-foreground">// Quick array shuffle function</div>
                  <div>
                    <span class="text-blue-500">function</span
                    ><span class="text-yellow-500"> shuffleArray</span>(<span
                      class="text-orange-500"
                      >array</span
                    >) {
                  </div>
                  <div class="pl-4">
                    <span class="text-purple-500">return</span> array.sort(() => Math.random() -
                    <span class="text-green-500">0.5</span>);
                  </div>
                  <div>}</div>
                </div>
                <div class="text-muted-foreground mt-4 flex items-center gap-4 text-sm">
                  <div class="flex items-center gap-1">
                    <HeartIcon class="h-4 w-4" />
                    24 likes
                  </div>
                  <div class="flex items-center gap-1">
                    <BookmarkIcon class="h-4 w-4" />
                    12 saves
                  </div>
                  <div class="flex items-center gap-1">
                    <TagIcon class="h-4 w-4" />
                    JavaScript
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="bg-muted/30 py-24">
      <div class="container mx-auto px-4">
        <div class="mx-auto mb-16 max-w-2xl text-center">
          <h2 class="text-foreground mb-4 text-3xl font-bold tracking-tight sm:text-4xl">
            Everything developers need to share code
          </h2>
          <p class="text-muted-foreground text-lg">
            A modern platform designed to make code sharing simple, organized, and enjoyable for
            developers.
          </p>
        </div>

        <div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
          <FeatureCard
            v-for="feature in features"
            :key="feature.title"
            :icon="feature.icon"
            :title="feature.title"
            :description="feature.description"
          />
        </div>
      </div>
    </section>

    <!-- Stats Section -->
    <section class="bg-background py-24">
      <div class="container mx-auto px-4">
        <div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
          <StatCard
            v-for="stat in stats"
            :key="stat.label"
            :value="stat.value"
            :label="stat.label"
            :description="stat.description"
          />
        </div>
      </div>
    </section>

    <!-- Testimonials Section -->
    <section class="bg-muted/30 py-24">
      <div class="container mx-auto px-4">
        <div class="mx-auto mb-16 max-w-2xl text-center">
          <h2 class="text-foreground mb-4 text-3xl font-bold tracking-tight sm:text-4xl">
            Loved by developers everywhere
          </h2>
          <p class="text-muted-foreground text-lg">
            See what fellow developers are saying about CodeShare.
          </p>
        </div>

        <div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
          <TestimonialCard
            v-for="testimonial in testimonials"
            :key="testimonial.name"
            :quote="testimonial.quote"
            :name="testimonial.name"
            :role="testimonial.role"
            :avatar="testimonial.avatar"
          />
        </div>
      </div>
    </section>

    <!-- Course Info Section -->
    <section class="bg-accent/50 py-16">
      <div class="container mx-auto px-4">
        <div class="mx-auto max-w-3xl text-center">
          <div class="mb-4 flex items-center justify-center">
            <GraduationCapIcon class="text-primary mr-3 h-8 w-8" />
            <h3 class="text-foreground text-xl font-semibold">Academic Project</h3>
          </div>
          <p class="text-muted-foreground mb-4">
            CodeShare was created as part of the course
            <strong class="text-primary"
              >Special Topic: Introduction to Full Stack Web Development</strong
            >
            at Johannes Kepler University.
          </p>
          <div class="text-muted-foreground flex items-center justify-center gap-4 text-sm">
            <div class="flex items-center gap-2">
              <MonitorIcon class="h-4 w-4" />
              Vue 3 + TypeScript
            </div>
            <div class="flex items-center gap-2">
              <PaletteIcon class="h-4 w-4" />
              Tailwind CSS
            </div>
            <div class="flex items-center gap-2">
              <DatabaseIcon class="h-4 w-4" />
              Go + SQLite
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="bg-primary text-primary-foreground py-24">
      <div class="container mx-auto px-4">
        <div class="mx-auto max-w-3xl text-center">
          <h2 class="mb-4 text-3xl font-bold tracking-tight sm:text-4xl">
            Ready to share your code?
          </h2>
          <p class="text-primary-foreground/80 mb-8 text-lg">
            Join the developer community and start sharing useful code snippets today.
          </p>
          <div class="flex flex-col gap-4 sm:flex-row sm:justify-center">
            <Button size="lg" variant="secondary" class="text-base" @click="handleGetStarted">
              <PlusIcon class="mr-2 h-4 w-4" />
              Share Your First Snippet
            </Button>
            <Button
              size="lg"
              variant="outline
            "
              class="text-base"
              @click="handleBrowseSnippets"
            >
              <LibraryIcon class="mr-2 h-4 w-4" />
              Explore Snippets
            </Button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import FeatureCard from './_components/FeatureCard.vue'
import TestimonialCard from './_components/TestimonialCard.vue'
import StatCard from './_components/StatCard.vue'
import router from '@/router'
import {
  Share2,
  Heart,
  Zap,
  Bell,
  Users,
  Globe,
  CodeIcon,
  FileCodeIcon,
  SearchIcon,
  CheckIcon,
  MonitorIcon,
  PaletteIcon,
  DatabaseIcon,
  PlusIcon,
  LibraryIcon,
  HeartIcon,
  BookmarkIcon,
  TagIcon,
  GraduationCapIcon,
} from 'lucide-vue-next'

// Features data
const features = [
  {
    icon: Share2,
    title: 'Share Code Snippets',
    description:
      'Easily share your useful code snippets with the developer community. Support for all programming languages.',
  },
  {
    icon: Heart,
    title: 'Like & Save Favorites',
    description:
      'Like snippets you find helpful and save your favorites for quick access later. Build your personal collection.',
  },
  {
    icon: Zap,
    title: 'Clean, Modern Interface',
    description:
      'Enjoy a retro-inspired design with crisp black borders and playful shadows for an engaging experience.',
  },
  {
    icon: Bell,
    title: 'Instant Feedback',
    description:
      'Get instant feedback with beautiful toast notifications. Know immediately when your actions are successful.',
  },
  {
    icon: Users,
    title: 'Community Driven',
    description:
      'Join a growing community of developers sharing knowledge and learning from each other.',
  },
  {
    icon: Globe,
    title: 'Built with Modern Tech',
    description:
      'Powered by Vue 3, TypeScript, and Tailwind CSS. A showcase of modern full-stack development.',
  },
]

// Stats data
const stats = [
  {
    value: '500+',
    label: 'Code Snippets',
    description: 'Shared by the community',
  },
  {
    value: '50+',
    label: 'Contributors',
    description: 'Active developers sharing code',
  },
  {
    value: '15+',
    label: 'Languages',
    description: 'Programming languages supported',
  },
  {
    value: '100%',
    label: 'Open Source',
    description: 'Free and accessible to all',
  },
]

// Testimonials data (styled as developer feedback)
const testimonials = [
  {
    quote:
      "Perfect for finding quick solutions! I've discovered so many useful snippets that saved me hours of coding time.",
    name: 'Alex Weber',
    role: 'Frontend Developer & JKU Student',
    avatar: 'AW',
  },
  {
    quote:
      'The interface is clean and intuitive. Sharing my code snippets has never been this easy and fun!',
    name: 'Florian Rohrauer',
    role: 'Full Stack Developer',
    avatar: 'FR',
  },
  {
    quote:
      'Great project for learning full-stack development. The retro design with modern tech stack is brilliant!',
    name: 'David Chen',
    role: 'Computer Science Student',
    avatar: 'DC',
  },
]

const handleGetStarted = () => {
  router.push({ name: 'signup' })
}

const handleBrowseSnippets = () => {
  router.push({ name: 'snippets' })
}
</script>
