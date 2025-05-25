<script lang="ts">
  import "./app.css";

  import { EventsOn } from "../wailsjs/runtime/runtime.js";
  import {
    ScrapeWebPage,
    SaveFile,
    CheckForUpdate,
    Update,
  } from "../wailsjs/go/main/App.js";

  import {
    ArrowDownToDotIcon,
    Info,
    Loader,
    MousePointerClick,
    Search,
    SearchIcon,
    XIcon,
  } from "lucide-svelte";
  import Button from "./lib/Button.svelte";
  import { onMount } from "svelte";
  import { scraper, main } from "../wailsjs/go/models";
  import { cn, compareVersion } from "./lib/utils";
  import { blur, fade, slide } from "svelte/transition";

  let url = $state("");
  let showURLForm = $state(true);
  let showSearchButton = $state(true);
  let showDownload = $state(false);
  let isError = $state(false);

  let isUpdateAvailable = $state(false);
  let versionInfo = $state(main.VersionInfo.createFrom());
  let updating = $state(false);
  let updateText = $state("New version available.");

  let res = $state(scraper.Response.createFrom());
  let data = $derived(res.Data);

  let processing = $state(false);
  let videoRef: HTMLVideoElement | null = $state(null);

  let logQueue: string[] = [];
  let logs: string[] = $state([]);

  onMount(() => {
    EventsOn("log", (message) => {
      logQueue.push(message);
      logDelay();
    });

    EventsOn("progressbar", (message) => {
      showSearchButton = false;
      logQueue.push(message);
      logDelay(0);
    });
  });

  const checkForUpdate = async () => {
    versionInfo = await CheckForUpdate();
    const isLatest = compareVersion(
      versionInfo.current_version,
      versionInfo.latest_version,
    );

    if (isLatest < 0) {
      isUpdateAvailable = true;
    }
  };

  onMount(checkForUpdate);

  const update = async () => {
    updating = true;

    const err = await Update(versionInfo.asset);
    if (err === null) {
      updating = false;
      checkForUpdate();
      updateText = "New version updated. Restart the application.";
      return;
    }
  };

  const logDelay = async (delay: number = 200) => {
    if (processing) return;

    processing = true;

    while (logQueue.length > 0) {
      const nextLog = logQueue.shift() || "";

      await new Promise((resolve) => setTimeout(resolve, delay));
      logs.push(nextLog);
      if (logs.length > 3) {
        logs.shift();
        if (nextLog === "DOWNLOAD DONE") {
          showSearchButton = true;
        }
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
      videoRef?.load();
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
        {#if !showURLForm}
          <p class="font-bold text-xl">Opps... No reel found</p>
        {/if}
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
        bind:this={videoRef}
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
        class="absolute flex flex-col top-0 left-0 justify-between items-center w-full h-full backdrop-blur-lg bg-black/85"
        role="dialog"
        in:blur
        out:blur
      >
        <div>
          {#if url || data?.Reel}
            <Button
              className="absolute top-0 left-0 translate-x-2 translate-y-2"
              clickHandler={() => (showURLForm = !showURLForm)}
            >
              <XIcon />
            </Button>
          {/if}
        </div>

        <div class="w-full flex flex-col justify-center items-center">
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
                disabled={!url || !showSearchButton}><Search /></button
              >
            </div>
          </form>
          {#if isUpdateAvailable}
            <div class="text-zinc-400 mt-1.5 w-4/5 flex justify-start text-sm">
              <p class="flex gap-1 justify-center items-center">
                <Info class="size-4.5 stroke-blue-400" />
                <span>{updateText}</span>
              </p>
            </div>
          {/if}
        </div>
        <footer
          class="text-zinc-400 pb-2.5 text-xs flex flex-col items-center gap-1"
        >
          <div class="flex items-center gap-1">
            <p>
              <span>{versionInfo.current_version}</span>
            </p>
            <p>•</p>
            <p>Made with ❤️</p>
            <p>•</p>
            <p>Natrium</p>
            {#if isUpdateAvailable}
              <p>•</p>
              <p class="flex items-center gap-1 font-semibold underline group">
                <span>
                  {#if updating}
                    Updating...
                  {:else}
                    <button class="inline cursor-pointer" onclick={update}
                      >Update</button
                    >
                  {/if}
                </span>
                <span>
                  {#if updating}
                    <Loader class="animate-spin" />
                  {:else}
                    <MousePointerClick
                      class="w-4.5 group-hover:animate-bounce"
                    />
                  {/if}
                </span>
              </p>
            {/if}
          </div>
        </footer>
      </div>
    {/if}
  </div>
</main>

<style></style>
