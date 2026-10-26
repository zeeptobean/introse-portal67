<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import DocumentList from "#lib/components/DocumentList.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<PageHeader title="Documents" subtitle="A clear record of your requests and decisions.">
  <a class="button primary" href="/student/documents/new">New request</a>
</PageHeader><DocumentList records={data.records} />
