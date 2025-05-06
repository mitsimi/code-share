<script setup lang="ts">
import { ref } from 'vue'
import AppHeader from '../components/AppHeader.vue'
import CardGrid from '../components/CardGrid.vue'
import AppFooter from '../components/AppFooter.vue'
import FloatingActionButton from '../components/FloatingActionButton.vue'
import SnippetModal from '../components/SnippetModal.vue'
import ToastNotification from '../components/ToastNotification.vue'

interface Card {
  id: number
  title: string
  code: string
  author: string
  likes: number
  isLiked: boolean
}

const showModal = ref(false)
const showToast = ref(false)
const toastMessage = ref('')
const toastTimeout = ref<number | null>(null)

const cards = ref<Card[]>([
  {
    id: 1,
    title: 'Vue 3 Composition API Example',
    code: "import { ref, computed } from 'vue'\nexport default {\n    setup() {\n        const count = ref(0)\n        const double = computed(() => count.value * 2)\n\n        function increment() {\n            count.value++\n        }\n\n        return {\n            count,\n            double,\n            increment\n        }\n    }\n}",
    author: 'John Doe',
    likes: 290,
    isLiked: true,
  },
  {
    id: 2,
    title: 'React Hooks: useEffect Example',
    code: "import React, { useState, useEffect } from 'react';\n\nfunction Example() {\n    const [count, setCount] = useState(0);\n\n    useEffect(() => {\n    document.title = `Count: ${count}`;\n    });\n\n    return (\n    <div>\n        <p>You clicked {count} times</p>\n        <button onClick={() => setCount(count + 1)}>Click me</button>\n    </div>\n    );\n}",
    author: 'Alice Johnson',
    likes: 104,
    isLiked: true,
  },
  {
    id: 3,
    title: 'Shell Script Example',
    code: '#!/bin/bash\n\necho "Hello, World!"',
    author: 'Paula Cyan',
    likes: 17,
    isLiked: false,
  },
  {
    id: 4,
    title: 'Java Basic Hello World',
    code: 'public class HelloWorld {\n    public static void main(String[] args) {\n        System.out.println("Hello, World!");\n    }\n}',
    author: 'Ivy Purple',
    likes: 3,
    isLiked: false,
  },
  {
    id: 5,
    title: 'Go Simple HTTP Server',
    code: 'package main\n\nimport (\n    "fmt"\n    "net/http"\n)\n\nfunc handler(w http.ResponseWriter, r *http.Request) {\n    fmt.Fprintf(w, "Hello, World!")\n}\n\nfunc main() {\n    http.HandleFunc("/", handler)\n    http.ListenAndServe(":8080", nil)\n}',
    author: 'Grace Yellow',
    likes: 561,
    isLiked: false,
  },
  {
    id: 6,
    title: 'Python Data Analysis with Pandas',
    code: "import pandas as pd\n\ndata = {'Name': ['Tom', 'Jerry', 'Mickey'], 'Age': [20, 21, 19]}\ndf = pd.DataFrame(data)\nprint(df.describe())",
    author: 'Guido van Rossum',
    likes: -666,
    isLiked: false,
  },
])

const toggleLike = (card: Card) => {
  card.isLiked = !card.isLiked
  card.likes += card.isLiked ? 1 : -1
}

const submitSnippet = (formData: { title: string; code: string; author: string }) => {
  const newId = Math.max(...cards.value.map((card) => card.id)) + 1

  const newCard: Card = {
    id: newId,
    title: formData.title,
    code: formData.code,
    author: formData.author,
    likes: 0,
    isLiked: false,
  }

  cards.value.unshift(newCard)
  showModal.value = false
  showToastMessage('"' + newCard.title + '" has been added successfully!')
}

const showToastMessage = (message: string) => {
  if (toastTimeout.value) {
    clearTimeout(toastTimeout.value)
  }

  toastMessage.value = message
  showToast.value = true
  toastTimeout.value = window.setTimeout(() => {
    showToast.value = false
  }, 3000)
}
</script>

<template>
  <AppHeader />

  <main class="mx-auto my-12 max-w-7xl px-4">
    <CardGrid :cards="cards" @toggle-like="toggleLike" />
  </main>

  <AppFooter />

  <FloatingActionButton @click="showModal = true" />

  <SnippetModal :show="showModal" @close="showModal = false" @submit="submitSnippet" />

  <ToastNotification :show="showToast" :message="toastMessage" />
</template>
