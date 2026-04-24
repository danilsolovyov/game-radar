<script lang="ts">
  import { onMount } from "svelte";
  import Check from "@lucide/svelte/icons/check";
  import Copy from "@lucide/svelte/icons/copy";
  import Card from "$lib/components/ui/card/card.svelte";
  import CardContent from "$lib/components/ui/card/card-content.svelte";
  import CardHeader from "$lib/components/ui/card/card-header.svelte";
  import CardTitle from "$lib/components/ui/card/card-title.svelte";
  import Button from "$lib/components/ui/button/button.svelte";
  import type { Translations } from "$lib/i18n";
  import { t } from "$lib/i18n";
  import { GetAppName, GetVersion } from "wailsjs/go/app/App";
  import { BrowserOpenURL } from "wailsjs/runtime/runtime";

  /** Donation links, synchronized with README.md */
  const donationPlatformLinks = [
    {
      labelKey: "infoDonationDonateStream" as const,
      url: "https://donate.stream/soldan-gameradar",
    },
  ] as const;

  type WalletLabelKey =
    | "infoWalletUsdtTrc"
    | "infoWalletUsdcPoly"
    | "infoWalletEth"
    | "infoWalletBtc"
    | "infoWalletLtc";

  type DonationLabelKey = "infoDonationDonateStream";

  const walletRows: {
    labelKey: WalletLabelKey;
    address: string;
    trustWalletUrl: string;
  }[] = [
    {
      labelKey: "infoWalletUsdtTrc",
      address: "TBxkDaADAbVk2VH3o3pTYamnnhCN3dSR1R",
      trustWalletUrl:
        "https://link.trustwallet.com/send?coin=195&address=TBxkDaADAbVk2VH3o3pTYamnnhCN3dSR1R&token_id=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
    },
    {
      labelKey: "infoWalletUsdcPoly",
      address: "0x6ceBb4a1EC0b50C3C68C8F5A09aA2ae4c944c4e0",
      trustWalletUrl:
        "https://link.trustwallet.com/send?coin=966&address=0x6ceBb4a1EC0b50C3C68C8F5A09aA2ae4c944c4e0&token_id=0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359",
    },
    {
      labelKey: "infoWalletEth",
      address: "0x6ceBb4a1EC0b50C3C68C8F5A09aA2ae4c944c4e0",
      trustWalletUrl:
        "https://link.trustwallet.com/send?coin=60&address=0x6ceBb4a1EC0b50C3C68C8F5A09aA2ae4c944c4e0",
    },
    {
      labelKey: "infoWalletBtc",
      address: "bc1qs3svtnv04tl23fyweq34l0jpny2ymftrv7tad9",
      trustWalletUrl:
        "https://link.trustwallet.com/send?coin=0&address=bc1qs3svtnv04tl23fyweq34l0jpny2ymftrv7tad9",
    },
    {
      labelKey: "infoWalletLtc",
      address: "ltc1qf7de8mtdeczmv8vz97u7ylfz8kh06ejnxhrczs",
      trustWalletUrl:
        "https://link.trustwallet.com/send?coin=2&address=ltc1qf7de8mtdeczmv8vz97u7ylfz8kh06ejnxhrczs",
    },
  ];

  const cryptobotUrl = "https://t.me/send?start=IVGukNPxmSM0";

  const contacts = [
    {
      label: "Issues",
      value: "github.com/danilsolovyov/game-radar/issues",
      url: "https://github.com/danilsolovyov/game-radar/issues",
    },
    {
      label: "Telegram chat",
      value: "t.me/game_radar_chat",
      url: "https://t.me/game_radar_chat",
    },
    {
      label: "Releases",
      value: "github.com/danilsolovyov/game-radar/releases",
      url: "https://github.com/danilsolovyov/game-radar/releases",
    },
  ] as const;

  function tr(
    messages: Translations,
    key: WalletLabelKey | DonationLabelKey,
  ): string {
    return messages[key];
  }

  let copiedWalletKey = $state("");
  let copiedWalletTimer: ReturnType<typeof setTimeout> | null = null;
  let appVersion = $state("dev");
  let appName = $state("Sound Radar");

  async function loadAppName() {
    try {
      const name = await GetAppName();
      appName = name?.trim() ? name : "Sound Radar";
    } catch (err) {
      console.error("Failed to get app name:", err);
      appName = "Sound Radar";
    }
  }

  async function loadAppVersion() {
    try {
      const version = await GetVersion();
      appVersion = version?.trim() ? version : "dev";
    } catch (err) {
      console.error("Failed to get app version:", err);
      appVersion = "dev";
    }
  }

  onMount(() => {
    void loadAppName();
    void loadAppVersion();
  });

  async function copyWalletAddress(walletKey: string, address: string) {
    try {
      await navigator.clipboard.writeText(address);
      copiedWalletKey = walletKey;

      if (copiedWalletTimer) clearTimeout(copiedWalletTimer);
      copiedWalletTimer = setTimeout(() => {
        if (copiedWalletKey === walletKey) copiedWalletKey = "";
      }, 1600);
    } catch (err) {
      console.error("Failed to copy wallet address:", err);
    }
  }

  function openExternalUrl(url: string) {
    try {
      BrowserOpenURL(url);
    } catch (err) {
      console.error("Failed to open external URL:", err);
    }
  }
