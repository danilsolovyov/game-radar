<script lang="ts">
  import { cn } from "$lib/utils.js";
  import ThemeParamHint from "./ThemeParamHint.svelte";

  interface Props {
    label: string;
    value: number;
    step?: number;
    min?: number;
    max?: number;
    placeholder?: string;
    showSlider?: boolean;
    /** Range hint shown below the controls */
    hint?: string;
    unit?: string;
    onUpdate: (value: number) => void;
    tooltipDescription?: string;
    tooltipAriaLabel?: string;
  }

  let {
    label,
    value,
    step = 1,
    min,
    max,
    placeholder,
    showSlider = false,
    hint,
    unit,
    onUpdate,
    tooltipDescription,
    tooltipAriaLabel,
  }: Props = $props();

  const inputClass = cn(
    "dark:bg-input/30 border-input focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 disabled:bg-input/50 dark:disabled:bg-input/80 h-8 rounded-lg border bg-transparent px-2.5 py-1 text-base transition-colors focus-visible:ring-3 aria-invalid:ring-3 md:text-sm text-foreground placeholder:text-muted-foreground w-full min-w-0 outline-none disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50",
  );

  function clamp(n: number): number {
    let v = n;
    if (min !== undefined) v = Math.max(min, v);
    if (max !== undefined) v = Math.min(max, v);
    return v;
  }

  function roundToStep(n: number, s: number): number {
    if (!Number.isFinite(s) || s <= 0) return n;
    const inv = 1 / s;
    return Math.round(n * inv) / inv;
  }

  function normalize(raw: number): number {
    let v = clamp(raw);
    v = roundToStep(v, step);
    if (min !== undefined) v = Math.max(min, v);
    if (max !== undefined) v = Math.min(max, v);
    return v;
  }

  function decimalsFromStep(s: number): number {
    const t = s.toString();
    if (t.includes("e") || t.includes("E")) return 8;
    const i = t.indexOf(".");
    if (i < 0) return 0;
    return t.length - i - 1;
  }

  /** One normalized value for input and slider; always in sync after parent updates. */
  const syncedValue = $derived(normalize(value));

  const displayValue = $derived.by(() => {
    const v = syncedValue;
    if (!Number.isFinite(v)) return "";
    const d = decimalsFromStep(step);
    if (d === 0) return String(Math.round(v));
    return v.toFixed(d);
  });

  function parseAndUpdate(rawValue: string) {
    const newValue = Number.parseFloat(rawValue);
    if (Number.isNaN(newValue)) return;
    onUpdate(normalize(newValue));
  }
</script>

<label class="space-y-1">
  <span class="flex items-center justify-between gap-2 text-sm">
    <span class="inline-flex items-center gap-1.5">
      <span>{label}</span>
      {#if tooltipDescription && tooltipAriaLabel}
        <ThemeParamHint
          description={tooltipDescription}
          ariaLabel={tooltipAriaLabel}
        />
      {/if}
    </span>
    {#if unit}
      <span class="text-xs text-muted-foreground">{unit}</span>
    {/if}
  </span>
  <div class="space-y-2">
    <input
      type="number"
      class={inputClass}
      {step}
      {min}
      {max}
      {placeholder}
      value={displayValue}
      oninput={(e) => parseAndUpdate((e.target as HTMLInputElement).value)}
    />
    {#if showSlider}
      <input
        class="w-full accent-primary"
        type="range"
        {min}
        {max}
        {step}
        value={syncedValue}
        oninput={(e) => parseAndUpdate((e.target as HTMLInputElement).value)}
      />
    {/if}
  </div>
  {#if hint}
    <span class="block text-xs text-muted-foreground">{hint}</span>
  {/if}
</label>
