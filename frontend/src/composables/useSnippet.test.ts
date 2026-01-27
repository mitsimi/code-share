import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { useSnippet } from './useSnippet'
import { snippetsService } from '@/services/snippets'
import { toast } from 'vue-sonner'
import type { Snippet } from '@/types'

// Mock dependencies
vi.mock('@/services/snippets', () => ({
  snippetsService: {
    getSnippet: vi.fn(),
    updateSnippet: vi.fn(),
    deleteSnippet: vi.fn(),
  },
}))

vi.mock('vue-sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

// Helper component to run the composable in a Vue context
const TestComponent = defineComponent({
  props: {
    snippetId: { type: String, required: true },
  },
  setup(props) {
    return useSnippet(props.snippetId)
  },
  template: '<div></div>',
})

describe('useSnippet', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    setActivePinia(createPinia())
    queryClient = new QueryClient({
      defaultOptions: {
        queries: {
          retry: false,
        },
      },
    })
    vi.clearAllMocks()
  })

  afterEach(() => {
    queryClient.clear()
  })

  const mockSnippet: Snippet = {
    id: '123',
    title: 'Test Snippet',
    content: 'console.log("hello")',
    language: 'typescript',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    author: {
      id: 'user1',
      username: '',
      avatar: '',
      email: '',
      createdAt: '',
      updatedAt: '',
    },
    likes: 0,
    views: 0,
    isLiked: false,
    isSaved: false,
  }

  it('fetches snippet and adds to store', async () => {
    // Setup mock
    vi.mocked(snippetsService.getSnippet).mockResolvedValue(mockSnippet)

    // Mount component with Vue Query plugin
    const wrapper = mount(TestComponent, {
      props: { snippetId: '123' },
      global: {
        plugins: [[VueQueryPlugin, { queryClient }]],
      },
    })

    // Wait for query to resolve
    await new Promise((resolve) => setTimeout(resolve, 0))
    // In a real scenario we might wait for isLoading to be false,
    // but with mocked immediate resolution, a tick is usually enough.
    // Better: use `flushPromises` if available, or wait for the query state.

    // Check if service was called
    expect(snippetsService.getSnippet).toHaveBeenCalledWith('123')

    // Check if data is available via the composable
    expect(wrapper.vm.snippet).toEqual(mockSnippet)
    expect(wrapper.vm.isLoading).toBe(false)
  })

  it('updates snippet successfully', async () => {
    vi.mocked(snippetsService.getSnippet).mockResolvedValue(mockSnippet)
    const updatedSnippet = { ...mockSnippet, title: 'Updated Title' }
    vi.mocked(snippetsService.updateSnippet).mockResolvedValue(updatedSnippet)

    const wrapper = mount(TestComponent, {
      props: { snippetId: '123' },
      global: {
        plugins: [[VueQueryPlugin, { queryClient }]],
      },
    })

    // Wait for initial fetch
    await new Promise((resolve) => setTimeout(resolve, 0))

    // Perform update
    wrapper.vm.updateSnippet({ title: 'Updated Title', content: '...' })

    // Wait for mutation
    await new Promise((resolve) => setTimeout(resolve, 0))

    expect(snippetsService.updateSnippet).toHaveBeenCalledWith('123', {
      title: 'Updated Title',
      content: '...',
    })
    expect(toast.success).toHaveBeenCalledWith('Snippet updated successfully')

    // Check if store was updated
    expect(wrapper.vm.snippet).toEqual(updatedSnippet)
  })

  it('deletes snippet successfully', async () => {
    vi.mocked(snippetsService.getSnippet).mockResolvedValue(mockSnippet)
    vi.mocked(snippetsService.deleteSnippet).mockResolvedValue()

    const wrapper = mount(TestComponent, {
      props: { snippetId: '123' },
      global: {
        plugins: [[VueQueryPlugin, { queryClient }]],
      },
    })

    // Wait for initial fetch
    await new Promise((resolve) => setTimeout(resolve, 0))

    // Perform delete
    wrapper.vm.deleteSnippet()

    // Wait for mutation
    await new Promise((resolve) => setTimeout(resolve, 0))

    expect(snippetsService.deleteSnippet).toHaveBeenCalledWith('123')
    expect(toast.success).toHaveBeenCalledWith('Snippet deleted successfully')

    // Check if removed from store (snippet should be undefined)
    expect(wrapper.vm.snippet).toBeUndefined()
  })
})
