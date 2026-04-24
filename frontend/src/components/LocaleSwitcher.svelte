<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import {
        DropdownMenu,
        DropdownMenuContent,
        DropdownMenuItem,
        DropdownMenuTrigger,
    } from "$lib/components/ui/dropdown-menu";
    import { locale, availableLocales, changeLocale } from "$lib/i18n";
    import type { LocaleOption } from "$lib/i18n/locales";

    function handleLocaleSelect(loc: LocaleOption) {
        changeLocale(loc.code);
    }
</script>

<div class="locale-switcher">
    <DropdownMenu>
        <DropdownMenuTrigger>
            <Button variant="ghost" size="sm" class="locale-button">
                <span class="current-flag">
                    {$availableLocales.find((l) => l.code === $locale)?.flag ??
                        "🌐"}
                </span>
                <span class="current-name">
                    {$availableLocales.find((l) => l.code === $locale)?.name ??
                        "Language"}
                </span>
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
            {#each $availableLocales as loc}
                <DropdownMenuItem
                    onclick={() => handleLocaleSelect(loc)}
                    class={loc.code === $locale ? "active-locale" : ""}
                >
                    <span class="locale-option">
                        <span class="locale-flag">{loc.flag}</span>
                        <span class="locale-name">{loc.name}</span>
                    </span>
                </DropdownMenuItem>
            {/each}
        </DropdownMenuContent>
    </DropdownMenu>
</div>

<style>
    .locale-switcher {
        min-width: 130px;
    }

    .locale-button {
        display: flex;
        align-items: center;
        gap: 0.375rem;
        font-size: 0.875rem;
    }

    .current-flag {
        font-size: 1.1em;
    }

    .current-name {
        font-size: 0.875rem;
    }

    .locale-option {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .locale-flag {
        font-size: 1.1em;
    }

    .locale-name {
        font-size: 0.875rem;
    }

    :global(.active-locale) {
        background-color: hsl(var(--accent));
    }
</style>
