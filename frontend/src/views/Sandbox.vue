<script setup lang="ts">
import CodeCardv2 from '@/components/CodeCard.vue'
import CardContent from '@/components/ui/card/CardContent.vue'
import CardFooter from '@/components/ui/card/CardFooter.vue'
import CardHeader from '@/components/ui/card/CardHeader.vue'
</script>

<template>
  <div class="grid grid-cols-1 gap-6 p-6 md:grid-cols-2 lg:grid-cols-3">
    <Card>
      <CardHeader>
        <CardTitle>Title</CardTitle>
        <CardDescription>Description</CardDescription>
      </CardHeader>
      <CardContent>Blablabla</CardContent>
      <CardFooter>Footer</CardFooter>
    </Card>
    <CodeCardv2
      :snippet="{
        id: '3',
        title: 'Python Function',
        content:
          'def calculate_sum(a, b):\n    return a + b\n\nresult = calculate_sum(5, 3)\nprint(result)',
        author: 'Alex Johnson',
        likes: 15,
        isLiked: false,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      }"
      @click="() => {}"
    />
  </div>
</template>
