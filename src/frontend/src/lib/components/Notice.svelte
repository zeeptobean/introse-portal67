<script lang="ts">
  import type { Snippet } from "svelte";
  let {
    tone = "info",
    children,
  }: { tone?: "info" | "error" | "success" | "caution"; children: Snippet } = $props();
</script>

<div class="notice {tone}" role={tone === "error" ? "alert" : "status"}>{@render children()}</div>
