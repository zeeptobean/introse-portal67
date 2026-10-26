<script lang="ts">
  import { page } from "$app/state";
  import { invalidateAll } from "$app/navigation";
  import PageHeader from "#lib/components/PageHeader.svelte";
</script>

<PageHeader
  title={page.status === 403
    ? "This page is not available for your account"
    : page.status === 404
      ? "We couldn’t find this record"
      : "Could not load this page"}
/>
<div class="panel">
  <p>{page.error?.message}</p>
  <div class="actions">
    <button onclick={() => invalidateAll()}>Retry</button>
    <a href="/">Return to your workspace</a>
  </div>
</div>
