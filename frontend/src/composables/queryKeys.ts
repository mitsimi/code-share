export const queryKeys = {
  all: () => ['snippets'] as const,
  lists: () => [...queryKeys.all(), 'list'] as const,
  list: (filters?: { language?: string; author?: string }) =>
    [...queryKeys.lists(), filters] as const,
  details: () => [...queryKeys.all(), 'detail'] as const,
  detail: (id: string) => [...queryKeys.details(), id] as const,
  my: () => [...queryKeys.all(), 'my'] as const,
  liked: () => [...queryKeys.all(), 'liked'] as const,
  saved: () => [...queryKeys.all(), 'saved'] as const,
  user: (id: string) => ['user', id] as const,
}

export type QueryKeys =
  | ReturnType<typeof queryKeys.all>
  | ReturnType<typeof queryKeys.lists>
  | ReturnType<typeof queryKeys.list>
  | ReturnType<typeof queryKeys.details>
  | ReturnType<typeof queryKeys.detail>
  | ReturnType<typeof queryKeys.my>
  | ReturnType<typeof queryKeys.liked>
  | ReturnType<typeof queryKeys.saved>
  | ReturnType<typeof queryKeys.user>
