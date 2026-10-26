<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import DocumentDetail from "#lib/components/DocumentDetail.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<a class="breadcrumb" href="/staff/documents">← Document requests</a>
<PageHeader
  title="Review request {data.record.id}"
  subtitle="Read the purpose, then give a clear decision."
/>
{#key data.record.id}
  <DocumentDetail record={data.record} staff decisionForm={data.form} />
{/key}