</script>

<Card class="w-full min-w-0 mx-auto overflow-hidden border-2">
  <CardHeader class="space-y-2 pb-3">
    <CardTitle>
      <h2 class="text-lg font-semibold">{$t.infoPanelTitle}</h2>
    </CardTitle>
    <p class="text-sm text-muted-foreground">
      {$t.infoPanelSubtitle}
    </p>
  </CardHeader>

  <CardContent class="space-y-4 overflow-auto">
    <section class="rounded-lg border bg-muted/10 p-4 space-y-4">
      <div class="space-y-3">
        <h3 class="text-sm font-semibold">{$t.infoAboutAppTitle}</h3>
        <p class="text-sm leading-relaxed text-muted-foreground">
          <strong class="font-semibold text-foreground">{appName}</strong>
          — {$t.infoReadmeTagline}
        </p>
        <p class="text-sm leading-relaxed text-muted-foreground">
          {$t.infoAboutAppText}
        </p>

        <div class="about-meta">
          <p class="text-sm text-muted-foreground">
            <strong class="font-semibold text-foreground">
              {$t.infoDeveloperLabel}:
            </strong>
            {$t.infoDeveloperText}
          </p>
          <p class="text-xs text-muted-foreground">
            {$t.infoAppVersionLabel}: {appVersion}
          </p>
        </div>
      </div>

      <div class="space-y-2 border-t border-border/70 pt-3">
        <h3 class="text-sm font-semibold">{$t.infoWhatItDoesTitle}</h3>
        <p class="text-sm leading-relaxed text-muted-foreground">
          {$t.infoWhatItDoesText}
        </p>
      </div>

      <div class="space-y-2 border-t border-border/70 pt-3">
        <h3 class="text-sm font-semibold">{$t.infoHowItWorksTitle}</h3>
        <ol class="list-decimal space-y-1.5 pl-4 text-sm text-muted-foreground">
          <li>{$t.infoHowItWorksStep1}</li>
          <li>{$t.infoHowItWorksStep2}</li>
          <li>{$t.infoHowItWorksStep3}</li>
          <li>{$t.infoHowItWorksStep4}</li>
        </ol>
      </div>

      <div class="space-y-2 border-t border-border/70 pt-3">
        <h3 class="text-sm font-semibold">{$t.infoUsageTitle}</h3>
        <p class="text-sm leading-relaxed text-muted-foreground">
          {$t.infoUsageText}
        </p>
      </div>

      <div class="space-y-2 border-t border-border/70 pt-3">
        <h3 class="text-sm font-semibold">{$t.infoAnticheatTitle}</h3>
        <p class="text-sm leading-relaxed text-muted-foreground">
          {$t.infoAnticheatLead}
        </p>
        <ul class="list-disc space-y-1 pl-4 text-sm text-muted-foreground">
          <li>{$t.infoAnticheatBullet1}</li>
          <li>{$t.infoAnticheatBullet2}</li>
          <li>{$t.infoAnticheatBullet3}</li>
        </ul>
        <p class="text-sm leading-relaxed text-muted-foreground">
          {$t.infoAnticheatFooter}
        </p>
      </div>

      <div class="space-y-2 border-t border-border/70 pt-3">
        <h3 class="text-sm font-semibold">{$t.contactsTitle}</h3>
        <div class="grid gap-2 sm:grid-cols-2">
          {#each contacts as contact (contact.url)}
            <Button
              type="button"
              variant="outline"
              class="justify-start min-w-0 external-link-button"
              onclick={() => openExternalUrl(contact.url)}
            >
              <span class="text-muted-foreground">{contact.label}:</span>
              <span class="truncate">{contact.value}</span>
            </Button>
          {/each}
        </div>
      </div>
    </section>

    <section class="rounded-lg border bg-muted/10 p-4 space-y-4">
      <div class="space-y-2">
        <h3 class="text-sm font-semibold">{$t.infoDonationsTitle}</h3>
        <p class="text-sm leading-relaxed text-muted-foreground">
          {$t.infoDonationsText}
        </p>
      </div>

      <div class="space-y-2 border-t border-border/70 pt-3">
        <h4
          class="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
        >
          {$t.infoDonationsRussiaTitle}
        </h4>
        <div class="flex flex-wrap gap-2">
          {#each donationPlatformLinks as link (link.url)}
            <Button
              type="button"
              variant="outline"
              class="external-link-button"
              onclick={() => openExternalUrl(link.url)}
            >
              {tr($t, link.labelKey)}
            </Button>
          {/each}
        </div>
      </div>

      <div class="space-y-3 border-t border-border/70 pt-3">
        <h4
          class="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
        >
          {$t.infoDonationsInternationalTitle}
        </h4>
        <p class="text-sm text-muted-foreground">
          {$t.infoDonationsCryptoIntro}
        </p>
        <div class="space-y-3">
          {#each walletRows as wallet (wallet.trustWalletUrl)}
            <div class="wallet-item rounded-md border bg-background/70 p-3">
              <div class="space-y-1">
                <p class="text-xs font-medium text-muted-foreground">
                  {tr($t, wallet.labelKey)}
                </p>
                <p class="wallet-address text-sm">{wallet.address}</p>
              </div>
              <div class="wallet-actions pt-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onclick={() =>
                    copyWalletAddress(wallet.trustWalletUrl, wallet.address)}
                >
                  {#if copiedWalletKey === wallet.trustWalletUrl}
                    <Check class="copy-icon copied-icon" />
                    {$t.copied}
                  {:else}
                    <Copy class="copy-icon" />
                    {$t.copyAddress}
                  {/if}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  class="external-link-button"
                  onclick={() => openExternalUrl(wallet.trustWalletUrl)}
                >
                  {$t.openTrustWallet}
                </Button>
              </div>
            </div>
          {/each}
        </div>
        <Button
          type="button"
          variant="outline"
          class="external-link-button w-full sm:w-auto"
          onclick={() => openExternalUrl(cryptobotUrl)}
        >
          {$t.infoDonationCryptobotAction}
        </Button>
        <p class="text-xs leading-relaxed text-muted-foreground">
          {$t.infoDonationsFallback}
        </p>
      </div>
    </section>
  </CardContent>
</Card>

<style>
  :global(.copy-icon) {
    width: 0.95rem;
    height: 0.95rem;
    color: hsl(var(--muted-foreground));
  }

  :global(.copied-icon) {
    color: hsl(var(--primary));
  }

  :global(.wallet-address) {
    color: hsl(var(--muted-foreground));
    overflow-wrap: anywhere;
    word-break: break-word;
    font-family: ui-monospace, "Cascadia Code", "SFMono-Regular", Menlo, Monaco,
      Consolas, "Liberation Mono", "Courier New", monospace;
  }

  .wallet-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .about-meta {
    display: grid;
    gap: 0.4rem;
    padding-top: 0.4rem;
    border-top: 1px solid hsl(var(--border) / 0.7);
  }

  :global(.external-link-button) {
    cursor: pointer;
  }
</style>
