<script lang="ts">
  import {
    clampByte,
    rgbaToHex,
    applyHexToColor,
    numberValue,
  } from "../utils/color";
  import type { color } from "../../wailsjs/go/models";
  import ThemeParamHint from "./ThemeParamHint.svelte";

  interface Props {
    label: string;
    colorValue: color.RGBA;
    onUpdate: (color: color.RGBA) => void;
    /** Alpha slider label (from i18n) */
    alphaLabel: string;
    /** Preview title attribute (from i18n) */
    previewTitle: string;
    tooltipDescription?: string;
    tooltipAriaLabel?: string;
  }

  let {
    label,
    colorValue,
    onUpdate,
    alphaLabel,
    previewTitle,
    tooltipDescription,
    tooltipAriaLabel,
  }: Props = $props();

  function handleHexChange(event: Event) {
    const hex = (event.target as HTMLInputElement).value;
    applyHexToColor(colorValue, hex);
    onUpdate(colorValue);
  }

  function handleAlphaInput(event: Event) {
    const raw = Number.parseInt((event.target as HTMLInputElement).value, 10);
    colorValue.A = clampByte(Number.isFinite(raw) ? raw : colorValue.A);
    onUpdate(colorValue);
  }

  const alphaNorm = $derived(
    Math.max(0, Math.min(1, numberValue(colorValue.A, 255) / 255)),
  );
</script>

<div class="space-y-2">
  <div class="flex items-center gap-1.5">
    <div class="text-sm font-medium">{label}</div>
    {#if tooltipDescription && tooltipAriaLabel}
      <ThemeParamHint
        description={tooltipDescription}
        ariaLabel={tooltipAriaLabel}
      />
    {/if}
  </div>
  <div class="flex flex-wrap items-center gap-3">
    <input
      type="color"
      value={rgbaToHex(colorValue)}
      onchange={handleHexChange}
      class="h-9 w-14 shrink-0 rounded border border-input bg-background p-1"
    />
    <div
      class="h-9 w-9 shrink-0 rounded border border-input"
      style={`background-color: rgba(${clampByte(colorValue.R)}, ${clampByte(colorValue.G)}, ${clampByte(colorValue.B)}, ${alphaNorm});`}
      title={previewTitle}
    ></div>
    <div class="flex min-w-[140px] flex-1 flex-col gap-1">
      <span class="text-xs text-muted-foreground">{alphaLabel}</span>
      <input
        class="w-full accent-primary"
        type="range"
        min="0"
        max="255"
        step="1"
        value={clampByte(colorValue.A)}
        oninput={handleAlphaInput}
      />
    </div>
  </div>
</div>
