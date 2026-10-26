<script lang="ts">
  import type { Person, Term } from "#lib/schemas/shared.ts";
  import { money } from "#lib/domain.ts";
  let {
    student,
    term,
    title,
    amount,
    before,
    after,
    description,
  }: {
    student: Person;
    term: Term;
    title: string;
    amount: string;
    before: number;
    after: number;
    description: string;
  } = $props();
</script>

<section class="review-box">
  <h2>{title}</h2>
  <p>{student.full_name} · {student.institutional_id ?? "ID not provided"}</p>
  <p>{term.name}</p>
  <dl class="review-values">
    <div>
      <dt>Amount</dt>
      <dd>{money(Number(amount))}</dd>
    </div>
    <div>
      <dt>Details</dt>
      <dd class="pre-wrap">{description || "No note"}</dd>
    </div>
    <div>
      <dt>Balance before</dt>
      <dd>{money(before)}</dd>
    </div>
    <div>
      <dt>Expected balance after</dt>
      <dd>{money(after)}{after < 0 ? " (credit)" : ""}</dd>
    </div>
  </dl>
</section>

<style>
  .review-values {
    margin: 0;
  }
  .review-values > div {
    display: grid;
    grid-template-columns: 160px minmax(0, 1fr);
    gap: 16px;
    padding: 16px 0;
    border-bottom: 1px solid var(--line);
  }
  @media (max-width: 620px) {
    .review-values > div {
      grid-template-columns: 1fr;
      gap: 8px;
    }
  }
</style>
