<script lang="ts">
  import * as Tooltip from "$lib/components/ui/tooltip/index.js";
  import CircleHelp from "@lucide/svelte/icons/circle-help";
  import { mergeProps } from "bits-ui";
  import { cn } from "$lib/utils.js";

  interface Props {
    /** Tooltip text (from i18n) */
    description: string;
    /** Label for button and screen readers */
    ariaLabel: string;
  }

  let { description, ariaLabel }: Props = $props();
</script>

<Tooltip.Root>
  <Tooltip.Trigger>
    {#snippet child({ props })}
      <button
        {...mergeProps(props, {
          type: "button" as const,
          "aria-label": ariaLabel,
          class: cn(
            "inline-flex size-7 shrink-0 items-center justify-center rounded-md text-muted-foreground",
            "hover:bg-muted hover:text-foreground",
            "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
          ),
        })}
      >
        <CircleHelp class="size-4 shrink-0" strokeWidth={2} />
      </button>
    {/snippet}
  </Tooltip.Trigger>
  <Tooltip.Content side="top" class="max-w-sm text-left leading-snug">
    <span class="block whitespace-pre-wrap">{description}</span>
  </Tooltip.Content>
</Tooltip.Root>
