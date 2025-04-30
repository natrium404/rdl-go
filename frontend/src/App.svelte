<script lang="ts">
  import "./app.css";

  import { EventsOn } from "../wailsjs/runtime/runtime.js";
  import { ScrapeWebPage, SaveFile } from "../wailsjs/go/main/App.js";

  import { ArrowDownToDotIcon, Search, SearchIcon, XIcon } from "lucide-svelte";
  import Button from "./lib/Button.svelte";
  import { onMount } from "svelte";
  import { scraper } from "../wailsjs/go/models";
  import { cn } from "./lib/utils";
  import { blur, fade, slide } from "svelte/transition";

  let url = $state("");
  let showURLForm = $state(true);
  let showDownload = $state(false);
  let isError = $state(false);

  let res = $state(scraper.Response.createFrom());
  let data = $derived(res.Data);

  let processing = $state(false);

  let logQueue: string[] = [];
  let logs: string[] = $state([]);

  onMount(() => {
    EventsOn("log", (message) => {
      logQueue.push(message);
      logDelay();
    });
  });

  const logDelay = async () => {
    if (processing) return;

    processing = true;

    while (logQueue.length > 0) {
      const nextLog = logQueue.shift() || "";

      await new Promise((resolve) => setTimeout(resolve, 200));
      logs.push(nextLog);
      if (logs.length > 3) {
        logs.shift();
        if (nextLog === "DONE") {
          showURLForm = !showURLForm;
        }
      }
    }

    processing = false;
  };

  const formSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    isError = false;

    res = await ScrapeWebPage(url);
    data = res.Data;
    if (res.Success) {
      url = "";
    }
    isError = !res.Success;
  };

  const setDynamicHeight = () => {
    const vh = window.innerHeight * 0.01;

    document.documentElement.style.setProperty("--vh", `${vh}px`);
  };

  window.addEventListener("resize", setDynamicHeight);
  window.addEventListener("load", setDynamicHeight);
</script>

<main class="max-w-2xl mx-auto h-svh flex justify-center items-center relative">
  <div
    class="bg-zinc-900 h-[calc(var(--vh,1vh)*100-4%)] w-3/4 rounded-md border border-zinc-800 p-1 overflow-hidden relative"
  >
    {#if !data?.Reel}
      <div class="w-full h-full flex justify-center items-center">
        <p class="font-bold text-xl">Opps... No reel found</p>
      </div>
    {:else}
      <video
        class="object-cover aspect-auto w-full h-full rounded-md"
        controls
        muted
        onloadstart={(e) => {
          const vid = e.target as HTMLVideoElement;
          vid.volume = 0.5;
        }}
      >
        <track kind="captions" />
        <source src={data.Reel} />
      </video>
    {/if}

    <div
      class="absolute top-0 right-0 translate-y-2 -translate-x-2 flex flex-col gap-1"
    >
      <!-- Search button -->

      {#if !showURLForm}
        <Button clickHandler={() => (showURLForm = !showURLForm)}>
          <SearchIcon />
        </Button>
      {/if}

      {#if data && data?.Reel}
        <!-- Download button -->
        <div class="relative">
          <Button clickHandler={() => (showDownload = !showDownload)}>
            <ArrowDownToDotIcon />
          </Button>

          {#if showDownload}
            <div
              class="absolute right-0 mt-2 max-w-md p-4 rounded-md bg-zinc-800 text-zinc-200 text-sm shadow-md"
            >
              <button
                class="bg-zinc-900 rounded-md p-1 border border-zinc-700 cursor-pointer hover:bg-zinc-950 active:scale-95 transition duration-100 ease-out"
                onclick={() => SaveFile(data.Reel, data.Code, ".mp4", "HD")}
                >Download Orginal</button
              >
              <div class="flex gap-6 mt-2">
                <div>
                  <p class="font-semibold">Videos</p>
                  <ul aria-label="videos" class="space-y-1 text-center">
                    {#each data.Videos as video}
                      <li aria-label="video-lable">
                        <button
                          class="underline cursor-pointer"
                          onclick={() =>
                            SaveFile(
                              video.URL,
                              data.Code,
                              video.MimeType.split("/")[1] || "mp4",
                              `${video.Quality}p_`,
                            )}>{video.Quality}</button
                        >
                      </li>
                    {/each}
                  </ul>
                </div>
                <div>
                  <p class="font-semibold">Audio</p>
                  <ul aria-label="audio" class="space-y-1 text-center">
                    <li class="list-item">
                      <button
                        class="underline cursor-pointer"
                        onclick={() =>
                          SaveFile(data.Audio.URL, data.Code, "mp3", "audio")}
                        >MP3</button
                      >
                    </li>
                  </ul>
                </div>
              </div>
              <p class="text-xs mt-3">
                * These are just seperated video and audio files.
              </p>
            </div>
          {/if}
        </div>
      {/if}
    </div>

    <!-- URL form -->
    {#if showURLForm}
      <div
        class="absolute flex top-0 left-0 justify-center items-center w-full h-full backdrop-blur-lg bg-black/85"
        role="dialog"
        in:blur
        out:blur
      >
        {#if url || data?.Reel}
          <Button
            className="absolute top-0 left-0 translate-x-2 translate-y-2"
            clickHandler={() => (showURLForm = !showURLForm)}
          >
            <XIcon />
          </Button>
        {/if}

        <!-- form -->
        <form class="w-4/5 relative" onsubmit={formSubmit}>
          <!-- Show logs form backend -->
          <div class="w-full flex max-w-xs relative items-end">
            <ul
              class="absolute list-disc ml-6 mb-2 *:transition-colors *:duration-700 transition duration-300 *:ease-out ease-in"
            >
              {#each logs as log}
                <li
                  class={cn(
                    "not-last:first:text-zinc-400/40 not-first:not-last:text-zinc-400/60 text-zinc-400 last:font-semibold",
                    isError && "text-red-800",
                  )}
                  in:slide
                  out:fade
                >
                  {log}
                </li>
              {/each}
            </ul>
          </div>

          <!-- URL input -->
          <div class="flex gap-1">
            <input
              bind:value={url}
              placeholder="Paste the reel url..."
              type="text"
              class={cn(
                "w-full rounded-md border bg-zinc-800  text-white placeholder:text-white ring-0 focus:border-blue-500 transition-colors duration-200 outline-none shadow-md focus:bg-zinc-900 border-zinc-700",
                isError && "border-red-500",
              )}
            />

            <button
              class={cn(
                "bg-zinc-800 px-2.5 rounded-full disabled:cursor-not-allowed cursor-pointer hover:bg-zinc-800/90",
                isError && "text-red-800",
              )}
              type="submit"
              disabled={!url}><Search /></button
            >
          </div>
        </form>
      </div>
    {/if}
  </div>
</main>

<style></style>
