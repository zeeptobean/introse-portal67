<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import DocumentDetail from "#lib/components/DocumentDetail.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<a class="breadcrumb" href="/student/documents">← Documents</a>
<PageHeader title="Request {data.record.id}" subtitle="Your request and its confirmed history." />
{#key data.record.id}
  <DocumentDetail record={data.record} />
{/key}
