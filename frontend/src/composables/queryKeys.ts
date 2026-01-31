export const queryKeys = {
  snippets: () => ['snippets'] as const,
  snippet: (id: string) => ['snippets', id] as const,
  mySnippets: () => ['my-snippets'] as const,
  likedSnippets: () => ['liked-snippets'] as const,
  savedSnippets: () => ['saved-snippets'] as const,
  user: (id: string) => ['user', id] as const,
}

// Helper type for query key types
export type QueryKeys = ReturnType<
  typeof queryKeys.snippets
> | ReturnType<
  typeof queryKeys.snippet
> | ReturnType<
  typeof queryKeys.mySnippets
> | ReturnType<
  typeof queryKeys.likedSnippets
> | ReturnType<
  typeof queryKeys.savedSnippets
> | ReturnType<
  typeof queryKeys.user
>
