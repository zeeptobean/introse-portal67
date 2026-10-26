<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import TermSelect from "#lib/components/TermSelect.svelte";
  import TuitionLedger from "#lib/components/TuitionLedger.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<PageHeader title="Tuition" subtitle="A balance you can trace back to every record.">
  <TermSelect terms={data.terms} value={data.termId} />
</PageHeader><TuitionLedger summary={data.summary} termId={data.termId} />
