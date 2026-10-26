<script lang="ts">
  import { page } from "$app/state";
</script>

<main class="sign-in">
  <!-- prettier-ignore -->
  <div class="brand">Portal<span>67</span></div>
  <h1>Portal67 is temporarily unavailable</h1>
  <p>{page.error?.message}</p>
  <button onclick={() => window.location.reload()}>Retry</button>
  <a href="/sign-in">Return to sign in</a>
</main>
