<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import DocumentList from "#lib/components/DocumentList.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<PageHeader
  title="Document requests"
  subtitle="Review requests, and give students a definite answer."
/><DocumentList records={data.records} staff />
