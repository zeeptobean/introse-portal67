<script lang="ts">
  import type { Snippet } from "svelte";
  let { title, subtitle, children }: { title: string; subtitle?: string; children?: Snippet } =
    $props();
</script>

<svelte:head><title>{title} · Portal67</title></svelte:head>
<header class="page-header">
  <div>
    <p class="eyebrow">Academic workspace</p>
    <h1>{title}</h1>
    {#if subtitle}
      <p class="muted">{subtitle}</p>
    {/if}
  </div>
  {#if children}
    <div class="header-actions">{@render children()}</div>
  {/if}
</header>

<style>
  .page-header {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 1px solid var(--line);
  }
  .page-header p:last-child {
    margin-bottom: 0;
  }
  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
  }
  @media (max-width: 900px) {
    .page-header {
      gap: 16px;
    }
    .header-actions {
      flex-shrink: 1;
    }
  }
  @media (max-width: 620px) {
    .page-header {
      display: block;
      margin-bottom: 24px;
      padding-bottom: 24px;
    }
    .page-header h1 {
      font-size: 26px;
    }
    .header-actions {
      margin-top: 24px;
    }
  }
</style>
