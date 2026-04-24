<script lang="ts">
  import Play from "@lucide/svelte/icons/play";
  import Square from "@lucide/svelte/icons/square";
  import {
    initAudioPreview,
    subscribeAudioPreview,
    toggleAudioPreview,
  } from "./audioPreview";
  import { onMount } from "svelte";
  import { Button } from "$lib/components/ui/button";
  import { t } from "$lib/i18n";

  let paused = $state(true);
  let ready = $state(false);

  function togglePlay() {
    toggleAudioPreview();
  }

  onMount(() => {
    initAudioPreview();
    const unsubscribe = subscribeAudioPreview((state) => {
      paused = state.paused;
      ready = state.ready;
    });
    return unsubscribe;
  });
</script>

<Button onclick={togglePlay} variant="outline" disabled={!ready}>
  {#if paused}
    <Play />
  {:else}
    <Square />
  {/if}
  {$t.testSound}
</Button>
