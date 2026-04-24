<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    activePanel: "audio" | "theme" | "info" | null;
    radarContent: Snippet;
    audioPanel: Snippet;
    themePanel: Snippet;
    infoPanel: Snippet;
  }

  const {
    activePanel,
    radarContent,
    audioPanel,
    themePanel,
    infoPanel,
  }: Props = $props();
</script>

<section class="workspace-grid" class:has-active-panel={activePanel !== null}>
  <section class="radar-workspace">
    <div class="radar-transparent-layer">
      {@render radarContent()}
    </div>
  </section>

  {#if activePanel === "audio"}
    <aside class="full-height-panel audio-panel">
      {@render audioPanel()}
    </aside>
  {:else if activePanel === "theme"}
    <aside class="full-height-panel theme-panel">
      {@render themePanel()}
    </aside>
  {:else if activePanel === "info"}
    <aside class="full-height-panel info-panel">
      {@render infoPanel()}
    </aside>
  {/if}
</section>

<style>
  .workspace-grid {
    min-height: 0;
    height: 100%;
    display: grid;
    gap: 1rem;
    align-items: stretch;
    grid-template-columns: minmax(360px, 540px) 0;
  }

  .workspace-grid.has-active-panel {
    grid-template-columns: minmax(360px, 540px) minmax(0, 1fr);
  }

  .radar-workspace {
    min-height: 0;
    display: grid;
    align-content: start;
    gap: 0.75rem;
  }

  .radar-transparent-layer {
    min-height: 444px;
    display: grid;
    place-items: center;
    background: transparent;
    border: none;
    border-radius: 0;
    box-shadow: none;
    padding: 0.25rem;
  }

  .radar-transparent-layer :global(.radar-container) {
    width: min(420px, 100%);
    display: grid;
    place-items: center;
    background: transparent;
  }

  .radar-transparent-layer :global(canvas) {
    max-width: min(420px, 100%);
    height: auto !important;
    background: transparent !important;
    box-shadow: none !important;
  }

  .full-height-panel {
    min-height: 0;
    height: 100%;
    overflow: hidden;
  }

  .audio-panel {
    display: grid;
    grid-template-rows: auto minmax(0, 1fr);
    gap: 0.8rem;
  }

  .audio-panel :global(.card) {
    height: 100%;
  }

  .theme-panel {
    display: grid;
    grid-template-rows: minmax(0, 1fr);
    overflow: hidden;
  }

  .theme-panel :global(.card) {
    height: 100%;
  }

  .info-panel {
    display: grid;
    grid-template-rows: minmax(0, 1fr);
    overflow: hidden;
  }

  .info-panel :global(.card) {
    height: 100%;
  }

  @media (max-width: 1200px) {
    .workspace-grid {
      grid-template-columns: minmax(0, 1fr);
      grid-template-rows: minmax(380px, 1fr) 0;
    }

    .workspace-grid.has-active-panel {
      grid-template-rows: minmax(380px, 1fr) minmax(0, 1fr);
    }
  }
</style>
