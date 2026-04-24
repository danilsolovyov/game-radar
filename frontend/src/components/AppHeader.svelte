<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import Headphones from "@lucide/svelte/icons/headphones";
  import Settings2 from "@lucide/svelte/icons/settings-2";
  import Play from "@lucide/svelte/icons/play";
  import Square from "@lucide/svelte/icons/square";
  import Info from "@lucide/svelte/icons/info";
  import LocaleSwitcher from "./LocaleSwitcher.svelte";
  import { t } from "$lib/i18n";

  interface Props {
    appName: string;
    activePanel: "audio" | "theme" | "info" | null;
    isRadarActive?: boolean;
    onStartRadar: () => void;
    onStopRadar?: () => void;
    onToggleAudio: () => void;
    onToggleTheme: () => void;
    onToggleInfo: () => void;
    onExitApp?: () => void;
  }

  const {
    appName,
    activePanel,
    isRadarActive = false,
    onStartRadar,
    onStopRadar,
    onToggleAudio,
    onToggleTheme,
    onToggleInfo,
  }: Props = $props();

  function handleRadarToggle() {
    if (isRadarActive && onStopRadar) {
      onStopRadar();
    } else {
      onStartRadar();
    }
  }
</script>

<header class="app-header">
  <div class="header-left">
    <div class="app-info">
      <h1 class="app-title">{appName}</h1>
      <p class="app-subtitle">{$t.appTitle}</p>
    </div>
  </div>

  <div class="header-center-action">
    <Button
      class="radar-toggle-button"
      variant={isRadarActive ? "destructive" : "default"}
      onclick={handleRadarToggle}
    >
      {#if isRadarActive}
        <Square class="button-icon" />
        <span>{$t.stopRadar}</span>
      {:else}
        <Play class="button-icon" />
        <span>{$t.startRadar}</span>
      {/if}
    </Button>
  </div>

  <div class="header-right">
    <LocaleSwitcher />

    <div class="app-actions">
      <Button
        variant={activePanel === "audio" ? "secondary" : "ghost"}
        onclick={onToggleAudio}
        title={$t.audioDevices}
      >
        <Headphones class="action-icon" />
        <span>{$t.audio}</span>
      </Button>
      <Button
        variant={activePanel === "theme" ? "secondary" : "ghost"}
        onclick={onToggleTheme}
        title={$t.themeEditor}
      >
        <Settings2 class="action-icon" />
        <span>{$t.themeEditor}</span>
      </Button>
      <Button
        variant={activePanel === "info" ? "secondary" : "ghost"}
        onclick={onToggleInfo}
        title={$t.info}
      >
        <Info class="action-icon" />
        <span>{$t.info}</span>
      </Button>
    </div>
  </div>
</header>

<style>
  .app-header {
    border: 1px solid hsl(var(--border));
    background: #fff;
    border-radius: 1rem;
    padding: 1rem 1.5rem;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
    align-items: center;
    gap: 1.5rem;
  }

  .header-left {
    flex: 1;
    min-width: 0;
  }

  .app-info {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }

  .app-title {
    margin: 0;
    font-size: clamp(1.3rem, 2.2vw, 1.65rem);
    line-height: 1.2;
    font-weight: 700;
    letter-spacing: -0.01em;
    color: hsl(var(--foreground));
  }

  .app-subtitle {
    margin: 0;
    font-size: 0.75rem;
    color: hsl(var(--muted-foreground));
    font-weight: 400;
    line-height: 1.4;
  }

  .header-center-action {
    display: flex;
    justify-content: center;
    width: 100%;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    justify-content: flex-end;
  }

  .app-actions {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.2rem;
    background: hsl(var(--muted) / 0.4);
    border-radius: 0.75rem;
  }

  :global(.action-icon) {
    width: 1.125rem;
    height: 1.125rem;
  }

  :global(.radar-toggle-button) {
    min-height: 2.75rem;
    padding-inline: 1.5rem;
    font-size: 0.9rem;
    font-weight: 600;
    gap: 0.5rem;
    border-radius: 0.75rem;
    justify-content: center;
  }

  :global(.radar-toggle-button) :global(.button-icon) {
    width: 1.1rem;
    height: 1.1rem;
  }

  @media (max-width: 768px) {
    .app-header {
      grid-template-columns: 1fr;
      align-items: stretch;
      padding: 1.25rem;
      gap: 1rem;
    }

    .header-left {
      text-align: center;
    }

    .header-right {
      flex-wrap: wrap;
      justify-content: center;
    }

    .app-actions {
      gap: 0.5rem;
    }
  }

  @media (max-width: 480px) {
    .app-title {
      font-size: 1.5rem;
    }

    :global(.radar-toggle-button) {
      flex: 1;
      justify-content: center;
    }
  }
</style>
