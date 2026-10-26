<script lang="ts">
  import { page } from "$app/state";
  import PageHeader from "#lib/components/PageHeader.svelte";
  import TermSelect from "#lib/components/TermSelect.svelte";
  import TuitionLedger from "#lib/components/TuitionLedger.svelte";
  import Notice from "#lib/components/Notice.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<a class="breadcrumb" href="/staff/tuition">← Student lookup</a>
<PageHeader title="Student tuition" subtitle="Confirm identity before changing a record.">
  <TermSelect terms={data.terms} value={data.termId} />
</PageHeader>
{#if page.url.searchParams.get("saved")}
  <Notice tone="success">
    Record {page.url.searchParams.get("saved")} saved. The confirmed ledger is shown below.
  </Notice>
{/if}
<TuitionLedger summary={data.summary} termId={data.termId} staff />
