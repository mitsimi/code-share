export interface Card {
  id: string
  title: string
  content: string
  author: string
  likes: number
  is_liked: boolean
  created_at: string
  updated_at: string
  user_likes: Record<string, boolean>
}
