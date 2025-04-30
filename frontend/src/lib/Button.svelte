<script lang="ts">
  const { children, className = "", clickHandler = null } = $props();
  import { cn } from "./utils";
</script>

<button
  class={cn(
    "cursor-pointer bg-zinc-700 p-2 rounded-full flex justify-center items-center shadow-md hover:bg-zinc-800 transition active:scale-95 duration-100 ease-out",
    className,
  )}
  type="button"
  onclick={clickHandler}
>
  {@render children?.()}
</button>

<style></style>
