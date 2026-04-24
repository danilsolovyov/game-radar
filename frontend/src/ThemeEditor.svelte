<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import { models } from "wailsjs/go/models";
  import {
    DeleteTheme,
    GetCurrentThemeName,
    GetThemes,
    SetTheme,
  } from "wailsjs/go/app/App";

  import Card from "$lib/components/ui/card/card.svelte";
  import CardContent from "$lib/components/ui/card/card-content.svelte";
  import CardHeader from "$lib/components/ui/card/card-header.svelte";
  import CardTitle from "$lib/components/ui/card/card-title.svelte";
  import Button from "$lib/components/ui/button/button.svelte";
  import Checkbox from "$lib/components/ui/checkbox/checkbox.svelte";
  import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
    AlertDialogAction,
    AlertDialogCancel,
  } from "$lib/components/ui/alert-dialog";
  import Save from "@lucide/svelte/icons/save";
  import Trash2 from "@lucide/svelte/icons/trash-2";
  import * as Tooltip from "$lib/components/ui/tooltip/index.js";

  import ThemeSelector from "./components/ThemeSelector.svelte";
  import ThemeParamHint from "./components/ThemeParamHint.svelte";
  import ColorPicker from "./components/ColorPicker.svelte";
  import NumberInput from "./components/NumberInput.svelte";
  import AudioPlayer from "./AudioPlayer.svelte";

  import { safeLogError, safeLogInfo } from "./utils/logger";
  import { cloneTheme } from "./utils/color";
  import { t } from "$lib/i18n";

  const dispatch = createEventDispatcher<{
    themechange: models.Theme;
  }>();

  let themes: models.Theme[] = $state([]);
  let selectedThemeName = $state("");
  let draftTheme: models.Theme | null = $state(null);
  let error = $state("");
  let loading = $state(false);
  let saving = $state(false);
  let deleteDialogOpen = $state(false);
  let showBlips = $state(false);

  function createEmptyThemeDraft(): models.Theme {
    return new models.Theme({
      name: "",
      background_color: { R: 8, G: 12, B: 18, A: 230 },
      radar_color: { R: 48, G: 151, B: 255, A: 255 },
      intensity_multiplier: 1,
      size: 320,
      pos_x: 220,
      pos_y: 220,
      section_count: 8,
      ring_count: 4,
      border_width: 2,
      border_opacity: 0.55,
      section_base_opacity: 0.25,
      section_bright_opacity: 0.85,
      section_timeout: 450,
      show_blips: true,
      blip_size: 8,
      blip_opacity: 0.95,
      blip_timeout: 650,
    });
  }

  async function loadThemes() {
    if (loading) return;
    loading = true;
    error = "";

    try {
      const [themesResult, currentName] = await Promise.all([
        GetThemes(),
        GetCurrentThemeName(),
      ]);
      themes = Array.isArray(themesResult)
        ? themesResult.map((item) => cloneTheme(item))
        : [];
      selectedThemeName =
        currentName && themes.some((theme) => theme.name === currentName)
          ? currentName
          : themes[0]?.name || "";
      const selected = themes.find((theme) => theme.name === selectedThemeName);
      draftTheme = selected ? cloneTheme(selected) : createEmptyThemeDraft();
      showBlips = draftTheme?.show_blips ?? false;
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
      themes = [];
      selectedThemeName = "";
      draftTheme = createEmptyThemeDraft();
      showBlips = draftTheme.show_blips ?? false;
      safeLogError("Failed to load themes", e);
    } finally {
      loading = false;
    }
  }

  async function onThemeSelect(name: string) {
    selectedThemeName = name;
    const selected = themes.find((theme) => theme.name === selectedThemeName);
    draftTheme = selected ? cloneTheme(selected) : createEmptyThemeDraft();
    showBlips = draftTheme?.show_blips ?? false;

    if (selected) {
      saving = true;
      error = "";
      try {
        await SetTheme(selected);
        safeLogInfo(`Theme activated: ${selected.name}`);
        dispatch("themechange", selected);
      } catch (e) {
        error = e instanceof Error ? e.message : String(e);
        safeLogError("Failed to activate theme", e);
      } finally {
        saving = false;
      }
    }
  }

  async function saveTheme() {
    if (!draftTheme) return;
    saving = true;
    error = "";

    try {
      const payload = cloneTheme(draftTheme);
      payload.name = payload.name.trim();
      if (!payload.name) {
        throw new Error($t.themeNameRequired);
      }
      await SetTheme(payload);
      safeLogInfo(`Theme saved: ${payload.name}`);
      dispatch("themechange", payload);
      await loadThemes();
      onThemeSelect(payload.name);
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
      safeLogError("Failed to save theme", e);
    } finally {
      saving = false;
    }
  }

  async function deleteSelectedTheme() {
    if (!selectedThemeName) return;
    saving = true;
    error = "";
    deleteDialogOpen = false;

    try {
      await DeleteTheme(selectedThemeName);
      safeLogInfo(`Theme deleted: ${selectedThemeName}`);
      await loadThemes();
      const selected =
        themes.find((theme) => theme.name === selectedThemeName) || themes[0];
      if (selected) {
        showBlips = selected.show_blips ?? false;
        dispatch("themechange", selected);
      }
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
      safeLogError("Failed to delete theme", e);
    } finally {
      saving = false;
    }
  }

  /** Create a new theme object so Svelte 5 detects updates (mutating Theme class fields does not trigger children). */
  function updateThemeField<K extends keyof models.Theme>(
    field: K,
    value: models.Theme[K],
  ) {
    if (!draftTheme) return;
    const next = cloneTheme(draftTheme);
    (next as unknown as Record<string, unknown>)[field as string] = value;
    draftTheme = next;
  }

  onMount(() => {
    loadThemes();
  });
