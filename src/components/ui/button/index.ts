import { cva, type VariantProps } from 'class-variance-authority'

export { default as Button } from './Button.vue'

export const buttonVariants = cva(
  "inline-flex items-center rounded-lg justify-center gap-2 whitespace-nowrap text-sm font-bold transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none border-4 border-black dark:border-white shadow-[4px_4px_0_0_#000] dark:shadow-[4px_4px_0_0_#fff] active:translate-x-1 active:translate-y-1 active:shadow-none",
  {
    variants: {
      variant: {
        default:
          'bg-primary text-primary-foreground hover:translate-x-1 hover:translate-y-1 hover:shadow-none',
        destructive:
          'bg-destructive text-white hover:translate-x-1 hover:translate-y-1 hover:shadow-none',
        outline:
          'bg-background hover:bg-accent hover:text-accent-foreground hover:translate-x-1 hover:translate-y-1 hover:shadow-none',
        secondary:
          'bg-secondary text-secondary-foreground hover:translate-x-1 hover:translate-y-1 hover:shadow-none',
        ghost:
          'border-none shadow-none hover:bg-accent hover:text-accent-foreground hover:border-4 hover:border-black dark:hover:border-white hover:shadow-[4px_4px_0_0_#000] dark:hover:shadow-[4px_4px_0_0_#fff]',
        link: 'border-none shadow-none text-primary underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-10 px-4 py-2 has-[>svg]:px-3',
        sm: 'h-8 px-3 has-[>svg]:px-2.5',
        lg: 'h-12 px-6 has-[>svg]:px-4',
        icon: 'size-10',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

export type ButtonVariants = VariantProps<typeof buttonVariants>
