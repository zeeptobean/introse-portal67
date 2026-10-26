<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import PaymentForm from "#lib/components/PaymentForm.svelte";
  import ChargeForm from "#lib/components/ChargeForm.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<a class="breadcrumb" href="/staff/tuition/students/{data.summary.student.id}?term={data.termId}">
  ← Student tuition
</a>
<PageHeader
  title={data.kind === "charges"
    ? "Add semester charge"
    : data.payment
      ? "Payment record"
      : "Record simulated payment"}
/>
<div class="form-width">
  <section class="identity-strip">
    <strong>{data.summary.student.full_name}</strong>
    <span>
      {data.summary.student.institutional_id ?? "ID not provided"} · {data.summary.student.email}
    </span>
    <span>{data.term.name} · {data.term.code}</span>
  </section>
  <p class="muted">Sample record maintenance. This does not initiate a money transfer.</p>
  {#key `${data.summary.student.id}-${data.termId}-${data.payment?.id ?? "new"}-${data.kind}`}
    {#if data.paymentForm}
      <PaymentForm
        summary={data.summary}
        term={data.term}
        payment={data.payment}
        initial={data.paymentForm}
      />
    {:else if data.chargeForm}
      <ChargeForm summary={data.summary} term={data.term} initial={data.chargeForm} />
    {/if}
  {/key}
</div>