</script>

<Tooltip.Provider delayDuration={0}>
  <Card class="w-full min-w-0 min-h-0 mx-auto overflow-hidden border-2">
    <CardHeader class="space-y-2 pb-3">
      <div class="space-y-3">
        <div>
          <CardTitle>
            <h2 class="text-lg font-semibold">
              {$t.themeEditorTitle}
            </h2>
          </CardTitle>
        </div>
        <div
          class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto_auto] lg:items-center"
        >
          <div class="min-w-0">
            <ThemeSelector
              {themes}
              {selectedThemeName}
              onSelect={onThemeSelect}
            />
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <Button
              onclick={saveTheme}
              variant="outline"
              class="shrink-0"
              title={$t.saveTheme}
              disabled={saving || !draftTheme?.name?.trim()}
            >
              <Save class="h-4 w-4" />
              <span>{$t.saveTheme}</span>
            </Button>
            <div class="h-6 w-px bg-border" aria-hidden="true"></div>
            <AlertDialog bind:open={deleteDialogOpen}>
              <AlertDialogTrigger>
                <Button
                  variant="destructive"
                  class="shrink-0"
                  title={$t.deleteTheme}
                  disabled={!selectedThemeName || saving}
                >
                  <Trash2 class="h-4 w-4" />
                  <span>{$t.deleteTheme}</span>
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>{$t.deleteTheme}</AlertDialogTitle>
                  <AlertDialogDescription>
                    {$t.confirmDeleteTheme}
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>{$t.cancel}</AlertDialogCancel>
                  <AlertDialogAction onclick={deleteSelectedTheme}>
                    {$t.delete}
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
          <div class="lg:justify-self-end">
            <AudioPlayer />
          </div>
        </div>
      </div>
    </CardHeader>

    <CardContent class="space-y-4 min-h-0 overflow-auto">
      {#if draftTheme}
        <div class="rounded-lg border bg-muted/10 p-3 space-y-5">
          {#if themes.length === 0}
            <p class="text-sm text-muted-foreground">
              {$t.noThemeSelected}
            </p>
          {/if}
          <section class="space-y-3">
            <h3 class="text-sm font-semibold">
              {$t.themeBasicSettings}
            </h3>
            <div class="grid gap-3 md:grid-cols-2">
              <label class="space-y-1">
                <span class="flex items-center gap-1.5 text-sm">
                  <span>{$t.themeName}</span>
                  <ThemeParamHint
                    description={$t.themeHintThemeName}
                    ariaLabel={$t.themeParamTooltipAria}
                  />
                </span>
                <input
                  type="text"
                  bind:value={draftTheme.name}
                  class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                />
              </label>
              <NumberInput
                label={$t.intensityMultiplier}
                value={draftTheme.intensity_multiplier}
                step={0.1}
                min={0.1}
                max={15}
                showSlider={true}
                hint="0.1 - 15.0"
                tooltipDescription={$t.themeHintIntensityMultiplier}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("intensity_multiplier", v)}
              />
            </div>
          </section>

          <section class="space-y-3">
            <h3 class="text-sm font-semibold">{$t.colorsSection}</h3>
            <div class="grid gap-3 md:grid-cols-2">
              <ColorPicker
                label={$t.backgroundColor}
                colorValue={draftTheme.background_color}
                alphaLabel={$t.themeColorAlpha}
                previewTitle={$t.themeColorPreview}
                tooltipDescription={$t.themeHintBackgroundColor}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(c) => updateThemeField("background_color", c)}
              />
              <ColorPicker
                label={$t.gridColor}
                colorValue={draftTheme.radar_color}
                alphaLabel={$t.themeColorAlpha}
                previewTitle={$t.themeColorPreview}
                tooltipDescription={$t.themeHintGridColor}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(c) => updateThemeField("radar_color", c)}
              />
            </div>
          </section>

          <section class="space-y-3">
            <h3 class="text-sm font-semibold">
              {$t.positionSizeSection}
            </h3>
            <div class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
              <NumberInput
                label={$t.sizeLabel}
                value={draftTheme.size}
                min={120}
                max={1200}
                step={10}
                showSlider={true}
                hint="120 - 1200"
                tooltipDescription={$t.themeHintSize}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("size", v)}
              />
              <NumberInput
                label={$t.posX}
                value={draftTheme.pos_x}
                tooltipDescription={$t.themeHintPosX}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("pos_x", v)}
              />
              <NumberInput
                label={$t.posY}
                value={draftTheme.pos_y}
                tooltipDescription={$t.themeHintPosY}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("pos_y", v)}
              />
            </div>
          </section>

          <section class="space-y-3">
            <h3 class="text-sm font-semibold">{$t.geometrySection}</h3>
            <div class="grid gap-3 md:grid-cols-3">
              <NumberInput
                label={$t.sectors}
                value={draftTheme.section_count}
                tooltipDescription={$t.themeHintSectors}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("section_count", v)}
              />
              <NumberInput
                label={$t.rings}
                value={draftTheme.ring_count}
                tooltipDescription={$t.themeHintRings}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("ring_count", v)}
              />
              <NumberInput
                label={$t.borderWidth}
                value={draftTheme.border_width}
                tooltipDescription={$t.themeHintBorderWidth}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("border_width", v)}
              />
            </div>
          </section>

          <section class="space-y-3">
            <h3 class="text-sm font-semibold">{$t.effectsSection}</h3>

            <div class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
              <NumberInput
                label={$t.borderOpacity}
                value={draftTheme.border_opacity}
                step={0.01}
                min={0}
                max={1}
                showSlider={true}
                hint="0 - 1"
                tooltipDescription={$t.themeHintBorderOpacity}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("border_opacity", v)}
              />
              <NumberInput
                label={$t.sectionBaseOpacity}
                value={draftTheme.section_base_opacity}
                step={0.01}
                min={0}
                max={1}
                showSlider={true}
                hint="0 - 1"
                tooltipDescription={$t.themeHintSectionBaseOpacity}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("section_base_opacity", v)}
              />
              <NumberInput
                label={$t.sectionBrightOpacity}
                value={draftTheme.section_bright_opacity}
                step={0.01}
                min={0}
                max={1}
                showSlider={true}
                hint="0 - 1"
                tooltipDescription={$t.themeHintSectionBrightOpacity}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("section_bright_opacity", v)}
              />
              <NumberInput
                label={$t.sectionTimeout}
                value={draftTheme.section_timeout}
                min={50}
                max={5000}
                step={50}
                showSlider={true}
                hint={`50 - 5000 ${$t.themeUnitMs}`}
                unit={$t.themeUnitMs}
                tooltipDescription={$t.themeHintSectionTimeout}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("section_timeout", v)}
              />
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                id="showBlips"
                bind:checked={showBlips}
                onCheckedChange={(v) => {
                  const newValue = v === true;
                  showBlips = newValue;
                  updateThemeField("show_blips", newValue);
                }}
              />
              <label for="showBlips" class="cursor-pointer text-sm">
                {$t.showBlips}
              </label>
              <ThemeParamHint
                description={$t.themeHintShowBlips}
                ariaLabel={$t.themeParamTooltipAria}
              />
            </div>
            <div
              class="grid gap-3 md:grid-cols-2 lg:grid-cols-3"
              class:opacity-50={!showBlips}
              class:pointer-events-none={!showBlips}
            >
              <NumberInput
                label={$t.blipSize}
                value={draftTheme.blip_size}
                min={1}
                max={40}
                step={1}
                showSlider={true}
                hint="1 - 40"
                tooltipDescription={$t.themeHintBlipSize}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("blip_size", v)}
              />
              <NumberInput
                label={$t.blipOpacity}
                value={draftTheme.blip_opacity}
                step={0.01}
                min={0}
                max={1}
                showSlider={true}
                hint="0 - 1"
                tooltipDescription={$t.themeHintBlipOpacity}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("blip_opacity", v)}
              />
              <NumberInput
                label={$t.blipTimeout}
                value={draftTheme.blip_timeout}
                min={50}
                max={5000}
                step={50}
                showSlider={true}
                hint={`50 - 5000 ${$t.themeUnitMs}`}
                unit={$t.themeUnitMs}
                tooltipDescription={$t.themeHintBlipTimeout}
                tooltipAriaLabel={$t.themeParamTooltipAria}
                onUpdate={(v) => updateThemeField("blip_timeout", v)}
              />
            </div>
          </section>
        </div>
      {/if}
    </CardContent>
  </Card>
</Tooltip.Provider>
