<script lang="ts">
  import AppHeader from "./components/AppHeader.svelte";
  import AppSidebar from "./components/AppSidebar.svelte";
  import AudioDevices from "./AudioDevices.svelte";
  import InfoPanel from "./components/InfoPanel.svelte";
  import Radar from "./Radar.svelte";
  import ThemeEditor from "./ThemeEditor.svelte";

  import {
    GetAppName,
    IsOverlayMode,
    UpdateOverlayStatus,
  } from "wailsjs/go/app/App";
  import { SetNormalWindow, SetRadarOverlay } from "wailsjs/go/app/App";
  import { EventsOn, Quit } from "wailsjs/runtime/runtime";

  import { safeLogInfo, safeLogError } from "./utils/logger";
  import { initLocaleSync } from "$lib/i18n";
  import { onMount } from "svelte";

  let overlayMode = $state(false);
  let appElement: HTMLElement | null = $state(null);
  let radarRenderKey = $state(0);
  let appName = $state("Game Radar");
  let activePanel: "audio" | "theme" | "info" | null = $state(null);

  async function enterRadarMode() {
    overlayMode = true;
    safeLogInfo("Switching to radar overlay mode");
    await UpdateOverlayStatus(true);
    await SetRadarOverlay();
  }

  async function exitRadarMode() {
    overlayMode = false;
    safeLogInfo("Switching to normal window mode");
    await UpdateOverlayStatus(false);
    await SetNormalWindow();
  }

  function refreshTheme() {
    radarRenderKey += 1;
  }

  function toggleAudioPanel() {
    activePanel = activePanel === "audio" ? null : "audio";
  }

  function toggleThemePanel() {
    activePanel = activePanel === "theme" ? null : "theme";
  }

  function toggleInfoPanel() {
    activePanel = activePanel === "info" ? null : "info";
  }

  $effect(() => {
    if (appElement) {
      if (overlayMode) {
        appElement.classList.add("radar-mode");
      } else {
        appElement.classList.remove("radar-mode");
      }
    }
  });

  onMount(() => {
    appElement = document.getElementById("app");
    let unsubscribeLocaleSync: (() => void) | undefined;

    GetAppName()
      .then((resolvedName) => {
        appName = String(resolvedName ?? "Game Radar");
      })
      .catch((err) => {
        safeLogInfo(`GetAppName failed, using fallback: ${String(err)}`);
      });

    const unsubscribeOverlayStatus = EventsOn(
      "overlay-status",
      (status: boolean) => {
        overlayMode = Boolean(status);
      },
    );

    IsOverlayMode()
      .then((status) => {
        overlayMode = Boolean(status);
      })
      .catch((err) => {
        safeLogInfo(`IsOverlayMode failed: ${String(err)}`);
      });

    initLocaleSync()
      .then((unsubscribe) => {
        unsubscribeLocaleSync = unsubscribe;
      })
      .catch((err) => {
        safeLogError("Failed to initialize locale sync", err);
      });

    return () => {
      unsubscribeOverlayStatus?.();
      unsubscribeLocaleSync?.();
    };
  });
</script>

<div class="app-shell" class:radar-mode={overlayMode}>
  {#if overlayMode}
    <div class="overlay-root">
      {#key radarRenderKey}
        <Radar />
      {/key}
    </div>
  {:else}
    <main class="main-layout">
      <AppHeader
        {appName}
        {activePanel}
        isRadarActive={overlayMode}
        onStartRadar={enterRadarMode}
        onStopRadar={exitRadarMode}
        onToggleAudio={toggleAudioPanel}
        onToggleTheme={toggleThemePanel}
        onToggleInfo={toggleInfoPanel}
      />

      <AppSidebar {activePanel}>
        {#snippet radarContent()}
          {#key radarRenderKey}
            <Radar />
          {/key}
        {/snippet}
        {#snippet audioPanel()}
          <AudioDevices />
        {/snippet}
        {#snippet themePanel()}
          <ThemeEditor on:themechange={refreshTheme} />
        {/snippet}
        {#snippet infoPanel()}
          <InfoPanel />
        {/snippet}
      </AppSidebar>
    </main>
  {/if}
</div>

<style>
  .app-shell {
    width: 100%;
    min-height: 100vh;
    color: hsl(var(--foreground));
  }

  .app-shell.radar-mode {
    min-height: 0;
    width: 100%;
    height: 100%;
  }

  .overlay-root {
    width: 100%;
    height: 100%;
    margin: 0;
    padding: 0;
    display: grid;
    place-items: center;
    background: transparent;
  }

  .main-layout {
    height: 100vh;
    padding: 1rem;
    display: grid;
    grid-template-rows: auto minmax(0, 1fr);
    gap: 1rem;
    overflow: hidden;
  }

  @media (max-width: 768px) {
    .main-layout {
      padding: 0.75rem;
      gap: 0.75rem;
    }
  }
</style>
