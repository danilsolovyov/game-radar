<script lang="ts">
    import {
        Select,
        SelectContent,
        SelectTrigger,
        SelectItem,
    } from "$lib/components/ui/select";
    import type { models } from "wailsjs/go/models";
    import { t } from "$lib/i18n";

    interface Props {
        themes: models.Theme[];
        selectedThemeName: string;
        onSelect: (name: string) => void;
    }

    const { themes, selectedThemeName, onSelect }: Props = $props();

    function handleSelect(value: string) {
        onSelect(value);
    }
</script>

<Select type="single" value={selectedThemeName} onValueChange={handleSelect}>
    <SelectTrigger class="w-full min-w-0">
        <span class="block truncate">
            {themes.find((t) => t.name === selectedThemeName)?.name ||
                $t.selectTheme}
        </span>
    </SelectTrigger>
    <SelectContent class="max-w-[min(80vw,36rem)]">
        {#each themes as theme (theme.name)}
            <SelectItem value={theme.name}>
                <div class="flex items-center gap-2 min-w-0 w-full">
                    <span class="truncate">{theme.name}</span>
                </div>
            </SelectItem>
        {/each}
    </SelectContent>
</Select>
