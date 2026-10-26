<script lang="ts">
  let { value }: { value: string } = $props();
  const labels: Record<string, string> = {
    no_charges: "No tuition charges",
    partial: "Partly paid",
    credit: "Credit balance",
    enrolled: "Registered",
  };
</script>

<span
  class="status"
  class:positive={["approved", "paid", "recorded", "enrolled"].includes(value)}
  class:negative={value === "rejected"}
  class:warning={["partial", "unpaid"].includes(value)}
>
  {labels[value] || value.replaceAll("_", " ")}
</span>

<style>
  .status {
    display: inline-flex;
    width: fit-content;
    align-items: center;
    background: var(--tint);
    color: var(--ink);
    padding: 4px 9px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
    text-transform: capitalize;
    line-height: 1.4;
  }
  .status.positive {
    background: #edf5f0;
    color: var(--success);
  }
  .status.negative {
    background: #fceff0;
    color: var(--error);
  }
  .status.warning {
    background: #fff4df;
    color: #7a4a0b;
  }
</style>
